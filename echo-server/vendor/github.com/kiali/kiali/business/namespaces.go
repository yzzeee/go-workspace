package business

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	osproject_v1 "github.com/openshift/api/project/v1"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/kiali/kiali/config"
	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/log"
	"github.com/kiali/kiali/models"
	"github.com/kiali/kiali/observability"
)

// NamespaceService deals with fetching k8sClients namespaces / OpenShift projects and convert to kiali model
type NamespaceService struct {
	userClients            map[string]kubernetes.ClientInterface
	kialiSAClients         map[string]kubernetes.ClientInterface
	homeClusterUserClient  kubernetes.ClientInterface
	hasProjects            bool
	isAccessibleNamespaces map[string]bool
}

type AccessibleNamespaceError struct {
	msg string
}

func (in *AccessibleNamespaceError) Error() string {
	return in.msg
}

func IsAccessibleError(err error) bool {
	_, isAccessibleError := err.(*AccessibleNamespaceError)
	return isAccessibleError
}

func NewNamespaceService(userClients map[string]kubernetes.ClientInterface, kialiSAClients map[string]kubernetes.ClientInterface) NamespaceService {
	var hasProjects bool
	conf := config.Get()

	homeClusterName := conf.KubernetesConfig.ClusterName
	if saClient, ok := kialiSAClients[homeClusterName]; ok && saClient.IsOpenShift() {
		hasProjects = true
	} else {
		hasProjects = false
	}

	ans := conf.Deployment.AccessibleNamespaces
	isAccessibleNamespaces := make(map[string]bool, len(ans))
	for _, ns := range ans {
		isAccessibleNamespaces[ns] = true
	}

	return NamespaceService{
		userClients:            userClients,
		kialiSAClients:         kialiSAClients,
		hasProjects:            hasProjects,
		homeClusterUserClient:  userClients[homeClusterName],
		isAccessibleNamespaces: isAccessibleNamespaces,
	}
}

// Returns a list of the given namespaces / projects
func (in *NamespaceService) GetNamespaces(ctx context.Context) ([]models.Namespace, error) {
	var end observability.EndFunc
	_, end = observability.StartSpan(ctx, "GetNamespaces",
		observability.Attribute("package", "business"),
	)
	defer end()

	if kialiCache != nil && in.homeClusterUserClient != nil {
		if ns := kialiCache.GetNamespaces(in.homeClusterUserClient.GetToken()); ns != nil {
			return ns, nil
		}
	}

	conf := config.Get()

	// determine what the discoverySelectors are by examining the Istio ConfigMap
	var discoverySelectors []*meta_v1.LabelSelector
	if kialiCache != nil {
		if icm, err := kialiCache.GetConfigMap(conf.IstioNamespace, conf.ExternalServices.Istio.ConfigMapName); err == nil {
			if ic, err2 := kubernetes.GetIstioConfigMap(icm); err2 == nil {
				discoverySelectors = ic.DiscoverySelectors
			} else {
				log.Errorf("Will not process discoverySelectors due to a failure to get the Istio ConfigMap: %v", err2)
			}
		} else {
			log.Errorf("Will not process discoverySelectors due to a failure to parse the Istio ConfigMap: %v", err)
		}
	}

	if len(discoverySelectors) > 0 {
		log.Tracef("Istio discovery selectors: %+v", discoverySelectors)
	} else {
		log.Tracef("No Istio discovery selectors defined.")
	}

	// Let's explain the four different filters along with accessible namespaces (aka AN).
	//
	// First, we look at AN. AN is either ["**"] or it is not.
	//
	// If AN is ["**"], then the entire cluster of namespaces is accessible to Kiali.
	// In this case, the user can further filter what namespaces this function should return using both includes and excludes.
	// 1. LabelSelectorInclude is used to obtain an initial set of namespaces, if specified.
	// 2. Added to that initial list will be the namespaces named in the Include list, if those namespaces actually exist.
	// 3. If no LabelSelectorInclude or Include list is specified, then all namespaces are in the list.
	// 4. Remove from that list those namespaces that match LabelSelectorExclude, as well as those namespaces found in the Exclude list.
	// (Side note: You might ask: Why have an Include list when we already have the AN list?
	// The difference is if you specify AN (not ["**"]), only those namespaces that exist __at install time__ will get a Role
	// and hence are accessible to Kiali. The Include list is evaluated at the time this function is called, thus it
	// allows Kiali to see those namespaces even if they are created after Kiali is installed).
	//
	// If AN is not ["**"], then only a subset of namespaces in the cluster is accessible to Kiali.
	// When installed by the operator, Kiali will be given access to a set of namespaces (as defined in AN) via Roles that
	// are created by the operator. Those namespaces that Kiali has access to (as defined in AN) will be labeled with the label
	// selector defined in LabelSelectorInclude (Kiali CR "spec.api.namespaces.label_selector_include").
	// 1. All of those namespaces are retrieved with the LabelSelectorInclude to obtain a set of namespaces.
	// 2. Remove from that list those namespaces that match LabelSelectorExclude, as well as those namespaces found in the Exclude list.
	// The Include option is ignored in this case - you cannot Include more namespaces over and above what AN specifies.
	// (Side note 1: It probably doesn't make sense to set LabelSelectorExclude and Excludes when AN is not ["**"]. This is because
	// you already have defined what namespaces you want to give Kiali access to (the AN list itself). However, for consistency,
	// this function will still use those additional filters to filter out namespaces. So it is possible this function returns
	// a subset of namespaces that are listed in AN.)
	// (Side note 2: Notice the difference here between when AN is set to ["**"] and when it is not. When AN is not set to ["**"],
	// LabelSelectorInclude does not tell the operator which namespaces are included - AN does that. Instead, the operator will
	// create that label as defined by LabelSelectorInclude on each namespace defined in AN. Thus, after the operator installs
	// Kiali, the Kiali Server can then use LabelSelectorInclude in this function in order to select all namespaces as defined in AN.
	// If installed via the server helm chart, none of that is done, and it is up to the user to ensure
	// LabelSelectorInclude (if defined) selects all namespaces in AN. It is a user-error if they do not configure that correctly.
	// The server helm chart will not assume they did it correctly. The user therefore normally should not set LabelSelectorInclude
	// if they also set AN to something not ["**"]. This is one reason why we recommend using the Kiali operator, and why we say
	// the server helm chart is only provided as a convenience.)
	// (Side note 3: The control plane namespace is always included via api.namespaces.include and
	// never excluded via api.namespaces.exclude or api.namespaces.label_selector_exclude.)

	// determine if we are to exclude namespaces by label - if so, set the label name and value for use later
	labelSelectorExclude := conf.API.Namespaces.LabelSelectorExclude
	var labelSelectorExcludeName string
	var labelSelectorExcludeValue string
	if labelSelectorExclude != "" {
		excludeLabelList := strings.Split(labelSelectorExclude, "=")
		if len(excludeLabelList) != 2 {
			return nil, fmt.Errorf("api.namespaces.label_selector_exclude is invalid: %v", labelSelectorExclude)
		}
		labelSelectorExcludeName = excludeLabelList[0]
		labelSelectorExcludeValue = excludeLabelList[1]
	}

	namespaces := []models.Namespace{}

	var clusterNames []string
	for c := range in.userClients {
		clusterNames = append(clusterNames, c)
	}

	wg := &sync.WaitGroup{}
	type result struct {
		cluster string
		ns      []models.Namespace
		err     error
	}
	resultsCh := make(chan result)

	// TODO: Use a context to define a timeout. The context should be passed to the k8s client
	go func() {
		for _, cluster := range clusterNames {
			wg.Add(1)
			go func(c string) {
				defer wg.Done()
				list, error := in.getNamespacesByCluster(c)
				if error != nil {
					resultsCh <- result{cluster: c, ns: nil, err: error}
				} else {
					resultsCh <- result{cluster: c, ns: list, err: nil}
				}
			}(cluster)
		}
		wg.Wait()
		close(resultsCh)
	}()

	// Combine namespace data
	for resultCh := range resultsCh {
		if resultCh.err != nil {
			if resultCh.cluster == conf.KubernetesConfig.ClusterName {
				log.Errorf("Error fetching Namespaces for local cluster [%s]: %s", resultCh.cluster, resultCh.err)
				return nil, resultCh.err
			} else {
				log.Infof("Error fetching Namespaces for cluster [%s]: %s", resultCh.cluster, resultCh.err)
				continue
			}
		}
		namespaces = append(namespaces, resultCh.ns...)
	}

	resultns := namespaces

	// Filter out those namespaces that do not match discoverySelectors.
	// Follow the semantics that Istio follows, which is:
	//   If there is no discoverySelectors section in the config, skip this entirely.
	//   If there is an empty discoverySelectors section, that means all namespaces are to be used.
	//   If there are one or more discoverySelectors specified, the filter namespaces based on what they select.
	if len(discoverySelectors) > 0 {
		// 1. convert LabelSelectors to Selectors
		selectors := make([]labels.Selector, 0)
		for _, selector := range discoverySelectors {
			ls, err := meta_v1.LabelSelectorAsSelector(selector)
			if err != nil {
				return nil, fmt.Errorf("error initializing discovery selectors filter, invalid discovery selector: %v", err)
			}
			selectors = append(selectors, ls)
		}

		// 2. range over all namespaces to get discovery namespaces, notice each selector result is ORed (as per Istio convention)
		selectedNamespaces := make([]models.Namespace, 0)
		for _, ns := range resultns {
			if ns.Name == conf.IstioNamespace {
				selectedNamespaces = append(selectedNamespaces, ns) // we always want to return the control plane namespace
			} else {
				for _, selector := range selectors {
					if selector.Matches(labels.Set(ns.Labels)) {
						selectedNamespaces = append(selectedNamespaces, ns)
						break
					}
				}
			}
		}
		namespaces = selectedNamespaces
		resultns = namespaces
	}

	// exclude namespaces that are:
	// 1. to be filtered out via the exclude list
	// 2. to be filtered out via the label selector
	// Note that the control plane namespace is never excluded
	excludes := conf.API.Namespaces.Exclude
	if len(excludes) > 0 || labelSelectorExclude != "" {
		resultns = []models.Namespace{}
	NAMESPACES:
		for _, namespace := range namespaces {
			if namespace.Name != conf.IstioNamespace {
				if len(excludes) > 0 {
					for _, excludePattern := range excludes {
						if match, _ := regexp.MatchString(excludePattern, namespace.Name); match {
							continue NAMESPACES
						}
					}
				}
				if labelSelectorExclude != "" {
					if namespace.Labels[labelSelectorExcludeName] == labelSelectorExcludeValue {
						continue NAMESPACES
					}
				}
			}
			resultns = append(resultns, namespace)
		}
	}

	// store only the filtered set of namespaces in cache for the token
	if kialiCache != nil && in.homeClusterUserClient != nil {
		// just get the home cluster token because it is assumed tokens are identical across all clusters
		kialiCache.SetNamespaces(in.homeClusterUserClient.GetToken(), resultns)
	}

	return resultns, nil
}

func (in *NamespaceService) getNamespacesByCluster(cluster string) ([]models.Namespace, error) {
	configObject := config.Get()

	labelSelectorInclude := configObject.API.Namespaces.LabelSelectorInclude

	namespaces := []models.Namespace{}
	_, queryAllNamespaces := in.isAccessibleNamespaces["**"]
	// If we are running in OpenShift, we will use the project names since these are the list of accessible namespaces
	if in.hasProjects {
		projects, err2 := in.userClients[cluster].GetProjects(labelSelectorInclude)
		if err2 == nil {
			// Everything is good, return the projects we got from OpenShift
			if queryAllNamespaces {
				namespaces = models.CastProjectCollection(projects, cluster)
				// add the namespaces explicitly included in the include list.
				includes := configObject.API.Namespaces.Include
				if len(includes) > 0 {
					var allNamespaces []models.Namespace
					var seedNamespaces []models.Namespace

					if labelSelectorInclude == "" {
						// we have already retrieved all the namespaces, but we want only those in the Include list
						allNamespaces = namespaces
						seedNamespaces = make([]models.Namespace, 0)
					} else {
						// we have already got those namespaces that match the LabelSelectorInclude - that is our seed list.
						// but we need ALL namespaces so we can look for more that match the Include list.
						if allProjects, err := in.userClients[cluster].GetProjects(""); err != nil {
							return nil, err
						} else {
							allNamespaces = models.CastProjectCollection(allProjects, cluster)
							seedNamespaces = namespaces
						}
					}
					namespaces = in.addIncludedNamespaces(allNamespaces, seedNamespaces)
				}
			} else {
				filteredProjects := make([]osproject_v1.Project, 0)
				for _, project := range projects {
					if _, isAccessible := in.isAccessibleNamespaces[project.Name]; isAccessible {
						filteredProjects = append(filteredProjects, project)
					}
				}
				namespaces = models.CastProjectCollection(filteredProjects, cluster)
			}
		}
	} else {
		// if the accessible namespaces define a distinct list of namespaces, use only those.
		// If accessible namespaces include the special "**" (meaning all namespaces) ask k8sClients for them.
		// Note that "**" requires cluster role permission to list all namespaces.
		accessibleNamespaces := configObject.Deployment.AccessibleNamespaces
		if queryAllNamespaces {

			nss, err := in.userClients[cluster].GetNamespaces(labelSelectorInclude)
			if err != nil {
				// Fallback to using the Kiali service account, if needed
				if errors.IsForbidden(err) {
					if nss, err = in.getNamespacesUsingKialiSA(cluster, labelSelectorInclude, err); err != nil {
						return nil, err
					}
				} else {
					return nil, err
				}
			}

			namespaces = models.CastNamespaceCollection(nss, cluster)

			// add the namespaces explicitly included in the includes list.
			includes := configObject.API.Namespaces.Include
			if len(includes) > 0 {
				var allNamespaces []models.Namespace
				var seedNamespaces []models.Namespace

				if labelSelectorInclude == "" {
					// we have already retrieved all the namespaces, but we want only those in the Include list
					allNamespaces = namespaces
					seedNamespaces = make([]models.Namespace, 0)
				} else {
					// we have already got those namespaces that match the LabelSelectorInclude - that is our seed list.
					// but we need ALL namespaces so we can look for more that match the Include list.
					allK8sNamespaces, errGetNs := in.userClients[cluster].GetNamespaces("")
					if errGetNs != nil {
						// Fallback to using the Kiali service account, if needed
						if errors.IsForbidden(errGetNs) {
							if allK8sNamespaces, errGetNs = in.getNamespacesUsingKialiSA(cluster, "", errGetNs); errGetNs != nil {
								return nil, errGetNs
							}
						} else {
							return nil, errGetNs
						}
					}
					allNamespaces = models.CastNamespaceCollection(allK8sNamespaces, cluster)
					seedNamespaces = namespaces
				}
				namespaces = in.addIncludedNamespaces(allNamespaces, seedNamespaces)
			}
		} else {
			k8sNamespaces := make([]core_v1.Namespace, 0)
			for _, ans := range accessibleNamespaces {
				k8sNs, err := in.userClients[cluster].GetNamespace(ans)
				if err != nil {
					if errors.IsNotFound(err) {
						// If a namespace is not found, then we skip it from the list of namespaces
						log.Warningf("Kiali has an accessible namespace [%s] which doesn't exist", ans)
					} else if errors.IsForbidden(err) {
						// Also, if namespace isn't readable, skip it.
						log.Warningf("Kiali has an accessible namespace [%s] which is forbidden", ans)
					} else {
						// On any other error, abort and return the error.
						return nil, err
					}
				} else {
					k8sNamespaces = append(k8sNamespaces, *k8sNs)
				}
			}
			namespaces = models.CastNamespaceCollection(k8sNamespaces, cluster)
		}
	}

	return namespaces, nil
}

// GetNamespacesForCluster is just a convenience routine that filters GetNamespaces for a particular cluster
func (in *NamespaceService) GetNamespacesForCluster(ctx context.Context, cluster string) ([]models.Namespace, error) {
	tokenNamespaces, err := in.GetNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	clusterNamespaces := []models.Namespace{}
	for _, ns := range tokenNamespaces {
		if ns.Cluster == cluster {
			clusterNamespaces = append(clusterNamespaces, ns)
		}
	}

	return clusterNamespaces, nil
}

// addIncludedNamespaces will look at all the namespaces and return all of them that match the Include list.
// The returned results will be guaranteed to include the namespaces found in the given seed list.
// There will be no duplicate namespaces in the returned list.
func (in *NamespaceService) addIncludedNamespaces(all []models.Namespace, seed []models.Namespace) []models.Namespace {
	var controlPlaneNamespace models.Namespace
	hasNamespace := make(map[string]bool, len(seed))
	results := make([]models.Namespace, 0, len(seed))
	configObject := config.Get()

	// seed with the initial set of namespaces - this ensures there are no duplicates in the seed list
	for _, ns := range seed {
		if _, exists := hasNamespace[ns.Name]; !exists {
			hasNamespace[ns.Name] = true
			results = append(results, ns)
		}
	}

	// go through the list of all namespaces and add to the results list those that match a regex found in the Include list
	includes := configObject.API.Namespaces.Include
NAMESPACES:
	for _, ns := range all {
		if _, exists := hasNamespace[ns.Name]; exists {
			continue
		}
		for _, includePattern := range includes {
			if match, _ := regexp.MatchString(includePattern, ns.Name); match {
				hasNamespace[ns.Name] = true
				results = append(results, ns)
				continue NAMESPACES
			}
		}
		if ns.Name == configObject.IstioNamespace {
			controlPlaneNamespace = ns // squirrel away the control plane namepace in case we need to add it
		}
	}

	// Kiali needs the control plane namespace, so it should always be included.
	// If the user did not configure the include list to explicitly include the control plane namespace, then we need to include it now.
	if _, exists := hasNamespace[configObject.IstioNamespace]; !exists {
		if controlPlaneNamespace.Name != "" {
			results = append(results, controlPlaneNamespace)
		} else {
			log.Errorf("Kiali needs to include the control plane namespace. Make sure you configured Kiali so it can access and include the namespace [%s].", configObject.IstioNamespace)
		}
	}
	return results
}

func (in *NamespaceService) isAccessibleNamespace(namespace string) bool {
	_, queryAllNamespaces := in.isAccessibleNamespaces["**"]
	if queryAllNamespaces {
		return true
	}
	_, isAccessible := in.isAccessibleNamespaces[namespace]
	return isAccessible
}

func (in *NamespaceService) isExcludedNamespace(namespace string) bool {
	configObject := config.Get()
	excludes := configObject.API.Namespaces.Exclude
	if len(excludes) == 0 {
		return false
	}
	if namespace == configObject.IstioNamespace {
		return false // the control plane namespace is never excluded
	}
	for _, excludePattern := range excludes {
		if match, _ := regexp.MatchString(excludePattern, namespace); match {
			return true
		}
	}
	return false
}

func (in *NamespaceService) isIncludedNamespace(namespace string) bool {
	_, queryAllNamespaces := in.isAccessibleNamespaces["**"]
	if !queryAllNamespaces {
		return true // Include list is ignored if accessible namespaces is not **; for our purposes, when ignored we assume the Include list includes all.
	}

	configObject := config.Get()
	if namespace == configObject.IstioNamespace {
		return true // the control plane namespace is always included
	}

	includes := configObject.API.Namespaces.Include
	if len(includes) == 0 {
		return true // if no Include list is specified, all namespaces are included
	}
	for _, includePattern := range includes {
		if match, _ := regexp.MatchString(includePattern, namespace); match {
			return true
		}
	}
	return false
}

// GetNamespace returns the definition of the specified namespace.
// TODO: Multicluster: We are going to need something else to identify the namespace, the cluster (OR Return a list/array/map)
func (in *NamespaceService) GetNamespace(ctx context.Context, namespace string) (*models.Namespace, error) {
	// TODO: Wrapper for MC while other services are not updated to propagate the cluster
	return in.GetNamespaceByCluster(ctx, namespace, "")
}

// GetNamespace returns the definition of the specified namespace.
// TODO: Multicluster: We are going to need something else to identify the namespace, the cluster (OR Return a list/array/map)
func (in *NamespaceService) GetNamespaceByCluster(ctx context.Context, namespace string, cluster string) (*models.Namespace, error) {
	var end observability.EndFunc
	ctx, end = observability.StartSpan(ctx, "GetNamespaceByCluster",
		observability.Attribute("package", "business"),
		observability.Attribute("namespace", namespace),
		observability.Attribute("cluster", cluster),
	)
	defer end()

	var err error

	// Cache already has included/excluded namespaces applied
	if kialiCache != nil && in.homeClusterUserClient != nil {
		if ns := kialiCache.GetNamespace(in.homeClusterUserClient.GetToken(), namespace, cluster); ns != nil {
			return ns, nil
		}
	}

	if !in.isAccessibleNamespace(namespace) {
		return nil, &AccessibleNamespaceError{msg: "Namespace [" + namespace + "] is not accessible for Kiali"}
	}

	if !in.isIncludedNamespace(namespace) {
		return nil, &AccessibleNamespaceError{msg: "Namespace [" + namespace + "] is not included for Kiali"}
	}

	if in.isExcludedNamespace(namespace) {
		return nil, &AccessibleNamespaceError{msg: "Namespace [" + namespace + "] is excluded for Kiali"}
	}

	var result models.Namespace
	if in.hasProjects {
		var project *osproject_v1.Project
		// TODO: MC
		if cluster == "" {
			var err2 error
			for cl := range in.userClients {
				project, err2 = in.userClients[cl].GetProject(namespace)
				if err2 == nil {
					result = models.CastProject(*project, cluster)
					break
				}
			}
			if err2 != nil {
				return nil, err2
			}
		} else {
			if _, ok := in.userClients[cluster]; !ok {
				return nil, fmt.Errorf("OCP Cluster [%s] is not found or is not accessible for Kiali", cluster)
			}
			project, errC := in.userClients[cluster].GetProject(namespace)
			if errC != nil {
				return nil, errC
			}
			result = models.CastProject(*project, cluster)
		}
	} else {
		// TODO: MC
		var ns *core_v1.Namespace
		var errC error
		if cluster == "" {
			for cl := range in.userClients {
				ns, errC = in.userClients[cl].GetNamespace(namespace)
				if errC == nil {
					// Namespace found, assign that cluster
					cluster = cl
					break
				}
			}
			if errC != nil {
				return nil, errC
			}
		} else {
			if _, ok := in.userClients[cluster]; !ok {
				return nil, fmt.Errorf("Cluster [%s] is not found or is not accessible for Kiali", cluster)
			}
			ns, errC = in.userClients[cluster].GetNamespace(namespace)
			if errC != nil {
				return nil, errC
			}
		}

		result = models.CastNamespace(*ns, cluster)
	}
	// Refresh cache in case of cache expiration
	if kialiCache != nil {
		if _, err = in.GetNamespaces(ctx); err != nil {
			return nil, err
		}
	}
	return &result, nil
}

func (in *NamespaceService) UpdateNamespace(ctx context.Context, namespace string, jsonPatch string, cluster string) (*models.Namespace, error) {
	var end observability.EndFunc
	ctx, end = observability.StartSpan(ctx, "UpdateNamespace",
		observability.Attribute("package", "business"),
		observability.Attribute("namespace", namespace),
		observability.Attribute("jsonPatch", jsonPatch),
	)
	defer end()

	// A first check to run the accessible/excluded logic and not run the Update operation on filtered namespaces
	_, err := in.GetNamespaceByCluster(ctx, namespace, cluster)
	if err != nil {
		return nil, err
	}

	_, err = in.userClients[cluster].UpdateNamespace(namespace, jsonPatch)
	if err != nil {
		return nil, err
	}

	// Cache is stopped after a Create/Update/Delete operation to force a refresh
	if kialiCache != nil && err == nil {
		kialiCache.Refresh(namespace)
		kialiCache.RefreshTokenNamespaces()
	}
	// Call GetNamespace to update the caching
	return in.GetNamespaceByCluster(ctx, namespace, cluster)
}

func (in *NamespaceService) getNamespacesUsingKialiSA(cluster string, labelSelector string, forwardedError error) ([]core_v1.Namespace, error) {
	// Check if we already are using the Kiali ServiceAccount token. If we are, no need to do further processing, since
	// this would just circle back to the same results.
	kialiToken := in.kialiSAClients[cluster].GetToken()
	if in.userClients[cluster].GetToken() == kialiToken {
		return nil, forwardedError
	}

	// Let's get the namespaces list using the Kiali Service Account
	nss, err := in.kialiSAClients[cluster].GetNamespaces(labelSelector)
	if err != nil {
		return nil, err
	}

	// Only take namespaces where the user has privileges
	var namespaces []core_v1.Namespace
	for _, item := range nss {
		if _, getNsErr := in.userClients[cluster].GetNamespace(item.Name); getNsErr == nil {
			// Namespace is accessible
			namespaces = append(namespaces, item)
		} else if !errors.IsForbidden(getNsErr) {
			// Since the returned error is NOT "forbidden", something bad happened
			return nil, getNsErr
		}
	}

	// Return the list of namespaces where the user has the 'get namespace' read privilege.
	return namespaces, nil
}

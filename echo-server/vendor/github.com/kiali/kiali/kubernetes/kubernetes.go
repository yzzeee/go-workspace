package kubernetes

import (
	"bytes"
	"context"
	goerrors "errors"
	"fmt"
	"io"
	"time"

	osapps_v1 "github.com/openshift/api/apps/v1"
	osproject_v1 "github.com/openshift/api/project/v1"
	osroutes_v1 "github.com/openshift/api/route/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	apps_v1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/authentication/v1"
	auth_v1 "k8s.io/api/authorization/v1"
	batch_v1 "k8s.io/api/batch/v1"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/kiali/kiali/config"
	"github.com/kiali/kiali/log"
	"github.com/kiali/kiali/observability"
	"github.com/kiali/kiali/util/httputil"
)

type K8SClientInterface interface {
	// Kube returns the underlying kubernetes client.
	Kube() kubernetes.Interface
	GetClusterServicesByLabels(labelsSelector string) ([]core_v1.Service, error)
	GetConfigMap(namespace, name string) (*core_v1.ConfigMap, error)
	GetCronJobs(namespace string) ([]batch_v1.CronJob, error)
	GetDaemonSet(namespace string, name string) (*apps_v1.DaemonSet, error)
	GetDaemonSets(namespace string) ([]apps_v1.DaemonSet, error)
	GetDeployment(namespace string, name string) (*apps_v1.Deployment, error)
	GetDeployments(namespace string) ([]apps_v1.Deployment, error)
	GetDeploymentConfig(namespace string, name string) (*osapps_v1.DeploymentConfig, error)
	GetDeploymentConfigs(namespace string) ([]osapps_v1.DeploymentConfig, error)
	GetEndpoints(namespace string, name string) (*core_v1.Endpoints, error)
	GetJobs(namespace string) ([]batch_v1.Job, error)
	GetNamespace(namespace string) (*core_v1.Namespace, error)
	GetNamespaces(labelSelector string) ([]core_v1.Namespace, error)
	GetPod(namespace, name string) (*core_v1.Pod, error)
	GetPods(namespace, labelSelector string) ([]core_v1.Pod, error)
	GetReplicationControllers(namespace string) ([]core_v1.ReplicationController, error)
	GetReplicaSets(namespace string) ([]apps_v1.ReplicaSet, error)
	GetSecret(namespace, name string) (*core_v1.Secret, error)
	GetSelfSubjectAccessReview(ctx context.Context, namespace, api, resourceType string, verbs []string) ([]*auth_v1.SelfSubjectAccessReview, error)
	GetService(namespace string, name string) (*core_v1.Service, error)
	GetServices(namespace string, selectorLabels map[string]string) ([]core_v1.Service, error)
	GetServicesByLabels(namespace string, labelsSelector string) ([]core_v1.Service, error)
	GetStatefulSet(namespace string, name string) (*apps_v1.StatefulSet, error)
	GetStatefulSets(namespace string) ([]apps_v1.StatefulSet, error)
	GetTokenSubject(authInfo *api.AuthInfo) (string, error)
	StreamPodLogs(namespace, name string, opts *core_v1.PodLogOptions) (io.ReadCloser, error)
	UpdateNamespace(namespace string, jsonPatch string) (*core_v1.Namespace, error)
	UpdateService(namespace string, name string, jsonPatch string, patchType string) error
	UpdateWorkload(namespace string, name string, workloadType string, jsonPatch string, patchType string) error
}

type OSClientInterface interface {
	GetProject(project string) (*osproject_v1.Project, error)
	GetProjects(labelSelector string) ([]osproject_v1.Project, error)
	GetRoute(namespace string, name string) (*osroutes_v1.Route, error)
	UpdateProject(project string, jsonPatch string) (*osproject_v1.Project, error)
}

func (in *K8SClient) forwardGetRequest(namespace, podName string, destinationPort int, path string) ([]byte, error) {
	localPort := httputil.Pool.GetFreePort()
	defer httputil.Pool.FreePort(localPort)

	f, err := in.getPodPortForwarder(namespace, podName, fmt.Sprintf("%d:%d", localPort, destinationPort))
	if err != nil {
		return nil, err
	}

	// Start the forwarding
	if err := f.Start(); err != nil {
		return nil, err
	}

	// Defering the finish of the port-forwarding
	defer f.Stop()

	// Ready to create a request
	resp, code, _, err := httputil.HttpGet(fmt.Sprintf("http://localhost:%d%s", localPort, path), nil, 10*time.Second, nil, nil)
	if code >= 400 {
		return resp, fmt.Errorf("error fetching %s from %s/%s. Response code: %d", path, namespace, podName, code)
	}

	return resp, err
}

func (in *K8SClient) Kube() kubernetes.Interface {
	return in.k8s
}

// GetClusterServicesByLabels fetches and returns all services in the whole cluster
// that match the optional labelSelector. This is using the cluster-wide call
// to fetch the services. The client will need to be created with an account that
// has cluster-wide privileges to list services.
func (in *K8SClient) GetClusterServicesByLabels(labelsSelector string) ([]core_v1.Service, error) {
	selector := meta_v1.ListOptions{LabelSelector: labelsSelector}
	if allServicesList, err := in.k8s.CoreV1().Services("").List(in.ctx, selector); err == nil {
		return allServicesList.Items, nil
	} else {
		return []core_v1.Service{}, err
	}
}

// GetConfigMap fetches and returns the specified ConfigMap definition
// from the cluster
func (in *K8SClient) GetConfigMap(namespace, name string) (*core_v1.ConfigMap, error) {
	configMap, err := in.k8s.CoreV1().ConfigMaps(namespace).Get(in.ctx, name, emptyGetOptions)
	if err != nil {
		return &core_v1.ConfigMap{}, err
	}

	return configMap, nil
}

// GetNamespace fetches and returns the specified namespace definition
// from the cluster
func (in *K8SClient) GetNamespace(namespace string) (*core_v1.Namespace, error) {
	ns, err := in.k8s.CoreV1().Namespaces().Get(in.ctx, namespace, emptyGetOptions)
	if err != nil {
		return &core_v1.Namespace{}, err
	}

	return ns, nil
}

// GetServerVersion fetches and returns information about the version Kubernetes that is running
func (in *K8SClient) GetServerVersion() (*version.Info, error) {
	return in.k8s.Discovery().ServerVersion()
}

// GetNamespaces returns a list of all namespaces of the cluster.
// It returns a list of all namespaces of the cluster.
// It returns an error on any problem.
func (in *K8SClient) GetNamespaces(labelSelector string) ([]core_v1.Namespace, error) {
	var listOptions meta_v1.ListOptions

	// Apply labelSelector filtering if specified
	if labelSelector != "" {
		listOptions = meta_v1.ListOptions{LabelSelector: labelSelector}
	} else {
		listOptions = emptyListOptions
	}

	namespaces, err := in.k8s.CoreV1().Namespaces().List(in.ctx, listOptions)
	if err != nil {
		return nil, err
	}

	return namespaces.Items, nil
}

// GetProject fetches and returns the definition of the project with
// the specified name by querying the cluster API. GetProject will fail
// if the underlying cluster is not Openshift.
func (in *K8SClient) GetProject(name string) (*osproject_v1.Project, error) {
	result := &osproject_v1.Project{}

	err := in.k8s.Discovery().RESTClient().Get().Prefix("apis", "project.openshift.io", "v1", "projects", name).Do(in.ctx).Into(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (in *K8SClient) GetProjects(labelSelector string) ([]osproject_v1.Project, error) {
	result := &osproject_v1.ProjectList{}

	request := in.k8s.Discovery().RESTClient().Get().Prefix("apis", "project.openshift.io", "v1", "projects")

	// Apply label selector filtering if specified
	if labelSelector != "" {
		request.Param("labelSelector", labelSelector)
	}

	err := request.Do(in.ctx).Into(result)
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

func (in *K8SClient) IsOpenShift() bool {
	if in.isOpenShift == nil {
		isOpenShift := false
		_, err := in.k8s.Discovery().RESTClient().Get().AbsPath("/apis/project.openshift.io").Do(in.ctx).Raw()
		if err == nil {
			isOpenShift = true
		}
		in.isOpenShift = &isOpenShift
	}
	return *in.isOpenShift
}

func (in *K8SClient) IsGatewayAPI() bool {
	if in.GatewayAPI() == nil {
		return false
	}
	if in.isGatewayAPI == nil {
		isGatewayAPI := false
		_, err := in.k8s.Discovery().RESTClient().Get().AbsPath("/apis/gateway.networking.k8s.io").Do(in.ctx).Raw()
		if err == nil {
			isGatewayAPI = true
		} else if !errors.IsNotFound(err) {
			log.Warningf("Error checking Kubernetes Gateway API configuration: %v", err)
		}
		in.isGatewayAPI = &isGatewayAPI
	}
	return *in.isGatewayAPI
}

// Is IstioAPI checks whether Istio API is installed or not
func (in *K8SClient) IsIstioAPI() bool {
	if in.Istio() == nil {
		return false
	}
	if in.isIstioAPI == nil {
		isIstioAPI := false
		_, err := in.k8s.Discovery().RESTClient().Get().AbsPath("/apis/networking.istio.io").Do(in.ctx).Raw()
		if err == nil {
			isIstioAPI = true
		} else if !errors.IsNotFound(err) {
			log.Warningf("Error checking Istio API configuration: %v", err)
		}
		in.isIstioAPI = &isIstioAPI
	}
	return *in.isIstioAPI
}

// GetServices returns a list of services for a given namespace.
// If selectorLabels is defined the list of services is filtered for those that matches Services selector labels.
// It returns an error on any problem.
// NOTE: The selectorLabels argument is NOT to find services matching the given labels. Assume selectorLabels are
// the labels of a Deployment. If this imaginary Deployment is selected by the Service (because of its Selector), then
// that service is returned; else it's omitted.
func (in *K8SClient) GetServices(namespace string, selectorLabels map[string]string) ([]core_v1.Service, error) {
	var allServices []core_v1.Service

	if allServicesList, err := in.k8s.CoreV1().Services(namespace).List(in.ctx, emptyListOptions); err == nil {
		allServices = allServicesList.Items
	} else {
		return []core_v1.Service{}, err
	}

	if selectorLabels == nil {
		return allServices, nil
	}
	var services []core_v1.Service
	for _, svc := range allServices {
		svcSelector := labels.Set(svc.Spec.Selector).AsSelector()
		if !svcSelector.Empty() && svcSelector.Matches(labels.Set(selectorLabels)) {
			services = append(services, svc)
		}
	}
	return services, nil
}

func (in *K8SClient) GetServicesByLabels(namespace string, labelsSelector string) ([]core_v1.Service, error) {
	selector := meta_v1.ListOptions{LabelSelector: labelsSelector}
	if allServicesList, err := in.k8s.CoreV1().Services(namespace).List(in.ctx, selector); err == nil {
		return allServicesList.Items, nil
	} else {
		return []core_v1.Service{}, err
	}
}

func (in *K8SClient) GetDaemonSet(namespace string, name string) (*apps_v1.DaemonSet, error) {
	return in.k8s.AppsV1().DaemonSets(namespace).Get(in.ctx, name, emptyGetOptions)
}

func (in *K8SClient) GetDaemonSets(namespace string) ([]apps_v1.DaemonSet, error) {
	if daeList, err := in.k8s.AppsV1().DaemonSets(namespace).List(in.ctx, emptyListOptions); err == nil {
		return daeList.Items, nil
	} else {
		return []apps_v1.DaemonSet{}, err
	}
}

// GetDeployment returns the definition of a specific deployment.
// It returns an error on any problem.
func (in *K8SClient) GetDeployment(namespace, name string) (*apps_v1.Deployment, error) {
	return in.k8s.AppsV1().Deployments(namespace).Get(in.ctx, name, emptyGetOptions)
}

// GetRoute returns the external URL endpoint of a specific route name.
// It returns an error on any problem.
func (in *K8SClient) GetRoute(namespace, name string) (*osroutes_v1.Route, error) {
	result := &osroutes_v1.Route{}
	err := in.k8s.Discovery().RESTClient().Get().Prefix("apis", "route.openshift.io", "v1").Namespace(namespace).Resource("routes").SubResource(name).Do(in.ctx).Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetDeployments returns an array of deployments for a given namespace.
// It returns an error on any problem.
func (in *K8SClient) GetDeployments(namespace string) ([]apps_v1.Deployment, error) {
	if depList, err := in.k8s.AppsV1().Deployments(namespace).List(in.ctx, emptyListOptions); err == nil {
		return depList.Items, nil
	} else {
		return []apps_v1.Deployment{}, err
	}
}

// GetDeployment returns the definition of a specific deployment.
// It returns an error on any problem.
func (in *K8SClient) GetDeploymentConfig(namespace, name string) (*osapps_v1.DeploymentConfig, error) {
	result := &osapps_v1.DeploymentConfig{}
	err := in.k8s.Discovery().RESTClient().Get().Prefix("apis", "apps.openshift.io", "v1").Namespace(namespace).Resource("deploymentconfigs").SubResource(name).Do(in.ctx).Into(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetDeployments returns an array of deployments for a given namespace.
// An empty labelSelector will fetch all Deployments for a namespace.
// It returns an error on any problem.
func (in *K8SClient) GetDeploymentConfigs(namespace string) ([]osapps_v1.DeploymentConfig, error) {
	result := &osapps_v1.DeploymentConfigList{}
	err := in.k8s.Discovery().RESTClient().Get().Prefix("apis", "apps.openshift.io", "v1").Namespace(namespace).Resource("deploymentconfigs").Do(in.ctx).Into(result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// GetReplicaSets returns the cached ReplicaSets for the namespace.  For any given RS for a given
// Owner (i.e. Deployment), only the most recent version of the RS will be included in the returned list.
// When an owning Deployment is configured with revisionHistoryLimit > 0, then k8s may return multiple
// versions of the RS for the same Deployment (current and older revisions). Note that it is still possible
// to have multiple RS for the same owner. In which case the most recent version of each is returned.
// see also: ../kubernetes.go
func (in *K8SClient) GetReplicaSets(namespace string) ([]apps_v1.ReplicaSet, error) {
	if rsList, err := in.k8s.AppsV1().ReplicaSets(namespace).List(in.ctx, emptyListOptions); err == nil {
		activeRSMap := map[string]apps_v1.ReplicaSet{}
		for _, rs := range rsList.Items {
			if len(rs.OwnerReferences) > 0 {
				for _, ownerRef := range rs.OwnerReferences {
					if ownerRef.Controller != nil && *ownerRef.Controller {
						key := fmt.Sprintf("%s_%s_%s", ownerRef.Name, rs.Name, rs.ResourceVersion)
						if currRS, ok := activeRSMap[key]; ok {
							if currRS.CreationTimestamp.Time.Before(rs.CreationTimestamp.Time) {
								activeRSMap[key] = rs
							}
						} else {
							activeRSMap[key] = rs
						}
					}
				}
			} else {
				// it is it's own controller
				activeRSMap[rs.Name] = rs
			}
		}

		result := make([]apps_v1.ReplicaSet, len(activeRSMap))
		i := 0
		for _, activeRS := range activeRSMap {
			result[i] = activeRS
			i = i + 1
		}
		return result, nil
	} else {
		return []apps_v1.ReplicaSet{}, err
	}
}

func (in *K8SClient) GetStatefulSet(namespace string, name string) (*apps_v1.StatefulSet, error) {
	return in.k8s.AppsV1().StatefulSets(namespace).Get(in.ctx, name, emptyGetOptions)
}

func (in *K8SClient) GetStatefulSets(namespace string) ([]apps_v1.StatefulSet, error) {
	if ssList, err := in.k8s.AppsV1().StatefulSets(namespace).List(in.ctx, emptyListOptions); err == nil {
		return ssList.Items, nil
	} else {
		return []apps_v1.StatefulSet{}, err
	}
}

func (in *K8SClient) GetReplicationControllers(namespace string) ([]core_v1.ReplicationController, error) {
	if rcList, err := in.k8s.CoreV1().ReplicationControllers(namespace).List(in.ctx, emptyListOptions); err == nil {
		return rcList.Items, nil
	} else {
		return []core_v1.ReplicationController{}, err
	}
}

// GetService returns the definition of a specific service.
// It returns an error on any problem.
func (in *K8SClient) GetService(namespace, name string) (*core_v1.Service, error) {
	return in.k8s.CoreV1().Services(namespace).Get(in.ctx, name, emptyGetOptions)
}

// GetEndpoints return the list of endpoint of a specific service.
// It returns an error on any problem.
func (in *K8SClient) GetEndpoints(namespace, name string) (*core_v1.Endpoints, error) {
	return in.k8s.CoreV1().Endpoints(namespace).Get(in.ctx, name, emptyGetOptions)
}

// GetPods returns the pods definitions for a given set of labels.
// An empty labelSelector will fetch all pods found per a namespace.
// It returns an error on any problem.
func (in *K8SClient) GetPods(namespace, labelSelector string) ([]core_v1.Pod, error) {
	// An empty selector is ambiguous in the go client, could mean either "select all" or "select none"
	// Here we assume empty == select all
	// (see also https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors)
	if pods, err := in.k8s.CoreV1().Pods(namespace).List(in.ctx, meta_v1.ListOptions{LabelSelector: labelSelector}); err == nil {
		return pods.Items, nil
	} else {
		return []core_v1.Pod{}, err
	}
}

// getPodPortForwarder returns a port-forwarder struct which represents an open server forwarding request to the
// requested pod and port
// namespace: name of the namespace where the pod lives in.
// name: name of the pod living in the namespace
// portMap: ports open by the forwarder. Local port and destination port. Format: "80:8080" (local:destination)
// It returns both a portforwarder and an error (if present)
func (in *K8SClient) getPodPortForwarder(namespace, name, portMap string) (httputil.PortForwarder, error) {
	// This branch is just used for testing.
	if in.getPodPortForwarderFunc != nil {
		return in.getPodPortForwarderFunc(namespace, name, portMap)
	}

	writer := new(bytes.Buffer)

	// First try whether the pod exist or not
	pod, err := in.GetPod(namespace, name)
	if err != nil {
		log.Errorf("Couldn't fetch the Pod: %v", err)
		return nil, err
	}

	// Prevent the forward if the pod is not running
	if pod.Status.Phase != core_v1.PodRunning {
		return nil, fmt.Errorf("error creating a pod forwarder for a non-running pod: %s/%s", namespace, name)
	}

	// Create a Port Forwarder
	restInterface := in.k8s.CoreV1().RESTClient()
	return httputil.NewPortForwarder(restInterface, in.restConfig,
		namespace, name, "localhost", portMap, writer)
}

// GetPod returns the pod definitions for a given pod name.
// It returns an error on any problem.
func (in *K8SClient) GetPod(namespace, name string) (*core_v1.Pod, error) {
	if pod, err := in.k8s.CoreV1().Pods(namespace).Get(in.ctx, name, emptyGetOptions); err != nil {
		return nil, err
	} else {
		return pod, nil
	}
}

// StreamPodLogs opens a connection to progressively fetch the logs of a pod. Callers must make sure to properly close the returned io.ReadCloser.
// It returns an error on any problem.
func (in *K8SClient) StreamPodLogs(namespace, name string, opts *core_v1.PodLogOptions) (io.ReadCloser, error) {
	req := in.k8s.CoreV1().RESTClient().Get().Namespace(namespace).Name(name).Resource("pods").SubResource("log").VersionedParams(opts, scheme.ParameterCodec)
	return req.Stream(in.ctx)
}

func (in *K8SClient) GetCronJobs(namespace string) ([]batch_v1.CronJob, error) {
	if cjList, err := in.k8s.BatchV1().CronJobs(namespace).List(in.ctx, emptyListOptions); err == nil {
		return cjList.Items, nil
	} else {
		return []batch_v1.CronJob{}, err
	}
}

func (in *K8SClient) GetJobs(namespace string) ([]batch_v1.Job, error) {
	if jList, err := in.k8s.BatchV1().Jobs(namespace).List(in.ctx, emptyListOptions); err == nil {
		return jList.Items, nil
	} else {
		return []batch_v1.Job{}, err
	}
}

// NewNotFound is a helper method to create a NotFound error similar as used by the kubernetes client.
// This method helps upper layers to send a explicit NotFound error without querying the backend.
func NewNotFound(name, group, resource string) error {
	return errors.NewNotFound(schema.GroupResource{Group: group, Resource: resource}, name)
}

// GetSelfSubjectAccessReview provides information on Kiali permissions
func (in *K8SClient) GetSelfSubjectAccessReview(ctx context.Context, namespace, api, resourceType string, verbs []string) ([]*auth_v1.SelfSubjectAccessReview, error) {
	if config.Get().Server.Observability.Tracing.Enabled {
		var span trace.Span
		ctx, span = otel.Tracer(observability.TracerName()).Start(ctx, "GetSelfSubjectAccessReview",
			trace.WithAttributes(
				attribute.String("package", "kubernetes"),
				attribute.String("namespace", namespace),
				attribute.String("api", api),
				attribute.String("resourceType", resourceType),
				attribute.StringSlice("verbs", verbs),
			),
		)
		defer span.End()
	}

	calls := len(verbs)
	ch := make(chan *auth_v1.SelfSubjectAccessReview, calls)
	errChan := make(chan error)
	for _, v := range verbs {
		go func(ctx context.Context, verb string) {
			res, err := in.k8s.AuthorizationV1().SelfSubjectAccessReviews().Create(ctx, &auth_v1.SelfSubjectAccessReview{
				Spec: auth_v1.SelfSubjectAccessReviewSpec{
					ResourceAttributes: &auth_v1.ResourceAttributes{
						Namespace: namespace,
						Verb:      verb,
						Group:     api,
						Resource:  resourceType,
					},
				},
			}, meta_v1.CreateOptions{})
			if err != nil {
				errChan <- err
			} else {
				ch <- res
			}
		}(ctx, v)
	}

	var err error
	result := []*auth_v1.SelfSubjectAccessReview{}
	for count := 0; count < calls; count++ {
		select {
		case res := <-ch:
			result = append(result, res)
		case err = <-errChan:
			// No op
		}
	}
	return result, err
}

func (in *K8SClient) UpdateWorkload(namespace string, workloadName string, workloadType string, jsonPatch string, patchType string) error {
	emptyPatchOptions := meta_v1.PatchOptions{}
	bytePatch := []byte(jsonPatch)
	var err error
	switch workloadType {
	case DeploymentType:
		_, err = in.k8s.AppsV1().Deployments(namespace).Patch(in.ctx, workloadName, GetPatchType(patchType), bytePatch, emptyPatchOptions)
	case ReplicaSetType:
		_, err = in.k8s.AppsV1().ReplicaSets(namespace).Patch(in.ctx, workloadName, GetPatchType(patchType), bytePatch, emptyPatchOptions)
	case ReplicationControllerType:
		_, err = in.k8s.CoreV1().ReplicationControllers(namespace).Patch(in.ctx, workloadName, GetPatchType(patchType), bytePatch, emptyPatchOptions)
	case DeploymentConfigType:
		if in.IsOpenShift() {
			result := &osapps_v1.DeploymentConfigList{}
			err = in.k8s.Discovery().RESTClient().Patch(GetPatchType(patchType)).Prefix("apis", "apps.openshift.io", "v1").Namespace(namespace).Resource("deploymentconfigs").SubResource(workloadName).Body(bytePatch).Do(in.ctx).Into(result)
		}
	case StatefulSetType:
		_, err = in.k8s.AppsV1().StatefulSets(namespace).Patch(in.ctx, workloadName, GetPatchType(patchType), bytePatch, emptyPatchOptions)
	case JobType:
		_, err = in.k8s.BatchV1().Jobs(namespace).Patch(in.ctx, workloadName, GetPatchType(patchType), bytePatch, emptyPatchOptions)
	case CronJobType:
		_, err = in.k8s.BatchV1().CronJobs(namespace).Patch(in.ctx, workloadName, GetPatchType(patchType), bytePatch, emptyPatchOptions)
	case PodType:
		_, err = in.k8s.CoreV1().Pods(namespace).Patch(in.ctx, workloadName, GetPatchType(patchType), bytePatch, emptyPatchOptions)
	case DaemonSetType:
		_, err = in.k8s.AppsV1().DaemonSets(namespace).Patch(in.ctx, workloadName, GetPatchType(patchType), bytePatch, emptyPatchOptions)
	default:
		err = fmt.Errorf("Workload type %s not found", workloadType)
	}
	return err
}

func (in *K8SClient) UpdateService(namespace string, name string, jsonPatch string, patchType string) error {
	emptyPatchOptions := meta_v1.PatchOptions{}
	bytePatch := []byte(jsonPatch)
	var err error
	_, err = in.k8s.CoreV1().Services(namespace).Patch(in.ctx, name, GetPatchType(patchType), bytePatch, emptyPatchOptions)
	if err != nil {
		log.Errorf("Error is %s", err.Error())
	}
	return err
}

func (in *K8SClient) UpdateNamespace(namespace string, jsonPatch string) (*core_v1.Namespace, error) {
	emptyPatchOptions := meta_v1.PatchOptions{}
	bytePatch := []byte(jsonPatch)
	ns, err := in.k8s.CoreV1().Namespaces().Patch(in.ctx, namespace, types.MergePatchType, bytePatch, emptyPatchOptions)
	if err != nil {
		return &core_v1.Namespace{}, err
	}

	return ns, nil
}

func (in *K8SClient) UpdateProject(namespace string, jsonPatch string) (*osproject_v1.Project, error) {
	result := &osproject_v1.Project{}
	bytePatch := []byte(jsonPatch)
	err := in.k8s.Discovery().RESTClient().Patch(types.MergePatchType).Prefix("apis", "project.openshift.io", "v1", "projects", namespace).Body(bytePatch).Do(in.ctx).Into(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetTokenSubject returns the subject of the authInfo using
// the TokenReview api
func (in *K8SClient) GetTokenSubject(authInfo *api.AuthInfo) (string, error) {
	tokenReview := &v1.TokenReview{}
	tokenReview.Spec.Token = authInfo.Token

	result, err := in.k8s.AuthenticationV1().TokenReviews().Create(in.ctx, tokenReview, meta_v1.CreateOptions{})

	if err != nil {
		return "", err
	} else if result.Status.Error != "" {
		return "", goerrors.New(result.Status.Error)
	} else {
		return result.Status.User.Username, nil
	}
}

package authorization

import (
	"fmt"

	api_security_v1beta "istio.io/api/security/v1beta1"
	networking_v1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	security_v1beta "istio.io/client-go/pkg/apis/security/v1beta1"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/models"
	"github.com/kiali/kiali/util/mtls"
)

const objectType = "authorizationpolicy"

type MtlsEnabledChecker struct {
	AuthorizationPolicies []*security_v1beta.AuthorizationPolicy
	MtlsDetails           kubernetes.MTLSDetails
	ServiceEntries        []networking_v1beta1.ServiceEntry
	RegistryServices      []*kubernetes.RegistryService
}

// Checks if mTLS is enabled, mark all Authz Policies with error
func (c MtlsEnabledChecker) Check() models.IstioValidations {
	validations := models.IstioValidations{}

	for _, ap := range c.AuthorizationPolicies {
		matchLabels := map[string]string{}
		if ap.Spec.Selector != nil {
			matchLabels = ap.Spec.Selector.MatchLabels
		}
		receiveMtlsTraffic := c.IsMtlsEnabledFor(matchLabels, ap.Namespace)
		if !receiveMtlsTraffic {
			if need, paths := needsMtls(ap); need {
				checks := make([]*models.IstioCheck, 0)
				key := models.BuildKey(objectType, ap.Name, ap.Namespace)

				for _, path := range paths {
					check := models.Build("authorizationpolicy.mtls.needstobeenabled", path)
					checks = append(checks, &check)
				}

				validations.MergeValidations(models.IstioValidations{key: &models.IstioValidation{
					Name:       ap.Namespace,
					ObjectType: objectType,
					Valid:      false,
					Checks:     checks,
				}})
			}
		}
	}

	return validations
}

func needsMtls(ap *security_v1beta.AuthorizationPolicy) (bool, []string) {
	paths := make([]string, 0)
	if len(ap.Spec.Rules) == 0 {
		return false, nil
	}

	for i, rule := range ap.Spec.Rules {
		if rule == nil {
			continue
		}
		if needs, fPaths := fromNeedsMtls(rule.From, i); needs {
			paths = append(paths, fPaths...)
		}
		if needs, cPaths := conditionNeedsMtls(rule.When, i); needs {
			paths = append(paths, cPaths...)
		}
	}
	return len(paths) > 0, paths
}

func fromNeedsMtls(froms []*api_security_v1beta.Rule_From, ruleNum int) (bool, []string) {
	paths := make([]string, 0)

	for _, from := range froms {
		if from == nil {
			continue
		}

		if from.Source == nil {
			continue
		}

		if len(from.Source.Principals) > 0 {
			paths = append(paths, fmt.Sprintf("spec/rules[%d]/source/principals", ruleNum))
		}
		if len(from.Source.NotPrincipals) > 0 {
			paths = append(paths, fmt.Sprintf("spec/rules[%d]/source/notPrincipals", ruleNum))
		}
		if len(from.Source.Namespaces) > 0 {
			paths = append(paths, fmt.Sprintf("spec/rules[%d]/source/namespaces", ruleNum))
		}
		if len(from.Source.NotNamespaces) > 0 {
			paths = append(paths, fmt.Sprintf("spec/rules[%d]/source/notNamespaces", ruleNum))
		}
	}
	return len(paths) > 0, paths
}

func conditionNeedsMtls(conditions []*api_security_v1beta.Condition, ruleNum int) (bool, []string) {
	var keysWithMtls = [3]string{"source.namespace", "source.principal", "connection.sni"}
	paths := make([]string, 0)

	for i, c := range conditions {
		if c == nil {
			continue
		}
		for _, key := range keysWithMtls {
			if c.Key == key {
				paths = append(paths, fmt.Sprintf("spec/rules[%d]/when[%d]", ruleNum, i))
			}
		}
	}
	return len(paths) > 0, paths
}

func (c MtlsEnabledChecker) IsMtlsEnabledFor(labels labels.Set, namespace string) bool {
	mtlsEnabledNamespaceLevel := c.hasMtlsEnabledForNamespace(namespace) == mtls.MTLSEnabled
	if labels == nil {
		return mtlsEnabledNamespaceLevel
	}

	workloadmTlsStatus := mtls.MtlsStatus{
		AutoMtlsEnabled:     c.MtlsDetails.EnabledAutoMtls,
		DestinationRules:    c.MtlsDetails.DestinationRules,
		MatchingLabels:      labels,
		PeerAuthentications: c.MtlsDetails.PeerAuthentications,
		RegistryServices:    c.RegistryServices,
	}.WorkloadMtlsStatus(namespace)

	if workloadmTlsStatus == mtls.MTLSEnabled {
		return true
	} else if workloadmTlsStatus == mtls.MTLSDisabled {
		return false
	} else if workloadmTlsStatus == mtls.MTLSNotEnabled {
		// need to check with ns-level and mesh-level status
		return mtlsEnabledNamespaceLevel
	}

	return false
}

func (c MtlsEnabledChecker) hasMtlsEnabledForNamespace(namespace string) string {
	mtlsStatus := mtls.MtlsStatus{
		AutoMtlsEnabled: c.MtlsDetails.EnabledAutoMtls,
	}.OverallMtlsStatus(c.namespaceMtlsStatus(namespace), c.meshWideMtlsStatus())

	// If there isn't any PeerAuthn or DestinationRule and AutoMtls is enabled,
	// then we can consider that the rule will be using mtls
	// Masthead icon won't be present in this case.
	if mtlsStatus == mtls.MTLSNotEnabled && c.MtlsDetails.EnabledAutoMtls {
		mtlsStatus = mtls.MTLSEnabled
	}

	return mtlsStatus
}

func (c MtlsEnabledChecker) meshWideMtlsStatus() mtls.TlsStatus {
	mtlsStatus := mtls.MtlsStatus{
		PeerAuthentications: c.MtlsDetails.MeshPeerAuthentications,
		DestinationRules:    c.MtlsDetails.DestinationRules,
		AutoMtlsEnabled:     c.MtlsDetails.EnabledAutoMtls,
		AllowPermissive:     true,
	}

	return mtlsStatus.MeshMtlsStatus()
}

func (c MtlsEnabledChecker) namespaceMtlsStatus(namespace string) mtls.TlsStatus {
	mtlsStatus := mtls.MtlsStatus{
		PeerAuthentications: c.MtlsDetails.PeerAuthentications,
		DestinationRules:    c.MtlsDetails.DestinationRules,
		AutoMtlsEnabled:     c.MtlsDetails.EnabledAutoMtls,
		AllowPermissive:     true,
	}

	return mtlsStatus.NamespaceMtlsStatus(namespace)
}

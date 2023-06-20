package peerauthentications

import (
	networking_v1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	security_v1beta "istio.io/client-go/pkg/apis/security/v1beta1"

	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/models"
)

type DisabledMeshWideChecker struct {
	PeerAuthn        *security_v1beta.PeerAuthentication
	DestinationRules []*networking_v1beta1.DestinationRule
}

func (c DisabledMeshWideChecker) Check() ([]*models.IstioCheck, bool) {
	validations := make([]*models.IstioCheck, 0)

	// Validation only affects to PeerAuthn disabling mTLS
	if _, mode := kubernetes.PeerAuthnHasMTLSEnabled(c.PeerAuthn); mode != "DISABLE" {
		return validations, true
	}

	for _, dr := range c.DestinationRules {
		if _, mode := kubernetes.DestinationRuleHasMeshWideMTLSEnabled(dr); mode == "ISTIO_MUTUAL" || mode == "MUTUAL" {
			check := models.Build("peerauthentications.mtls.disablemeshdestinationrulemissing", "spec/mtls")
			return append(validations, &check), false
		}
	}

	return validations, true
}

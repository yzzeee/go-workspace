package checkers

import (
	"istio.io/client-go/pkg/apis/telemetry/v1alpha1"

	"github.com/kiali/kiali/models"
)

const TelemetryCheckerType = "telemetry"

type TelemetryChecker struct {
	Namespaces  models.Namespaces
	Telemetries []*v1alpha1.Telemetry
}

// An Object Checker runs all checkers for an specific object type (i.e.: pod, route rule,...)
// It run two kinds of checkers:
// 1. Individual checks: validating individual objects.
// 2. Group checks: validating behaviour between configurations.
func (in TelemetryChecker) Check() models.IstioValidations {
	validations := models.IstioValidations{}

	return validations
}

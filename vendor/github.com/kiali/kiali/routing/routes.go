package routing

import (
	"net/http"

	"github.com/kiali/kiali/handlers"
)

// Route describes a single route
type Route struct {
	Name          string
	Method        string
	Pattern       string
	HandlerFunc   http.HandlerFunc
	Authenticated bool
}

// Routes holds an array of Route. A note on swagger documentation. The path variables and query parameters
// are defined in ../doc.go.  YOu need to manually associate params and routes.
type Routes struct {
	Routes []Route
}

// NewRoutes creates and returns all the API routes
func NewRoutes() (r *Routes) {
	r = new(Routes)

	r.Routes = []Route{
		// swagger:route GET /healthz kiali healthz
		// ---
		// Endpoint to get the health of Kiali
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		// responses:
		//		500: internalError
		//		200
		{
			"Healthz",
			"GET",
			"/healthz",
			handlers.Healthz,
			false,
		},
		// swagger:route GET / kiali root
		// ---
		// Endpoint to get the status of Kiali
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		// responses:
		//      500: internalError
		//      200: statusInfo
		{
			"Root",
			"GET",
			"/api",
			handlers.Root,
			false,
		},
		// swagger:route GET /authenticate auth authenticate
		// ---
		// Endpoint to authenticate the user
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		//    Security:
		//     authorization: user, password
		//
		// responses:
		//      500: internalError
		//      200: userSessionData
		{
			"Authenticate",
			"GET",
			"/api/authenticate",
			handlers.Authenticate,
			false,
		},
		// swagger:route POST /authenticate auth openshiftCheckToken
		// ---
		// Endpoint to check if a token from Openshift is working correctly
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      200: userSessionData
		{
			"OpenshiftCheckToken",
			"POST",
			"/api/authenticate",
			handlers.Authenticate,
			false,
		},
		// swagger:route GET /logout auth logout
		// ---
		// Endpoint to logout an user (unset the session cookie)
		//
		//     Schemes: http, https
		//
		// responses:
		//      204: noContent
		{
			"Logout",
			"GET",
			"/api/logout",
			handlers.Logout,
			false,
		},
		// swagger:route GET /auth/info auth authenticationInfo
		// ---
		// Endpoint to get login info, such as strategy, authorization endpoints
		// for OAuth providers and so on.
		//
		//     Consumes:
		//     - application/json
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      200: authenticationInfo
		{
			"AuthenticationInfo",
			"GET",
			"/api/auth/info",
			handlers.AuthenticationInfo,
			false,
		},
		// swagger:route GET /status status getStatus
		// ---
		// Endpoint to get the status of Kiali
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      200: statusInfo
		{
			"Status",
			"GET",
			"/api/status",
			handlers.Root,
			true,
		},
		// swagger:route GET /config kiali getConfig
		// ---
		// Endpoint to get the config of Kiali
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      200: statusInfo
		{
			"Config",
			"GET",
			"/api/config",
			handlers.Config,
			true,
		},
		// swagger:route GET /crippled kiali getCrippledFeatures
		// ---
		// Endpoint to get the crippled features of Kiali
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      200: statusInfo
		{
			"Crippled",
			"GET",
			"/api/crippled",
			handlers.CrippledFeatures,
			true,
		},
		// swagger:route GET /istio/permissions config getPermissions
		// ---
		// Endpoint to get the caller permissions on new Istio Config objects
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      200: istioConfigPermissions
		{
			"IstioConfigPermissions",
			"GET",
			"/api/istio/permissions",
			handlers.IstioConfigPermissions,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/istio config istioConfigList
		// ---
		// Endpoint to get the list of Istio Config of a namespace
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      200: istioConfigList
		//
		{
			"IstioConfigList",
			"GET",
			"/api/namespaces/{namespace}/istio",
			handlers.IstioConfigList,
			true,
		},
		// swagger:route GET /istio config istioConfigListAll
		// ---
		// Endpoint to get the list of Istio Config of all namespaces
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      200: istioConfigList
		//
		{
			"IstioConfigListAll",
			"GET",
			"/api/istio/config",
			handlers.IstioConfigList,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/istio/{object_type}/{object} config istioConfigDetails
		// ---
		// Endpoint to get the Istio Config of an Istio object
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      404: notFoundError
		//      500: internalError
		//      200: istioConfigDetailsResponse
		//
		{
			"IstioConfigDetails",
			"GET",
			"/api/namespaces/{namespace}/istio/{object_type}/{object}",
			handlers.IstioConfigDetails,
			true,
		},
		// swagger:route DELETE /namespaces/{namespace}/istio/{object_type}/{object} config istioConfigDelete
		// ---
		// Endpoint to delete the Istio Config of an (arbitrary) Istio object
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      404: notFoundError
		//      500: internalError
		//      200
		//
		{
			"IstioConfigDelete",
			"DELETE",
			"/api/namespaces/{namespace}/istio/{object_type}/{object}",
			handlers.IstioConfigDelete,
			true,
		},
		// swagger:route PATCH /namespaces/{namespace}/istio/{object_type}/{object} config istioConfigUpdate
		// ---
		// Endpoint to update the Istio Config of an Istio object used for templates and adapters using Json Merge Patch strategy.
		//
		//     Consumes:
		//	   - application/json
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      404: notFoundError
		//      500: internalError
		//      200: istioConfigDetailsResponse
		//
		{
			"IstioConfigUpdate",
			"PATCH",
			"/api/namespaces/{namespace}/istio/{object_type}/{object}",
			handlers.IstioConfigUpdate,
			true,
		},
		// swagger:route POST /namespaces/{namespace}/istio/{object_type} config istioConfigCreate
		// ---
		// Endpoint to create an Istio object by using an Istio Config item
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      404: notFoundError
		//      500: internalError
		//		202
		//		201: istioConfigDetailsResponse
		//      200: istioConfigDetailsResponse
		//
		{
			"IstioConfigCreate",
			"POST",
			"/api/namespaces/{namespace}/istio/{object_type}",
			handlers.IstioConfigCreate,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/services services serviceList
		// ---
		// Endpoint to get the details of a given service
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      200: serviceListResponse
		//
		{
			"ServiceList",
			"GET",
			"/api/namespaces/{namespace}/services",
			handlers.ServiceList,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/services/{service} services serviceDetails
		// ---
		// Endpoint to get the details of a given service
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      404: notFoundError
		//      500: internalError
		//      200: serviceDetailsResponse
		//
		{
			"ServiceDetails",
			"GET",
			"/api/namespaces/{namespace}/services/{service}",
			handlers.ServiceDetails,
			true,
		},
		// swagger:route PATCH /namespaces/{namespace}/services/{service} services serviceUpdate
		// ---
		// Endpoint to update the Service configuration using Json Merge Patch strategy.
		//
		//     Consumes:
		//	   - application/json
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      404: notFoundError
		//      500: internalError
		//      200: serviceDetailsResponse
		//
		{
			"ServiceUpdate",
			"PATCH",
			"/api/namespaces/{namespace}/services/{service}",
			handlers.ServiceUpdate,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/apps/{app}/spans traces appSpans
		// ---
		// Endpoint to get Jaeger spans for a given app
		//
		//		Produces:
		//		- application/json
		//
		//		Schemes: http, https
		//
		// responses:
		// 		500: internalError
		//		200: spansResponse
		{
			"AppSpans",
			"GET",
			"/api/namespaces/{namespace}/apps/{app}/spans",
			handlers.AppSpans,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/workloads/{workload}/spans traces workloadSpans
		// ---
		// Endpoint to get Jaeger spans for a given workload
		//
		//		Produces:
		//		- application/json
		//
		//		Schemes: http, https
		//
		// responses:
		// 		500: internalError
		//		200: spansResponse
		{
			"WorkloadSpans",
			"GET",
			"/api/namespaces/{namespace}/workloads/{workload}/spans",
			handlers.WorkloadSpans,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/services/{service}/spans traces serviceSpans
		// ---
		// Endpoint to get Jaeger spans for a given service
		//
		//		Produces:
		//		- application/json
		//
		//		Schemes: http, https
		//
		// responses:
		// 		500: internalError
		//		200: spansResponse
		{
			"ServiceSpans",
			"GET",
			"/api/namespaces/{namespace}/services/{service}/spans",
			handlers.ServiceSpans,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/apps/{app}/traces traces appTraces
		// ---
		// Endpoint to get the traces of a given app
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      404: notFoundError
		//      500: internalError
		//      200: traceDetailsResponse
		//
		{
			"AppTraces",
			"GET",
			"/api/namespaces/{namespace}/apps/{app}/traces",
			handlers.AppTraces,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/services/{service}/traces traces serviceTraces
		// ---
		// Endpoint to get the traces of a given service
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      404: notFoundError
		//      500: internalError
		//      200: traceDetailsResponse
		//
		{
			"ServiceTraces",
			"GET",
			"/api/namespaces/{namespace}/services/{service}/traces",
			handlers.ServiceTraces,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/workloads/{workload}/traces traces workloadTraces
		// ---
		// Endpoint to get the traces of a given workload
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      404: notFoundError
		//      500: internalError
		//      200: traceDetailsResponse
		//
		{
			"WorkloadTraces",
			"GET",
			"/api/namespaces/{namespace}/workloads/{workload}/traces",
			handlers.WorkloadTraces,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/apps/{app}/errortraces traces errorTraces
		// ---
		// Endpoint to get the number of traces in error for a given service
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      404: notFoundError
		//      500: internalError
		//      200: errorTracesResponse
		//
		{
			"ErrorTraces",
			"GET",
			"/api/namespaces/{namespace}/apps/{app}/errortraces",
			handlers.ErrorTraces,
			true,
		},
		// swagger:route GET /traces/{traceID} traces traceDetails
		// ---
		// Endpoint to get a specific trace from ID
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      404: notFoundError
		//      500: internalError
		//      200: traceDetailsResponse
		//
		{
			"TracesDetails",
			"GET",
			"/api/traces/{traceID}",
			handlers.TraceDetails,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/workloads workloads workloadList
		// ---
		// Endpoint to get the list of workloads for a namespace
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      200: workloadListResponse
		//
		{
			"WorkloadList",
			"GET",
			"/api/namespaces/{namespace}/workloads",
			handlers.WorkloadList,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/workloads/{workload} workloads workloadDetails
		// ---
		// Endpoint to get the workload details
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      404: notFoundError
		//      200: workloadDetails
		//
		{
			"WorkloadDetails",
			"GET",
			"/api/namespaces/{namespace}/workloads/{workload}",
			handlers.WorkloadDetails,
			true,
		},
		// swagger:route PATCH /namespaces/{namespace}/workloads/{workload} workloads workloadUpdate
		// ---
		// Endpoint to update the Workload configuration using Json Merge Patch strategy.
		//
		//     Consumes:
		//	   - application/json
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      404: notFoundError
		//      500: internalError
		//      200: workloadDetails
		//
		{
			"WorkloadUpdate",
			"PATCH",
			"/api/namespaces/{namespace}/workloads/{workload}",
			handlers.WorkloadUpdate,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/apps apps appList
		// ---
		// Endpoint to get the list of apps for a namespace
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      200: appListResponse
		//
		{
			"AppList",
			"GET",
			"/api/namespaces/{namespace}/apps",
			handlers.AppList,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/apps/{app} apps appDetails
		// ---
		// Endpoint to get the app details
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      404: notFoundError
		//      200: appDetails
		//
		{
			"AppDetails",
			"GET",
			"/api/namespaces/{namespace}/apps/{app}",
			handlers.AppDetails,
			true,
		},
		// swagger:route GET /namespaces namespaces namespaceList
		// ---
		// Endpoint to get the list of the available namespaces
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      200: namespaceList
		//
		{
			"NamespaceList",
			"GET",
			"/api/namespaces",
			handlers.NamespaceList,
			true,
		},
		// swagger:route PATCH /namespaces/{namespace} namespaces namespaceUpdate
		// ---
		// Endpoint to update the Namespace configuration using Json Merge Patch strategy.
		//
		//     Consumes:
		//	   - application/json
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      404: notFoundError
		//      500: internalError
		//      200: namespaceResponse
		//
		{
			"NamespaceUpdate",
			"PATCH",
			"/api/namespaces/{namespace}",
			handlers.NamespaceUpdate,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/services/{service}/metrics services serviceMetrics
		// ---
		// Endpoint to fetch metrics to be displayed, related to a single service
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      503: serviceUnavailableError
		//      200: metricsResponse
		//
		{
			"ServiceMetrics",
			"GET",
			"/api/namespaces/{namespace}/services/{service}/metrics",
			handlers.ServiceMetrics,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/aggregates/{aggregate}/{aggregateValue}/metrics aggregates aggregateMetrics
		// ---
		// Endpoint to fetch metrics to be displayed, related to a single aggregate
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      503: serviceUnavailableError
		//      200: metricsResponse
		//
		{
			"AggregateMetrics",
			"GET",
			"/api/namespaces/{namespace}/aggregates/{aggregate}/{aggregateValue}/metrics",
			handlers.AggregateMetrics,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/apps/{app}/metrics apps appMetrics
		// ---
		// Endpoint to fetch metrics to be displayed, related to a single app
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      503: serviceUnavailableError
		//      200: metricsResponse
		//
		{
			"AppMetrics",
			"GET",
			"/api/namespaces/{namespace}/apps/{app}/metrics",
			handlers.AppMetrics,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/workloads/{workload}/metrics workloads workloadMetrics
		// ---
		// Endpoint to fetch metrics to be displayed, related to a single workload
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      503: serviceUnavailableError
		//      200: metricsResponse
		//
		{
			"WorkloadMetrics",
			"GET",
			"/api/namespaces/{namespace}/workloads/{workload}/metrics",
			handlers.WorkloadMetrics,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/services/{service}/dashboard services serviceDashboard
		// ---
		// Endpoint to fetch dashboard to be displayed, related to a single service
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      503: serviceUnavailableError
		//      200: dashboardResponse
		//
		{
			"ServiceDashboard",
			"GET",
			"/api/namespaces/{namespace}/services/{service}/dashboard",
			handlers.ServiceDashboard,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/apps/{app}/dashboard apps appDashboard
		// ---
		// Endpoint to fetch dashboard to be displayed, related to a single app
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      503: serviceUnavailableError
		//      200: dashboardResponse
		//
		{
			"AppDashboard",
			"GET",
			"/api/namespaces/{namespace}/apps/{app}/dashboard",
			handlers.AppDashboard,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/workloads/{workload}/dashboard workloads workloadDashboard
		// ---
		// Endpoint to fetch dashboard to be displayed, related to a single workload
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      503: serviceUnavailableError
		//      200: dashboardResponse
		//
		{
			"WorkloadDashboard",
			"GET",
			"/api/namespaces/{namespace}/workloads/{workload}/dashboard",
			handlers.WorkloadDashboard,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/customdashboard/{dashboard} dashboards customDashboard
		// ---
		// Endpoint to fetch a custom dashboard
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      503: serviceUnavailableError
		//      200: dashboardResponse
		//
		{
			"CustomDashboard",
			"GET",
			"/api/namespaces/{namespace}/customdashboard/{dashboard}",
			handlers.CustomDashboard,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/metrics namespaces namespaceMetrics
		// ---
		// Endpoint to fetch metrics to be displayed, related to a namespace
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      503: serviceUnavailableError
		//      200: metricsResponse
		//
		{
			"NamespaceMetrics",
			"GET",
			"/api/namespaces/{namespace}/metrics",
			handlers.NamespaceMetrics,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/health namespaces namespaceHealth
		// ---
		// Get health for all objects in the given namespace
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      200: namespaceAppHealthResponse
		//      400: badRequestError
		//      500: internalError
		//
		{
			"NamespaceHealth",
			"GET",
			"/api/namespaces/{namespace}/health",
			handlers.NamespaceHealth,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/validations namespaces namespaceValidations
		// ---
		// Get validation summary for all objects in the given namespace
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      200: namespaceValidationSummaryResponse
		//      400: badRequestError
		//      500: internalError
		//
		{
			"NamespaceValidationSummary",
			"GET",
			"/api/namespaces/{namespace}/validations",
			handlers.NamespaceValidationSummary,
			true,
		},
		// swagger:route GET /istio/validations namespaces namespacesValidations
		// ---
		// Get validation summary for all objects in the given namespaces
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      200: namespaceValidationSummaryResponse
		//      400: badRequestError
		//      500: internalError
		//
		{
			"ConfigValidationSummary",
			"GET",
			"/api/istio/validations",
			handlers.ConfigValidationSummary,
			true,
		},
		// swagger:route GET /mesh/tls tls meshTls
		// ---
		// Get TLS status for the whole mesh
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      200: meshTlsResponse
		//      400: badRequestError
		//      500: internalError
		//
		{
			"NamespaceTls",
			"GET",
			"/api/mesh/tls",
			handlers.MeshTls,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/tls tls namespaceTls
		// ---
		// Get TLS status for the given namespace
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      200: namespaceTlsResponse
		//      400: badRequestError
		//      500: internalError
		//
		{
			"NamespaceTls",
			"GET",
			"/api/namespaces/{namespace}/tls",
			handlers.NamespaceTls,
			true,
		},
		// swagger:route GET /istio/status status istioStatus
		// ---
		// Get the status of each components needed in the control plane
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      200: istioStatusResponse
		//      400: badRequestError
		//      500: internalError
		//
		{
			"IstioStatus",
			"GET",
			"/api/istio/status",
			handlers.IstioStatus,
			true,
		},
		// swagger:route GET /istio/certs certs istioCerts
		// ---
		// Get certificates (internal) information used by Istio
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      200: certsInfoResponse
		//      500: internalError
		//
		{
			"IstioCerts",
			"GET",
			"/api/istio/certs",
			handlers.IstioCerts,
			true,
		},
		// swagger:route GET /namespaces/graph graphs graphNamespaces
		// ---
		// The backing JSON for a namespaces graph.
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      500: internalError
		//      200: graphResponse
		//
		{
			"GraphNamespaces",
			"GET",
			"/api/namespaces/graph",
			handlers.GraphNamespaces,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/aggregates/{aggregate}/{aggregateValue}/graph graphs graphAggregate
		// ---
		// The backing JSON for an aggregate node detail graph. (supported graphTypes: app | versionedApp | workload)
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      500: internalError
		//      200: graphResponse
		//
		{

			"GraphAggregate",
			"GET",
			"/api/namespaces/{namespace}/aggregates/{aggregate}/{aggregateValue}/graph",
			handlers.GraphNode,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/aggregates/{aggregate}/{aggregateValue}/{service}/graph graphs graphAggregateByService
		// ---
		// The backing JSON for an aggregate node detail graph, specific to a service. (supported graphTypes: app | versionedApp | workload)
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      500: internalError
		//      200: graphResponse
		//
		{

			"GraphAggregateByService",
			"GET",
			"/api/namespaces/{namespace}/aggregates/{aggregate}/{aggregateValue}/{service}/graph",
			handlers.GraphNode,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/applications/{app}/versions/{version}/graph graphs graphAppVersion
		// ---
		// The backing JSON for a versioned app node detail graph. (supported graphTypes: app | versionedApp)
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      500: internalError
		//      200: graphResponse
		//
		{

			"GraphAppVersion",
			"GET",
			"/api/namespaces/{namespace}/applications/{app}/versions/{version}/graph",
			handlers.GraphNode,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/applications/{app}/graph graphs graphApp
		// ---
		// The backing JSON for an app node detail graph. (supported graphTypes: app | versionedApp)
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      500: internalError
		//      200: graphResponse
		//
		{
			"GraphApp",
			"GET",
			"/api/namespaces/{namespace}/applications/{app}/graph",
			handlers.GraphNode,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/services/{service}/graph graphs graphService
		// ---
		// The backing JSON for a service node detail graph.
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      500: internalError
		//      200: graphResponse
		//
		{
			"GraphService",
			"GET",
			"/api/namespaces/{namespace}/services/{service}/graph",
			handlers.GraphNode,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/workloads/{workload}/graph graphs graphWorkload
		// ---
		// The backing JSON for a workload node detail graph.
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      400: badRequestError
		//      500: internalError
		//      200: graphResponse
		//
		{
			"GraphWorkload",
			"GET",
			"/api/namespaces/{namespace}/workloads/{workload}/graph",
			handlers.GraphNode,
			true,
		},
		// swagger:route GET /grafana integrations grafanaInfo
		// ---
		// Get the grafana URL and other descriptors
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      503: serviceUnavailableError
		//      200: grafanaInfoResponse
		//      204: noContent
		//
		{
			"GrafanaURL",
			"GET",
			"/api/grafana",
			handlers.GetGrafanaInfo,
			true,
		},
		// swagger:route GET /jaeger integrations jaegerInfo
		// ---
		// Get the jaeger URL and other descriptors
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      404: notFoundError
		//      406: notAcceptableError
		//      200: jaegerInfoResponse
		//
		{
			"JaegerURL",
			"GET",
			"/api/jaeger",
			handlers.GetJaegerInfo,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/pods/{pod} pods podDetails
		// ---
		// Endpoint to get pod details
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      404: notFoundError
		//      200: workloadDetails
		//
		{
			"PodDetails",
			"GET",
			"/api/namespaces/{namespace}/pods/{pod}",
			handlers.PodDetails,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/pods/{pod}/logs pods podLogs
		// ---
		// Endpoint to get pod logs
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      404: notFoundError
		//      200: workloadDetails
		//
		{
			"PodLogs",
			"GET",
			"/api/namespaces/{namespace}/pods/{pod}/logs",
			handlers.PodLogs,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/pods/{pod}/config_dump pods podProxyDump
		// ---
		// Endpoint to get pod proxy dump
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      404: notFoundError
		//      200: configDump
		//
		{
			"PodConfigDump",
			"GET",
			"/api/namespaces/{namespace}/pods/{pod}/config_dump",
			handlers.ConfigDump,
			true,
		},
		// swagger:route GET /namespaces/{namespace}/pods/{pod}/config_dump/{resource} pods podProxyResource
		// ---
		// Endpoint to get pod logs
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      404: notFoundError
		//      200: configDumpResource
		//
		{
			"PodConfigDump",
			"GET",
			"/api/namespaces/{namespace}/pods/{pod}/config_dump/{resource}",
			handlers.ConfigDumpResourceEntries,
			true,
		},
		// swagger:route POST /namespaces/{namespace}/pods/{pod}/logging pods podProxyLogging
		// ---
		// Endpoint to set pod proxy log level
		//
		//     Produces:
		//     - application/json
		//
		//     Schemes: http, https
		//
		// responses:
		//      500: internalError
		//      404: notFoundError
		//      400: badRequestError
		//      200: noContent
		//
		{
			"PodProxyLogging",
			"POST",
			"/api/namespaces/{namespace}/pods/{pod}/logging",
			handlers.LoggingUpdate,
			true,
		},

		// swagger:route POST /stats/metrics stats metricsStats
		// ---
		// Produces metrics statistics
		//
		// 		Produces:
		//		- application/json
		//
		//		Schemes: http, https
		//
		// responses:
		//    400: badRequestError
		//    503: serviceUnavailableError
		//		500: internalError
		//		200: metricsStatsResponse
		{
			Name:          "MetricsStats",
			Method:        "POST",
			Pattern:       "/api/stats/metrics",
			HandlerFunc:   handlers.MetricsStats,
			Authenticated: true,
		},
		// swagger:route GET /api/clusters
		// ---
		// Endpoint to get the list of the clusters that are hosting the service mesh.
		//              Produces:
		//              - application/json
		//
		//              Schemes: http, https
		//
		// responses:
		//              500: internalError
		//              200: clustersResponse
		{
			"GetClusters",
			"GET",
			"/api/clusters",
			handlers.GetClusters,
			true,
		},
		// swagger:route GET /api/mesh/outbound_traffic_policy/mode
		// ---
		// Endpoint to get the OutboundTrafficPolicy Mode configured in the service mesh.
		//              Produces:
		//              - application/json
		//
		//              Schemes: http, https
		//
		// responses:
		//              500: internalError
		//              200: clustersResponse
		{
			"OutboundTrafficPolicyMode",
			"GET",
			"/api/mesh/outbound_traffic_policy/mode",
			handlers.OutboundTrafficPolicyMode,
			true,
		},
		// swagger:route GET /api/mesh/resources/thresholds
		// ---
		// Endpoint to get the IstiodResourceThresholds.
		//              Produces:
		//              - application/json
		//
		//              Schemes: http, https
		//
		// responses:
		//              500: internalError
		//              200: istiodResourceThresholds
		{
			"IstiodResourceThresholds",
			"GET",
			"/api/mesh/resources/thresholds",
			handlers.IstiodResourceThresholds,
			true,
		},
		// swagger:route GET /api/mesh/canaries/status
		// ---
		// Endpoint to get the IstiodCanariesStatus.
		//              Produces:
		//              - application/json
		//
		//              Schemes: http, https
		//
		// responses:
		//              500: internalError
		//              200: istiodCanariesStatus
		{
			"IstiodCanariesStatus",
			"GET",
			"/api/mesh/canaries/status",
			handlers.IstiodCanariesStatus,
			true,
		},
	}

	return
}

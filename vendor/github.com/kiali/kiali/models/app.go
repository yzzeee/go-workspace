package models

type AppList struct {
	// Namespace where the apps live in
	// required: true
	// example: bookinfo
	Namespace Namespace `json:"namespace"`

	// Cluster where the apps live in
	// required: true
	// example: east
	Cluster string `json:"cluster"`

	// Applications for a given namespace
	// required: true
	Apps []AppListItem `json:"applications"`
}

// AppListItem has the necessary information to display the console app list
type AppListItem struct {
	// Name of the application
	// required: true
	// example: reviews
	Name string `json:"name"`

	// Cluster of the application
	// required: true
	// example: reviews
	Cluster string `json:"cluster"`

	// Define if all Pods related to the Workloads of this app has an IstioSidecar deployed
	// required: true
	// example: true
	IstioSidecar bool `json:"istioSidecar"`

	// Define if any pod has the Ambient annotation
	// required: true
	// example: true
	IstioAmbient bool `json:"istioAmbient"`

	// Labels for App
	Labels map[string]string `json:"labels"`

	// Istio References
	IstioReferences []*IstioValidationKey `json:"istioReferences"`

	// Health
	Health AppHealth `json:"health,omitempty"`
}

type WorkloadItem struct {
	// Name of a workload member of an application
	// required: true
	// example: reviews-v1
	WorkloadName string `json:"workloadName"`

	// Define if all Pods related to the Workload has an IstioSidecar deployed
	// required: true
	// example: true
	IstioSidecar bool `json:"istioSidecar"`

	// Define if belongs to a namespace labeled as ambient
	// required: true
	// example: true
	IstioAmbient bool `json:"istioAmbient"`

	// Labels for Workload
	Labels map[string]string `json:"labels"`

	// List of service accounts involved in this application
	// required: true
	ServiceAccountNames []string `json:"serviceAccountNames"`
}

type App struct {
	// Namespace where the app lives in
	// required: true
	// example: bookinfo
	Namespace Namespace `json:"namespace"`

	// Name of the application
	// required: true
	// example: reviews
	Name string `json:"name"`

	// Cluster of the application
	// required: false
	// example: east
	Cluster string `json:"cluster"`

	// Workloads for a given application
	// required: true
	Workloads []WorkloadItem `json:"workloads"`

	// List of service names linked with an application
	// required: true
	ServiceNames []string `json:"serviceNames"`

	// Runtimes and associated dashboards
	Runtimes []Runtime `json:"runtimes"`

	// Health
	Health AppHealth `json:"health"`
}

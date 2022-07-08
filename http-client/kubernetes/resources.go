package kubernetes

import (
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/json"
)

func getV1GroupVersionResource(resourceGroup string, resourceVersion string, resourceKind string) schema.GroupVersionResource {
	var groupVersionResource schema.GroupVersionResource
	if resourceGroup != "" && resourceVersion != "" {
		groupVersionResource = schema.GroupVersionResource{
			Group:    resourceGroup,
			Version:  resourceVersion,
			Resource: resourceKind,
		}
	} else {
		groupVersionResource = schema.GroupVersionResource{
			Version:  "v1",
			Resource: resourceKind,
		}
	}
	return groupVersionResource
}

// GetLimitRange 리소스의 limit range 정보를 가져온다
func GetLimitRange() ([]byte, error) {
	limitRanges, err := ClientSettings.CoreV1().LimitRanges("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return json.Marshal(limitRanges.Items[0].Spec)
}

// GetPipelines 리소스의 pipelines 정보를 가져온다
func GetPipelines() ([]byte, error) {
	groupVersionResource := getV1GroupVersionResource("tekton.dev", "v1beta1", "pipelines")
	var unstructuredResult *unstructured.UnstructuredList
	var err error
	unstructuredResult, err = DynamicClient.Resource(groupVersionResource).Namespace("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	resourceList := unstructuredResult.UnstructuredContent()

	return json.Marshal(resourceList)
}

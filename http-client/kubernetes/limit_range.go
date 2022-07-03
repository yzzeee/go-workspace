package kubernetes

import (
	"context"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

// GetLimitRange 리소스의 limit range 정보를 가져온다
func GetLimitRange() ([]byte, error) {
	limitRanges, err := ClientSettings.CoreV1().LimitRanges("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var items interface{}
	if len(limitRanges.Items) != 0 {
		items = limitRanges.Items[0].Spec
	}
	return json.Marshal(items)
}

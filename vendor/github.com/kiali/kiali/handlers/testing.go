/*
  This file contains testing helpers for the handlers package.
*/

package handlers

import (
	"fmt"
	"time"

	osproject_v1 "github.com/openshift/api/project/v1"
	core_v1 "k8s.io/api/core/v1"

	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/util"
)

func mockClock() {
	clockTime := time.Date(2017, 0o1, 15, 0, 0, 0, 0, time.UTC)
	util.Clock = util.ClockMock{Time: clockTime}
}

type noPrivClient struct {
	kubernetes.ClientInterface
}

func (n *noPrivClient) GetProjects(labelSelector string) ([]osproject_v1.Project, error) {
	return nil, fmt.Errorf("Rejecting")
}

func (n *noPrivClient) GetProject(name string) (*osproject_v1.Project, error) {
	return nil, fmt.Errorf("Rejecting")
}

func (n *noPrivClient) GetNamespace(namespace string) (*core_v1.Namespace, error) {
	return nil, fmt.Errorf("Rejecting")
}

func (n *noPrivClient) GetNamespaces(labelSelector string) ([]core_v1.Namespace, error) {
	return nil, fmt.Errorf("Rejecting")
}

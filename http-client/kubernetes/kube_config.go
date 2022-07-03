package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeConfigPath kube config 설정파일의 위치
var KubeConfigPath string

// ClientSettings kube config 설정을 담고 있다
var ClientSettings *kubernetes.Clientset

// IgnoreTLSVerification 클라이언트가 Kubernetes API 에 접속할때 TLS 인증을 무시할지 여부
var IgnoreTLSVerification bool

func InitConfig() error {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", KubeConfigPath)
	if err != nil {
		return err
	}

	if IgnoreTLSVerification {
		kubeConfig.Insecure = true
		kubeConfig.TLSClientConfig.CAData = nil
	}

	ClientSettings, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}

	return nil
}

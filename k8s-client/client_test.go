package k8s_client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	v12 "k8s.io/api/authentication/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"testing"
)

var clientset *kubernetes.Clientset

func init() {
	config, err := clientcmd.BuildConfigFromFlags(
		"", "/home/yzzeee/.kube/config",
	)
	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

func TestGetSecret(t *testing.T) {
	secret, _ := clientset.CoreV1().Secrets("istio-system").Get(context.TODO(), "kiali-signing-key", v1.GetOptions{})
	fmt.Printf("%+v", secret)
}

func TestCreateToken(t *testing.T) {
	token, _ := clientset.CoreV1().ServiceAccounts("istio-system").
		CreateToken(context.TODO(), "kiali-service-account", &v12.TokenRequest{}, v1.CreateOptions{})
	fmt.Println(token.Status.Token)
}

func TestGetAuthentication(t *testing.T) {
	token, _ := clientset.CoreV1().ServiceAccounts("istio-system").
		CreateToken(context.TODO(), "kiali-service-account", &v12.TokenRequest{}, v1.CreateOptions{})
	fmt.Println(token.Status.Token)

	requestURL := "http://192.168.67.100/kiali/api/authenticate"
	res, _ := req.C().R().
		SetContentType("application/x-www-form-urlencoded").
		SetBody(bytes.NewBufferString("token=eyJhbGciOiJSUzI1NiIsImtpZCI6IkdsbXNpT1laRzNTd1hQSzNkN2Z1eWdJRVREVmFmVWtOZ2Vnd0FKWGlkVHMifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNjkwMjgwMTA5LCJpYXQiOjE2OTAyNzY1MDksImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJpc3Rpby1zeXN0ZW0iLCJzZXJ2aWNlYWNjb3VudCI6eyJuYW1lIjoia2lhbGktc2VydmljZS1hY2NvdW50IiwidWlkIjoiYWE2Mzc2ODAtZWY1ZS00YmFhLTgyMzMtMGQ4MDBlYWE3Y2QzIn19LCJuYmYiOjE2OTAyNzY1MDksInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDppc3Rpby1zeXN0ZW06a2lhbGktc2VydmljZS1hY2NvdW50In0.f0ezhpzByieP4K1TyNC9Nl-WLKpaLVh4ZyxcoZ1dAwcUoI19hbD50e13-lwCe28RT3gdOC3xsZ2axn3vVo_Kr9EYGXiEbgqEGQajFBIPZGmDFFhmuDWZG_SpbLy1-HQZdCKMCb7OTVgZFaOurwgBjDX9pYCGIH1XAjuFdTImX0YlmbuQHF-1zSH_u_OxBdZuxqp8ycvWWe7aKQ3iVWQfx3aneyZIX7zMJDFIZImK509qv-BnLIPTXDXKxkYlRzsiCzNlo3emY6C8hXOhLTU7FIyK9tjR-_dMKCkdux623NylacsnFxONDf7Am2aSvOuZ8aXAvldFVm0TRohaYWvg-Q")).
		Post(requestURL)

	fmt.Println(res.StatusCode, token.Status.Token, "<----------------")

	for _, cookie := range res.Cookies() {
		fmt.Println("cookie: ", cookie)
	}
}

func TestCreateVirtualService(t *testing.T) {
	token, _ := clientset.CoreV1().ServiceAccounts("istio-system").
		CreateToken(context.TODO(), "kiali-service-account", &v12.TokenRequest{}, v1.CreateOptions{})
	fmt.Println(token.Status.Token)

	authRequestURL := "http://192.168.67.100/kiali/api/authenticate"
	authRes, _ := req.C().R().
		SetContentType("application/x-www-form-urlencoded").
		SetBody(bytes.NewBufferString("token=eyJhbGciOiJSUzI1NiIsImtpZCI6IkdsbXNpT1laRzNTd1hQSzNkN2Z1eWdJRVREVmFmVWtOZ2Vnd0FKWGlkVHMifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNjkwMjg1NTMyLCJpYXQiOjE2OTAyODE5MzIsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJpc3Rpby1zeXN0ZW0iLCJzZXJ2aWNlYWNjb3VudCI6eyJuYW1lIjoia2lhbGktc2VydmljZS1hY2NvdW50IiwidWlkIjoiYWE2Mzc2ODAtZWY1ZS00YmFhLTgyMzMtMGQ4MDBlYWE3Y2QzIn19LCJuYmYiOjE2OTAyODE5MzIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDppc3Rpby1zeXN0ZW06a2lhbGktc2VydmljZS1hY2NvdW50In0.razvgOfFBkJL88MsvrzWxLtSRMTTmEEaZy9hT8RSi0uiRrJm5UmEYX9Jwiqokv9KKSf1Yli-gfIERcTuV08qg-k908-39rPeW4dluPZ_4Y4ag9iI6JgrHlSJ_UATsVOZgFtU-Uc_jFDeuC_KDmIivsSaQKEqMYxBBoXWsA8BJekZumNAHG82a638WfK7ObmIH9H9Lo395dvudEA-tKoOmSaJV_NjoopHWS1zEkjZbi6QoUYhdk4UijaWuycW4V-oCGXOahznZJw9gsqlIIRjcMaPaUJ3_7h4Qh9l4TuGrkmOZiNEeZI01nNQIzbANnAM8MeVkPLFw081dRlW3raLwA")).
		Post(authRequestURL)

	fmt.Println(authRes.StatusCode, token.Status.Token, "<----------------")
	client := req.C().R()
	client.SetCookies(authRes.Cookies()...)

	requestURL := "http://192.168.67.100/kiali/api/namespaces/default/istio/virtualservices"
	res, _ := client.
		SetContentType("application/x-www-form-urlencoded").
		SetBody(bytes.NewBufferString("{\"kind\":\"VirtualService\",\"apiVersion\":\"networking.istio.io/v1beta1\",\"metadata\":{\"namespace\":\"default\",\"name\":\"fleetma이름n-position-tracker\",\"labels\":{\"kiali_wizard\":\"fault_injection\"}},\"spec\":{\"http\":[{\"route\":[{\"destination\":{\"host\":\"fleetman-position-tracker.default.svc.cluster.local\",\"subset\":\"latest\"},\"weight\":100}],\"fault\":{\"delay\":{\"percentage\":{\"value\":100},\"fixedDelay\":\"5s\"}}}],\"hosts\":[\"fleetman-position-tracker.default.svc.cluster.local\"],\"gateways\":null}}")).
		Post(requestURL)

	t.Log(res.StatusCode)

	for _, cookie := range res.Cookies() {
		fmt.Println("cookie: ", cookie)
	}
}

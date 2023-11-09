package k8s

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restClient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8sClient struct {
	clientset *kubernetes.Clientset
}

func NewK8sClient() (*K8sClient, error) {
	c, err := resolveK8sClient()
	if err != nil {
		return nil, err
	}
	return &K8sClient{clientset: c}, nil
}

func resolveK8sClient() (*kubernetes.Clientset, error) {
	config, err := getK8sConfig()
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func getK8sConfig() (*restClient.Config, error) {
	config, err := restClient.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// fmt.Println("error creating client InCluster, fallback to kubeconfig")

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err == nil {
		return config, nil
	}

	return nil, err
}

func (k *K8sClient) ListEvents() error {
	ev, err := k.clientset.CoreV1().Events("default").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	fmt.Println(ev)
	return nil
}

// ListPods returns a list of all pending pods all namespaces.
func (k *K8sClient) ListPendingPods(namespace string) (*v1.PodList, error) {
	// get all pending pods
	podList, err := k.clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		FieldSelector: "status.phase=Pending",
	})
	if err != nil {
		return nil, err
	}

	return podList, nil
}

func (k *K8sClient) ListNodesByLabel(label string) (*v1.NodeList, error) {
	n, err := k.clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
		LabelSelector: label,
	})

	if err != nil {
		return nil, err
	}

	return n, nil
}

// func (k *K8sClient) DrainNode(node *v1.Node) error {
// 	k.clientset.CoreV1().Nodes().Update(context.Background(), node, metav1.UpdateOptions{})
// }

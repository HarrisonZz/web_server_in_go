package kubernetes

import (
	"context"
	"fmt"
	"os"

	"github.com/HarrisonZz/web_server_in_go/internal/logger"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ─── NODE 資訊取得 ───────────────────────────────────────────────────
type Node struct {
	Name       string            `json:"name"`
	InternalIP string            `json:"internal_ip"`
	Labels     map[string]string `json:"labels"`
	CPU        string            `json:"cpu"`
	Memory     string            `json:"memory"`
	OS         string            `json:"os"`
	Arch       string            `json:"arch"`
	Kernel     string            `json:"kernel"`
	OSImage    string            `json:"os_image"`
}

var (
	clientSet *kubernetes.Clientset
	NodeInfo  *Node
)

func init() {
	clientSet = newClient()
	node_name := os.Getenv("NODE_NAME")

	info, err := getNodeInfo(clientSet, node_name)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to preload NodeInfo: %v", err))
		return
	}

	NodeInfo = info
	logger.Info(fmt.Sprintf("NodeInfo initialized: %s (%s)", NodeInfo.Name, NodeInfo.InternalIP))
}

func newClient() *kubernetes.Clientset {

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil
	}
	clientSet, err := kubernetes.NewForConfig(config)

	return clientSet
}

func getNodeInfo(clientset *kubernetes.Clientset, nodeName string) (*Node, error) {
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	info := &Node{
		Name:       node.Name,
		InternalIP: getNodeIP(node, v1.NodeInternalIP),
		Labels:     node.Labels,
		CPU:        node.Status.Allocatable.Cpu().String(),
		Memory:     node.Status.Allocatable.Memory().String(),
		OS:         node.Status.NodeInfo.OperatingSystem,
		Arch:       node.Status.NodeInfo.Architecture,
		Kernel:     node.Status.NodeInfo.KernelVersion,
		OSImage:    node.Status.NodeInfo.OSImage,
	}
	return info, nil
}

func getNodeIP(node *v1.Node, ipType v1.NodeAddressType) string {
	for _, addr := range node.Status.Addresses {
		if addr.Type == ipType {
			return addr.Address
		}
	}
	return ""
}

package ks

import (
	"context"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Node Name -> Node
var NodesInfo map[string]Node
var NextUpdate4NodesInfo time.Time

type Node struct {
	Ip   string `json:"ip"`
	Zone string `json:"zone,omitempty"`
}

type NodeInterface interface {
	GetNodes() error
	GetZone(name string) string
	GetIP(name string) string
}

// Get Nodes information from Kubernetes API
func (n *Node) GetNodes() error {

	if NodesInfo == nil || NextUpdate4NodesInfo.After(time.Now()) {
		NodesInfo = make(map[string]Node)

		// creates the in-cluster config
		config, err := rest.InClusterConfig()
		if err != nil {
			log.Error().Interface("err", err).Msg("rest.InClusterConfig at GetNodes()")
			log.Info().Msg("Outside of cluster and try reading kubeconfig")
			config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
			if err != nil {
				log.Error().Interface("err", err).Msg("clientcmd.BuildConfigFromFlags")
				return err
			}
		}
		// creates the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Error().Interface("err", err).Msg("kubernetes.NewForConfig")
			return err
		}

		nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			nodes, err = clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Error().Interface("err", err).Msg("clientset.CoreV1.Nodes.List")
				return err
			}
		}
		for _, nd := range nodes.Items {
			name := nd.GetName()
			ip := ""
			for _, addr := range nd.Status.Addresses {
				if addr.Type == "InternalIP" {
					ip = addr.Address
					break
				}
			}
			if zone, ok := nd.GetLabels()["topology.kubernetes.io/zone"]; ok {
				NodesInfo[name] = Node{
					Ip:   ip,
					Zone: zone,
				}
			} else {
				log.Error().Interface("err", err).Fields(nd.GetLabels()).Msg("Failed to find label - topology.kubernetes.io/zone")
				NodesInfo[name] = Node{
					Ip: ip,
				}
			}
		}
		log.Info().Interface("nodes", NodesInfo).Send()
		NextUpdate4NodesInfo = time.Now().Add(300 * time.Second)
	}

	return nil
}

// Get Zone by Node name
func (n *Node) GetZone(name string) string {
	if z, ok := NodesInfo[name]; ok {
		return z.Zone
	}
	return ""

}

// Get node IP by Node name
func (n *Node) GetIP(name string) string {
	if z, ok := NodesInfo[name]; ok {
		return z.Ip
	}
	return ""
}

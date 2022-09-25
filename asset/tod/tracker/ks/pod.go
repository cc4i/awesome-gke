package ks

import (
	"net"
	"os"

	"github.com/rs/zerolog/log"
)

// TODO: Retrieve pod info by calling APIs.
// Pod info
type Pod struct {
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
	Image      string `json:"image,omitempty"`
	NodeName   string `json:"node_name,omitempty"`
	NodeIp     string `json:"node_ip,omitempty"`
	Zone       string `json:"zone,omitempty"`
	PodIp      string `json:"pod_ip,omitempty"`
	UpCaller   string `json:"up_caller,omitempty"`
	NextCallee string `json:"next_callee,omitempty"`
}

type PodInterface interface {
	GetLocalIP() string
	GetPodName() string
	GetNodeName() string
	BuildResponse() string
	BuildPod() Pod
}

func (p *Pod) BuildPod() *Pod {
	node := p.GetNodeName()
	return &Pod{
		Name:      p.GetPodName(),
		Namespace: p.GetPodNamespace(),
		NodeName:  node,
		NodeIp:    p.GetNodeIP(node),
		Zone:      p.GetZone(node),
		PodIp:     p.GetLocalIP(),
	}

}

// Get local/Pod IP
func (p *Pod) GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// Get Pod name from env variable
func (p *Pod) GetPodName() string {
	return os.Getenv("POD_NAME")
}

// Get the namepspace where the Pod was in
func (p *Pod) GetPodNamespace() string {
	return os.Getenv("POD_NAMESPACE")
}

// Get the Node name from env variable
func (p *Pod) GetNodeName() string {
	return os.Getenv("POD_NODE_NAME")
}

func (p *Pod) BuildResponse() string {
	return "StatusOK from " + p.GetPodName()
}

func (p *Pod) GetNodeIP(nodeName string) string {
	n := &Node{}
	if err := n.GetNodes(); err != nil {
		log.Error().Interface("err", err).Msg("GetNodes")
	}
	return n.GetIP(nodeName)
}

func (p *Pod) GetZone(nodeName string) string {
	n := &Node{}
	if err := n.GetNodes(); err != nil {
		log.Error().Interface("err", err).Msg("GetNodes")
	}
	return n.GetZone(nodeName)
}

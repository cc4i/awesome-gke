package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Node Name -> Node IP
var nodeIps map[string]string

// Node Name -> Zone
var nodeZones map[string]string

// Invoking chain with reverse order
type TripDetail []P2p

type P2p struct {
	Number      int    `json:"number"`
	Source      Point  `json:"source"`
	Destination Point  `json:"destination"`
	Method      string `json:"rethod"`
	RequestURI  string `json:"request_uri,omitempty"`
	Reqest      string `json:"reqest,omitempty"`
	Response    string `json:"response"`
}

type Point struct {
	Ip  string `json:"ip"`
	Pod *Pod   `json:"pod,omitempty"`
}

type Pod struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	NodeName  string `json:"node_name"`
	NodeIp    string `json:"node_ip"`
	Zone      string `json:"zone"`
}

// Get Nodes information from Kubernetes API
func getNodes() {
	nodeIps = make(map[string]string)
	nodeZones = make(map[string]string)

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal().Interface("err", err).Msg("rest.InClusterConfig")
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal().Interface("err", err).Msg("kubernetes.NewForConfig")
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		time.Sleep(500 * time.Millisecond)
		nodes, err = clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatal().Interface("err", err).Msg("clientset.CoreV1.Nodes.List")
		}
	}
	for _, n := range nodes.Items {
		name := n.GetName()
		for _, addr := range n.Status.Addresses {
			if addr.Type == "InternalIP" {
				nodeIps[name] = addr.Address
				break
			}
		}
		if zone, ok := n.GetLabels()["topology.kubernetes.io/zone"]; ok {
			nodeZones[name] = zone
		} else {
			log.Error().Interface("err", err).Fields(n.GetLabels()).Msg("Failed to find label - topology.kubernetes.io/zone")
		}
	}
	log.Info().Interface("node_ips", nodeIps).Send()
	log.Info().Interface("node_zones", nodeZones).Send()

}

// Get Zone by Node name
func getZone(name string) string {
	return nodeZones[name]
}

// Get node IP by Node name
func getNodeIP(name string) string {
	return nodeIps[name]
}

// Get local/Pod IP
func getLocalIP() string {
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
func getPodName() string {
	return os.Getenv("POD_NAME")
}

// Get the namepspace where the Pod was in
func getPodNamespace() string {
	return os.Getenv("POD_NAMESPACE")
}

// Get the Node name from env variable
func getNodeName() string {
	return os.Getenv("POD_NODE_NAME")
}

// Get the URL for next call from env variable
func getNextCall() string {
	return os.Getenv("NEXT_CALL")
}

func buildResponse() string {
	return "StatusOK from " + getLocalIP()
}

// Call /trip API and return TripDetail
func callTrip(url string) TripDetail {
	log.Info().Str("next_call", url).Send()
	if url != "null" && url != "" {
		client := http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Error().Interface("err", err).Msg("http.NewRequest")
			return nil
		}
		req.Header.Add("x-pod-name", getPodName())
		req.Header.Add("x-pod-namespace", getPodNamespace())
		req.Header.Add("x-pod-ip", getLocalIP())
		res, err := client.Do(req)
		if err != nil {
			log.Error().Interface("err", err).Msg("client.Do")
			return nil
		}
		buf, _ := ioutil.ReadAll(res.Body)
		var ttd TripDetail
		json.Unmarshal(buf, &ttd)
		log.Info().Interface("remote_return", ttd).Send()
		return ttd

	}
	return nil
}

// /trip API
func trip(c *gin.Context) {
	var td TripDetail
	whoami := c.Param("whoami")

	log.Info().Str("next_call", getNextCall()).Send()
	if ttd := callTrip(getNextCall()); ttd != nil {
		td = append(td, ttd...)
	}

	myself := &Pod{
		Name:      getPodName(),
		Namespace: getPodNamespace(),
		NodeName:  getNodeName(),
		NodeIp:    getNodeIP(getNodeName()),
		Zone:      getZone(getNodeName()),
	}

	src := Point{
		Ip: c.Request.Header.Get("x-pod-ip"),
	}

	if whoami == "pod" {
		src.Pod = &Pod{
			Name:      c.Request.Header.Get("x-pod-name"),
			Namespace: c.Request.Header.Get("x-pod-namespace"),
		}
	}
	p2p := P2p{
		Number: len(td) + 1,
		Source: src,
		Destination: Point{
			Ip:  getLocalIP(),
			Pod: myself,
		},
		Method:     c.Request.Method,
		RequestURI: c.Request.RequestURI,
		Response:   buildResponse(),
	}
	log.Info().Interface("return", p2p).Send()
	td = append(td, p2p)
	c.JSON(http.StatusOK, &td)
}

func router(ctx context.Context, r *gin.Engine) *gin.Engine {
	r.GET("/trip", trip)
	r.GET("/trip/:whoami", trip)
	return r
}

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Info().Str("local_ip", getLocalIP()).Str("status", "ready").Msg("Tracker service")

	//Retrieve Nodes inforamtion
	getNodes()

	gin.DisableConsoleColor()
	server := gin.Default()
	port := os.Getenv("POD_PORT")
	if port != "" {
		log.Fatal().Err(router(context.Background(), server).Run("0.0.0.0:" + port))
	} else {
		log.Fatal().Err(router(context.Background(), server).Run("0.0.0.0:8000"))
	}

}

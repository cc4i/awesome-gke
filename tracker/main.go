package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
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
	Method      string `json:"method,omitempty"`
	RequestURI  string `json:"request_uri,omitempty"`
	Reqest      string `json:"reqest,omitempty"`
	Response    string `json:"response,omitempty"`
	Latency     int64  `json:"latency,omitempty"`
}

type Point struct {
	Ip  string `json:"ip"`
	Pod *Pod   `json:"pod,omitempty"`
}

type Pod struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	NodeName  string `json:"node_name,omitempty"`
	NodeIp    string `json:"node_ip,omitempty"`
	Zone      string `json:"zone,omitempty"`
	PodIp     string `json:"pod_ip,omitempty"`
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
	if z, ok := nodeZones[name]; ok {
		return z
	} else {
		getNodes()
	}
	return nodeZones[name]
}

// Get node IP by Node name
func getNodeIP(name string) string {
	if ip, ok := nodeIps[name]; ok {
		return ip
	} else {
		getNodes()
	}
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

// Inter-call API (/trip) and return TripDetail
func callTrip(url string) TripDetail {
	log.Info().Str("url", url).Send()
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
		req.Header.Add("x-node-name", getNodeName())
		req.Header.Add("x-node-ip", getNodeIP(getNodeName()))
		req.Header.Add("x-zone", getZone(getNodeName()))
		req.Header.Add("x-request-start", strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10))
		res, err := client.Do(req)
		if err != nil {
			log.Error().Interface("err", err).Msg("client.Do")
			return nil
		}
		log.Info().Str("http_code", res.Status).Msg("return code from " + url)
		buf, _ := ioutil.ReadAll(res.Body)
		var ttd TripDetail
		json.Unmarshal(buf, &ttd)
		log.Info().Interface("remote_return", ttd).Send()
		return ttd

	}
	return nil
}

// API (/trip) - get TripDetail
func trip(c *gin.Context) {
	var td TripDetail
	whoami := c.Param("whoami")

	log.Info().Str("next_call", getNextCall()).Send()
	if ttd := callTrip(getNextCall()); ttd != nil {
		td = append(td, ttd...)
	}

	//Get one-way trip latency: A->B / put time into header & calculate
	start, _ := strconv.ParseInt(c.Request.Header.Get("x-request-start"), 10, 64)
	end := time.Now().UnixNano() / int64(time.Millisecond)
	// Get remote client IP if it's first call
	cltIp := c.Request.Header.Get("x-pod-ip")
	if cltIp == "" {
		//TODO: Get external IP when call from inside Pod
		cltIp = c.ClientIP()
	}
	src := Point{
		Ip: cltIp,
	}

	// Build source pod from headers
	if whoami == "pod" {
		src.Pod = &Pod{
			Name:      c.Request.Header.Get("x-pod-name"),
			Namespace: c.Request.Header.Get("x-pod-namespace"),
			NodeName:  c.Request.Header.Get("x-node-name"),
			NodeIp:    c.Request.Header.Get("x-node-ip"),
			Zone:      c.Request.Header.Get("x-zone"),
			PodIp:     cltIp,
		}
	}
	p2p := P2p{
		Number: len(td) + 1,
		Source: src,
		Destination: Point{
			Ip: getLocalIP(),
			Pod: &Pod{
				Name:      getPodName(),
				Namespace: getPodNamespace(),
				NodeName:  getNodeName(),
				NodeIp:    getNodeIP(getNodeName()),
				Zone:      getZone(getNodeName()),
				PodIp:     getLocalIP(),
			},
		},
		Method:     c.Request.Method,
		RequestURI: c.Request.RequestURI,
		Response:   buildResponse(),
		Latency:    end - start,
	}
	td = append(td, p2p)
	log.Info().Interface("return", td).Msg("trip()")
	c.JSON(http.StatusOK, &td)
}

func contains(s []string, str string) bool {

	for _, v := range s {
		if strings.HasPrefix(str, v) {
			log.Debug().Str("str", str).Str("prefix", v).Msg("true")
			return true
		}
	}
	log.Debug().Interface("all_prefix", s).Str("str", str).Msg("false")
	return false
}

// API (/initial) - get all pods as per request
func getInitialPods(c *gin.Context) {
	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Interface("err", err).Msg("Read request body")
	}
	log.Info().Str("request_str", string(buf)).Send()

	strs := strings.Split(string(buf), "::")
	ns := strs[0]
	prefixs := strings.Split(strs[1], ",")
	log.Info().Strs("prefixs", prefixs).Send()

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

	pods, err := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		time.Sleep(500 * time.Millisecond)
		pods, err = clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatal().Interface("err", err).Msg("clientset.CoreV1.Nodes.List")
		}
	}
	//Scanning all pods
	allPods := make(map[string]Pod)
	for _, p := range pods.Items {
		log.Info().Str("pod", p.Name).Str("pod_status_podip", p.Status.PodIP).Send()
		if contains(prefixs, p.Name) && p.Status.PodIP != "" {
			node := p.Spec.NodeName
			if _, ok := nodeZones[node]; !ok {
				getNodes()
			}
			log.Info().Str("pod", p.Name).Msg("Added into allPods")
			allPods[p.Name] = Pod{
				Namespace: p.Namespace,
				Name:      p.Name,
				NodeName:  node,
				NodeIp:    nodeIps[node],
				Zone:      nodeZones[node],
				PodIp:     p.Status.PodIP,
			}
		}
	}

	//Build possible links as per pods
	var td TripDetail
	for i := 0; i < len(prefixs); i++ {
		for _, pSrc := range allPods {
			if strings.HasPrefix(pSrc.Name, prefixs[i]) && (i+1) < len(prefixs) {
				for _, pDst := range allPods {
					if strings.HasPrefix(pDst.Name, prefixs[i+1]) {
						src := pSrc
						dst := pDst
						p2p := P2p{
							Number: 0,
							Source: Point{
								Ip:  src.PodIp,
								Pod: &src,
							},
							Destination: Point{
								Ip:  dst.PodIp,
								Pod: &dst,
							},
						}
						log.Info().Str("src", pSrc.Name).Str("dst", pDst.Name).Msg("checking out src->dst")
						td = append(td, p2p)
					}
				}
			}
		}
	}
	log.Info().Interface("return", td).Msg("getInitialPods()")
	c.JSON(http.StatusOK, &td)
}

func router(ctx context.Context, r *gin.Engine) *gin.Engine {
	log.Info().Interface("ctx", ctx).Msg("context.Context pairs")
	r.GET("/trip", trip)
	r.GET("/trip/:whoami", trip)
	r.POST("/initial", getInitialPods)

	// ko add static assets under ./kodata - https://github.com/google/ko#static-assets
	if staticDir := os.Getenv("KO_DATA_PATH"); staticDir != "" {
		log.Info().Str("KO_DATA_PATH", staticDir).Send()
		r.Static("tracker-ui", staticDir)
	} else {
		r.Static("tracker-ui", "./kodata")
	}

	return r
}

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Info().Str("local_ip", getLocalIP()).Str("status", "ready").Msg("Tracker service")

	gin.DisableConsoleColor()
	server := gin.Default()

	//setup CORS policies
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	server.Use(cors.Default())

	if port := os.Getenv("POD_PORT"); port != "" {
		log.Fatal().Err(router(context.Background(), server).Run("0.0.0.0:" + port))
	} else {
		log.Fatal().Err(router(context.Background(), server).Run("0.0.0.0:8000"))
	}

}

package trip

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"tracker/ks"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Invoking chain with reverse order
type TripDetail struct {
	Id     string `json:"id"`
	Detail []P2p  `json:"detail"`
}

// Single call source -> destination
type P2p struct {
	Number      int    `json:"number"`
	Source      Point  `json:"source"`
	Destination Point  `json:"destination"`
	Method      string `json:"method,omitempty"`
	RequestURI  string `json:"request_uri,omitempty"`
	Request     string `json:"request,omitempty"`
	Response    string `json:"response,omitempty"`
	Latency     int64  `json:"latency,omitempty"`
}

// Point for source/destination
type Point struct {
	Ip  string  `json:"ip"`
	Pod *ks.Pod `json:"pod,omitempty"`
}

type TripInterface interface {
	GetInitialPods(from string, ns string, prefixs []string) error
	CallTrip(kp *ks.Pod, url string) error
	TripHistory() error
	GoTrip(whoami string, headers map[string]string, clientIp string, reqMethod string, reqUri string, nextCall string, whereami string) error
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

// Get intial relations of pods under specificed namespace & call chain
func (td *TripDetail) GetInitialPods(from string, ns string, prefixs []string) error {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Error().Interface("err", err).Msg("rest.InClusterConfig")
		return err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error().Interface("err", err).Msg("kubernetes.NewForConfig")
		return err
	}

	pods, err := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		time.Sleep(500 * time.Millisecond)
		pods, err = clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Error().Interface("err", err).Msg("clientset.CoreV1.Nodes.List")
			return err
		}
	}
	//Scanning all pods
	allPods := make(map[string]ks.Pod)
	for _, p := range pods.Items {
		log.Info().Str("pod", p.Name).Str("pod_status_podip", p.Status.PodIP).Send()
		if contains(prefixs, p.Name) && p.Status.PodIP != "" {
			log.Info().Str("pod", p.Name).Msg("Added into allPods")

			image := ""
			for _, img := range p.Spec.Containers {
				if strings.Contains(img.Image, "tracker") {
					image = img.Image
				}
			}

			pp := &ks.Pod{
				Namespace: p.Namespace,
				Name:      p.Name,
				Image:     image,
				NodeName:  p.Spec.NodeName,
				PodIp:     p.Status.PodIP,
			}
			pp.NodeIp = pp.GetNodeIP(pp.NodeName)
			pp.Zone = pp.GetZone(pp.NodeName)
			allPods[p.Name] = *pp
		}
	}

	//Build possible links as per pods

	for i := 0; i < len(prefixs); i++ {
		//add from_nodes
		if i == 0 {
			for _, pSrc := range allPods {
				if strings.HasPrefix(pSrc.Name, prefixs[i]) && (i+1) < len(prefixs) {
					fsrc := pSrc
					fP2p := P2p{
						Number: 0,
						Source: Point{
							Ip: from,
						},
						Destination: Point{
							Ip:  fsrc.PodIp,
							Pod: &fsrc,
						},
					}
					td.Detail = append(td.Detail, fP2p)
				}
			}
		}
		for _, pSrc := range allPods {
			if strings.HasPrefix(pSrc.Name, prefixs[i]) && (i+1) < len(prefixs) {
				// add rest of nodes
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
						td.Detail = append(td.Detail, p2p)
					}
				}
			}
		}
	}
	log.Info().Interface("return", td).Msg("GetInitialPods()")
	return nil
}

// Inter-call API (/trip) and return TripDetail
func (td *TripDetail) CallTrip(kp *ks.Pod, url string) error {

	log.Info().Str("url", url).Send()
	if url != "null" && url != "" {
		client := http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Error().Interface("err", err).Msg("http.NewRequest")
			return err
		}
		req.Header.Add("x-pod-name", kp.GetPodName())
		req.Header.Add("x-pod-namespace", kp.GetPodNamespace())
		req.Header.Add("x-pod-ip", kp.GetLocalIP())
		req.Header.Add("x-node-name", kp.GetNodeName())
		req.Header.Add("x-node-ip", kp.GetNodeIP(kp.GetNodeName()))
		req.Header.Add("x-zone", kp.GetZone(kp.GetNodeName()))
		req.Header.Add("x-request-start", strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10))
		res, err := client.Do(req)
		if err != nil {
			log.Error().Interface("err", err).Msg("client.Do")
			return err
		}
		log.Info().Str("http_code", res.Status).Msg("return code from " + url)
		buf, _ := ioutil.ReadAll(res.Body)
		log.Info().Int("buf_num", len(buf)).Send()
		log.Info().Interface("buf", buf).Send()
		var p2ps []P2p
		if err = json.Unmarshal(buf, &p2ps); err != nil {
			log.Error().Interface("err", err).Msg("json.Unmarshal")
			return err
		} else {
			log.Info().Int("p2ps", len(p2ps)).Send()
			td.Detail = append(td.Detail, p2ps...)
		}
		log.Info().Interface("remote_return", td.Detail).Send()
		return nil

	}
	return nil
}

func (td *TripDetail) TripHistory() error {
	if maps, err := QueryAllTds4Redis(); err != nil {
		return err
	} else {

		for _, val := range maps {
			var ttd []P2p
			json.Unmarshal([]byte(val), &ttd)
			td.Detail = append(td.Detail, ttd...)
		}
	}
	return nil
}

// Go through trip chain/call one by one and return all data
//
// whoami - indicate for origal call or call from Pods
// headers - all headers from HTTP request
// clientIp - remote ip of client
// reqMethod - request method of HTTP
// reqUri - URL of http call
// nextCall - URL of next http call
// whereami - where the Tracker was deployed
func (td *TripDetail) GoTrip(whoami string, headers map[string]string, clientIp string, reqMethod string, reqUri string, nextCall string, whereami string) error {

	// Build source pod from headers if call from pod
	srcPod := &ks.Pod{}
	if whoami == "pod" {
		srcPod.Name = headers["x-pod-name"]
		srcPod.Namespace = headers["x-pod-namespace"]
		srcPod.NodeName = headers["x-node-name"]
		srcPod.NodeIp = headers["x-node-ip"]
		srcPod.Zone = headers["x-zone"]
		srcPod.PodIp = headers["x-pod-ip"]
	}

	if err := td.CallTrip(srcPod, nextCall); err != nil {
		log.Error().Interface("err", err).Msg("tp.CallTrip")
		return err
	}
	log.Info().Int("p2p_num", len(td.Detail)).Send()

	//Get one-way trip latency: A->B / put time into header & calculate
	start, _ := strconv.ParseInt(headers["x-request-start"], 10, 64)
	end := time.Now().UnixNano() / int64(time.Millisecond)
	// One time call and latency is 0
	if headers["x-request-start"] == "" {
		start = end
	}

	//Source
	src := Point{
		Ip: clientIp,
	}
	if whoami == "pod" {
		src.Pod = srcPod
	}

	wp := Point{
		Ip: whereami,
	}
	var swp Point
	if headers["x-pod-ip"] == "" {
		fp2p := P2p{
			Number:      0,
			Source:      src,
			Destination: wp,
		}
		td.Detail = append(td.Detail, fp2p)
		swp = wp
	} else {
		swp = src
	}

	//Destination
	dstPod := &ks.Pod{}
	p2p := P2p{
		Number: len(td.Detail),
		Source: swp,
		Destination: Point{
			Ip:  dstPod.GetLocalIP(),
			Pod: dstPod.BuildPod(),
		},
		Method:     reqMethod,
		RequestURI: reqUri,
		Response:   dstPod.BuildResponse(),
		Latency:    end - start,
	}
	td.Detail = append(td.Detail, p2p)
	log.Info().Interface("return", td.Detail).Msg("startTrip()")
	// Save to redis, only once
	if whoami != "pod" {
		buf, _ := json.Marshal(td.Detail)
		if err := SaveTd2Redis(td.Id, buf); err != nil {
			log.Error().Interface("error", err).Msg("fail to trip.SaveTd2Redis()")
		}
	}
	return nil
}

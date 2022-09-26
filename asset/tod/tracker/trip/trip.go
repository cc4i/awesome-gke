package trip

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"tracker/ks"

	"github.com/rs/zerolog/log"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
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
	GetInitialPods(from string, ns string) error
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
func (td *TripDetail) GetInitialPods(from string, ns string) error {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Error().Interface("err", err).Msg("rest.InClusterConfig at GetInitialPods()")
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

	pods, err := clientset.CoreV1().Pods(ns).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		time.Sleep(500 * time.Millisecond)
		pods, err = clientset.CoreV1().Pods(ns).List(context.TODO(), v1.ListOptions{})
		if err != nil {
			log.Error().Interface("err", err).Msg("clientset.CoreV1.Nodes.List")
			return err
		}
	}

	//Scanning all pods with Label :: "app:tracker"
	allPods := make(map[string]ks.Pod)
	for _, p := range pods.Items {
		log.Info().Str("pod", p.Name).Str("pod_status_podip", p.Status.PodIP).Send()
		lable := p.Labels["app"]
		if lable == "tracker" && p.Status.PodIP != "" {
			log.Info().Str("pod", p.Name).Msg("Added into allPods")

			var image, upCaller, nextCallee string
			for _, img := range p.Spec.Containers {
				if strings.Contains(img.Image, "tracker") {
					image = img.Image
					for _, env := range img.Env {
						if env.Name == "UP_CALLER" {
							upCaller = env.Value
						}
						if env.Name == "NEXT_CALLEE" {
							nextCallee = env.Value
						}
					}
				}
			}
			pp := &ks.Pod{
				Namespace:  p.Namespace,
				Name:       p.Name,
				Image:      image,
				NodeName:   p.Spec.NodeName,
				PodIp:      p.Status.PodIP,
				UpCaller:   upCaller,
				NextCallee: nextCallee,
			}
			pp.NodeIp = pp.GetNodeIP(pp.NodeName)
			pp.Zone = pp.GetZone(pp.NodeName)
			allPods[p.Name] = *pp
		}
	}

	//Build possible links as per pods
	//BO:
	no := 0
	log.Info().Int("allPods", len(allPods)).Msgf("The number of pods in %s", ns)
	for _, pod := range allPods {

		// Starting nodes
		if pod.UpCaller == "" && pod.NextCallee != "" {
			dstPod := pod
			fP2p := P2p{
				Number: no,
				Source: Point{
					Ip: from,
				},
				Destination: Point{
					Ip:  dstPod.PodIp,
					Pod: &dstPod,
				},
			}
			td.Detail = append(td.Detail, fP2p)
			no += 1
		}

		if pod.NextCallee != "" {
			for _, callee := range strings.Split(pod.NextCallee, ",") {
				srcPod := pod
				for _, dstPod := range allPods {
					if strings.HasPrefix(dstPod.Name, callee) {
						fP2p := P2p{
							Number: no,
							Source: Point{
								Ip:  srcPod.PodIp,
								Pod: &srcPod,
							},
							Destination: Point{
								Ip:  dstPod.PodIp,
								Pod: &dstPod,
							},
						}
						td.Detail = append(td.Detail, fP2p)
						no += 1
					}
				}
			}
		}
	}
	//EO:

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

	//Multiple services to call
	nextCallServices := strings.Split(nextCall, ",")
	for _, ncs := range nextCallServices {
		if ncs != "" {
			// TODO: pods namespace/port
			url := fmt.Sprintf("http://%s:%s/trip/pod", ncs, "8000")

			if err := td.CallTrip(srcPod, url); err != nil {
				log.Error().Interface("err", err).Msg("tp.CallTrip")
				return err
			}
			log.Info().Int("p2p_num", len(td.Detail)).Send()

		}
	}

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

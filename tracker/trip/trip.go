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
	From   string `json:"from,omitempty"`
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
	GetInitialPods(ns string, prefixs []string) error
	CallTrip(url string) error
	TripHistory() error
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
func (td *TripDetail) GetInitialPods(ns string, prefixs []string) error {

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

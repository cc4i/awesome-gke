package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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
}

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

func getPodName() string {
	return os.Getenv("POD_NAME")
}

func getPodNamespace() string {
	return os.Getenv("POD_NAMESPACE")
}

func getNextCall() string {
	return os.Getenv("NEXT_CALL")
}

func buildResponse() string {
	return "StatusOK from " + getLocalIP()
}

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
	}

	src := Point{
		Ip: c.ClientIP(),
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

	gin.DisableConsoleColor()
	server := gin.Default()
	port := os.Getenv("POD_PORT")
	if port != "" {
		log.Fatal().Err(router(context.Background(), server).Run("0.0.0.0:" + port))
	} else {
		log.Fatal().Err(router(context.Background(), server).Run("0.0.0.0:8000"))
	}

}

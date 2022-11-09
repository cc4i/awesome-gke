package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"tracker/tcp"
	"tracker/trip"

	"github.com/google/uuid"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// where Tracker is deployed
var whereami string

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func simpleWs(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()
	for {
		//Read Message from client
		mt, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		//If client message is ping will return pong
		if string(message) == "ping" {
			message = []byte("pong")
		}
		//Response message to client
		err = ws.WriteMessage(mt, message)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

// Get the services list of next call from environment variable
// eg: svc1,svc2,svc3
func getNextCall() string {
	nextCallee := os.Getenv("NEXT_CALLEE")
	log.Info().Str("NEXT_CALLEE", nextCallee).Send()
	return nextCallee
}

// API (/trip) - get TripDetail
func goTrip(c *gin.Context) {
	// Initial TripDetail with UUID
	td := &trip.TripDetail{
		Id: uuid.New().String(),
	}
	headers := make(map[string]string)

	headers["x-pod-name"] = c.Request.Header.Get("x-pod-name")
	headers["x-pod-namespace"] = c.Request.Header.Get("x-pod-namespace")
	headers["x-node-name"] = c.Request.Header.Get("x-node-name")
	headers["x-node-ip"] = c.Request.Header.Get("x-node-ip")
	headers["x-zone"] = c.Request.Header.Get("x-zone")
	headers["x-pod-ip"] = c.Request.Header.Get("x-pod-ip")
	headers["x-request-start"] = c.Request.Header.Get("x-request-start")

	// Get remote client IP if it's first call
	clientIp := c.Request.Header.Get("x-pod-ip")
	if clientIp == "" {
		//TODO: Get external IP when call from inside Pod/? Call API?
		clientIp = c.ClientIP()
	}
	whoami := c.Param("whoami")

	if err := td.GoTrip(whoami, headers, clientIp, c.Request.Method, c.Request.RequestURI, getNextCall(), whereami); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, td.Detail)
	}

}

// API (/initial) - get all pods as per request
func getInitialPods(c *gin.Context) {
	//Get initial cloud provider & namespace, then read /CRD::TrackerTop/ inside the namespace
	from := c.Param("from")
	whereami = from
	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Interface("err", err).Msg("Read request body")
	}
	log.Info().Str("request_str", string(buf)).Send()

	//reqeust string format : <provider>::<host:port>::<namespace>
	strs := strings.Split(string(buf), "::")
	ns := strs[2]

	// Initial TripDetail with UUID
	tp := &trip.TripDetail{
		Id: uuid.New().String(),
	}

	if err = tp.GetInitialPods(whereami, ns); err != nil {
		log.Error().Interface("error", err).Msg("getInitialPods()")
		c.JSON(http.StatusInternalServerError, err)
	} else {
		log.Info().Interface("return", tp.Detail).Msg("getInitialPods()")
		c.JSON(http.StatusOK, tp.Detail)
	}

}

// Load all trips from storage
func allTrips(c *gin.Context) {
	// Initial TripDetail with UUID
	tp := &trip.TripDetail{
		Id: uuid.New().String(),
	}
	if err := tp.TripHistory(); err != nil {
		log.Error().Interface("error", err).Msg("allTrips()")
		c.JSON(http.StatusInternalServerError, err)
	} else {
		log.Info().Interface("return", tp.Detail).Msg("allTrips()")
		c.JSON(http.StatusOK, tp.Detail)
	}

}

// Clear trip history in redis
func clearTrips(c *gin.Context) {
	// trip.S2Redis
	tp := &trip.TripDetail{}
	if err := tp.ClearTripHistory(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.String(http.StatusOK, "Cleared")
	}

}

// Test panic and reboot Pod
func doPanic(c *gin.Context) {
	log.Info().Str("exit_code", "255").Msg("Received panic request for demo")
	os.Exit(255)
}

// Echo function for general usage and return stuff as per request methods
func echo(c *gin.Context) {
	if c.Request.Method == "GET" {
		get := map[string]interface{}{
			"url":     c.Request.RequestURI,
			"method":  c.Request.Method,
			"headers": c.Request.Header,
			"queries": c.Request.URL.Query(),
		}
		c.JSON(http.StatusOK, get)
	} else if c.Request.Method == "POST" {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		post := map[string]interface{}{
			"url":     c.Request.RequestURI,
			"method":  c.Request.Method,
			"headers": c.Request.Header,
			"queries": string(buf),
		}
		c.JSON(http.StatusOK, post)
	} else {
		other := map[string]interface{}{
			"url":     c.Request.RequestURI,
			"method":  c.Request.Method,
			"headers": c.Request.Header,
		}
		c.JSON(http.StatusOK, other)
	}
}

func router(ctx context.Context, r *gin.Engine) *gin.Engine {
	log.Info().Interface("ctx", ctx).Msg("context.Context pairs")
	r.POST("/initial/:from", getInitialPods)
	r.GET("/trip", goTrip)
	r.GET("/trip/:whoami", goTrip)
	r.GET("/all-trips", allTrips)
	r.GET("/clear-trips", clearTrips)

	// ko add static assets under ./kodata - https://github.com/google/ko#static-assets
	if staticDir := os.Getenv("KO_DATA_PATH"); staticDir != "" {
		log.Info().Str("KO_DATA_PATH", staticDir).Send()
		r.Static("tracker-ui", staticDir)
	} else {
		r.Static("tracker-ui", "./kodata")
	}

	r.Any("/echo", echo)
	r.GET("/ws", simpleWs)
	r.GET("/panic", doPanic)

	return r
}

func httpSrv() {
	gin.DisableConsoleColor()
	server := gin.Default()

	//setup CORS policies
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	server.Use(cors.Default())

	//default port=8000
	if port := os.Getenv("POD_PORT"); port != "" {
		log.Fatal().Err(router(context.Background(), server).Run("0.0.0.0:" + port))
	} else {
		log.Fatal().Err(router(context.Background(), server).Run("0.0.0.0:8000"))
	}

}

func tcpSrv() {
	//default port=8008
	if port := os.Getenv("POD_TCP_PORT"); port != "" {
		tcp.Run(port)
	} else {
		tcp.Run("8008")
	}

}

func gameSrv() {
	tcp.RunBoringGame()
}

func main() {
	whereami = "gcp"
	//Configure zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	//Start TCP backend
	go tcpSrv()
	//Start HTTP backend
	go httpSrv()
	//Start boring Game Server based on Agones
	go gameSrv()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	s := <-shutdown
	log.Info().Msgf("Signal is %s", s)
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer handleTermination(cancel)

}

func handleTermination(cancel context.CancelFunc) {

}

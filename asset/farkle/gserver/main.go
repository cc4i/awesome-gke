// The game server follows rules from https://cardgames.io/farkle/, which
// uses gameing frameworks to demonstrate their capbilities on GKE.
//
// :: Game rules ::
//
//	Ones: Any die depicting a one. Worth 100 points each.
//	Fives: Any die depicting a five. Worth 50 points each.
//	Three Ones: A set of three dice depicting a one. worth 1000 points
//	Three Twos: A set of three dice depicting a two. worth 200 points
//	Three Threes: A set of three dice depicting a three. worth 300 points
//	Three Fours: A set of three dice depicting a four. worth 400 points
//	Three Fives: A set of three dice depicting a five. worth 500 points
//	Three Sixes: A set of three dice depicting a six. worth 600 points
//	Four of a kind: Any set of four dice depicting the same value. Worth 1000 points
//	Five of a kind: Any set of five dice depicting the same value. Worth 2000 points
//	Six of a kind: Any set of six dice depicting the same value. Worth 3000 points
//	Three Pairs: Any three sets of two pairs of dice. Includes having a four of a kind plus a pair. Worth 1500 points
//	Run: Six dice in a sequence (1,2,3,4,5,6). Worth 2500 points
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"gserver/game"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	coresdk "agones.dev/agones/pkg/sdk"
	"agones.dev/agones/pkg/util/signals"
	sdk "agones.dev/agones/sdks/go"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// EO: ///////////
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//BO: //////////

func main() {
	port := flag.String("port", "7654", "The port to listen to traffic on")
	passthrough := flag.Bool("passthrough", false, "Get listening port from the SDK, rather than use the 'port' value")
	readyOnStart := flag.Bool("ready", true, "Mark this GameServer as Ready on startup")
	shutdownDelayMin := flag.Int("automaticShutdownDelayMin", 0, "[Deprecated] If greater than zero, automatically shut down the server this many minutes after the server becomes allocated (please use automaticShutdownDelaySec instead)")
	shutdownDelaySec := flag.Int("automaticShutdownDelaySec", 0, "If greater than zero, automatically shut down the server this many seconds after the server becomes allocated (cannot be used if automaticShutdownDelayMin is set)")
	readyDelaySec := flag.Int("readyDelaySec", 0, "If greater than zero, wait this many seconds each time before marking the game server as ready")
	readyIterations := flag.Int("readyIterations", 0, "If greater than zero, return to a ready state this number of times before shutting down")
	enablePlayerTracking := flag.Bool("player-tracking", true, "If true, player tracking will be enabled.")
	playerCapacity := flag.Int64("player-capacity", 10, "The capacity of gameserver, default is 10 (Alpha, enable player tracking)")
	flag.Parse()

	// 1. Listen to singal
	go doSignal()

	// 2. Initial a SDK instance
	log.Print("Creating SDK instance")
	s, err := sdk.NewSDK()
	if err != nil {
		log.Fatalf("Could not connect to sdk: %v", err)
	}

	// 3. Send health ping to SDK
	log.Print("Starting Health Ping")
	ctx, cancel := context.WithCancel(context.Background())
	go doHealth(s, ctx)

	// 4. Player Tracking
	if *enablePlayerTracking {
		if err = s.Alpha().SetPlayerCapacity(*playerCapacity); err != nil {
			log.Fatalf("could not set play count: %v", err)
		}
	}

	// 5. Create a Game Server
	if *passthrough {
		var gs *coresdk.GameServer
		gs, err = s.GameServer()
		if err != nil {
			log.Fatalf("Could not get gameserver port details: %s", err)
		}

		p := strconv.FormatInt(int64(gs.Status.Ports[0].Port), 10)
		port = &p
	}

	// 6. The Game Server starts to listen
	go httpServe(port, s, cancel)

	// 7. Configure shutdown deplay
	if *shutdownDelaySec > 0 {
		shutdownAfterNAllocations(s, *readyIterations, *shutdownDelaySec)
	} else if *shutdownDelayMin > 0 {
		shutdownAfterNAllocations(s, *readyIterations, *shutdownDelayMin*60)
	}

	// 8. Waiting to ready
	if *readyOnStart {
		if *readyDelaySec > 0 {
			log.Printf("Waiting %d seconds before moving to ready", *readyDelaySec)
			time.Sleep(time.Duration(*readyDelaySec) * time.Second)
		}
		log.Print("Marking this server as ready")
		ready(s)
	}

	// Prevent the program from quitting as the server is listening on goroutines.
	for {
	}

	fmt.Println("Dedicated Game Server for Farkel")
}

// doSignal shutsdown on SIGTERM/SIGKILL
func doSignal() {
	ctx := signals.NewSigKillContext()
	<-ctx.Done()
	log.Println("Exit signal received. Shutting down.")
	os.Exit(0)
}

// doHealth sends the regular Health Pings
func doHealth(sdk *sdk.SDK, ctx context.Context) {
	tick := time.Tick(5 * time.Second)
	for {
		log.Printf("Health Ping")
		err := sdk.Health()
		if err != nil {
			log.Fatalf("Could not send health ping, %v", err)
		}
		select {
		case <-ctx.Done():
			log.Print("Stopped health pings")
			return
		case <-tick:
		}
	}
}

// shutdownAfterNAllocations creates a callback to automatically shut down
// the server a specified number of seconds after the server becomes
// allocated the Nth time.
//
// The algorithm is:
//
//  1. Move the game server back to ready N times after it is allocated
//  2. Shutdown the game server after the Nth time is becomes allocated
//
// This follows the integration pattern documented on the website at
// https://agones.dev/site/docs/integration-patterns/reusing-gameservers/
func shutdownAfterNAllocations(s *sdk.SDK, readyIterations, shutdownDelaySec int) {
	gs, err := s.GameServer()
	if err != nil {
		log.Fatalf("Could not get game server: %v", err)
	}
	log.Printf("Initial game Server state = %s", gs.Status.State)

	m := sync.Mutex{} // protects the following two variables
	lastAllocated := gs.ObjectMeta.Annotations["agones.dev/last-allocated"]
	remainingIterations := readyIterations

	if err := s.WatchGameServer(func(gs *coresdk.GameServer) {
		m.Lock()
		defer m.Unlock()
		la := gs.ObjectMeta.Annotations["agones.dev/last-allocated"]
		log.Printf("Watch Game Server callback fired. State = %s, Last Allocated = %q", gs.Status.State, la)
		if lastAllocated != la {
			log.Println("Game Server Allocated")
			lastAllocated = la
			remainingIterations--
			// Run asynchronously
			go func(iterations int) {
				time.Sleep(time.Duration(shutdownDelaySec) * time.Second)

				if iterations > 0 {
					log.Println("Moving Game Server back to Ready")
					readyErr := s.Ready()
					if readyErr != nil {
						log.Fatalf("Could not set game server to ready: %v", readyErr)
					}
					log.Println("Game Server is Ready")
					return
				}

				log.Println("Moving Game Server to Shutdown")
				if shutdownErr := s.Shutdown(); shutdownErr != nil {
					log.Fatalf("Could not shutdown game server: %v", shutdownErr)
				}
				// The process will exit when Agones removes the pod and the
				// container receives the SIGTERM signal
				return
			}(remainingIterations)
		}
	}); err != nil {
		log.Fatalf("Could not watch Game Server events, %v", err)
	}
}

func handleResponse(txt string, s *sdk.SDK, cancel context.CancelFunc) (response string, responseError error) {
	parts := strings.Split(strings.TrimSpace(txt), " ")
	response = txt
	responseError = nil

	switch parts[0] {
	// All actions from Farkle client
	case "FARKLE_ACTION":
		if len(parts) != 2 {
			response = "Invalid FARKLE_ACTION, should have 1 argument"
			responseError = fmt.Errorf("Invalid FARKLE_ACTION, should have 1 argument")
			return
		}
		response, _ = game.FarkleHandler(parts[1])
		return

	// shuts down the gameserver
	case "EXIT":
		// handle elsewhere, as we respond before exiting
		return

	// turns off the health pings
	case "UNHEALTHY":
		cancel()

	case "GAMESERVER":
		response = gameServerName(s)

	case "READY":
		ready(s)

	case "ALLOCATE":
		allocate(s)

	case "RESERVE":
		if len(parts) != 2 {
			response = "Invalid RESERVE, should have 1 argument"
			responseError = fmt.Errorf("Invalid RESERVE, should have 1 argument")
		}
		if dur, err := time.ParseDuration(parts[1]); err != nil {
			response = fmt.Sprintf("%s\n", err)
			responseError = err
		} else {
			reserve(s, dur)
		}

	case "WATCH":
		watchGameServerEvents(s)

	case "LABEL":
		switch len(parts) {
		case 1:
			// legacy format
			setLabel(s, "timestamp", strconv.FormatInt(time.Now().Unix(), 10))
		case 3:
			setLabel(s, parts[1], parts[2])
		default:
			response = "Invalid LABEL command, must use zero or 2 arguments"
			responseError = fmt.Errorf("Invalid LABEL command, must use zero or 2 arguments")
		}

	case "CRASH":
		log.Print("Crashing.")
		os.Exit(1)
		return "", nil

	case "ANNOTATION":
		switch len(parts) {
		case 1:
			// legacy format
			setAnnotation(s, "timestamp", time.Now().UTC().String())
		case 3:
			setAnnotation(s, parts[1], parts[2])
		default:
			response = "Invalid ANNOTATION command, must use zero or 2 arguments"
			responseError = fmt.Errorf("Invalid ANNOTATION command, must use zero or 2 arguments")
		}

	case "PLAYER_CAPACITY":
		switch len(parts) {
		case 1:
			response = getPlayerCapacity(s)
		case 2:
			if cap, err := strconv.Atoi(parts[1]); err != nil {
				response = fmt.Sprintf("%s", err)
				responseError = err
			} else {
				setPlayerCapacity(s, int64(cap))
			}
		default:
			response = "Invalid PLAYER_CAPACITY, should have 0 or 1 arguments"
			responseError = fmt.Errorf("Invalid PLAYER_CAPACITY, should have 0 or 1 arguments")
		}

	case "PLAYER_CONNECT":
		if len(parts) < 2 {
			response = "Invalid PLAYER_CONNECT, should have 1 arguments"
			responseError = fmt.Errorf("Invalid PLAYER_CONNECT, should have 1 arguments")
			return
		}
		playerConnect(s, parts[1])

	case "PLAYER_DISCONNECT":
		if len(parts) < 2 {
			response = "Invalid PLAYER_DISCONNECT, should have 1 arguments"
			responseError = fmt.Errorf("Invalid PLAYER_DISCONNECT, should have 1 arguments")
			return
		}
		playerDisconnect(s, parts[1])

	case "PLAYER_CONNECTED":
		if len(parts) < 2 {
			response = "Invalid PLAYER_CONNECTED, should have 1 arguments"
			responseError = fmt.Errorf("Invalid PLAYER_CONNECTED, should have 1 arguments")
			return
		}
		response = playerIsConnected(s, parts[1])

	case "GET_PLAYERS":
		response = getConnectedPlayers(s)

	case "PLAYER_COUNT":
		response = getPlayerCount(s)
	}

	return
}

// BO: Websocket listener ////////////
func wsListener(s *sdk.SDK, cancel context.CancelFunc) gin.HandlerFunc {

	fn := func(c *gin.Context) {
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

			response, err := handleResponse(string(message), s, cancel)
			if err != nil {
				response = "ERROR: " + response + "\n"
			} else {
				response = "ACK TCP: " + response + "\n"

			}
			err = ws.WriteMessage(mt, []byte(response))
			if err != nil {
				fmt.Println(err)
				break
			}
			if string(message) == "EXIT" {
				exit(s)
			}

		}

	}
	return gin.HandlerFunc(fn)

}

func router(r *gin.Engine, s *sdk.SDK, cancel context.CancelFunc) *gin.Engine {
	r.GET("/ws", wsListener(s, cancel))
	return r
}

func httpServe(port *string, s *sdk.SDK, cancel context.CancelFunc) {
	gin.DisableConsoleColor()
	server := gin.Default()

	//setup CORS policies
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	server.Use(cors.Default())

	log.Fatal(router(server, s, cancel).Run(":" + *port))

}

// EO: ////////////
// ready attempts to mark this gameserver as ready
func ready(s *sdk.SDK) {
	err := s.Ready()
	if err != nil {
		log.Fatalf("Could not send ready message")
	}
}

// allocate attempts to allocate this gameserver
func allocate(s *sdk.SDK) {
	err := s.Allocate()
	if err != nil {
		log.Fatalf("could not allocate gameserver: %v", err)
	}
}

// reserve for 10 seconds
func reserve(s *sdk.SDK, duration time.Duration) {
	if err := s.Reserve(duration); err != nil {
		log.Fatalf("could not reserve gameserver: %v", err)
	}
}

// exit shutdowns the server
func exit(s *sdk.SDK) {
	log.Printf("Received EXIT command. Exiting.")
	// This tells Agones to shutdown this Game Server
	shutdownErr := s.Shutdown()
	if shutdownErr != nil {
		log.Printf("Could not shutdown")
	}
	// The process will exit when Agones removes the pod and the
	// container receives the SIGTERM signal
}

// gameServerName returns the GameServer name
func gameServerName(s *sdk.SDK) string {
	var gs *coresdk.GameServer
	gs, err := s.GameServer()
	if err != nil {
		log.Fatalf("Could not retrieve GameServer: %v", err)
	}
	var j []byte
	j, err = json.Marshal(gs)
	if err != nil {
		log.Fatalf("error mashalling GameServer to JSON: %v", err)
	}
	log.Printf("GameServer: %s \n", string(j))
	return "NAME: " + gs.ObjectMeta.Name + "\n"
}

// watchGameServerEvents creates a callback to log when
// gameserver events occur
func watchGameServerEvents(s *sdk.SDK) {
	err := s.WatchGameServer(func(gs *coresdk.GameServer) {
		j, err := json.Marshal(gs)
		if err != nil {
			log.Fatalf("error mashalling GameServer to JSON: %v", err)
		}
		log.Printf("GameServer Event: %s \n", string(j))
	})
	if err != nil {
		log.Fatalf("Could not watch Game Server events, %v", err)
	}
}

// setAnnotation sets a given annotation
func setAnnotation(s *sdk.SDK, key, value string) {
	log.Printf("Setting annotation %v=%v", key, value)
	err := s.SetAnnotation(key, value)
	if err != nil {
		log.Printf("could not set annotation: %v", err)
	}
}

// setLabel sets a given label
func setLabel(s *sdk.SDK, key, value string) {
	log.Printf("Setting label %v=%v", key, value)
	// label values can only be alpha, - and .
	err := s.SetLabel(key, value)
	if err != nil {
		log.Printf("could not set label: %v", err)
	}
}

// setPlayerCapacity sets the player capacity to the given value
func setPlayerCapacity(s *sdk.SDK, capacity int64) {
	log.Printf("Setting Player Capacity to %d", capacity)
	if err := s.Alpha().SetPlayerCapacity(capacity); err != nil {
		log.Printf("could not set capacity: %v", err)
	}
}

// getPlayerCapacity returns the current player capacity as a string
func getPlayerCapacity(s *sdk.SDK) string {
	log.Print("Getting Player Capacity")
	capacity, err := s.Alpha().GetPlayerCapacity()
	if err != nil {
		log.Printf("could not get capacity: %v", err)
	}
	return strconv.FormatInt(capacity, 10) + "\n"
}

// playerConnect connects a given player
func playerConnect(s *sdk.SDK, id string) {
	log.Printf("Connecting Player: %s", id)
	if _, err := s.Alpha().PlayerConnect(id); err != nil {
		log.Printf("could not connect player: %v", err)
	}
}

// playerDisconnect disconnects a given player
func playerDisconnect(s *sdk.SDK, id string) {
	log.Printf("Disconnecting Player: %s", id)
	if _, err := s.Alpha().PlayerDisconnect(id); err != nil {
		log.Printf("could not disconnect player: %v", err)
	}
}

// playerIsConnected returns a bool as a string if a player is connected
func playerIsConnected(s *sdk.SDK, id string) string {
	log.Printf("Checking if player %s is connected", id)

	connected, err := s.Alpha().IsPlayerConnected(id)
	if err != nil {
		log.Printf("could not retrieve if player is connected: %v", err)
	}

	return strconv.FormatBool(connected) + "\n"
}

// getConnectedPlayers returns a comma delimeted list of connected players
func getConnectedPlayers(s *sdk.SDK) string {
	log.Print("Retrieving connected player list")
	list, err := s.Alpha().GetConnectedPlayers()
	if err != nil {
		log.Printf("could not retrieve connected players: %s", err)
	}

	return strings.Join(list, ",") + "\n"
}

// getPlayerCount returns the count of connected players as a string
func getPlayerCount(s *sdk.SDK) string {
	log.Print("Retrieving connected player count")
	count, err := s.Alpha().GetPlayerCount()
	if err != nil {
		log.Printf("could not retrieve player count: %s", err)
	}
	return strconv.FormatInt(count, 10) + "\n"
}

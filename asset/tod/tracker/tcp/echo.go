package tcp

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/panjf2000/gnet/v2"
)

type echoServer struct {
	gnet.BuiltinEventEngine

	eng       gnet.Engine
	addr      string
	multicore bool
	sessions  map[string]string
}

func (es *echoServer) OnBoot(eng gnet.Engine) gnet.Action {
	es.eng = eng
	log.Info().Msgf("TCP Server with multi-core=%t is listening on %s\n", es.multicore, es.addr)
	return gnet.None
}

func (es *echoServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	remote := c.RemoteAddr()
	log.Info().Msgf("Connection from %s", remote.String())
	if es.sessions == nil {
		es.sessions = make(map[string]string)
	}
	es.sessions[remote.String()] = remote.Network()
	return []byte{}, gnet.None
}

func (es *echoServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf, _ := c.Next(-1)
	log.Info().Str("msg", string(buf)).Msg("Received message")
	if string(buf) == "bye\r\n" {
		log.Info().Msg("Received break signal and disconnected with client.")
		c.Close()
	}
	switch string(buf) {
	case "bye":
	case "bye\r\n":
		log.Info().Msg("Received break signal and disconnected with client.")
		c.Close()
		return gnet.None
	case "list":
	case "list\r\n":
		log.Info().Msg("All clients:")
		buf, _ := json.Marshal(es.sessions)
		c.Write(buf)
		return gnet.None

	}
	// built response message
	node := os.Getenv("POD_NODE_NAME")
	pod := os.Getenv("POD_NAME")
	ver := os.Getenv("TRACKER_VERSION")
	rs := fmt.Sprintf("[%s]:[%s] in [%s] has received: [%s]", pod, ver, node, buf)
	//

	c.Write([]byte(rs))
	log.Info().Str("msg", string(buf)).Msg("Responsed message")
	return gnet.None
}

func Run(port string) {
	echo := &echoServer{addr: fmt.Sprintf("tcp://:%s", port), multicore: true}
	log.Fatal().Err(gnet.Run(echo, echo.addr, gnet.WithMulticore(true)))
}

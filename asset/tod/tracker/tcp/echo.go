package tcp

import (
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
}

func (es *echoServer) OnBoot(eng gnet.Engine) gnet.Action {
	es.eng = eng
	log.Info().Msgf("echo server with multi-core=%t is listening on %s\n", es.multicore, es.addr)
	return gnet.None
}

func (es *echoServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf, _ := c.Next(-1)
	log.Info().Str("msg", string(buf)).Msg("Received message")
	if string(buf) == "bye\r\n" {
		log.Info().Msg("Received break signal and disconnected with client.")
		c.Close()
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

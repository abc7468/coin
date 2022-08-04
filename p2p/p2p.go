package p2p

import (
	"coin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func Upgrade(c *gin.Context) {
	openPort := c.Query("openPort")
	ip := utils.GetSplitedStrings(c.Request.RemoteAddr, ":", 0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return ip != "" && openPort != ""
	}
	fmt.Printf("%s want to Upgarde\n", openPort)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	utils.HandleErr(err)

	initPeer(conn, ip, openPort)

}

func AddPeer(address, port, openPort string) {
	fmt.Printf("%s want to connect to port %s\n", openPort, port)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil)
	utils.HandleErr(err)
	p := initPeer(conn, address, port)
	sendNewestBlock(p)
}

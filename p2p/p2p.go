package p2p

import (
	"coin/blockchain"
	"coin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func Upgrade(c *gin.Context) {
	openPort := c.Query("openPort")
	ip := utils.Splitter(c.Request.RemoteAddr, ":", 0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return ip != "" && openPort != ""
	}
	fmt.Printf("%s want to Upgarde\n", openPort)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	utils.HandleErr(err)

	initPeer(conn, ip, openPort)

}

func AddPeer(address, port, openPort string, broadcast bool) {
	fmt.Printf("%s want to connect to port %s\n", openPort, port)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil)
	utils.HandleErr(err)
	p := initPeer(conn, address, port)
	if broadcast {
		BroadcastNewPeer(p)
	}
	sendNewestBlock(p)

}
func notifyNewPeer(address string, p *peer) {
	m := makeMessage(MessageNewPeerNotify, address)
	p.inbox <- m
}
func BroadcastNewBlock(b *blockchain.Block) {
	Peers.mu.Lock()
	defer Peers.mu.Unlock()
	for _, p := range Peers.v {
		notifyNewBlock(b, p)
	}
}

func BroadcastNewTx(tx *blockchain.Tx) {
	for _, p := range Peers.v {
		notifyNewTx(tx, p)
	}
}

func BroadcastNewPeer(newPeer *peer) {
	for key, p := range Peers.v {
		if key != newPeer.key {
			payload := fmt.Sprintf("%s:%s", newPeer.key, p.port)
			notifyNewPeer(payload, p)
		}
	}
}

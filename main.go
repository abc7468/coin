package main

import (
	"coin/blockchain"
	"github.com/gin-gonic/gin"
	"net/http"
)

const port string = ":4000"

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func showHomePage(c *gin.Context) {
	data := homeData{"roy coin", blockchain.GetBlockchain().AllBlocks()}
	c.HTML(http.StatusOK, "home", data)
}
func showAddPage(c *gin.Context) {
	c.HTML(http.StatusOK, "add", nil)
}
func addBlock(c *gin.Context) {
	c.Request.ParseForm()
	data := c.PostForm("data")
	blockchain.GetBlockchain().AddBlock(data)
	c.Redirect(308, "http://localhost"+port)
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")
	r.GET("/", showHomePage)
	r.GET("/add", showAddPage)
	r.POST("/add", addBlock)
	r.Run(port)
}

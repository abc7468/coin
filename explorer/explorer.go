package explorer

import (
	"coin/blockchain"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

var port string

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func showHomePage(c *gin.Context) {
	//data := homeData{"roy coin", blockchain.GetBlockchain().AllBlocks()}
	c.HTML(http.StatusOK, "home", nil)
}
func showAddPage(c *gin.Context) {
	c.HTML(http.StatusOK, "add", nil)
}

// Welcome godoc
// @Summary blockchain에 새로운 block을 추가
// @Description 자세한 설명은 이곳에 적습니다.
// @name add-block
// @Accept  json
// @Produce  json
// @Param data path string true "block's data"
// @Router /add [POST]
// @Success 301
func addBlock(c *gin.Context) {
	c.Request.ParseForm()
	data := c.PostForm("data")
	blockchain.Blockchain().AddBlock(data)
	c.Redirect(http.StatusMovedPermanently, "home")
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")

	r.GET("/home", showHomePage)
	r.GET("/blocks", showAddPage)
	r.POST("/blocks", addBlock)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(port)
}

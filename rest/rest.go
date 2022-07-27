package rest

import (
	"coin/blockchain"
	"coin/docs"
	"coin/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

type addBlockBody struct {
	Message string `json:"message"`
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

// Welcome godoc
// @Summary 현재 blockchain의 모든 Block을 출력.
// @Description blockchain의 현 상태를 출력.
// @name show-blocks
// @Accept  json
// @Produce  json
// @Router /blocks [GET]
// @Success 200
func showBlocks(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	json.NewEncoder(c.Writer).Encode(blockchain.Blockchain().Blocks())
}

// Welcome godoc
// @Summary blockchain에 새로운 block을 추가
// @Description 자세한 설명은 이곳에 적습니다.
// @name add-block
// @Accept  json
// @Produce  json
// @Param message body string true "추가하고자 하는 Block의 Data"
// @Router /blocks [POST]
// @Success 201
func addBlocks(c *gin.Context) {
	var addBlockBody addBlockBody
	err := json.NewDecoder(c.Request.Body).Decode(&addBlockBody)
	utils.HandleErr(err)
	blockchain.Blockchain().AddBlock(addBlockBody.Message)
	c.Writer.WriteHeader(http.StatusCreated)
}

func getBlock(c *gin.Context) {
	hash := c.Param("hash")
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(c.Writer)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", "http://localhost:4000")
		c.Header("Access-Control-Allow-Methods", "GET, DELETE, POST")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func Start(port int) {

	r := gin.Default()
	docs.SwaggerInfo.Description = "This is a sample server for Swagger."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", port)
	docs.SwaggerInfo.Title = "Swagger Example Test"

	r.Use(CORSMiddleware())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/blocks", showBlocks)
	r.POST("/blocks", addBlocks)
	r.GET("/blocks/:hash", getBlock)
	r.Run(fmt.Sprintf(":%d", port))
}

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

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type addTxPayload struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
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
	blockchain.Blockchain().AddBlock()
	c.Writer.WriteHeader(http.StatusCreated)

}

// Welcome godoc
// @Summary 하나의 블록 데이터를 출력
// @Description 입력 해쉬값을 key로 가진 하나의 블록 데이터를 출력합니다.
// @name show-block
// @Accept  json
// @Produce  json
// @Param hash path string true "추가하고자 하는 Block의 Data"
// @Router /blocks/{hash} [Get]
// @Success 200
func getBlock(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	hash := c.Param("hash")
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(c.Writer)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

// Welcome godoc
// @Summary 블록체인의 현 상태 출력
// @Description 블록체인에 포함된 블록들의 정보를 출력합니다.
// @name show-blocks
// @Accept  json
// @Produce  json
// @Router /status [Get]
// @Success 200
func showBlockchain(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	json.NewEncoder(c.Writer).Encode(blockchain.Blockchain())
}

// Welcome godoc
// @Summary Address를 통해 balance를 출력
// @Description 입력한 Address의 Balance를 출력합니다.
// @name balance
// @Accept  json
// @Produce  json
// @Param address path string true "Balance를 확인하고자하는 address"
// @Router /balance/{address} [Get]
// @Success 200
func getBalance(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	address := c.Param("address")
	encoder := json.NewEncoder(c.Writer)
	total := c.Query("total")
	switch total {
	case "true":
		amount := blockchain.Blockchain().BalanceByAddress(address)
		json.NewEncoder(c.Writer).Encode(balanceResponse{address, amount})
	default:
		utils.HandleErr(encoder.Encode(blockchain.Blockchain().UTxOutsByAddress(address)))
	}

}

// Welcome godoc
// @Summary Address를 통해 balance를 출력
// @Description 입력한 Address의 Balance를 출력합니다.
// @name balance
// @Accept  json
// @Produce  json
// @Param address path string true "Balance를 확인하고자하는 address"
// @Router /mempool [Get]
// @Success 200
func mempool(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	utils.HandleErr(json.NewEncoder(c.Writer).Encode(blockchain.Mempool.Txs))
}

func transactions(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var payload addTxPayload
	utils.HandleErr(json.NewDecoder(c.Request.Body).Decode(&payload))
	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		json.NewEncoder(c.Writer).Encode(errorResponse{"not enough funds"})
	}
	c.Writer.WriteHeader(http.StatusCreated)

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
	r.GET("/status", showBlockchain)
	r.GET("/blocks/:hash", getBlock)
	r.GET("/balance/:address", getBalance)
	r.GET("/mempool", mempool)
	r.POST("/transactions", transactions)
	r.POST("/blocks", addBlocks)

	r.Run(fmt.Sprintf(":%d", port))
}

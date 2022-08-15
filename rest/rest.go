package rest

import (
	"bytes"
	"coin/blockchain"
	"coin/p2p"
	"coin/utils"
	"coin/wallet"
	"github.com/whatap/go-api/httpc"
	"strings"

	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/whatap/go-api/instrumentation/github.com/gin-gonic/gin/whatapgin"
	"github.com/whatap/go-api/method"
	"github.com/whatap/go-api/trace"
	"io/ioutil"
	"net/http"
	"time"
)

var port string

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

type walletRes struct {
	Address string `json:"address"`
}

type addPeerPayload struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

func getUser(ctx context.Context) {
	methodCtx, _ := method.Start(ctx, "getUser")
	defer method.End(methodCtx, nil)
	time.Sleep(time.Duration(1) * time.Second)
}

func httpGet(callUrl string) (int, string, error) {
	fmt.Println("httpGet ", callUrl)
	// GET 호출
	if resp, err := http.Get(callUrl); err == nil {
		defer resp.Body.Close()
		fmt.Println("status=", resp.StatusCode)

		// 결과 출력
		if data, err := ioutil.ReadAll(resp.Body); err == nil {
			return resp.StatusCode, string(data), err
		} else {
			return resp.StatusCode, "", err
		}

	} else {
		fmt.Println(err)
		return -1, "", err
	}
}
func httpWithRequest(method string, callUrl string, body string, headers http.Header) (int, string, error) {
	fmt.Println("httpGetWithRequest ", method, ", ", callUrl, ", ", body, ", ", headers)
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	if req, err := http.NewRequest(strings.ToUpper(method), callUrl, bytes.NewBufferString(body)); err == nil {
		if headers != nil {
			for key, _ := range headers {
				req.Header.Add(key, headers.Get(key))
			}
		}
		if resp, err := client.Do(req); err == nil {
			defer resp.Body.Close()
			if data, err := ioutil.ReadAll(resp.Body); err == nil {
				fmt.Println("status=", resp.StatusCode)
				return resp.StatusCode, string(data), err
			} else {
				fmt.Println("Read response Error ", err)
				return resp.StatusCode, "", err
			}
		} else {
			fmt.Println("client.Do Error ", err)
			return -2, "", err
		}

	} else {
		fmt.Println("NewRequest Error ", err)
		return -1, "", err
	}
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

	blockchain.Status(blockchain.Blockchain(), c)
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
	newBlock := blockchain.Blockchain().AddBlock()
	p2p.BroadcastNewBlock(newBlock)
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
	fmt.Println(c.Request.Host)
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
		amount := blockchain.BalanceByAddress(address, blockchain.Blockchain())
		json.NewEncoder(c.Writer).Encode(balanceResponse{address, amount})
	default:
		utils.HandleErr(encoder.Encode(blockchain.UTxOutsByAddress(address, blockchain.Blockchain())))
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
	utils.HandleErr(json.NewEncoder(c.Writer).Encode(blockchain.Mempool().Txs))
}

func transactions(c *gin.Context) {
	var payload addTxPayload
	utils.HandleErr(json.NewDecoder(c.Request.Body).Decode(&payload))
	tx, err := blockchain.Mempool().AddTx(payload.To, payload.Amount)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(errorResponse{err.Error()})
	}
	p2p.BroadcastNewTx(tx)
	c.Writer.WriteHeader(http.StatusCreated)

}

func getWallet(c *gin.Context) {
	address := wallet.Wallet().Address
	json.NewEncoder(c.Writer).Encode(walletRes{address})
}

func addPeer(c *gin.Context) {
	var payload addPeerPayload
	utils.HandleErr(json.NewDecoder(c.Request.Body).Decode(&payload))

	p2p.AddPeer(payload.Address, payload.Port, port[1:], true)
	c.Writer.WriteHeader(http.StatusOK)
}

func showPeers(c *gin.Context) {
	json.NewEncoder(c.Writer).Encode(p2p.AllPeers(&p2p.Peers))
}

func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(c.Request.URL)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, DELETE, POST")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func Start(aPort int) {

	port = fmt.Sprintf(":%d", aPort)
	r := gin.Default()
	//docs.SwaggerInfo.Description = "This is a sample server for Swagger."
	//docs.SwaggerInfo.Version = "1.0"
	//docs.SwaggerInfo.Host = fmt.Sprintf("localhost%s", port)
	//docs.SwaggerInfo.Title = "Swagger Example Test"
	config := make(map[string]string)
	trace.Init(config)
	defer trace.Shutdown()

	r.Use(loggerMiddleware(), whatapgin.Middleware())
	r.GET("/", func(c *gin.Context) {
		ctx := c.Request.Context()
		fmt.Println("Request -", c.Request)

		callUrl := "http://localhost:4000/blocks"
		httpcCtx, _ := httpc.Start(ctx, callUrl)
		var buffer bytes.Buffer
		if statusCode, data, err := httpWithRequest("GET", callUrl, "", httpc.GetMTrace(httpcCtx)); err == nil {
			httpc.End(httpcCtx, statusCode, "", nil)
			buffer.WriteString(fmt.Sprintln("httpc callUrl=", callUrl, ", statuscode=", statusCode, ", data=", data))
		} else {
			httpc.End(httpcCtx, -1, "", err)
			buffer.WriteString(fmt.Sprintln("httpc Error callUrl=", callUrl, ", err=", err))
		}

		trace.Step(ctx, "Text Message 2", "Message2", 6, 6)
		c.JSON(http.StatusOK, gin.H{
			"message": string(buffer.Bytes()),
		})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/blocks", showBlocks)
	r.GET("/status", showBlockchain)
	r.GET("/blocks/:hash", getBlock)
	r.GET("/balance/:address", getBalance)
	r.GET("/mempool", mempool)
	r.GET("/wallet", getWallet)
	r.GET("/ws", p2p.Upgrade)
	r.GET("/peers", showPeers)
	r.POST("/transactions", transactions)
	r.POST("/blocks", addBlocks)
	r.POST("/peers", addPeer)

	r.Run(port)
}

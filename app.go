package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/checksum0/go-electrum/electrum"
	"github.com/gin-gonic/gin"
)

func setupServer() *gin.Engine {

	paymentAddress := os.Getenv("KEVA_PAYMENT_ADDRESS")
	minPayment := os.Getenv("KEVA_MIN_PAYMENT")

	if len(paymentAddress) == 0 || len(minPayment) == 0 {
		log.Fatalln("Environment variable KEVA_PAYMENT_ADDRESS and KEVA_MIN_PAYMENT required.")
	}

	minPaymentVal, err := strconv.ParseFloat(minPayment, 64)
	if err != nil {
		log.Fatalln("Invalid minimal payment value.")
	}

	electrumServer := electrum.NewServer()
	if err = electrumServer.ConnectTCP("127.0.0.1:50001"); err != nil {
		log.Fatal(err)
	}

	// Timed "server.ping" call to prevent disconnection.
	go func() {
		for {
			if err = electrumServer.Ping(); err != nil {
				log.Fatal(err)
			}
			time.Sleep(60 * time.Second)
		}
	}()

	router := gin.New()
	v1 := router.Group("/v1")
	{
		// Get payment info
		v1.GET("/payment_info", func(c *gin.Context) {
			getPaymentInfo(c, paymentAddress, minPaymentVal)
		})

		// Upload media
		v1.POST("/media", func(c *gin.Context) {
			uploadMedia(c)
		})

		// Add to IPFS
		v1.POST("/pin", func(c *gin.Context) {
			publishMediaIPFS(electrumServer, c, paymentAddress, minPaymentVal)
		})
	}

	return router
}

func main() {
	router := setupServer()
	port := os.Getenv("KEVA_IPFS_PORT")
	if len(port) == 0 {
		log.Fatalln("Invalid port.")
	}
	router.Run(":" + port)
}

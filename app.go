package main

import (
	"crypto/tls"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/checksum0/go-electrum/electrum"
	"github.com/gin-gonic/gin"
)

func setupServer(port string, isSSL bool) *gin.Engine {

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

	if isSSL {
		conf := &tls.Config{}
		if err = electrumServer.ConnectSSL("127.0.0.1:"+port, conf); err != nil {
			log.Fatal(err)
		}
	} else {
		if err = electrumServer.ConnectTCP("127.0.0.1:" + port); err != nil {
			log.Fatal(err)
		}
	}

	// Timed "server.ping" call to prevent disconnection.
	go func() {
		for {
			if err = electrumServer.Ping(); err != nil {
				log.Println(err)
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
	portSSLStr := os.Getenv("KEVA_ELECTRUM_SSL_PORT")
	portTCPStr := os.Getenv("KEVA_ELECTRUM_TCP_PORT")
	if len(portSSLStr) == 0 && len(portTCPStr) == 0 {
		log.Fatalln("Either KEVA_ELECTRUM_SSL_PORT or KEVA_ELECTRUM_TCP_PORT must be set.")
	}
	var port int
	var router *gin.Engine
	if len(portSSLStr) > 0 {
		router = setupServer(portSSLStr, true)
		port, _ = strconv.Atoi(portSSLStr)
	} else {
		router = setupServer(portTCPStr, false)
		port, _ = strconv.Atoi(portTCPStr)
	}
	// The port used by the server is the eletrumx port plus 10.
	port += 10

	tlsEnabled := 0
	tlsEnabled, _ = strconv.Atoi(os.Getenv("KEVA_TLS_ENABLED"))
	if tlsEnabled != 0 {
		log.Println("Using TLS/SSL")
	} else {
		log.Println("**Warning: TLS/SSL not enabled. Set KEVA_TLS_ENABLED to 1 to enable TLS/SSL.")
	}

	if tlsEnabled != 0 {
		serverCert := os.Getenv("KEVA_TLS_CERT")
		serverKey := os.Getenv("KEVA_TLS_KEY")
		if len(serverCert) == 0 || len(serverKey) == 0 {
			log.Fatalln("Environment variables KEVA_TLS_CERT and KEVA_TLS_KEY required.")
		}
		router.RunTLS(":"+strconv.Itoa(port), serverCert, serverKey)
	} else {
		router.Run(":" + strconv.Itoa(port))
	}
}

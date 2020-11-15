package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/checksum0/go-electrum/electrum"
	"github.com/gin-gonic/gin"
)

const (
	sessionKey = "session"
)

func getPaymentInfo(c *gin.Context, paymentAddr string, minPayment float64) {
	paymentInfo := PaymentInfo{
		PaymentAddress: paymentAddr,
		MinPayment:     fmt.Sprintf("%f", minPayment),
	}
	c.JSON(200, &paymentInfo)
	return
}

func uploadMedia(c *gin.Context) {
	f, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}
	defer f.Close()

	tmpfile, err := ioutil.TempFile("", "CID")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create tmp file",
			"error":   true,
		})
		return
	}

	r := bufio.NewReader(f)
	w := bufio.NewWriter(tmpfile)
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := w.Write(buf[:n]); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error writing to tmp file",
				"error":   true,
			})
			return
		}
	}

	if err = w.Flush(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error flushing to tmp file",
			"error":   true,
		})
		return
	}

	// Get the CID of the tmp file.
	tmpfileName := tmpfile.Name()
	cid, err := exec.Command("ipfs", "add", "-n", tmpfileName).Output()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get CID",
			"error":   true,
		})
		return
	}
	reg := regexp.MustCompile(`added\s+([[:alnum:]]+)\s+`)
	cidStr := string(reg.FindSubmatch(cid)[1])

	err = os.Rename(tmpfileName, "/tmp/"+cidStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to rename file",
			"error":   true,
		})
		return
	}

	c.JSON(200, gin.H{
		"CID": cidStr,
	})
	return
}

func extractCID(keyStr string) string {
	reg := regexp.MustCompile(`\{\{([[:ascii:]]+)\}\}`)
	cidStr := string(reg.FindSubmatch([]byte(keyStr))[1])
	cidAndMime := strings.Split(cidStr, "|")
	return cidAndMime[0]
}

func publishMediaIPFS(server *electrum.Server, c *gin.Context, paymentAddr string, minPayment float64) {
	var pinMedia PinMedia
	c.BindJSON(&pinMedia)

	tx, err := server.GetTransaction(pinMedia.Tx)
	if err != nil {
		var count = 0
		// Try mutiple times - the transaction may not be in the mempool yet.
		for (tx == nil) && (count < 6) {
			time.Sleep(1 * time.Second)
			tx, err = server.GetTransaction(pinMedia.Tx)
			count++
		}
	}

	if tx == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Transaction",
			"error":   true,
		})
		return
	}

	paymentOK := false
	for _, v := range tx.Vout {
		if v.ScriptPubkey.Addresses[0] == paymentAddr {
			if v.Value < minPayment {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "Insuffcient payment",
					"error":   true,
				})
				return
			}
			paymentOK = true
			break
		}
	}

	if !paymentOK {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "No payment",
			"error":   true,
		})
		return
	}

	for _, v := range tx.Vout {
		if strings.HasPrefix(v.ScriptPubkey.Asm, "OP_KEVA_PUT") {
			items := strings.Split(v.ScriptPubkey.Asm, " ")
			key := items[2]
			dst := make([]byte, hex.DecodedLen(len(key)))
			_, err := hex.Decode(dst, []byte(key))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "No payment",
					"error":   true,
				})
				return
			}
			fmt.Println(string(dst))
			cid := extractCID(string(dst))
			cidResult, err := exec.Command("ipfs", "add", "/tmp/"+cid).Output()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Failed to get CID",
					"error":   true,
				})
				return
			}
			reg := regexp.MustCompile(`added\s+([[:alnum:]]+)\s+`)
			cidStr := string(reg.FindSubmatch(cidResult)[1])
			if cidStr != cid {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Incorrect CID",
					"error":   true,
				})
				return
			}
		}
	}

	c.JSON(200, gin.H{
		"message": "OK",
	})
	return
}

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
)

func handleIPFSFileUpload(c *gin.Context, sh *shell.Shell) {

	f, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	defer f.Close()
	cid, err := sh.Add(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}
	c.JSON(200, gin.H{
		"CID": cid,
	})
}

func setupIPFS() {
	// IPFS
	/*
		sh := shell.NewShell("localhost:5001")
		router.POST("/upload_ipfs", func(c *gin.Context) {
			handleIPFSFileUpload(c, sh)
		})
	*/
	// curl "https://ipfs.infura.io:5001/api/v0/add" -X POST -H "Content-Type: multipart/form-data" -F file=@"the_return.mp4"
}

func addFileToIPFS(f io.Reader) (string, error) {
	tmpFile, err := ioutil.TempFile("", "_ipfs_tmp_")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	w := bufio.NewWriter(tmpFile)
	io.Copy(w, f)
	w.Flush()

	cid, err := exec.Command("ipfs", "add", tmpFile.Name()).Output()
	reg := regexp.MustCompile(`added\s+([[:alnum:]]+)\s+`)
	cidStr := string(reg.FindSubmatch(cid)[1])

	if err != nil {
		return "", err
	}
	return cidStr, nil
}

// InfuraResponse response from Infura
type InfuraResponse struct {
	Name string `json:"Name"`
	Hash string `json:"Hash"`
	Size string `json:"Size"`
}

func addFileToInfura(f io.Reader) (string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormField("file")
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(fw, f); err != nil {
		return "", err
	}
	w.Close()

	url := "https://ipfs.infura.io:5001/api/v0/add"
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return "", err
	}

	// Submit the request
	req.Header.Set("Content-Type", w.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// Check the response
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Infura bad status: %s", resp.Status)
		return "", err
	}

	var infuraResp InfuraResponse
	json.NewDecoder(resp.Body).Decode(&infuraResp)
	defer resp.Body.Close()
	return infuraResp.Hash, nil
}

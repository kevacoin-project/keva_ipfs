package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"testing"
)

func createBody(t *testing.T, fileName string, title, mediaType string) (body *bytes.Buffer, contentType string, err error) {
	// Add a media file to the event.
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	fi, err := file.Stat()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	file.Close()

	body = new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	t.Log(fi.Name())
	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	part.Write(fileContents)

	// Additional parameters.
	writer.WriteField("title", title)
	writer.WriteField("type", mediaType)
	contentType = writer.FormDataContentType()

	err = writer.Close()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	return body, contentType, nil
}

func TestIPFS(t *testing.T) {
	ts := httptest.NewServer(setupServer())
	defer ts.Close()

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	client := &http.Client{
		Jar: jar,
	}

	// Get payment info
	url := fmt.Sprintf("%s/v1/payment_info", ts.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	respPayment, err := client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	var paymentInfo PaymentInfo
	json.NewDecoder(respPayment.Body).Decode(&paymentInfo)

	defer respPayment.Body.Close()
	fmt.Printf("Payment address: %s, minimal payment: %s \n", paymentInfo.PaymentAddress, paymentInfo.MinPayment)

	// Upload media file.
	body, formDataContentType, err := createBody(t, "./testdata/keva_logo.png", "Kevacoin logo", "logo")
	url = fmt.Sprintf("%s/v1/media", ts.URL)
	req, err = http.NewRequest("POST", url, body)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Set("Content-Type", formDataContentType)
	respMedia, err := client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if respMedia.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", respMedia.StatusCode)
	}
	var mediaResp MediaResponse
	err = json.NewDecoder(respMedia.Body).Decode(&mediaResp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Pin media
	url = fmt.Sprintf("%s/v1/pin", ts.URL)

	pinMedia := PinMedia{
		Tx: "1c9634d94bd832cb7ba402c4e4e96fa5ada7199f713627e7c804ce52a4a65b8d",
	}
	pinBytes, err := json.Marshal(&pinMedia)
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(pinBytes))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

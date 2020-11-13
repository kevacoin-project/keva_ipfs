package main

// MediaResponse the reponse from uploadMedia
type MediaResponse struct {
	CID string `json:"CID"`
}

// PinMedia the data body for pinning media to IPFS
type PinMedia struct {
	Tx string `json:"tx"`
}

// PaymentInfo payment information for the server.
type PaymentInfo struct {
	PaymentAddress string `json:"payment_address"`
	MinPayment     string `json:"min_payment"`
}

package main

import (
	"strconv"

	"github.com/kevacoin-project/go-bitcoind"
)

func kevaCreateNamespace(bc *bitcoind.Bitcoind, displayName string) (namespace string, txID string, err error) {
	result, err := bc.KevaCreateNamespace(displayName)
	if err != nil {
		return "", "", err
	}
	return result.NamespaceID, result.TxID, nil
}

func kevaPutValue(bc *bitcoind.Bitcoind, namespaceID string, key string, value string) (string, error) {
	result, err := bc.KevaPut(namespaceID, key, value)
	if err != nil {
		return "", err
	}
	return result.TxID, nil
}

func kevaGetValue(bc *bitcoind.Bitcoind, namespaceID string, key string) (value string, height int, err error) {
	result, err := bc.KevaGet(namespaceID, key)
	if err != nil {
		return "", -1, err
	}
	return result.TxID, result.Height, nil
}

func kevaGetShortcode(bc *bitcoind.Bitcoind, namespaceID string, txID string) (shortcode int, err error) {
	transcation, err := bc.GetTransaction(txID)
	if err != nil {
		return 0, err
	}

	if len(transcation.BlockHash) == 0 {
		return -1, nil
	}

	block, err := bc.GetBlock(transcation.BlockHash)
	if err != nil {
		return -1, err
	}

	height := int(block.Height)
	if height <= 0 {
		return -1, nil
	}

	heightStr := strconv.Itoa(height)
	shortCodeStr := strconv.Itoa(len(heightStr)) + heightStr + strconv.Itoa(int(transcation.BlockIndex))
	shortCode, _ := strconv.Atoi(shortCodeStr)
	return shortCode, nil
}

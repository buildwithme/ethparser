package rpcfetch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	TRANSACTION_BLOCK = "blockNumber"
	TRANSACTION_HASH  = "hash"
	TRANSACTION_FROM  = "from"
	TRANSACTION_TO    = "to"
	TRANSACTION_VALUE = "value"
)

type (
	// BlockByNumberResultResponse captures the response from the Ethereum node.
	BlockByNumberResultResponse struct {
		Number       string           `json:"number"`
		Transactions []map[string]any `json:"transactions"`
	}

	// BlockByNumberResponse captures the response from the Ethereum node.
	BlockByNumberResponse struct {
		JSONRPC string                      `json:"jsonrpc"`
		ID      int                         `json:"id"`
		Result  BlockByNumberResultResponse `json:"result"`
	}
)

// fetchBlock fetches a block from the Ethereum node and returns transactions.
func (p *ethFetcher) FetchBlock(ctx context.Context, blockNum int) (*BlockResult, error) {
	hexBlock := fmt.Sprintf("0x%X", blockNum)
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{hexBlock, true},
		"id":      1,
	}

	resp, err := p.postRequest(ctx, payload)
	if err != nil {
		return nil, err
	}

	var response BlockByNumberResponse
	err = p.decode(resp, &response)
	if err != nil {
		return nil, err
	}

	var txs []*BlockTransaction
	for _, raw := range response.Result.Transactions {
		txs = append(txs, &BlockTransaction{
			BlockNumber: int(getTrxIntValue(raw, TRANSACTION_BLOCK)),
			Hash:        getTrxStringValue(raw, TRANSACTION_HASH),
			From:        getTrxStringValue(raw, TRANSACTION_FROM),
			To:          getTrxStringValue(raw, TRANSACTION_TO),
			Value:       getTrxStringValue(raw, TRANSACTION_VALUE),
		})
	}

	return &BlockResult{BlockNumber: blockNum, Transactions: txs}, nil
}

// getTrxStringValue returns the string value from the raw transaction data.
func getTrxStringValue(raw map[string]any, value string) string {
	data, ok := raw[value].(string)
	if !ok {
		return ""
	}

	return data
}

// getTrxIntValue returns the integer value from the raw transaction data.
func getTrxIntValue(raw map[string]any, value string) int64 {
	bn, _ := strconv.ParseInt(getTrxStringValue(raw, value), 0, 64)
	return bn
}

// fetchLatestBlock gets the chain tip from Ethereum node.
func (p *ethFetcher) postRequest(ctx context.Context, payload any) (*http.Response, error) {
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetLatestBlock returns the latest block number from the Ethereum node.
func (p *ethFetcher) decode(resp *http.Response, output any) error {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		return err
	}

	return nil
}

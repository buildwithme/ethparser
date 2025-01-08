package rpcfetch

import (
	"context"
	"strconv"
)

const BLOCK_NUMBER_METHOD = "eth_blockNumber"

type (
	RPCRequest struct {
		JSONRPC string `json:"jsonrpc"`
		Method  string `json:"method"`
		Params  []interface{}
		ID      int `json:"id"`
	}

	BlockNumberResponse struct {
		Result string `json:"result"`
	}
)

// GetLatestBlock returns the latest block number from the Ethereum node.
func (p *ethFetcher) GetLatestBlock(ctx context.Context) (int, error) {
	payload := RPCRequest{
		JSONRPC: JSON_RPC_VERSION,
		Method:  BLOCK_NUMBER_METHOD,
		Params:  []interface{}{},
		ID:      1,
	}

	resp, err := p.postRequest(ctx, payload)
	if err != nil {
		return 0, err
	}

	var response BlockNumberResponse
	err = p.decode(resp, &response)
	if err != nil {
		return 0, err
	}

	val, err := strconv.ParseInt(response.Result, 0, 64)
	if err != nil {
		return 0, err
	}

	return int(val), nil
}

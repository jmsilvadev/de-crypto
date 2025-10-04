package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jmsilvadev/de-crypto/pkg/utils"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Ethereum struct {
	cliUrl     string
	httpClient HTTPClient
}

var _ JsonRpcClient = &Ethereum{}

func NewEthereum(cliUrl string, httpClient HTTPClient) *Ethereum {
	return &Ethereum{
		cliUrl:     cliUrl,
		httpClient: httpClient,
	}
}

func (e *Ethereum) GetCurrentBlockNumber(ctx context.Context) (uint64, error) {
	payload := `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.cliUrl, strings.NewReader(payload))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result rpcResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	if result.Error != nil {
		return 0, fmt.Errorf("rpc error %d: %s", result.Error.Code, result.Error.Message)
	}

	var block string
	if err := json.Unmarshal(result.Result, &block); err != nil {
		return 0, fmt.Errorf("unexpected result type for eth_blockNumber: %w", err)
	}

	return utils.ParseHexUint64(block)
}

func (e *Ethereum) GetBlockByNumber(ctx context.Context, blockNumber uint64) (*Block, error) {
	payload := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x%x", true],"id":1}`, blockNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.cliUrl, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result rpcResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Error != nil {
		return nil, fmt.Errorf("rpc error %d: %s", result.Error.Code, result.Error.Message)
	}

	var block Block
	if err := json.Unmarshal(result.Result, &block); err != nil {
		return nil, fmt.Errorf("unexpected result type for eth_getBlockByNumber: %w", err)
	}

	return &block, nil
}

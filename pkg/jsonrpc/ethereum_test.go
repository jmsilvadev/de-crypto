package jsonrpc

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockHTTPClient struct {
	response *http.Response
	err      error
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func TestNewEthereum(t *testing.T) {
	t.Run("create_ethereum_client", func(t *testing.T) {
		client := &http.Client{Timeout: 5 * time.Second}
		ethereum := NewEthereum("https://test-rpc.com", client)

		assert.NotNil(t, ethereum)
		assert.Equal(t, "https://test-rpc.com", ethereum.cliUrl)
		assert.Equal(t, client, ethereum.httpClient)
	})
}

func TestEthereum_GetCurrentBlockNumber(t *testing.T) {
	t.Run("get_current_block_number_success", func(t *testing.T) {
		responseBody := `{"jsonrpc":"2.0","result":"0x1a2b","id":1}`
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(responseBody)),
			Header:     make(http.Header),
		}

		mockClient := &mockHTTPClient{response: mockResponse}
		ethereum := &Ethereum{
			cliUrl:     "https://test-rpc.com",
			httpClient: mockClient,
		}

		ctx := context.Background()
		blockNumber, err := ethereum.GetCurrentBlockNumber(ctx)

		assert.NoError(t, err)
		assert.Equal(t, uint64(6699), blockNumber)
	})

	t.Run("get_current_block_number_rpc_error", func(t *testing.T) {
		responseBody := `{"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"},"id":1}`
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(responseBody)),
			Header:     make(http.Header),
		}

		mockClient := &mockHTTPClient{response: mockResponse}
		ethereum := &Ethereum{
			cliUrl:     "https://test-rpc.com",
			httpClient: mockClient,
		}

		ctx := context.Background()
		blockNumber, err := ethereum.GetCurrentBlockNumber(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rpc error -32601")
		assert.Equal(t, uint64(0), blockNumber)
	})

	t.Run("get_current_block_number_http_error", func(t *testing.T) {
		mockClient := &mockHTTPClient{err: assert.AnError}
		ethereum := &Ethereum{
			cliUrl:     "https://test-rpc.com",
			httpClient: mockClient,
		}

		ctx := context.Background()
		blockNumber, err := ethereum.GetCurrentBlockNumber(ctx)

		assert.Error(t, err)
		assert.Equal(t, uint64(0), blockNumber)
	})

	t.Run("get_current_block_number_invalid_json", func(t *testing.T) {
		responseBody := `invalid json`
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(responseBody)),
			Header:     make(http.Header),
		}

		mockClient := &mockHTTPClient{response: mockResponse}
		ethereum := &Ethereum{
			cliUrl:     "https://test-rpc.com",
			httpClient: mockClient,
		}

		ctx := context.Background()
		blockNumber, err := ethereum.GetCurrentBlockNumber(ctx)

		assert.Error(t, err)
		assert.Equal(t, uint64(0), blockNumber)
	})
}

func TestEthereum_GetBlockByNumber(t *testing.T) {
	t.Run("get_block_by_number_success", func(t *testing.T) {
		responseBody := `{"jsonrpc":"2.0","result":{"number":"0x1a2b","hash":"0xabc","transactions":[]},"id":1}`
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(responseBody)),
			Header:     make(http.Header),
		}

		mockClient := &mockHTTPClient{response: mockResponse}
		ethereum := &Ethereum{
			cliUrl:     "https://test-rpc.com",
			httpClient: mockClient,
		}

		ctx := context.Background()
		block, err := ethereum.GetBlockByNumber(ctx, 6699)

		assert.NoError(t, err)
		assert.NotNil(t, block)
		assert.Equal(t, "0x1a2b", block.Number)
		assert.Equal(t, "0xabc", block.Hash)
	})

	t.Run("get_block_by_number_rpc_error", func(t *testing.T) {
		responseBody := `{"jsonrpc":"2.0","error":{"code":-32602,"message":"Invalid params"},"id":1}`
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(responseBody)),
			Header:     make(http.Header),
		}

		mockClient := &mockHTTPClient{response: mockResponse}
		ethereum := &Ethereum{
			cliUrl:     "https://test-rpc.com",
			httpClient: mockClient,
		}

		ctx := context.Background()
		block, err := ethereum.GetBlockByNumber(ctx, 6699)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rpc error -32602")
		assert.Nil(t, block)
	})

	t.Run("get_block_by_number_http_error", func(t *testing.T) {
		mockClient := &mockHTTPClient{err: assert.AnError}
		ethereum := &Ethereum{
			cliUrl:     "https://test-rpc.com",
			httpClient: mockClient,
		}

		ctx := context.Background()
		block, err := ethereum.GetBlockByNumber(ctx, 6699)

		assert.Error(t, err)
		assert.Nil(t, block)
	})

	t.Run("get_block_by_number_with_transactions", func(t *testing.T) {
		responseBody := `{"jsonrpc":"2.0","result":{"number":"0x1a2b","hash":"0xabc","transactions":[{"hash":"0xtx1","from":"0xfrom","to":"0xto","value":"0x100"}]},"id":1}`
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(responseBody)),
			Header:     make(http.Header),
		}

		mockClient := &mockHTTPClient{response: mockResponse}
		ethereum := &Ethereum{
			cliUrl:     "https://test-rpc.com",
			httpClient: mockClient,
		}

		ctx := context.Background()
		block, err := ethereum.GetBlockByNumber(ctx, 6699)

		assert.NoError(t, err)
		assert.NotNil(t, block)
		assert.Len(t, block.Transactions, 1)
		assert.Equal(t, "0xtx1", block.Transactions[0].Hash)
	})
}

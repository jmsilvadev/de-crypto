package jsonrpc

import (
	"context"
)

type JsonRpcClient interface {
	GetCurrentBlockNumber(context.Context) (uint64, error)
	GetBlockByNumber(context.Context, uint64) (*Block, error)
}

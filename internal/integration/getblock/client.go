package getblock

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/restlesswhy/eth-balance-searcher/config"
	"github.com/restlesswhy/eth-balance-searcher/internal/models"
	"github.com/ybbus/jsonrpc/v3"
)

type getBlockRPC struct {
	client jsonrpc.RPCClient
}

func New(cfg *config.Config) *getBlockRPC {
	rpcClient := jsonrpc.NewClientWithOpts("https://eth.getblock.io/mainnet/", &jsonrpc.RPCClientOpts{
		HTTPClient: http.DefaultClient,
		CustomHeaders: map[string]string{
			"x-api-key": cfg.APIToken,
		},
	})

	return &getBlockRPC{
		client: rpcClient,
	}
}

func (g *getBlockRPC) GetLastBlockNumber(ctx context.Context) (string, error) {
	numb := ""

	err := g.client.CallFor(ctx, &numb, "eth_blockNumber")
	if err != nil {
		return "", errors.Wrap(err, "call rpc service error")
	}

	return numb, nil
}

func (g *getBlockRPC) GetBlockByNumber(ctx context.Context, number string) (*models.Block, error) {
	block := &models.Block{}

	err := g.client.CallFor(ctx, block, "eth_getBlockByNumber", number, true) // send with Authorization-Header
	if err != nil {
		return nil, errors.Wrap(err, "get block by number error")
	}

	return block, nil
}

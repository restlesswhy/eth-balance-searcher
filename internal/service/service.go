package service

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/pkg/errors"
	"github.com/restlesswhy/eth-balance-searcher/internal/models"
	"github.com/restlesswhy/eth-balance-searcher/internal/utils"
	"github.com/restlesswhy/eth-balance-searcher/pkg/logger"
)

type GetBlockRPC interface {
	GetLastBlockNumber(ctx context.Context) (string, error)
	GetBlockByNumber(ctx context.Context, number string) (*models.Block, error)
}

type service struct {
	log    logger.Logger
	client GetBlockRPC
}

func New(log logger.Logger, client GetBlockRPC) *service {
	return &service{log: log, client: client}
}

func (s *service) GetAddress(ctx context.Context) (string, error) {
	lastBlock, err := s.client.GetLastBlockNumber(ctx)
	if err != nil {
		return "", errors.Wrap(err, "get last block number error")
	}

	num, err := utils.HexToUInt64(lastBlock)
	if err != nil {
		return "", errors.Wrap(err, "convert hex to uint64 error")
	}

	seq, err := utils.GetSequence(num-99, num)
	if err != nil {
		return "", errors.Wrap(err, "get sequence error")
	}

	type asd struct {
		wg sync.WaitGroup
		mu sync.Mutex
		ma map[string]*big.Int
	}

	aas := asd{
		ma: make(map[string]*big.Int),
		wg: sync.WaitGroup{},
	}

	for _, i := range seq {
		aas.wg.Add(1)
		go func(d uint64) {
			defer aas.wg.Done()
			block, err := s.client.GetBlockByNumber(ctx, utils.UInt64ToHex(d))
			if err != nil {
				s.log.Error(err)
			}

			for _, v := range block.Transactions {
				val, err := utils.BigIntFromHex(v.Value)
				if err != nil {
					log.Fatal(err)
					continue
				}

				aas.mu.Lock()
				if from, ok := aas.ma[v.From]; !ok {
					res := new(big.Int)
					res.Sub(res, val)
					aas.ma[v.From] = res
				} else {
					res := new(big.Int)
					res.Sub(from, val)
					aas.ma[v.From] = res
				}

				if to, ok := aas.ma[v.To]; !ok {
					aas.ma[v.To] = val
				} else {
					res := new(big.Int)
					res.Add(to, val)
					aas.ma[v.From] = res
				}
				aas.mu.Unlock()
			}
		}(i)
	}
	aas.wg.Wait()

	as := new(big.Int)
	ag := ""
	for k, v := range aas.ma {
		if v.CmpAbs(as) == 1 {
			as = v
			ag = k
		}
	}

	fmt.Println(as.String())
	return ag, nil
}

func (s *service) GetAddress1(ctx context.Context) (string, error) {
	lastBlock, err := s.client.GetLastBlockNumber(ctx)
	if err != nil {
		return "", errors.Wrap(err, "get last block number error")
	}

	num, err := utils.HexToUInt64(lastBlock)
	if err != nil {
		return "", errors.Wrap(err, "convert hex to uint64 error")
	}

	seq, err := utils.GetSequence(num-99, num)
	if err != nil {
		return "", errors.Wrap(err, "get sequence error")
	}

	blockGetter := func(ctx context.Context, integers ...uint64) <-chan *models.Block {
		blockStream := make(chan *models.Block, 100)
		go func() {
			defer close(blockStream)
			wg := sync.WaitGroup{}
			for _, i := range integers {
				wg.Add(1)
				go func(d uint64) {
					defer wg.Done()
					block, err := s.client.GetBlockByNumber(ctx, utils.UInt64ToHex(d))
					if err != nil {
						s.log.Error(err)
					}
					blockStream <- block
				}(i)
			}
			wg.Wait()
		}()
		return blockStream
	}

	calculation := func(
		ctx context.Context,
		blockStream <-chan *models.Block,
	) <-chan map[string]*big.Int {
		res := make(chan map[string]*big.Int)
		go func() {
			defer close(res)
			data := make(map[string]*big.Int)

			for i := range blockStream {
				select {
				case <-ctx.Done():
					return
				default:
					for _, v := range i.Transactions {
						val, err := utils.BigIntFromHex(v.Value)
						if err != nil {
							log.Fatal(err)
							continue
						}

						if from, ok := data[v.From]; !ok {
							res := new(big.Int)
							res.Sub(res, val)
							data[v.From] = res
						} else {
							res := new(big.Int)
							res.Sub(from, val)
							data[v.From] = res
						}

						if to, ok := data[v.To]; !ok {
							data[v.To] = val
						} else {
							res := new(big.Int)
							res.Add(to, val)
							data[v.From] = res
						}
					}
				}
			}
			res <- data
		}()
		return res
	}

	res := calculation(ctx, blockGetter(ctx, seq...))

	as := new(big.Int)
	ag := ""
	for k, v := range <-res {
		if v.CmpAbs(as) == 1 {
			as = v
			ag = k
		}
	}

	fmt.Println(as.String())
	return ag, nil
}

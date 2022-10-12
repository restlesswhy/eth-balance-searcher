package service

import (
	"context"
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

type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) (bool, error)
	Set(ctx context.Context, key string, value interface{}) error
}

type service struct {
	log    logger.Logger
	client GetBlockRPC
	cache  Cache
}

func New(log logger.Logger, client GetBlockRPC, cache Cache) *service {
	return &service{log: log, client: client, cache: cache}
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

	data := models.NewSearchData()

	wg := sync.WaitGroup{}
	wg.Add(len(seq))
	for _, i := range seq {
		go func(d uint64) {
			defer wg.Done()
			hexNum := utils.UInt64ToHex(d)

			block := &models.Block{}
			exist, err := s.cache.Get(ctx, hexNum, block)
			if err != nil {
				s.log.Warn(err)
			}

			if !exist {
				block, err = s.client.GetBlockByNumber(ctx, hexNum)
				if err != nil {
					s.log.Error(err)
					return
				}
				
				if err := s.cache.Set(ctx, hexNum, block); err != nil {
					s.log.Warn(err)
				}
			}

			for _, v := range block.Transactions {
				val, err := utils.BigIntFromHex(v.Value)
				if err != nil {
					s.log.Error(err)
					continue
				}

				if from, ok := data.Get(v.From); !ok {
					res := new(big.Int)
					res.Sub(res, val)
					data.Set(v.From, res)
				} else {
					res := new(big.Int)
					res.Sub(from, val)
					data.Set(v.From, res)
				}

				if to, ok := data.Get(v.To); !ok {
					data.Set(v.To, val)
				} else {
					res := new(big.Int)
					res.Add(to, val)
					data.Set(v.To, res)
				}
			}
		}(i)
	}
	wg.Wait()

	largerValue := new(big.Int)
	address := ""
	for k, v := range data.GetAll() {
		if v.CmpAbs(largerValue) == 1 {
			largerValue = v
			address = k
		}
	}

	s.log.Infof("largest value - %s", largerValue.String())
	return address, nil
}

// Вариант через пайплайны
// func (s *service) GetAddress(ctx context.Context) (string, error) {
// 	lastBlock, err := s.client.GetLastBlockNumber(ctx)
// 	if err != nil {
// 		return "", errors.Wrap(err, "get last block number error")
// 	}

// 	num, err := utils.HexToUInt64(lastBlock)
// 	if err != nil {
// 		return "", errors.Wrap(err, "convert hex to uint64 error")
// 	}

// 	seq, err := utils.GetSequence(num-99, num)
// 	if err != nil {
// 		return "", errors.Wrap(err, "get sequence error")
// 	}

// 	blockGetter := func(ctx context.Context, integers ...uint64) <-chan *models.Block {
// 		blockStream := make(chan *models.Block, 100)
// 		go func() {
// 			defer close(blockStream)
// 			wg := sync.WaitGroup{}
// 			for _, i := range integers {
// 				wg.Add(1)
// 				go func(d uint64) {
// 					defer wg.Done()
// 					block, err := s.client.GetBlockByNumber(ctx, utils.UInt64ToHex(d))
// 					if err != nil {
// 						s.log.Error(err)
// 					}
// 					blockStream <- block
// 				}(i)
// 			}
// 			wg.Wait()
// 		}()
// 		return blockStream
// 	}

// 	calculation := func(
// 		ctx context.Context,
// 		blockStream <-chan *models.Block,
// 	) <-chan map[string]*big.Int {
// 		res := make(chan map[string]*big.Int)
// 		go func() {
// 			defer close(res)
// 			data := make(map[string]*big.Int)

// 			for i := range blockStream {
// 				select {
// 				case <-ctx.Done():
// 					return
// 				default:
// 					for _, v := range i.Transactions {
// 						val, err := utils.BigIntFromHex(v.Value)
// 						if err != nil {
// 							log.Fatal(err)
// 							continue
// 						}

// 						if from, ok := data[v.From]; !ok {
// 							res := new(big.Int)
// 							res.Sub(res, val)
// 							data[v.From] = res
// 						} else {
// 							res := new(big.Int)
// 							res.Sub(from, val)
// 							data[v.From] = res
// 						}

// 						if to, ok := data[v.To]; !ok {
// 							data[v.To] = val
// 						} else {
// 							res := new(big.Int)
// 							res.Add(to, val)
// 							data[v.From] = res
// 						}
// 					}
// 				}
// 			}
// 			res <- data
// 		}()
// 		return res
// 	}

// 	res := calculation(ctx, blockGetter(ctx, seq...))

// 	as := new(big.Int)
// 	ag := ""
// 	for k, v := range <-res {
// 		if v.CmpAbs(as) == 1 {
// 			as = v
// 			ag = k
// 		}
// 	}

// 	fmt.Println(as.String())
// 	return ag, nil
// }

package utils

import (
	"math/big"

	"github.com/pkg/errors"
	"github.com/ubiq/go-ubiq/common/hexutil"
)

func BigIntFromHex(s string) (*big.Int, error) {
	res, err := hexutil.DecodeBig(s)
	if err != nil {
		return nil, errors.Wrap(err, "decode hex to big int error")
	}
	
	return res, nil
}

func UInt64ToHex(n uint64) string {
	return hexutil.EncodeUint64(n)
}

func HexToUInt64(hex string) (uint64, error) {
	res, err := hexutil.DecodeUint64(hex)
	if err != nil {
		return 0, errors.New("convert hex to decimal failed")
	}

	return res, nil
}

func GetSequence(min, max uint64) ([]uint64, error) {
	if min > max {
		return nil, errors.New("first num must be larger then second num")
	}

	a := make([]uint64, max-min+1)
	for i := range a {
		a[i] = min + uint64(i)
	}

	return a, nil
}

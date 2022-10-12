package utils_test

import (
	"testing"

	"github.com/restlesswhy/eth-balance-searcher/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestBigIntFromHex(t *testing.T) {
	res, err := utils.BigIntFromHex("0x2386f26fc10000")
	assert.NoError(t, err)
	assert.Equal(t, "10000000000000000", res.String())
}

func TestUInt64ToHex(t *testing.T) {
	res := utils.UInt64ToHex(15725513)
	assert.Equal(t, "0xeff3c9", res)
}

func TestHexToUInt64(t *testing.T) {
	res, err := utils.HexToUInt64("0xeff3c9")
	assert.NoError(t, err)
	assert.Equal(t, uint64(15725513), res)
}

func TestGetSequence(t *testing.T) {
	res, err := utils.GetSequence(1001, 1010)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{1001, 1002, 1003, 1004, 1005, 1006, 1007, 1008, 1009, 1010}, res)
}

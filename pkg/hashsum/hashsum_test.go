package hashsum_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/chiyutianyi/git-hashsum/pkg/hashsum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSum(t *testing.T) {
	var (
		sum [32]byte
		old = "d294bfa287a8df15375a37f71aea40ad8f6b1a66915efa1093878a5812789cd2"
	)
	oldSum, err := hex.DecodeString(old)
	require.NoError(t, err)
	for i := 0; i < 32; i++ {
		sum[i] = oldSum[i]
	}

	// update ref 0000000000000000000000000000000000000000 f3aaf7021071e311d06454317951a943f0233ccd refs/heads/test
	sum = hashsum.Sum(sum, "0000000000000000000000000000000000000000", "refs/heads/test")
	sum = hashsum.Sum(sum, "f3aaf7021071e311d06454317951a943f0233ccd", "refs/heads/test")
	assert.Equal(t, "328749b3b9adf2c7ac6c1d74dcd4d6cabe8eb34fd7b3889aa83d99e443ec1520", fmt.Sprintf("%x", sum))

	// update ref f3aaf7021071e311d06454317951a943f0233ccd 0000000000000000000000000000000000000000 refs/heads/test
	sum = hashsum.Sum(sum, "f3aaf7021071e311d06454317951a943f0233ccd", "refs/heads/test")
	sum = hashsum.Sum(sum, "0000000000000000000000000000000000000000", "refs/heads/test")
	assert.Equal(t, old, fmt.Sprintf("%x", sum))
}

package hashsum

import (
	"crypto/sha256"
	"fmt"
)

func Sum(old [32]byte, oid, ref string) [32]byte {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s %s", ref, oid)))

	return Xor(old, sum)
}

func Xor(a, b [32]byte) (r [32]byte) {
	for i := range a {
		r[i] = a[i] ^ b[i]
	}
	return r
}

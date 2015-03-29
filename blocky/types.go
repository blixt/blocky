package blocky

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type Id string

func NewId() Id {
	v, _ := rand.Int(rand.Reader, big.NewInt(0).Exp(big.NewInt(2), big.NewInt(80), nil))
	return Id(fmt.Sprintf("%020x", v))
}

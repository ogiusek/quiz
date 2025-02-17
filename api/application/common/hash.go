package common

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

type Hasher interface {
	Hash([]byte) string
}

type hasher struct {
	hasher func() hash.Hash
}

func (config *hasher) Hash(encoded []byte) string {
	hasher := config.hasher()
	hasher.Write(encoded)
	return hex.EncodeToString(hasher.Sum(nil))
}

func NewHasher() Hasher {
	return &hasher{
		hasher: sha256.New,
	}
}

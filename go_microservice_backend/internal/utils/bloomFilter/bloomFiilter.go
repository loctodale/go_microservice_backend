package bloomFilter

import (
	"context"
	"crypto/sha256"
	"fmt"
	"go_microservice_backend_api/global"
	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/sha3"
	"hash"
	"math/big"
)

func getHashFunction() []hash.Hash {
	return []hash.Hash{
		sha256.New(),
		sha3.New224(),
		sha3.New384(),
		sha3.New256(),
		md4.New(),
	}
}

func computeHashes(value string, numBits int) []int {
	hashes := getHashFunction()
	indices := make([]int, numBits)
	for i, h := range hashes {
		h.Reset()
		h.Write([]byte(value))
		hashValue := big.NewInt(0).SetBytes(h.Sum(nil))
		indices[i] = int(hashValue.Mod(hashValue, big.NewInt(int64(numBits))).Int64())
	}
	return indices
}

func AddToBloomFilter(key string, value string) bool {
	indices := computeHashes(value, global.Config.Redis.Numbits)
	var err error
	for _, index := range indices {
		err = global.Rdb.SetBit(context.Background(), key, int64(index), 1).Err()
	}
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func CheckBloomFilter(key string, value string) bool {
	indices := computeHashes(value, global.Config.Redis.Numbits)
	for _, index := range indices {
		bit, err := global.Rdb.GetBit(context.Background(), key, int64(index)).Result()
		if err != nil || bit == 0 {
			return false
		}
	}
	return true
}

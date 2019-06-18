package domain

import "context"

//provides flexibility to change the hash function used by worker manager
type HashFunction interface {
	Hash(ctx context.Context,key string)(hash int64,err error)
}

type HashFunc int

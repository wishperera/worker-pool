package domain

import "context"

type HashFunction interface {
	Hash(ctx context.Context,key string)(hash int64,err error)
}

type HashFunc int

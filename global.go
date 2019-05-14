package worker_pool

import "github.com/google/uuid"

var(
	PoolMap map[uuid.UUID](*Pool)
)


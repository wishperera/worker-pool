package worker_pool

import "github.com/google/uuid"


//Reduntant for now.Will be useful for a multiple pool scenario
var (
	PoolMap map[uuid.UUID](*Pool)
)

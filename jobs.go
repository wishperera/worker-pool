package worker_pool

import "github.com/google/uuid"

type job struct {
	id uuid.UUID
	input interface{}
	output interface{}
	err    error
}



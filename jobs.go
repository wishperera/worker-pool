package worker_pool

import (
	"context"
	"github.com/google/uuid"
)

type Job struct {
	ctx    context.Context
	key    interface{}   // key must be of type int or string
	id     uuid.UUID
	input  interface{}
	output interface{}
	err    error
}


//Returns a pointer to a Job object with provided input.
func newJob(ctx context.Context, in,key interface{}) (j Job) {
	return Job{
		id:     uuid.New(),
		input:  in,
		key:key,
		output: nil,
		err:    nil,
	}
}

// get the output and the error returned by the process function
func (j Job) GetOutput() (output interface{}, err error) {
	return j.output, j.err
}

// get the input passed to the process function
func (j Job) GetInput() (input interface{}) {
	return j.input
}

// returns a id that can uniquely address the job
func (j Job) GetID() uuid.UUID {
	return j.id
}

// return the key associated with the job
func (j Job)GetKey() interface{}{
	return j.key
}
package worker_pool

import (
	"github.com/go-errors/errors"
	"strconv"
	"context"
	"sync"
	"github.com/google/uuid"
)

type Pool struct {
	id uuid.UUID
	input chan Job
	output chan Job
	workers int  // maximum number of workers
	bufferSize int //  buffer size for all channels
	processFunc func(ctx context.Context,in interface{})(out interface{},err error)
	wGroup sync.WaitGroup
	closeWorkers chan bool
}


/*
	Func NewPol(maxWorkers,buffersize int)(p *Pool,err error)
	---------------------------------------------------------
	returns a pointer to pool object with the input,output and
    error channels size set to the buffersize and routines count
    limited to workers
 */
func NewPool(maxWorkers,buffersize int)(p *Pool,err error){

	if (buffersize < 0){
		return nil,errors.New("buffersize must be non negative,provided:"+strconv.Itoa(buffersize))
	}

	if (maxWorkers < 0){
		return nil,errors.New("maxWorkers must be non negative,provided:"+strconv.Itoa(maxWorkers))
	}

	return &Pool{
		input: make(chan Job,buffersize),
		output: make(chan Job,buffersize),
		closeWorkers: make(chan bool),
		workers: maxWorkers,
		bufferSize: buffersize,
	},nil
}

/*
	Func Init(ctx context.Context,processFunc func(ctx context.Context,in interface{})(out interface{},err error))
	--------------------------------------------------------------------------------------------------------------
	intiailize the pool with a process function that accepts a context and the function parameters as a interface.
    parameter can be a single value or a structure in case of multiple expected inputs, same goes for output.
 */
func (p *Pool)Init(ctx context.Context,processFunc func(ctx context.Context,in interface{})(out interface{},err error)){
	p.processFunc = processFunc
}
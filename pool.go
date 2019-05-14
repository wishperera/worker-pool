package worker_pool

import (
	"github.com/go-errors/errors"
	"strconv"
	"context"
)

type Pool struct {
	input chan interface{}
	output chan interface{}
	error  chan error
	workers int  // maximum number of workers
	bufferSize int //  buffer size for all channels
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
		input: make(chan interface{},buffersize),
		output: make(chan interface{},buffersize),
		error: make(chan error,buffersize),
		workers: maxWorkers,
		bufferSize: buffersize,
	},nil
}


func (p *Pool)Init(ctx context.Context){

}
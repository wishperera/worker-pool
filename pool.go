package worker_pool

import (
	"github.com/go-errors/errors"
	"strconv"
	"context"
	"sync"
	"github.com/google/uuid"
	"fmt"
)

type Pool struct {
	id uuid.UUID
	input chan Job
	Output chan Job
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
		Output: make(chan Job,buffersize),
		closeWorkers: make(chan bool),
		workers: maxWorkers,
		bufferSize: buffersize,
	},nil
}

/*
	Func (p *Pool)Init(ctx context.Context,processFunc func(ctx context.Context,in interface{})(out interface{},err error))
	--------------------------------------------------------------------------------------------------------------
	intiailize the pool with a process function that accepts a context and the function parameters as a interface.
    parameter can be a single value or a structure in case of multiple expected inputs, same goes for output.
 */
func (p *Pool)Init(ctx context.Context,processFunc func(ctx context.Context,in interface{})(out interface{},err error)){
	p.processFunc = processFunc
	p.id = uuid.New()
	for i := 0; i < p.workers; i ++{
		worker := &worker{
			id: uuid.New(),
			pool: p,
		}

		worker.run()
		p.wGroup.Add(1)
	}

	fmt.Println("worker pool successfully initialized with, pool_id: ",p.id, "workers_count: ",p.workers, "buffer_size:",p.bufferSize)

}



/*
	Func (p *Pool)Close(ctx context.Context)
	---------------------------------------
    shut down the pool gracefully after waiting all worker routines to close.
 */
func (p *Pool)Close(ctx context.Context){
	for i := 0; i < p.workers; i++{
		p.closeWorkers <- true
	}

	p.wGroup.Wait()
	close(p.closeWorkers)
	close(p.input)
	close(p.Output)
	fmt.Println("worker pool gracefully shut down, pool_id:",p.id)

}


/*
	FuncAddNewJob(ctx context.Context,input interface{})(jobID uuid.UUID)
	---------------------------------------------------------------------
    adds a new job to the process queue. will panic if the pool is not
    initialized using pool.Init(). Returns the job id for future use.
 */
func (p *Pool)AddNewJob(ctx context.Context,input interface{})(jobID uuid.UUID){
	if p.processFunc == nil{
		panic("process function empty in pool,please initialize using pool.Init()")
	}
	jb := newJob(ctx,input)
	p.input <- jb
	return jb.id
}
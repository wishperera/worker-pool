//worker pool is a simple go library to intergrate a  limited routine process in to your
//golang application.
package worker_pool

import (
	"context"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/wishperera/worker-pool/domain"
	"github.com/wishperera/worker-pool/hasher"
	"log"
	"strconv"
	"sync"
)

type Pool struct {
	id           uuid.UUID
	Output       chan Job
	workers      int64 // maximum number of workers
	workerBuffer int64 // buffer size for worker channels
	processFunc  func(ctx context.Context, in interface{}) (out interface{}, err error)
	wGroup       sync.WaitGroup
	closeWorkers chan bool
	manager      workerManager
}

//returns a pointer to pool object with the input,output and
//error channels size set to the buffersize and routines count
//limited to workers
func NewPool(maxWorkers , workerBufferSize int64, hashFunc domain.HashFunc) (p *Pool, err error) {

	if workerBufferSize < 0 {
		return nil, errors.New("[pool] workerBufferSize must be non negative,provided:" + strconv.FormatInt(workerBufferSize,10))
	}

	if maxWorkers < 0 {
		return nil, errors.New("[pool] workerBufferSize must be non negative,provided:" + strconv.FormatInt(maxWorkers,10))
	}

	hf,err := selectHash(hashFunc)
	if err != nil{
		return nil,err
	}

	return &Pool{
		Output:       make(chan Job),
		closeWorkers: make(chan bool),
		workers:      maxWorkers,
		workerBuffer:   workerBufferSize,
		manager: workerManager{
			buckets: maxWorkers,
			hashfunc: hf,
			workers: make(map[int64]*worker,0),
		},
	}, nil
}


//intiailize the pool with a process function that accepts a context and the function parameters as a interface.
//parameter can be a single value or a structure in case of multiple expected inputs, same goes for output.
func (p *Pool) Init(ctx context.Context, processFunc func(ctx context.Context, in interface{}) (out interface{}, err error)) {
	p.processFunc = processFunc
	p.id = uuid.New()
	var i int64
	for i = 0; i < p.workers; i++ {
		worker := &worker{
			id:   uuid.New(),
			pool: p,
			buffer: make(chan Job,p.workerBuffer),
		}
		p.manager.pool = p
		p.manager.addWorker(i,worker)
		worker.run()
		p.wGroup.Add(1)
	}
	log.Println("[pool] worker pool successfully initialized with, pool_id: ", p.id, "workers_count: ", p.workers, "worker_buffer_size:", p.workerBuffer)

}


//shut down the pool gracefully after waiting all worker routines to close.
func (p *Pool) Close(ctx context.Context) {
	p.manager.stopWorkers()
	p.wGroup.Wait()
	close(p.closeWorkers)
	close(p.Output)
	log.Println("worker pool gracefully shut down, pool_id:", p.id)

}


//adds a new job to the process queue. will panic if the pool is not
//initialized using pool.Init(). Returns the job id for future use.
func (p *Pool) AddNewJob(ctx context.Context, input,key interface{}) (jobID uuid.UUID) {
	if p.processFunc == nil {
		log.Fatal("process function empty in pool,please initialize using pool.Init()")
	}
	jb := newJob(ctx, input,key)
	p.manager.assignJobToWorkers(jb)
	return jb.id
}

//returns the implementation of the hash fucntion given the id
// currently supports SHA256,SHA512,MD5
func selectHash(hashFunc domain.HashFunc)(function domain.HashFunction,err error){
	switch hashFunc {
	case SHA256:
		return hasher.SHA256{},nil
	default:
		return function,errors.New("[pool] unsupported hash function provided")
	}
}
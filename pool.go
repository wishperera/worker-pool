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


type poolConfig struct {
	Workers int64
	WorkerBufferSize int64
	Metrics struct{
		NameSpace string
		SubSystem string
	}
	HashFunc domain.HashFunc
}

// default pool  config
var DefaultConfig poolConfig

func init(){
	DefaultConfig = poolConfig{
		Workers: 100,
		WorkerBufferSize: 10,
		HashFunc:SHA256,
	}
	DefaultConfig.Metrics.NameSpace = "worker_pool"
	DefaultConfig.Metrics.SubSystem = "sub_system"
}

// return an empty pool config
func NewPoolConfig()poolConfig{
	return poolConfig{}
}


type Pool struct {
	id           uuid.UUID
	Output       chan Job
	conf         poolConfig
	processFunc  func(ctx context.Context, in interface{}) (out interface{}, err error)
	wGroup       sync.WaitGroup
	closeWorkers chan bool
	manager      workerManager
}

//returns a pointer to pool object with the input,output and
//error channels size set to the buffersize and routines count
//limited to workers
func NewPool(conf poolConfig) (p *Pool, err error) {

	err = validatePoolConfig(conf)
	if err != nil{
		return nil,err
	}

	hf,err := selectHash(conf.HashFunc)
	if err != nil{
		return nil,err
	}

	return &Pool{
		Output:       make(chan Job),
		closeWorkers: make(chan bool),
		conf: conf,
		manager: workerManager{
			buckets: conf.Workers,
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

	initMetrics(p)

	for i = 0; i < p.conf.Workers; i++ {
		worker := &worker{
			id:   uuid.New(),
			pool: p,
			buffer: make(chan Job,p.conf.WorkerBufferSize),
		}
		p.manager.pool = p
		p.manager.addWorker(i,worker)
		worker.run()
		p.wGroup.Add(1)
	}
	log.Println("[pool] worker pool successfully initialized with, pool_id: ", p.id, "workers_count: ", p.conf.Workers, "worker_buffer_size:", p.conf.WorkerBufferSize)

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
	activeJobs.Add(1)
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


//validate the config information
func validatePoolConfig(cf poolConfig)(err error){
	if cf.WorkerBufferSize < 0 {
		return errors.New("[pool] number of workers must be non negative,provided:" + strconv.FormatInt(cf.Workers,10))
	}

	if cf.Workers < 0 {
		return  errors.New("[pool] workerBufferSize must be non negative,provided:" + strconv.FormatInt(cf.WorkerBufferSize,10))
	}


	if !(len(cf.Metrics.NameSpace) > 0) {
		return errors.New("[pool] metric namespace not provided")
	}

	if !(len(cf.Metrics.SubSystem) > 0) {
		return errors.New("[pool] metric subsystem not provided")
	}


	return nil
}
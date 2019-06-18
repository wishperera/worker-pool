package worker_pool

import (
	"context"
	"github.com/go-errors/errors"
	"github.com/wishperera/worker-pool/domain"
	"reflect"
	"strconv"
)

type workerManager struct {
	hashfunc domain.HashFunction
	workers map[int64](*worker)
	buckets int64
	pool *Pool
}

func (w *workerManager)addWorker(id int64,wkr *worker)(err error){
	if w.workers == nil{
		return errors.New("[manager] pool is not initialized, please invoke pool.Init()")
	}
	if _,ok := w.workers[id]; ok{
		return errors.New("[manager] duplicate register request for worker,already exists"+strconv.FormatInt(id,10))
	}
	w.workers[id] = wkr
	return nil
}

func (w *workerManager)stopWorkers()(err error){
	if w.workers == nil{
		return errors.New("[manager] pool is not initialized, please invoke pool.Init()")
	}

	for id := range w.workers{
		wkr,ok := w.workers[id]
		if ok{
			close(wkr.buffer)
		}
	}
	return nil
}

func (w  *workerManager)assignJobToWorkers(job Job){
	bucketId,err := w.getHashBucketID(job.ctx,job.key)
	if err != nil{
		w.pool.Output <- Job{
			ctx:    job.ctx,
			input:  job.input,
			id:     job.id,
			key:    job.key,
			output: nil,
			err:    err,
		}
	}
	wkr,ok := w.workers[bucketId]
	if !ok{
		panic("[manager] worker not registered for bucketID")
	}

	wkr.buffer <- job
}

func (w *workerManager)	getHashBucketID(ctx context.Context,key interface{})(bucketId int64,err error){
	var str string
	switch reflect.TypeOf(key).Kind() {
	case reflect.Int64:
		str = strconv.FormatInt(key.(int64),10)
	case reflect.Int32:
		str = strconv.FormatInt(int64(key.(int32)),10)
	case reflect.Int16:
		str = strconv.FormatInt(int64(key.(int16)),10)
	case reflect.Int8:
		str = strconv.FormatInt(int64(key.(int8)),10)
	case reflect.Int:
		str = strconv.Itoa(key.(int))
	case reflect.String:
		str = key.(string)
	default:
		return 0, errors.New("[manager] key must be either a string or an integer")

	}

	h,err := w.hashfunc.Hash(ctx,str)
	if err != nil{
		return 0,err
	}

	return h % w.buckets,nil

}


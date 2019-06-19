package tests

import (
	"context"
	"github.com/go-errors/errors"
	"github.com/wishperera/worker-pool"
	"reflect"
	"testing"
	"time"
)

func TestUnitPoolInvalidWorkerSize(t *testing.T){

	conf := worker_pool.NewPoolConfig()
	conf.Workers = -1
	conf.WorkerBufferSize = 10
	conf.Metrics.NameSpace = "test"
	conf.Metrics.SubSystem = "test"
	conf.HashFunc = worker_pool.SHA256

	_,err := worker_pool.NewPool(conf)
	if err != nil{
		t.Log("[test_case]:[invalid_worker_count]: passed")
		return
	}

	t.Error("[test_case]:[invalid_worker_count]: failed")

}

func TestUnitPoolInvalidBufferSize(t *testing.T){

	conf := worker_pool.NewPoolConfig()
	conf.Workers = 100
	conf.WorkerBufferSize = -10
	conf.Metrics.NameSpace = "test"
	conf.Metrics.SubSystem = "test"
	conf.HashFunc = worker_pool.SHA256

	_,err := worker_pool.NewPool(conf)
	if err != nil{
		t.Log("[test_case]:[invalid_buffer_size]: passed")
		return
	}

	t.Error("[test_case]:[invalid_buffer_size]: failed")

}


func TestUnitPool(t *testing.T){

	endTest := make(chan bool)

	testInput := make([]interface{},0)
	testInput = append(testInput, 1)
	testInput = append(testInput, "abc")
	testInput = append(testInput, 0)

	pool,err := worker_pool.NewPool(worker_pool.DefaultConfig)
	if err != nil{
		t.Errorf("[test_case]:[invalid_buffer_size]: failed with error: %v",err)
		return
	}

	//
	pool.Init(context.Background(), func(ctx context.Context, in interface{}) (out interface{}, err error) {
		v,ok := in.(int)
		if !ok{
			if in == nil{
				return nil,errors.New("expected: [int] recieved:[nil]")
			}else{
				return nil,errors.New("expected: [int] recieved:"+reflect.TypeOf(in).String())
			}
		}
		return v,nil
	})

	for _,val := range testInput{
		t.Log("[test] job id:",pool.AddNewJob(context.Background(),val,123123))
	}

	go func() {
		for {
			select {
				case out := <-pool.Output:
					o, e := out.GetOutput()
					t.Log("[test] id:", out.GetID(),"input:", out.GetInput(),"output:", o, "error:",e)

				case <-endTest:
					return
				}
		}
	}()


	time.Sleep(time.Second*8)

	endTest <- true

	pool.Close(context.Background())

}

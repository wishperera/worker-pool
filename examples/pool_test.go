package examples

import (
	"testing"
	"github.com/wishperera/worker-pool"
	"context"
	"github.com/go-errors/errors"
	"reflect"
)

func TestUnitPoolInvalidWorkerSize(t *testing.T){
	_,err := worker_pool.NewPool(-1,10)
	if err != nil{
		t.Log("[test_case]:[invalid_worker_count]: passed")
		return
	}

	t.Error("[test_case]:[invalid_worker_count]: failed")

}

func TestUnitPoolInvalidBufferSize(t *testing.T){
	_,err := worker_pool.NewPool(10,-2)
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

	pool,err := worker_pool.NewPool(10,10)
	if err != nil{
		t.Errorf("[test_case]:[invalid_buffer_size]: failed with error: %v",err)
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


	go func() {
		for {
			select {
				case out := <- pool.Output:
					o,e := out.GetOutput()
					t.Logf("id: %v ,input: %v, output: %v , err: %v",out.GetID(),out.GetInput(),o,e)

				case <- endTest:
					return
			}
		}

	}()
	for _,val := range testInput{
		pool.AddNewJob(context.Background(),val)
	}

	endTest <- true

	pool.Close(context.Background())

	select {

	}

}
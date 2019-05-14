## worker-pool

#### intro:
worker pool is a simple go library to intergrate a  limited routine process in to your 
golang application.

#### usage:

- import the library to development environment using go-get
    
    ```bash
        go-get github.com/wishperera/worker-pool
    ```
- creating a new pool and initializing with a process function

    ```go
    var workerCount = 10;
    var bufferSize  = 10;
    pool,err := worker_pool.NewPool(workerCount,bufferSize)
    if err != nil{
  	    //todo- handle error as preffered
    }
  
    pool.Init(context.Background(), func(ctx context.Context, in interface{}) (out interface{}, err error) {
        //todo - process function body
        return out,nil
    })
  
    ```
- `workerCount` refers to the number of go routines running parellely, and `bufferSize` refers to the size of 
   the input and output channels of the pool.
   
-  adding jobs to the pool 
    
    ```go
     val := 300
     jobID := pool.AddNewJob(context.Background(),val)
    ```
-  value is the input to the process function, which can be a primitive type or a struct if multiple input
   parameters are required  
   
-  output can be retreived from the `pool.Output` channel which returns a `Job` type structure containg the 
   original input,output and the possible errors. The `GetID` method returns a unique id to the Job that 
   can be used to align the input with output.
   

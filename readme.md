## worker-pool

#### intro:
worker pool is a simple go library to intergrate a  limited routine process in to your 
golang application. Jobs added to the worker pool will be assigned to a worker depending 
on a preffered hash algorithm calculated by the key. This ensures that the Jobs with same
key will be processed by the same worker, thus ensuring race conditions in concurrent  processes.

> `Warning:` if shared resources are used by the process function, the library does not gurantee the above claims,
in which case you must implement your own logic to avoid race conditions.  

#### usage:

- import the library to development environment using go-get
    
    ```bash
        go-get github.com/wishperera/worker-pool
    ```
- creating a new pool and initializing with a process function

    ```go
    var workerCount = 10
    var workerBufferSize  = 10
    var hashFunction = SHA256
  
    pool,err := worker_pool.NewPool(workerCount,workerBufferSize,hashFunction)
    if err != nil{
  	    //todo- handle error as preffered
    }
  
    pool.Init(context.Background(), func(ctx context.Context, in interface{}) (out interface{}, err error) {
        //todo - process function body
        return out,nil
    })
  
    ```
- `workerCount` refers to the number of workers(spawned in go routines) running parellely, and `workerBufferSize` refers to the size of 
   the workers buffered channel.
   
-  adding jobs to the pool 
    
    ```go
     val := 300
     jobID := pool.AddNewJob(context.Background(),val,key)
    ```
-  value is the input to the process function, which can be a primitive type or a struct if multiple input
   parameters are required  
   
-  key is can be a string or an integer (8,16,32,64). Jobs with same key will be processed by the same worker.
   
-  output can be retreived from the `pool.Output` channel which returns a `Job` type structure containg the 
   original input,output and the possible errors. The `GetID` method returns a unique id to the Job that 
   can be used to align the input with output.
   
- more documentation can be found [here](/doc/index.html)
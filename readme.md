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
    
    conf := worker_pool.NewPoolConfig()
    conf.Workers = 100
    conf.WorkerBufferSize = 10
    conf.Metrics.NameSpace = "test_namespace"
    conf.Metrics.SubSystem = "test_subsystem"
    conf.HashFunc = worker_pool.SHA256
  
    pool,err := worker_pool.NewPool(conf)
    if err != nil{
  	    //todo- handle error as preffered
    }
  
    pool.Init(context.Background(), func(ctx context.Context, in interface{}) (out interface{}, err error) {
        //todo - process function body
        return out,nil
    })
  
    ```
- `conf.Workers` refers to the number of workers(spawned in go routines) running in parallel, and `conf.WorkerBufferSize` refers to the size of 
   the workers buffered channel.
   
- `conf.HashFunc` refers to the bucket hashing algorithm used by job manager
- `conf.Metrics` is used to set the prometheus metric name space and subsystem for pool metrics
- `worker_pool.DefaultConfig` can be passed as the config for worker_pool.NewPool() function in which case the 
  defaults will be as follows
    ```
        Workers           = 100
        WorkerBufferSize  = 10
        Metrics.NameSpace = "worker_pool"
        Metrics.SubSystem = "worker_pool"
        HashFunc          =  worker_pool.SHA256
    ```
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
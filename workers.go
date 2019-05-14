package worker_pool

import (
	"github.com/google/uuid"
	"fmt"
)

type worker struct {
	id uuid.UUID
	pool *Pool
}


/*
	Func run()
	----------
	run the Worker as a go routine untill the pool is closed.
 */
func (w *worker)run(){
	go func() {
		fmt.Println("starting worker, poolID:", w.pool.id,"workerID:",w.id)
		for{
			select{
			case job := <- w.pool.input:
				res,err := w.pool.processFunc(job.ctx,job.input)
				w.pool.Output <- Job{
					ctx: job.ctx,
					id:  job.id,
					output: res,
					err: err,
				}

			case _ = <-w.pool.closeWorkers:
				fmt.Println("shutting down worker, poolID:",w.pool.id, "workerID:" ,w.id)
				w.pool.wGroup.Done()
				return
			}
		}

	}()
}


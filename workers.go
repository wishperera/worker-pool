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
func (w *worker)Run(){
	go func() {
		fmt.Printf("starting worker, poolID: %v, workerID: %v",w.pool.id,w.id)
		for{
			select{
			case job := <- w.pool.input:
				res,err := w.pool.processFunc(job.ctx,job.input)
				w.pool.output <- Job{
					ctx: job.ctx,
					id:  job.id,
					output: res,
					err: err,
				}

			case _ = <-w.pool.closeWorkers:
				fmt.Printf("closing down worker, poolID: %v , workerID: %v" ,w.pool.id,w.id)
				w.pool.wGroup.Done()
			}
		}

	}()
}


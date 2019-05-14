package worker_pool

import (
	"fmt"
	"github.com/google/uuid"
)

type worker struct {
	id   uuid.UUID
	pool *Pool
}


//run the Worker as a go routine untill the pool is closed.
func (w *worker) run() {
	go func() {
		fmt.Println("starting worker, poolID:", w.pool.id, "workerID:", w.id)
		for {
			select {
			case job := <-w.pool.input:
				res, err := w.pool.processFunc(job.ctx, job.input)
				w.pool.Output <- Job{
					ctx:    job.ctx,
					input:  job.input,
					id:     job.id,
					output: res,
					err:    err,
				}

			case _ = <-w.pool.closeWorkers:
				fmt.Println("shutting down worker, poolID:", w.pool.id, "workerID:", w.id)
				w.pool.wGroup.Done()
				return
			}
		}

	}()
}

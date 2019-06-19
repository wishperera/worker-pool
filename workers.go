package worker_pool

import (
	"log"
	"github.com/google/uuid"
	"time"
)

type worker struct {
	id   uuid.UUID
	pool *Pool
	buffer chan Job
}


//run the Worker as a go routine untill the pool is closed.
func (w *worker) run() {
	go func() {
		log.Println("[worker] starting worker, poolID:", w.pool.id, "workerID:", w.id)
		for job := range w.buffer{
				st := time.Now()
				res, err := w.pool.processFunc(job.ctx, job.input)
				processLatency.Observe(float64(time.Now().Sub(st).Nanoseconds())/1e3)

				log.Println("[worker] processed job, job_id:",job.id,"worker_id:",w.id)
				w.pool.Output <- Job{
					ctx:    job.ctx,
					input:  job.input,
					id:     job.id,
					key:    job.key,
					output: res,
					err:    err,
				}
			}

		log.Println("[worker] shutting down worker, poolID:", w.pool.id, "workerID:", w.id)
		w.pool.wGroup.Done()
		return
	}()
}

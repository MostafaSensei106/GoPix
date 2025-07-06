package worker

import (
	"context"
	"golang.org/x/time/rate"
	"sync"
    "github.com/mostafasensei106/gopix/internal/convert"
)

type Job struct {
	Path   string
	Format string
}

type WorkerPool struct {
	workers   uint8
	jobs      chan Job
	results   chan *converter.ConversionResult
	converter *converter.ImageConverter
	limiter   *rate.Limiter
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewWorkerPool(workers uint8, converter *converter.ImageConverter, rateLimit float64) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	var limiter *rate.Limiter
	if rateLimit > 0 {
		limiter = rate.NewLimiter(rate.Limit(rateLimit), 1)
	}

	return &WorkerPool{
		workers:   workers,
		jobs:      make(chan Job, workers*2),
		results:   make(chan *converter.ConversionResult, workers*2),
		converter: converter,
		limiter:   limiter,
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (wp *WorkerPool) Start() {
	for i := uint8(0); i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
	wp.cancel()
}

func (wp *WorkerPool) AddJob(job Job) {
	select {
	case wp.jobs <- job:
	case <-wp.ctx.Done():
	}
}

func (wp *WorkerPool) Results() <-chan *converter.ConversionResult {
	return wp.results
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for {
		select {
		case job, ok := <-wp.jobs:
			if !ok {
				return
			}

			// Apply rate limiting if configured
			if wp.limiter != nil {
				if err := wp.limiter.Wait(wp.ctx); err != nil {
					continue
				}
			}

			result := wp.converter.Convert(job.Path, job.Format)

			select {
			case wp.results <- result:
			case <-wp.ctx.Done():
				return
			}

		case <-wp.ctx.Done():
			return
		}
	}
}

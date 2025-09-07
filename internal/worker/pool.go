package worker

import (
	"context"
	"sync"

	"golang.org/x/time/rate"

	conv "github.com/MostafaSensei106/GoPix/internal/converter"
)

type Job struct {
	Path       string
	Format     string
	OutputPath string // Optional custom output path for batch processing
}

type WorkerPool struct {
	workers   uint8
	jobs      chan Job
	results   chan *conv.ConversionResult
	converter *conv.ImageConverter
	limiter   *rate.Limiter
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

// NewWorkerPool creates a new WorkerPool with the specified number of workers,
// an ImageConverter for handling image conversion jobs, and an optional rate
// limit to control the processing rate. The function initializes a context with
// cancellation, sets up job and result channels, and configures rate limiting
// if a positive rateLimit is provided.

func NewWorkerPool(workers uint8, converter *conv.ImageConverter, rateLimit float64) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	var limiter *rate.Limiter
	if rateLimit > 0 {
		limiter = rate.NewLimiter(rate.Limit(rateLimit), 1)
	}

	return &WorkerPool{
		workers:   workers,
		jobs:      make(chan Job, workers*2),
		results:   make(chan *conv.ConversionResult, workers*2),
		converter: converter,
		limiter:   limiter,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Start initializes the worker pool by spawning the specified number of
// worker goroutines. Each worker will process jobs from the job channel
// until the channel is closed or the context is cancelled. This function
// should be called before adding jobs to ensure workers are ready to process.

func (wp *WorkerPool) Start() {
	for i := uint8(0); i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

// Stop gracefully shuts down the worker pool by closing the job channel,
// waiting for all ongoing tasks to complete, and then closing the results
// channel. It also cancels the context, signaling that no further processing
// should occur. This ensures that all resources are released properly and
// no new jobs are processed.

func (wp *WorkerPool) Stop() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
	wp.cancel()
}

// AddJob adds a job to the job channel. If the context is cancelled, it will not add the job and return immediately.
func (wp *WorkerPool) AddJob(job Job) {
	select {
	case wp.jobs <- job:
	case <-wp.ctx.Done():
	}
}

// Results returns a receive-only channel of ConversionResult pointers.
// This channel provides the results of processed jobs. It can be used
// to retrieve conversion results as they become available.

func (wp *WorkerPool) Results() <-chan *conv.ConversionResult {
	return wp.results
}

// worker is a goroutine function that continuously processes jobs from the job channel.
// It applies rate limiting if a limiter is configured and handles job cancellations gracefully.
// Upon processing each job, it sends the conversion result to the results channel.
// The function exits when the job channel is closed or the context is cancelled.

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

			var result *conv.ConversionResult
			if job.OutputPath != "" {
				result = wp.converter.ConvertWithOutputPath(job.Path, job.Format, job.OutputPath)
			} else {
				result = wp.converter.Convert(job.Path, job.Format)
			}

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

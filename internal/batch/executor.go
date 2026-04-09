package batch

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Task represents a single unit of work in a batch operation
type Task struct {
	ID   string
	Name string
	Fn   func(ctx context.Context) error
}

// Result holds the outcome of a batch execution
type Result struct {
	Total     int
	Succeeded int
	Failed    int
	Errors    []TaskError
	Duration  time.Duration
}

// TaskError records a failure
type TaskError struct {
	TaskID   string
	TaskName string
	Err      error
}

func (e TaskError) Error() string {
	return fmt.Sprintf("task %s (%s): %v", e.TaskID, e.TaskName, e.Err)
}

// ProgressFunc is called after each task completes
type ProgressFunc func(completed, total int, lastTaskName string, lastErr error)

// Executor runs batch tasks with concurrency control
type Executor struct {
	Concurrency int
	OnProgress  ProgressFunc
	SkipErrors  bool
}

// NewExecutor creates a new batch executor
func NewExecutor(concurrency int) *Executor {
	if concurrency < 1 {
		concurrency = 1
	}
	return &Executor{
		Concurrency: concurrency,
	}
}

// Execute runs all tasks with the configured concurrency
func (e *Executor) Execute(ctx context.Context, tasks []Task) *Result {
	start := time.Now()
	result := &Result{Total: len(tasks)}

	if len(tasks) == 0 {
		return result
	}

	var (
		completed int64
		mu        sync.Mutex
		wg        sync.WaitGroup
		sem       = make(chan struct{}, e.Concurrency)
	)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, task := range tasks {
		select {
		case <-ctx.Done():
			break
		case sem <- struct{}{}:
		}

		wg.Add(1)
		go func(t Task) {
			defer wg.Done()
			defer func() { <-sem }()

			err := t.Fn(ctx)

			mu.Lock()
			if err != nil {
				result.Failed++
				result.Errors = append(result.Errors, TaskError{
					TaskID:   t.ID,
					TaskName: t.Name,
					Err:      err,
				})
				if !e.SkipErrors {
					cancel()
				}
			} else {
				result.Succeeded++
			}
			mu.Unlock()

			c := atomic.AddInt64(&completed, 1)
			if e.OnProgress != nil {
				e.OnProgress(int(c), len(tasks), t.Name, err)
			}
		}(task)
	}

	wg.Wait()
	result.Duration = time.Since(start)
	return result
}

// FormatResult returns a human-readable summary
func FormatResult(r *Result) string {
	return fmt.Sprintf("完成: 成功 %d, 失败 %d, 总计 %d (耗时 %s)",
		r.Succeeded, r.Failed, r.Total, r.Duration.Round(time.Millisecond))
}

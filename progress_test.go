package progress

import (
	"sync"
	"testing"
	"time"
)

func TestProgress(t *testing.T) {
	progress := New()
	progress.Start()
	wg := &sync.WaitGroup{}
	wg.Add(3)
	go worker("thread1", progress, wg)
	time.Sleep(100 * time.Millisecond)
	go worker("thread2", progress, wg)
	time.Sleep(100 * time.Millisecond)
	go worker("thread3", progress, wg)
	wg.Wait()
	progress.Stop()
}

func worker(name string, progress *Progress, wg *sync.WaitGroup) {
	defer wg.Done()
	bar := progress.AddBar(name, 100)
	for i := 1; i <= 100; i++ {
		bar.Set(i)
		time.Sleep(100 * time.Millisecond)
	}
}

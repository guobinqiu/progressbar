package progress

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type Progress struct {
	bars      []*Bar
	stopCh    chan bool
	interval  time.Duration
	out       io.Writer
	lineCount int
	mu        *sync.Mutex
}

func New() *Progress {
	return &Progress{
		stopCh:    make(chan bool, 1),
		interval:  10 * time.Millisecond,
		out:       os.Stdout,
		mu:        &sync.Mutex{},
		lineCount: 0,
	}
}

func (p *Progress) Start() {
	go p.listen()
}

func (p *Progress) listen() {
	t := time.NewTicker(p.interval)
	for {
		select {
		case <-t.C:
			p.print()
		case <-p.stopCh:
			close(p.stopCh)
			t.Stop()
			return
		}
	}
}

func (p *Progress) Stop() {
	p.stopCh <- true
}

func (p *Progress) print() {
	p.mu.Lock()
	defer p.mu.Unlock()

	var buf bytes.Buffer
	for _, bar := range p.bars {
		buf.WriteString(bar.getPercentString())
	}

	//光标回退到上一次输出位置
	fmt.Print(strings.Repeat("\033[1A", p.lineCount))

	//本次输出
	fmt.Print(buf.String())

	//记录本次输出行数供下次光标回退用
	p.lineCount = len(p.bars)
}

func (p *Progress) AddBar(name string, total int) *Bar {
	p.mu.Lock()
	defer p.mu.Unlock()

	bar := NewBar(name, total)
	p.bars = append(p.bars, bar)

	return bar
}

type Bar struct {
	name    string
	total   int
	current int
	mu      *sync.Mutex
}

func NewBar(name string, total int) *Bar {
	return &Bar{
		name:  name,
		total: total,
		mu:    &sync.Mutex{},
	}
}

func (b *Bar) Set(current int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.current = current
}

func (b *Bar) getPercentString() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return fmt.Sprintf("%s:%3.f%%\n", b.name, float64(b.current)/float64(b.total)*100)
}

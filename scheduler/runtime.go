package main

import "fmt"

const maxRunDur = 100

type gStatus string

var (
	statusRunnable gStatus = "runnable"
	statusRunning  gStatus = "running"
	statusWaiting  gStatus = "waiting"
	statusDead     gStatus = "dead"
)

type Goroutine struct {
	id     int
	runDur int
	status gStatus
}

func (g *Goroutine) Block() {
	if g.status != statusRunning {
		panic("invalid status for block")
	}
	g.status = statusWaiting
}

func (g *Goroutine) Unblock() {
	if g.status != statusWaiting {
		panic("invalid status for unblock")
	}
	g.status = statusRunnable
}

func (g *Goroutine) Done() {
	if g.status != statusRunning {
		panic("invalid status for done")
	}
	g.status = statusDead
}

type Thread struct {
	id   int
	goro *Goroutine
}

type RuntimeState struct {
	dur      int
	threads  map[int]int
	runnable []int
	running  []int
	waiting  []int
	dead     []int
}

type Runtime struct {
	dur      int
	nGoro    int
	threads  []*Thread
	runnable []*Goroutine
	running  []*Goroutine
	waiting  []*Goroutine
	dead     []*Goroutine
}

func NewRuntime(gomaxprocs int) *Runtime {
	threads := make([]*Thread, gomaxprocs)
	for i := 0; i < gomaxprocs; i++ {
		threads[i] = &Thread{id: i + 1}
	}
	return &Runtime{threads: threads}
}

func (r *Runtime) Go() *Goroutine {
	r.nGoro++
	g := &Goroutine{id: r.nGoro, status: statusRunnable}
	r.runnable = append(r.runnable, g)
	return g
}

func (r *Runtime) Forward(dur int) {
	r.dur += dur
	for _, t := range r.threads {
		if g := t.goro; g != nil && g.status == statusRunning {
			g.runDur += dur
		}
	}
}

func (r *Runtime) Schedule() {
	for _, t := range r.threads {
		if g := t.goro; g != nil {
			switch g.status {
			case statusRunning:
				if g.runDur >= maxRunDur {
					g.runDur = 0
					g.status = statusRunnable
					r.runnable = append(r.runnable, g)
					t.goro = nil
				}
			case statusWaiting:
				r.waiting = append(r.waiting, g)
				t.goro = nil
			case statusDead:
				r.dead = append(r.dead, g)
				t.goro = nil
			}
		}
	}

	var nw []*Goroutine
	for _, g := range r.waiting {
		if g.status == statusRunnable {
			r.runnable = append(r.runnable, g)
		} else {
			nw = append(nw, g)
		}
	}
	r.waiting = nw

	for _, t := range r.threads {
		if t.goro == nil && len(r.runnable) > 0 {
			g := r.runnable[0]
			r.runnable = r.runnable[1:]
			g.status = statusRunning
			g.runDur = 0
			t.goro = g
		}
	}

	r.running = r.running[:0]
	for _, t := range r.threads {
		if t.goro != nil {
			r.running = append(r.running, t.goro)
		}
	}
}

func (r *Runtime) State() RuntimeState {
	threads := make(map[int]int)
	for _, t := range r.threads {
		if t.goro != nil {
			threads[t.id] = t.goro.id
		} else {
			threads[t.id] = 0
		}
	}
	runnable := make([]int, len(r.runnable))
	for i, g := range r.runnable {
		runnable[i] = g.id
	}
	running := make([]int, len(r.running))
	for i, g := range r.running {
		running[i] = g.id
	}
	waiting := make([]int, len(r.waiting))
	for i, g := range r.waiting {
		waiting[i] = g.id
	}
	dead := make([]int, len(r.dead))
	for i, g := range r.dead {
		dead[i] = g.id
	}
	return RuntimeState{
		dur:      r.dur,
		threads:  threads,
		runnable: runnable,
		running:  running,
		waiting:  waiting,
		dead:     dead,
	}
}

func main() {
	r := NewRuntime(2)
	g1 := r.Go()
	g2 := r.Go()
	r.Go()
	r.Go()
	r.Schedule()
	r.Forward(10)
	g1.Done()
	g2.Block()
	r.Schedule()
	state := r.State()
	fmt.Printf("%+v\n", state)
}

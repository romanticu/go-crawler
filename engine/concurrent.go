package engine

import "fmt"

type ConcurrentEngine struct {
	Scheduler Scheduler
	WorkerCount int
	ItemChan chan interface{}
}

type Scheduler interface {
	Submit(Request)
	ConfigureMasterWorkerChan(chan Request)
	WorkerReady(chan Request)
	Run()
}
func (e *ConcurrentEngine) Run(seeds ...Request)  {
	out := make(chan ParseResult)
	e.Scheduler.Run()
	for i := 0; i < e.WorkerCount; i++ {
		createWorker(out, e.Scheduler)
	}
	for _,r := range seeds {
		e.Scheduler.Submit(r)
	}
	itemCount := 0
	for {
		result := <- out
		for _, item := range result.Items  {
			fmt.Printf("Got #: %d item %v \n",itemCount, item)
			itemCount++
			go func() {
				e.ItemChan <- item
			}()
		}
		for _, request := range result.Requests  {
			e.Scheduler.Submit(request)
		}
	}
}

func createWorker(out chan ParseResult, s Scheduler) {
	in := make(chan Request)
	go func() {
		for {
			s.WorkerReady(in)
			request := <-in
			result, err := worker(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}
//Implement the dining philosopher’s problem with the following constraints/modifications.
//
//There should be 5 philosophers sharing chopsticks, with one chopstick between each adjacent pair of philosophers.
//
//Each philosopher should eat only 3 times (not in an infinite loop as we did in lecture)
//
//The philosophers pick up the chopsticks in any order, not lowest-numbered first (which we did in lecture).
//
//In order to eat, a philosopher must get permission from a host which executes in its own goroutine.
//
//The host allows no more than 2 philosophers to eat concurrently.
//
//Each philosopher is numbered, 1 through 5.
//
//When a philosopher starts eating (after it has obtained necessary locks) it prints “starting to eat <number>” on a line by itself, where <number> is the number of the philosopher.
//
//When a philosopher finishes eating (before it has released its locks) it prints “finishing eating <number>” on a line by itself, where <number> is the number of the philosopher.

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type stick struct {
	stickMutex sync.Mutex
}

type philosopher struct {
	index         int
	leftS, rightS *stick
	h             *host
}

type host struct {
	ch chan struct{}
}

func (p *philosopher) eat(i int) {
	p.h.askPermission()
	fmt.Printf("Philosopher %d is starting to eat its meal %d\n\n", p.index, i+1)
	time.Sleep(2 * time.Second)
	fmt.Printf("Philosopher %d finished eating its meal %d\n\n", p.index, i+1)
	if i == 2 {
		fmt.Printf("**** Philosopher %d is done! *****\n\n", p.index)
	}
	p.h.communicateEnd()
}
func (p *philosopher) start(wg *sync.WaitGroup) {
	for i := 0; i < 3; i++ {
		choice := rand.Intn(2)
		if choice == 0 {
			p.leftS.stickMutex.Lock()
			p.rightS.stickMutex.Lock()
		} else {
			p.rightS.stickMutex.Lock()
			p.leftS.stickMutex.Lock()
		}
		p.eat(i)
		choice = rand.Intn(2)
		if choice == 0 {
			p.leftS.stickMutex.Unlock()
			p.rightS.stickMutex.Unlock()
		} else {
			p.rightS.stickMutex.Unlock()
			p.leftS.stickMutex.Unlock()
		}
	}
	wg.Done()
}
func (h *host) askPermission() {
	h.ch <- struct{}{}
}
func (h *host) communicateEnd() {
	<-h.ch
}
func (h *host) start(workersDone *sync.WaitGroup, hostStart *sync.WaitGroup, hostDone *sync.WaitGroup) {
	fmt.Println("Host has started\n")
	hostStart.Done()
	workersDone.Wait()
	fmt.Println("Host is finished")
	hostDone.Done()
}

func main() {
	var philosophers []philosopher
	var sticks []stick
	mainHost := host{make(chan struct{}, 2)}
	for i := 0; i < 5; i++ {
		sticks = append(sticks, stick{})
	}
	for i := 0; i < 5; i++ {
		philosophers = append(philosophers, philosopher{
			index:  i + 1,
			leftS:  &sticks[i],
			rightS: &sticks[(i+1)%5],
			h:      &mainHost,
		})
	}
	workDone := sync.WaitGroup{}
	workDone.Add(5)
	hostStart := sync.WaitGroup{}
	hostStart.Add(1)
	hostDone := sync.WaitGroup{}
	hostDone.Add(1)
	go mainHost.start(&workDone, &hostStart, &hostDone)
	hostStart.Wait()

	for i := 0; i < 5; i++ {
		go philosophers[i].start(&workDone)
	}
	hostDone.Wait()
}

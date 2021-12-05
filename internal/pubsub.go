package internal

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Message[T any] struct {
	Topic   string
	Meta    map[string]string
	Ctx     context.Context
	Content T
}

type Reg[T any] struct {
	topic   string
	channel chan Message[T]
}

type PubSub[T any] struct {
	mu     sync.RWMutex
	regs   map[string]map[Reg[T]]struct{}
	closed bool
}

func NewPubsub[T any]() *PubSub[T] {
	ps := &PubSub[T]{}
	ps.regs = make(map[string]map[Reg[T]]struct{})
	return ps
}

func (ps *PubSub[T]) Subscribe(topic string) Reg[T] {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	channel := make(chan Message[T])
	reg := Reg[T]{topic, channel}
	_, found := ps.regs[topic]
	if !found {
		ps.regs[topic] = make(map[Reg[T]]struct{})
	}
	ps.regs[topic][reg] = struct{}{}
	return reg
}

func removeAll[T any](ch chan T) {
	for {
		select {
		case <-ch: 
			fmt.Println("remove")
			continue
		default:
			fmt.Println("no  activity")
			return
		}
	}
}

func (ps *PubSub[T]) Unsubscribe(reg Reg[T]) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	regs, found := ps.regs[reg.topic]
	if found {
		fmt.Println("UNSUBSCRIBE", reg.topic)
		delete(regs, reg)
		if len(regs) == 0 {
			delete(ps.regs, reg.topic)
		}
		removeAll[Message[T]](reg.channel)
		fmt.Println("DONE UNSUBSCRIBE", reg.topic)
		close(reg.channel)
	}
}

func (ps *PubSub[T]) Publish(topic string, content T) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}

	for reg, _ := range ps.regs[topic] {
		reg.channel <- Message[T]{Topic: topic, Content: content}
	}
}

func (ps *PubSub[T]) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.closed {
		ps.closed = true
		for _, regs := range ps.regs {
			for reg, _ := range regs {
				close(reg.channel)
			}
		}
	}
}

const topic = "TOPIC"

func worker(name string, ps *PubSub[string], done func()) {
	defer done()

	reg := ps.Subscribe(topic)
	fmt.Println("Subscribed", name)
	for {
		msg := <-reg.channel
		fmt.Println(name, msg.Topic, msg.Content)
		if msg.Content == "stop" {
			ps.Unsubscribe(reg)
			fmt.Println("Unsubscribe", name)
			return
		}
	}
}

func publish(ps *PubSub[string], text string) {
	fmt.Println("Publis", text)
	ps.Publish(topic, text)
}

func PubSubDemo() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(3)

	ps := NewPubsub[string]()
	for i := 1; i <= 3; i++ {
		name := fmt.Sprintf("sub%d", i)
		go worker(name, ps, waitGroup.Done)
	}
	time.Sleep(200 * time.Millisecond)
	publish(ps, "one")
	publish(ps, "two")
	publish(ps, "stop")
	//time.Sleep(200 * time.Millisecond)
	for i := 0; i < 100; i++ {
		ps.Publish(topic, fmt.Sprintf("msg-%d", i))
	}
	waitGroup.Wait()
}

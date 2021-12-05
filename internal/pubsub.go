package internal

import (
	"fmt"
	"sync"
)

type Message[T any] struct {
	Topic   string
	Content T
}

type Reg[T any] struct {
	topic   string
	channel chan Message[T]
}

type PubSub[T any] struct {
	mu     sync.RWMutex
	subs   map[string][]Reg[T]
	closed bool
}

func NewPubsub[T any]() *PubSub[T] {
	ps := &PubSub[T]{}
	ps.subs = make(map[string][]Reg[T])
	return ps
}

func (ps *PubSub[T]) Subscribe(topic string, channel chan Message[T]) Reg[T] {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	reg := Reg[T]{topic, channel}
	ps.subs[topic] = append(ps.subs[topic], reg)
	return reg
}

func (ps *PubSub[T]) Unsubscribe(reg Reg[T]) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	regs, found := ps.subs[reg.topic]
	if found {
		// maybe close channel?
		ps.subs[reg.topic] = Filter[Reg[T]](regs, func(item Reg[T]) bool { return item != reg })
	}
}

func (ps *PubSub[T]) Publish(topic string, content T) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}

	for _, reg := range ps.subs[topic] {
		reg.channel <- Message[T]{topic, content}
	}
}

func (ps *PubSub[T]) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.closed {
		ps.closed = true
		for _, subs := range ps.subs {
			for _, reg := range subs {
				close(reg.channel)
			}
		}
	}
}

const topic = "TOPIC"

func worker(name string, ps *PubSub[string], done func()) {
	defer done()

	channel := make(chan Message[string])
	reg := ps.Subscribe(topic, channel)
	for {
		msg := <-channel
		fmt.Println(name, msg.Topic, msg.Content)
		if msg.Content == "stop" {
			ps.Unsubscribe(reg)
			fmt.Println("Unsubscribe", name)
			return
		}
	}
}

func PubSubDemo() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(3)

	ps := NewPubsub[string]()
	for i := 1; i <= 3; i++ {
		name := fmt.Sprintf("sub%d", i)
		go worker(name, ps, waitGroup.Done)
	}
	ps.Publish(topic, "one")
	ps.Publish(topic, "two")
	ps.Publish(topic, "stop")
	ps.Publish(topic, "three")

	waitGroup.Wait()
}

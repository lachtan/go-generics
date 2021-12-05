package internal

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Message conusmer
// if return true unsubsribe this registration
type Consumer[T any] func(message Message[T]) bool

type Message[T any] struct {
	Topic   string
	Meta    map[string]string
	Ctx     context.Context
	Content T
}

type Reg[T any] struct {
	topic   string
	id      int64
	consume Consumer[T]
}

type Unreg[T any] struct {
	reg Reg[T]
}

type PubSub[T any] struct {
	queue     chan interface{}
	consumers map[string][]Reg[T]
	done      chan struct{}
	closed    bool
}

func NewPubsub[T any]() *PubSub[T] {
	ps := &PubSub[T]{queue: make(chan interface{}, 100), done: make(chan struct{})}
	ps.consumers = make(map[string][]Reg[T])
	go ps.mainLoop()
	return ps
}

func (ps *PubSub[T]) register(reg Reg[T]) {
	fmt.Println("Register", reg)
	ps.consumers[reg.topic] = append(ps.consumers[reg.topic], reg)
}

func (ps *PubSub[T]) unregister(reg Reg[T]) {
	fmt.Println("Unregister", reg)
	consumers, found := ps.consumers[reg.topic]
	if found {
		index := Find[Reg[T]](consumers, func(x Reg[T]) bool { return x.id == reg.id })
		if index > 0 {
			ps.consumers[reg.topic] = RemoveIndex[Reg[T]](consumers, index)
		}
	}
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

func (ps *PubSub[T]) send(msg Message[T]) {
	consumers, found := ps.consumers[msg.Topic]
	if found {
		for _, consumer := range consumers {
			unregister := consumer.consume(msg)
			if unregister {
				ps.unregister(consumer)
			}
		}
	}
}

func (ps *PubSub[T]) process(msg interface{}) {
	switch v := msg.(type) {
	case Reg[T]: ps.register(v)
	case Unreg[T]: ps.unregister(v.reg)
	case Message[T]: ps.send(v)
	default: panic("Unknown message type")
	}
}

func (ps *PubSub[T]) close() {
	ps.closed = true
	close(ps.queue)
	ps.consumers = nil
}

func (ps *PubSub[T]) mainLoop() {
	for {
		select {
		case msg := <-ps.queue:
			ps.process(msg)
		case <-ps.done:
			ps.close()
			return
		}
	}
}

func (ps *PubSub[T]) Subscribe(topic string, consume Consumer[T]) Reg[T] {
	reg := Reg[T]{topic: topic, id: time.Now().UnixNano(), consume: consume}
	ps.queue <- reg
	return reg
}

func (ps *PubSub[T]) Unsubscribe(reg Reg[T]) {
	ps.queue <- Unreg[T]{reg}
}

func (ps *PubSub[T]) Publish(topic string, content T) {
	if !ps.closed {
		ps.queue <- Message[T]{Topic: topic, Content: content}
	}
}

func (ps *PubSub[T]) Close() {
	ps.done <- struct{}{}
}

const topic = "TOPIC"

func worker(name string, ps *PubSub[string], stop string, done func()) {
	reg := ps.Subscribe(topic, func(msg Message[string]) bool { 
		fmt.Println(name, "Received", msg.Topic, msg.Content)
		quit := msg.Content == stop
		if  quit{
			done()
		}
		return quit
	})
	fmt.Println("Subscribed", name, reg)
}

func publish(ps *PubSub[string], text string) {
	fmt.Println("Publis", text)
	ps.Publish(topic, text)
}

func PubSubDemo() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(3)

	ps := NewPubsub[string]()
	go worker("sub-1", ps, "stop", waitGroup.Done)
	go worker("sub-2", ps, "two", waitGroup.Done)
	go worker("sub-3", ps, "one", waitGroup.Done)
	time.Sleep(200 * time.Millisecond)
	publish(ps, "one")
	publish(ps, "two")
	publish(ps, "stop")
	//time.Sleep(200 * time.Millisecond)
	for i := 0; i < 10; i++ {
		ps.Publish(topic, fmt.Sprintf("msg-%d", i))
	}
	waitGroup.Wait()
}

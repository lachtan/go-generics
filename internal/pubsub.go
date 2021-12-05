package internal

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Consumer[T any] interface {
	Consume(message Message[T]) bool
}

type Message[T any] struct {
	Topic   string
	Meta    map[string]string
	Ctx     context.Context
	Content T
}

// Register request and unregistre token
type Reg[T any] struct {
	topic    string
	messages chan Message[T]
}

// Unregister command
type Unreg[T any] struct {
	reg Reg[T]
}

type PubSub[T any] struct {
	// commans or messages
	queue     chan interface{}
	consumers map[string]map[Reg[T]]struct{}
	closed    bool
}

func NewPubsub[T any]() *PubSub[T] {
	ps := &PubSub[T]{queue: make(chan interface{}, 100)}
	ps.consumers = make(map[string]map[Reg[T]]struct{})
	go ps.mainLoop()
	return ps
}

func (ps *PubSub[T]) register(reg Reg[T]) {
	_, found := ps.consumers[reg.topic]
	if !found {
		ps.consumers[reg.topic] = make(map[Reg[T]]struct{})
	}
	ps.consumers[reg.topic][reg] = struct{}{}
}

func (ps *PubSub[T]) unregister(reg Reg[T]) {
	regs, found := ps.consumers[reg.topic]
	if found {
		fmt.Println("UNSUBSCRIBE", reg.topic)
		delete(regs, reg)
		if len(regs) == 0 {
			delete(ps.consumers, reg.topic)
		}
		removeAll[Message[T]](reg.messages)
		fmt.Println("DONE UNSUBSCRIBE", reg.topic)
		close(reg.messages)
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
		for consumer, _ := range consumers {
			consumer.messages <- msg 
		}
	}
}

func (ps *PubSub[T]) mainLoop() {
	for {
		msg := <-ps.queue
		switch v := msg.(type) {
		case Reg[T]:
			ps.register(v)
		case Unreg[T]:
			ps.unregister(v.reg)
		case Message[T]:
			ps.send(v)
		default:
			panic("Unknown message type")
		}
	}
}

func (ps *PubSub[T]) Subscribe(topic string) Reg[T] {
	channel := make(chan Message[T], 1)
	reg := Reg[T]{topic, channel}
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
	// Close request
	// if !ps.closed {
	// 	ps.closed = true
	// 	for _, regs := range ps.regs {
	// 		for reg, _ := range regs {
	// 			close(reg.messages)
	// 		}
	// 	}
	// }
}

const topic = "TOPIC"

func worker(name string, ps *PubSub[string], done func()) {
	defer done()

	reg := ps.Subscribe(topic)
	fmt.Println("Subscribed", name)
	for {
		// mel by to vycitat s timeoutem
		msg := <-reg.messages
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
	for i := 0; i < 10; i++ {
		ps.Publish(topic, fmt.Sprintf("msg-%d", i))
	}
	waitGroup.Wait()
}

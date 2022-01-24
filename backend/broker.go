package main

import "sync"

// The broker is used to fan out messages to multiple clients currently subscribed
// via server sent events.
type Broker struct {
	subscriptions   []chan string
	subscribeChan   chan (chan string)
	unsubscribeChan chan (chan string)
	publishChan     chan string
}

func NewBroker() *Broker {
	b := &Broker{}
	b.subscribeChan = make(chan (chan string))
	b.unsubscribeChan = make(chan (chan string))
	b.publishChan = make(chan string, 100)
	return b
}

func (b *Broker) Run() *sync.WaitGroup {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		for {
			select {
			case s := <-b.subscribeChan:
				b.subscriptions = append(b.subscriptions, s)
			case u := <-b.unsubscribeChan:
				for i, ch := range b.subscriptions {
					if ch == u {
						b.subscriptions = append(b.subscriptions[:i], b.subscriptions[i+1:]...)
						close(u)
						break
					}
				}
			case msg := <-b.publishChan:
				for _, ch := range b.subscriptions {
					ch <- msg
				}
			}
		}
	}()
	return wg
}

func (b *Broker) Subscribe() chan string {
	ch := make(chan string)
	b.subscribeChan <- ch
	return ch
}

func (b *Broker) Unsubscribe(ch chan string) {
	b.unsubscribeChan <- ch
}

func (b *Broker) Publish(message string) {
	b.publishChan <- message
}

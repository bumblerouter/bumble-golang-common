package envelope

import ()

type Queue struct {
	channel chan *Envelope
}

func NewQueue() *Queue {
	q := new(Queue)
	q.channel = make(chan *Envelope)
	return q
}

func (q *Queue) Add(e *Envelope) {
	q.channel <- e
}

func (q *Queue) Channel() chan *Envelope {
	return q.channel
}

package painter

import (
	"image"
	"log"
	"sync"

	"golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циклі подій.
type Receiver interface {
	Update(t screen.Texture)
}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
	Receiver Receiver

	next screen.Texture // текстура, яка зараз формується
	prev screen.Texture // текстура, яка була відправлення останнього разу у Receiver

	mq messageQueue

	stop    chan struct{}
	stopReq bool
}

var size = image.Pt(400, 400)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start(s screen.Screen) {
	var err error
	l.next, err = s.NewTexture(size)
	if err != nil {
		log.Fatalf("failed to create texture: %v", err)
	}
	l.prev, err = s.NewTexture(size)
	if err != nil {
		log.Fatalf("failed to create texture: %v", err)
	}

	l.stop = make(chan struct{})
	go func() {
		for !l.stopReq || !l.mq.empty() {
			message := l.mq.pull()
			if update := message.Do(l.next); update {
				l.Receiver.Update(l.next)
				l.next, l.prev = l.prev, l.next
			}
		}
		close(l.stop)
	}()
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {
	l.mq.push(op)
}

// StopAndWait сигналізує про необхідність завершити цикл та блокується до моменту його повної зупинки.
func (l *Loop) StopAndWait() {
	l.Post(OperationFunc(func(t screen.Texture) {
		l.stopReq = true
	}))
	<-l.stop
}

type messageQueue struct {
	ops  []Operation
	mu   sync.Mutex
	cond *sync.Cond
}

func (mq *messageQueue) push(op Operation) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	mq.ops = append(mq.ops, op)
	if mq.cond != nil {
		mq.cond.Signal()
	}
}

func (mq *messageQueue) pull() Operation {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	for len(mq.ops) == 0 {
		if mq.cond == nil {
			mq.cond = sync.NewCond(&mq.mu)
		}
		mq.cond.Wait()
	}
	op := mq.ops[0]
	mq.ops[0] = nil
	mq.ops = mq.ops[1:]
	return op
}

func (mq *messageQueue) empty() bool {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	return len(mq.ops) == 0
}

package painter

import (
	"image"
	"log"
	"sync"

	"golang.org/x/exp/shiny/screen"
)

const (
	ImageSize = 800
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

	Mq messageQueue

	stop    chan struct{}
	stopReq bool
}

var size = image.Pt(ImageSize, ImageSize)

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

	l.Mq = messageQueue{}
	go func() {
		for !l.stopReq || !l.Mq.Empty() {
			message := l.Mq.Pull()
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
	l.Mq.Push(op)
}

// StopAndWait сигналізує про необхідність завершити цикл та блокується до моменту його повної зупинки.
func (l *Loop) StopAndWait() {
	l.Post(OperationFunc(func(t screen.Texture) {
		l.stopReq = true
	}))
	<-l.stop
}

type messageQueue struct {
	Ops  []Operation
	mu   sync.Mutex
	cond *sync.Cond
}

func (mq *messageQueue) Push(op Operation) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	mq.Ops = append(mq.Ops, op)
	if mq.cond != nil {
		mq.cond.Signal()
	}
}

func (mq *messageQueue) Pull() Operation {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	for len(mq.Ops) == 0 {
		if mq.cond == nil {
			mq.cond = sync.NewCond(&mq.mu)
		}
		mq.cond.Wait()
	}
	op := mq.Ops[0]
	mq.Ops[0] = nil
	mq.Ops = mq.Ops[1:]
	return op
}

func (mq *messageQueue) Empty() bool {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	return len(mq.Ops) == 0
}

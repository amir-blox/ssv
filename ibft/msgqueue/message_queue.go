package msgqueue

import (
	"sync"

	"github.com/bloxapp/ssv/ibft/proto"
)

// MessageQueue is a broker of messages for the ibft instance to process.
type MessageQueue struct {
	msgMutex          sync.Mutex
	currentRound      uint64
	currentRoundQueue []*proto.SignedMessage
	futureRoundQueue  []*proto.SignedMessage
}

func New() *MessageQueue {
	return &MessageQueue{
		msgMutex:          sync.Mutex{},
		currentRoundQueue: make([]*proto.SignedMessage, 0),
		futureRoundQueue:  make([]*proto.SignedMessage, 0),
	}
}

// AddMessage adds a message the queue based on the message round.
// AddMessage is thread safe
func (q *MessageQueue) AddMessage(msg *proto.SignedMessage) {
	q.msgMutex.Lock()
	defer q.msgMutex.Unlock()

	if msg.Message.Round < q.currentRound {
		return // not adding previous round messages
	}

	if q.currentRound == msg.Message.Round {
		q.currentRoundQueue = append(q.currentRoundQueue, msg)
	} else {
		q.futureRoundQueue = append(q.futureRoundQueue, msg)
	}
}

// PopMessage returns and removes a msg from the queue, FIFO
// PopMessage is thread safe
// returns nil if no messages are in the queue
func (q *MessageQueue) PopMessage() *proto.SignedMessage {
	q.msgMutex.Lock()
	defer q.msgMutex.Unlock()

	var ret *proto.SignedMessage
	if len(q.currentRoundQueue) > 0 {
		ret, q.currentRoundQueue = q.currentRoundQueue[0], q.currentRoundQueue[1:]
	}
	return ret
}

func (q *MessageQueue) SetRound(newRound uint64) {
	q.msgMutex.Lock()
	defer q.msgMutex.Unlock()

	// set round
	q.currentRound = newRound

	// empty current round
	q.currentRoundQueue = make([]*proto.SignedMessage, 0)

	// move from future round to current round, also remove dead messages
	newFutureRoundQueue := make([]*proto.SignedMessage, 0)
	for _, msg := range q.futureRoundQueue {
		if msg.Message.Round < q.currentRound {
			// do nothing, will delete this message
		} else if msg.Message.Round == q.currentRound {
			q.currentRoundQueue = append(q.currentRoundQueue, msg)
		} else { //  msg.Message.Round > q.currentRound
			newFutureRoundQueue = append(newFutureRoundQueue, msg)
		}
	}
}

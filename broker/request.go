package broker

import (
	"fmt"
	"math/rand"
)

type requestEvent struct {
	Name      string      `json:"name"`
	Msg       interface{} `json:"msg"`
	RequestId string      `json:"request_id"`
}

func (requestEvent requestEvent) WithRandomRequestId() requestEvent {
	requestEvent.RequestId = fmt.Sprint(rand.Int63n(10000000000))
	return requestEvent
}

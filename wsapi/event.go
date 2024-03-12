package wsapi

type Event struct {
	Name      string      `json:"name"`
	Msg       interface{} `json:"msg"`
	Version   string      `json:"version"`
	RequestId string      `json:"request_id"`
}

func NewEvent(name string, msg map[string]interface{}, requestId string) *Event {
	return &Event{
		Name:      name,
		Msg:       msg,
		Version:   "1.0",
		RequestId: requestId,
	}
}

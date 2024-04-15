package btypes

type RequestEvent struct {
	Name      string      `json:"name"`
	Msg       interface{} `json:"msg"`
	RequestId string      `json:"request_id"`
}

package wsapi

func CreateSubscription(ws *Socket, msg map[string]interface{}, requestId string) *RequestEvent {
	return &RequestEvent{
		Name:      "subscribeMessage",
		Msg:       msg,
		RequestId: requestId,
	}
}

func CreateUnsubscription(subscription RequestEvent) *RequestEvent {
	subscription.Name = "unsubscribeMessage"

	return &subscription
}
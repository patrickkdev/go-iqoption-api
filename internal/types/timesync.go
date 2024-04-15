package types

import (
	"time"
)

type TimesyncEvent struct {
	ClientSessionID string `json:"client_session_id"`
	Msg             int64  `json:"msg"`
	Name            string `json:"name"`
	RequestID       string `json:"request_id"`
	SessionID       string `json:"session_id"`
}

type Timesync struct {
	Name            string
	ServerTimestamp int64
	ExpirationTime  int
}

func (t *Timesync) GetServerTimestamp() int64 {
	return t.ServerTimestamp / 1000
}

func (t *Timesync) GetExpirationDatetime() time.Time {
	return time.Now().Add(time.Minute * time.Duration(t.ExpirationTime))
}

func (t *Timesync) SetExpirationTime(timeInMinutes int) {
	t.ExpirationTime = timeInMinutes
}

func (t *Timesync) ExpirationTimestamp() float64 {
	return float64(t.GetExpirationDatetime().UnixMicro())
}

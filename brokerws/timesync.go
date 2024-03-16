package brokerws

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
	serverTimestamp int64
	expirationTime  int
}

func NewTimesync() *Timesync {
	return &Timesync{
		Name:            "timeSync",
		serverTimestamp: time.Now().UnixNano(),
		expirationTime:  1,
	}
}

func (t *Timesync) GetServerTimestamp() int64 {
	return t.serverTimestamp / 1000
}

func (t *Timesync) GetServerDatetime() time.Time {
	return time.Unix(int64(t.serverTimestamp), 0)
}

func (t *Timesync) SetServerTimestamp(timestamp int64) {
	t.serverTimestamp = timestamp
}

func (t *Timesync) GetExpirationTime() int {
	return t.expirationTime
}

func (t *Timesync) GetExpirationDatetime() time.Time {
	return time.Now().Add(time.Minute * time.Duration(t.expirationTime))
}

func (t *Timesync) SetExpirationTime(timeInMinutes int) {
	t.expirationTime = timeInMinutes
}

func (t *Timesync) ExpirationTimestamp() float64 {
	return float64(t.GetExpirationDatetime().UnixMicro())
}

package wsapi

import (
	"time"
)

type Timesync struct {
	Name            string
	serverTimestamp float64
	expirationTime  int
}

func NewTimesync() *Timesync {
	return &Timesync{
		Name:            "timeSync",
		serverTimestamp: float64(time.Now().UnixMicro()),
		expirationTime:  1,
	}
}

func (t *Timesync) GetServerTimestamp() float64 {
	return t.serverTimestamp / 1000
}

func (t *Timesync) GetServerDatetime() time.Time {
	return time.Unix(int64(t.serverTimestamp), 0)
}

func (t *Timesync) SetServerTimestamp(timestamp float64) {
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

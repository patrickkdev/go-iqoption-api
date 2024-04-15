package types

type AuthenticationResponse struct {
	Name            string `json:"name"`
	Msg             bool   `json:"msg"`
	ClientSessionID string `json:"client_session_id"`
	RequestID       string `json:"request_id"`
}

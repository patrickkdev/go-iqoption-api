package data

import "fmt"

type Host struct {
	Host       string
	LoginURL   string
	HTTPAPIURL string
	WSAPIURL   string
}

func GetHostData(hostName string) *Host {
	urls := Host{
		Host:       hostName,
		LoginURL:   fmt.Sprintf("https://auth.%s.com/api/v2/login", hostName),
		HTTPAPIURL: fmt.Sprintf("https://%s/api", hostName),
		WSAPIURL:   fmt.Sprintf("hwss://%s/echo/websocket", hostName),
	}

	return &urls
}

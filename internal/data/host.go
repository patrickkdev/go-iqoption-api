package data

import "fmt"

type Host struct {
	Host       string
	LoginURL   string
	LogoutURL  string
	HTTPAPIURL string
	WSAPIURL   string
}

func GetHostData(hostName string) Host {
	urls := Host{
		Host:       hostName,
		LoginURL:   fmt.Sprintf("https://auth.%s/api/v2/login", hostName),
		LogoutURL:  fmt.Sprintf("https://auth.%s/api/v2/logout", hostName),
		HTTPAPIURL: fmt.Sprintf("https://%s/api", hostName),
		WSAPIURL:   fmt.Sprintf("wss://%s/echo/websocket", hostName),
	}

	return urls
}

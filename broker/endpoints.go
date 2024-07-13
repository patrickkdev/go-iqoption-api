package broker

import "fmt"

type brokerEndpoints struct {
	Domain     string
	LoginURL   string
	LogoutURL  string
	HTTPAPIURL string
	WSAPIURL   string
}

// Get broker endpoints by domain like 'iqoption.com'
func getEndpointsByBrokerDomain(brokerDomain string) brokerEndpoints {
	urls := brokerEndpoints{
		Domain:     brokerDomain,
		LoginURL:   fmt.Sprintf("https://auth.%s/api/v2/login", brokerDomain),
		LogoutURL:  fmt.Sprintf("https://auth.%s/api/v2/logout", brokerDomain),
		HTTPAPIURL: fmt.Sprintf("https://%s/api", brokerDomain),
		WSAPIURL:   fmt.Sprintf("wss://%s/echo/websocket", brokerDomain),
	}

	return urls
}

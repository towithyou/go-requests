# go-requests

## Installation
```
go get github.com/towithyou/go-requests
```

## Quick Start
```go
// reference requests_test.go
import (
	github.com/towithyou/go-requests
)

type RespData struct {
	Args    map[string]string `json:"args"`
	Data    string            `json:"data"`
	Headers struct {
		Accept                  string `json:"Accept"`
		AcceptEncoding          string `json:"Accept-Encoding"`
		AcceptLanguage          string `json:"Accept-Language"`
		Host                    string `json:"Host"`
		UpgradeInsecureRequests string `json:"Upgrade-Insecure-Requests"`
		UserAgent               string `json:"User-Agent"`
		XAmznTraceId            string `json:"X-Amzn-Trace-Id"`
		JWT                     string `json:"JWT"`
	} `json:"headers"`
	Origin string            `json:"origin"`
	Json   map[string]string `json:"json"`
	Url    string            `json:"url"`
}

func main() {
	var resp RespData
	_, err := requests.GetAndParse("http://httpbin.org/get", &resp)
	if err != nil {
		return
	}
	fmt.Printf("%+v", resp)
}
```
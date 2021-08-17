package requests

import (
	"io/ioutil"
	"net/http"
	"testing"
)

const (
	getUrl    = "http://httpbin.org/get"
	postUrl   = "http://httpbin.org/post"
	putUrl    = "http://httpbin.org/put"
	deleteUrl = "http://httpbin.org/delete"
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

func TestGetAndParse(t *testing.T) {
	var resp RespData
	_, err := GetAndParse(getUrl, &resp)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("%+v", resp)
}

func TestGetAndParse2(t *testing.T) {
	c := NewClient()
	c.AddQuery("search", "mysql").AddQuery("order", "id")
	c.AddHeader("JWT", "1qasdfsddf")

	var resp RespData
	if _, err := c.GetAndParseJson(getUrl, &resp); err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", resp)
}

func TestPost(t *testing.T) {
	//data := map[string]string{
	//	"name": "golang",
	//}

	data := struct {
		Name string
	}{
		Name: "golang",
	}

	var resp RespData
	_, err := PostAndParse(postUrl, &data, &resp)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("%+v", resp)
}

func TestPut(t *testing.T) {
	data := struct {
		Name string
	}{
		Name: "python",
	}

	var resp RespData
	_, err := PutAndParse(putUrl, &data, &resp)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("%+v", resp)
}

func TestDelete(t *testing.T) {
	c := NewClient()
	c.AddQuery("search", "mysql").AddQuery("action", "delete")
	c.AddHeader("JWT", "1qasdfsddf")

	var resp RespData
	_, err := c.DeleteAndParseJson(deleteUrl, &resp)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("%+v", resp)
}

func TestCustomClient(t *testing.T) {
	// custom client
	httpClient := &http.Client{}
	c := NewNativeClient(httpClient)

	resp, err := c.Get(getUrl)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(bytes))
}

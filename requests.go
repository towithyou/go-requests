package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type respCallBack func(resp *http.Response, err error)

type Method string

func (m Method) String() string {
	return string(m)
}

const (
	POST   Method = "POST"
	GET    Method = "GET"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
	PATCH  Method = "PATCH"
)

type headers map[string]string

func (h headers) Add(k, v string) {
	h[k] = v
}

type query map[string]string

func (q query) Add(k, v string) {
	q[k] = v
}

type Client struct {
	headers  headers
	urlQuery query
	client   *http.Client
}

var defaultClient = &Client{
	client: &http.Client{
		Timeout: 5 * time.Second,
	},
	headers: headers{
		"Content-Type": "application/json",
	},
	urlQuery: query{},
}

func NewClient() *Client {
	return &Client{
		client:   defaultClient.client,
		headers:  defaultClient.headers,
		urlQuery: defaultClient.urlQuery,
	}
}

func NewNativeClient(c *http.Client) *Client {
	return &Client{
		client:   c,
		headers:  defaultClient.headers,
		urlQuery: defaultClient.urlQuery,
	}
}

func (c *Client) parserQuery(u string) (fullUrl string, err error) {
	parseUrl, err := url.Parse(u)
	if err != nil {
		return
	}

	params, err := url.ParseQuery(parseUrl.RawQuery)
	if err != nil {
		return
	}

	for k, v := range c.urlQuery {
		params.Set(k, v)
	}

	parseUrl.RawQuery = params.Encode()
	fullUrl = parseUrl.String()

	return
}

func (c *Client) ClearQuery(u string) *Client {
	c.urlQuery = query{}
	return c
}

func (c *Client) AddQuery(key, val string) *Client {
	c.urlQuery.Add(key, val)
	return c
}

func (c *Client) AddHeader(key, val string) *Client {
	c.headers.Add(key, val)
	return c
}

func (c *Client) ClearHeader() *Client {
	c.headers = headers{}
	return c
}

func (c *Client) delete(u string) (resp *http.Response, err error) {
	req, err := c.NewRequest(DELETE, u, nil)
	if err != nil {
		return
	}
	return c.Do(req)
}

func (c *Client) Delete(u string) (resp *http.Response, err error) {
	return c.delete(u)
}

func (c *Client) DeleteAndParseJson(url string, v interface{}) (resp *http.Response, err error) {
	resp, err = c.delete(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	err = c.deserialization(resp, v)

	if err != nil {
		return
	}

	return
}

func (c *Client) DeleteAndParseJsonAsync(url string, v interface{}, callback respCallBack) {
	go callback(c.DeleteAndParseJson(url, v))
}

func (c *Client) Get(u string) (resp *http.Response, err error) {
	return c.get(u)
}

func (c *Client) get(u string) (resp *http.Response, err error) {
	req, err := c.NewRequest(GET, u, nil)
	if err != nil {
		return
	}
	return c.Do(req)
}

func (c *Client) GetAndParseJson(url string, v interface{}) (resp *http.Response, err error) {
	resp, err = c.get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	err = c.deserialization(resp, v)

	if err != nil {
		return
	}

	return
}

func (c *Client) GetAndParseJsonAsync(url string, v interface{}, callback respCallBack) {
	go callback(c.GetAndParseJson(url, v))
}

func (c *Client) PostAndParseJson(url string, reqData interface{}, v interface{}) (resp *http.Response, err error) {
	resp, err = c.post(url, reqData)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	err = c.deserialization(resp, v)

	if err != nil {
		return
	}

	return
}

func (c *Client) PostAndParseJsonAsync(url string, reqData interface{}, v interface{}, callback respCallBack) {
	go callback(c.PostAndParseJson(url, reqData, v))
}

func (c *Client) Post(u string, reqData interface{}) (resp *http.Response, err error) {
	return c.post(u, reqData)
}

func (c *Client) post(u string, reqData interface{}) (resp *http.Response, err error) {
	return c.withBody(POST, u, reqData)
}

func (c *Client) Put(u string, reqData interface{}) (resp *http.Response, err error) {
	return c.put(u, reqData)
}

func (c *Client) put(u string, reqData interface{}) (resp *http.Response, err error) {
	return c.withBody(PUT, u, reqData)
}

func (c *Client) PutAndParseJson(url string, reqData interface{}, v interface{}) (resp *http.Response, err error) {
	resp, err = c.put(url, reqData)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	err = c.deserialization(resp, v)

	if err != nil {
		return
	}

	return
}

func (c *Client) PutAndParseJsonAsync(url string, reqData interface{}, v interface{}, callback respCallBack) {
	go callback(c.PutAndParseJson(url, reqData, v))
}

func (c *Client) patch(u string, reqData interface{}) (resp *http.Response, err error) {
	return c.withBody(PATCH, u, reqData)
}

func (c *Client) Patch(u string, reqData interface{}) (resp *http.Response, err error) {
	return c.patch(u, reqData)
}

func (c *Client) PatchAndParseJson(url string, reqData interface{}, v interface{}) (resp *http.Response, err error) {
	resp, err = c.patch(url, reqData)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	err = c.deserialization(resp, v)

	if err != nil {
		return
	}

	return
}

func (c *Client) PatchAndParseJsonAsync(url string, reqData interface{}, v interface{}, callback respCallBack) {
	go callback(c.PatchAndParseJson(url, reqData, v))
}

func (c *Client) deserialization(resp *http.Response, v interface{}) error {
	body, err := c.read(resp)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) withBody(m Method, u string, reqData interface{}) (resp *http.Response, err error) {
	var r []byte
	if d, ok := reqData.([]byte); ok {
		r = d
	} else {
		if r, err = json.Marshal(reqData); err != nil {
			return
		}
	}

	req, err := c.NewRequest(m, u, bytes.NewReader(r))
	if err != nil {
		return
	}
	return c.Do(req)
}

func (c *Client) NewRequest(method Method, url string, body io.Reader) (req *http.Request, err error) {
	fullUrl, err := c.parserQuery(url)
	if err != nil {
		return
	}
	return http.NewRequest(method.String(), fullUrl, body)
}

func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
	for k, v := range c.headers {
		req.Header.Add(k, v)
	}
	return c.client.Do(req)
}

func (c *Client) read(resp *http.Response) (b []byte, err error) {
	return ioutil.ReadAll(resp.Body)
}

func (c *Client) DoAsync(req *http.Request, back respCallBack) {
	go back(c.Do(req))
}

func Get(u string) (*http.Response, error) {
	return defaultClient.get(u)
}
func GetAndParse(url string, v interface{}) (*http.Response, error) {
	return defaultClient.GetAndParseJson(url, v)
}

func Post(u string, data interface{}) (*http.Response, error) {
	return defaultClient.post(u, data)
}

func PostAndParse(url string, data, val interface{}) (*http.Response, error) {
	return defaultClient.PostAndParseJson(url, data, val)
}

func Put(u string, data interface{}) (*http.Response, error) {
	return defaultClient.put(u, data)
}

func PutAndParse(url string, data, val interface{}) (*http.Response, error) {
	return defaultClient.PutAndParseJson(url, data, val)
}

func Patch(u string, data interface{}) (*http.Response, error) {
	return defaultClient.patch(u, data)
}

func PatchAndParse(url string, data, val interface{}) (*http.Response, error) {
	return defaultClient.PatchAndParseJson(url, data, val)
}

func Delete(u string) (*http.Response, error) {
	return defaultClient.delete(u)
}

func DeleteAndParse(u string, v interface{}) (*http.Response, error) {
	return defaultClient.DeleteAndParseJson(u, v)
}

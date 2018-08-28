package jsonrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/magento-mcom/go-common/rpc"
	"github.com/magento-mcom/go-common/rpc/errors"
)

// NewClient constructs new JSONRPC client.
func NewClient(url string, transport rpc.Transport, mw ...rpc.TransportMiddleware) *Client {
	if transport == nil {
		transport = http.DefaultClient
	}

	for _, m := range mw {
		transport = m(transport)
	}

	return &Client{
		url:       url,
		transport: transport,
	}
}

// Client implements JSONRPC client.
type Client struct {
	url       string
	transport rpc.Transport
}

// Call RPC method
func (c *Client) Call(req rpc.Request) (rpc.Data, error) {
	params := req.Raw()
	if params == nil {
		params = struct{}{}
	}

	bodyreq, err := json.Marshal(requestDataEncode{
		Version: version,
		ID:      "1",
		Method:  req.Method,
		Params:  params,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to compose JSONRPC request: %v", err)
	}

	httpreq, err := http.NewRequest(http.MethodPost, c.url, bytes.NewReader(bodyreq))
	if err != nil {
		return nil, fmt.Errorf("failed to build HTTP request: %v", err)
	}

	httpreq = httpreq.WithContext(req.Context)
	httpreq.Header.Set("Content-Type", "application/json")
	httpreq.Header.Set("Connection", "close")
	httpreq.Close = true

	httpresp, err := c.transport.Do(httpreq)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %v", err)
	}

	defer httpresp.Body.Close()

	if int(httpresp.StatusCode/100) != 2 {
		return nil, errors.NewHTTPResponseError(
			fmt.Sprintf("failed to make HTTP request: response status code %v", httpresp.StatusCode),
			httpresp.StatusCode,
		)
	}

	bodyresp, err := ioutil.ReadAll(httpresp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP response: %v", err)
	}

	resp := &responseData{}

	if err := json.Unmarshal(bodyresp, resp); err != nil {
		return nil, errors.NewDecodingResponseError(
			fmt.Sprintf("server response is not a valid JSONRPC message: %v", err),
			bodyresp,
		)
	}

	if resp.Error != nil {
		return nil, errors.New(resp.Error.Code, resp.Error.Message, resp.Error.Data)
	}

	return rpc.NewJSONData(resp.Result), nil
}

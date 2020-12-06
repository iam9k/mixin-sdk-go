package mixin

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	Endpoint = "https://mixin-api.zeromesh.net"
)

var httpClient = resty.New().
	SetHeader("Accept", "application/json").
	SetHostURL(Endpoint).
	SetTimeout(300 * time.Millisecond)

func Request(ctx context.Context) *resty.Request {
	return httpClient.R().SetContext(ctx)
}

func DecodeResponse(resp *resty.Response) ([]byte, error) {
	var body struct {
		Error
		Data json.RawMessage `json:"data,omitempty"`
	}

	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return nil, err
	}

	if body.Error.Code > 0 {
		return nil, &body.Error
	}
	return body.Data, nil
}

func UnmarshalResponse(resp *resty.Response, v interface{}) error {
	data, err := DecodeResponse(resp)
	if err != nil {
		return err
	}

	if v != nil {
		return json.Unmarshal(data, v)
	}

	return nil
}

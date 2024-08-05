package be

import (
	"encoding/json"
	"net/http"
)

type RespCode int

type Resp struct {
	http.ResponseWriter `json:"-"`

	Code RespCode        `json:"code"`
	Data json.RawMessage `json:"data"`
}

// Sends an HTTP response with the given code and no data; see status defaults
// in `resp.SendData`.
func (resp *Resp) Send(code RespCode) error {
	return resp.SendData(code, nil)
}

// Sends an HTTP response with the given code and data; the status defaults to
// 400 for error codes, 200 for success codes, and the code itself if it's
// equal to an HTTP status.
func (resp *Resp) SendData(code RespCode, data any) error {
	if code >= 2000 {
		return resp.SendStatus(code, data, http.StatusBadRequest)
	} else if code >= 1000 {
		return resp.SendStatus(code, data, http.StatusOK)
	} else {
		return resp.SendStatus(code, data, int(code))
	}
}

// Sends an HTTP response with the given app code, data and status.
func (resp *Resp) SendStatus(
	code RespCode,
	data any,
	httpStatus int,
) error {
	rawData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	resp.Code = RespCode(code)
	resp.Data = rawData
	rawResp, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return err
	}

	resp.WriteHeader(httpStatus)
	resp.Write(rawResp)

	return nil
}

package gosnowflake

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/mailru/easyjson"
)

const (
	// limit http response to be 100MB to avoid overwhelming the scheduler
	ResponseBodyLimit = 100 * 1024 * 1024
)

var ErrResponseTooLarge = fmt.Errorf("response is too large")

func decodeResponse(body io.ReadCloser, resp interface{}) error {
	lr := io.LimitReader(body, ResponseBodyLimit)
	var err error
	if v, is := resp.(easyjson.Unmarshaler); is {
		err = easyjson.UnmarshalFromReader(lr, v)
	} else {
		err = json.NewDecoder(lr).Decode(resp)
	}
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return ErrResponseTooLarge
	}
	return err
}

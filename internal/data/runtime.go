package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvaldRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	unsupportedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvaldRuntimeFormat
	}

	parts := strings.Split(unsupportedJSONValue, " ")
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvaldRuntimeFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvaldRuntimeFormat
	}

	*r = Runtime(i)

	return nil
}

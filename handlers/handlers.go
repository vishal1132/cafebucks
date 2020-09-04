package handlers

import (
	"context"
	"fmt"

	"github.com/valyala/fastjson"
)

type ctxKey uint8

const (
	ctxKeyReqID ctxKey = iota
)

// CtxRequestID extracts the request id if set
func CtxRequestID(ctx context.Context) (string, bool) {
	v := ctx.Value(ctxKeyReqID)
	if v == nil {
		return "", false
	}

	rid, ok := v.(string)
	if !ok {
		return "", false
	}

	return rid, true
}

// GetJSONString is wisott
func GetJSONString(document *fastjson.Value, key string) (string, error) {
	if !document.Exists(key) {
		return "", fmt.Errorf("failed to get field %s: key does not exist", key)
	}

	v, err := document.Get(key).StringBytes()
	if err != nil {
		return "", fmt.Errorf("failed to get field %s: %w", key, err)
	}

	s := make([]byte, len(v))

	copy(s, v)

	return string(s), nil
}

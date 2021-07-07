package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// APIError 微信支付 API v3 标准错误结构
type APIError struct {
	StatusCode int         // 应答报文的 HTTP 状态码
	Header     http.Header // 应答报文的 Header 信息
	Body       string      // 应答报文的 Body 原文
	Code       string      `json:"code"`             // 应答报文的 Body 解析后的错误码信息，仅不符合预期/发生系统错误时存在
	Message    string      `json:"message"`          // 应答报文的 Body 解析后的文字说明信息，仅不符合预期/发生系统错误时存在
	Detail     interface{} `json:"detail,omitempty"` // 应答报文的 Body 解析后的详细信息，仅不符合预期/发生系统错误时存在
}

func (e *APIError) Error() string {
	var buf bytes.Buffer
	_, _ = fmt.Fprintf(&buf, "error http response:[StatusCode: %d Code: \"%s\"", e.StatusCode, e.Code)
	if e.Message != "" {
		_, _ = fmt.Fprintf(&buf, "\nMessage: %s", e.Message)
	}
	if e.Detail != nil {
		var detailBuf bytes.Buffer
		enc := json.NewEncoder(&detailBuf)
		enc.SetIndent("", "  ")
		if err := enc.Encode(e.Detail); err == nil {
			_, _ = fmt.Fprint(&buf, "\nDetail:")
			_, _ = fmt.Fprintf(&buf, "\n%s", strings.TrimSpace(detailBuf.String()))
		}
	}
	if len(e.Header) > 0 {
		_, _ = fmt.Fprint(&buf, "\nHeader:")
		for key, value := range e.Header {
			_, _ = fmt.Fprintf(&buf, "\n - %v=%v", key, value)
		}
	}
	_, _ = fmt.Fprintf(&buf, "]")
	return buf.String()
}

// ParameterError 请求参数错误
type ParameterError struct {
	FieldName string
	Message   string
}

func (e *ParameterError) Error() string {
	return fmt.Sprintf("`%s` %s", e.FieldName, e.Message)
}

package service

import (
	"testing"
)

func TestWxPay(t *testing.T) {
	ws := NewWxPayService()

	//ws.PrePay()

	ws.PrePay1(100, "12345567")
}

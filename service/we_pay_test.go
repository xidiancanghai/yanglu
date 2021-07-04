package service

import (
	"fmt"
	"testing"
)

func TestWxPay(t *testing.T) {
	ws, err := NewWxPayService()
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	ws.PrePay()
}

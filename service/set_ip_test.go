package service

import (
	"fmt"
	"testing"
)

func TestSetIp(t *testing.T) {
	ss := NewSetIpService("172.0.0.1")
	err := ss.SetIp()
	fmt.Println("err = ", err)
}

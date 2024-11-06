package models

import "net"

type ConnInfo struct {
	IsTls  bool
	IPType uint8
	Port   int
	IP     net.IP
}

package tcpx

import (
	"unsafe"
)

// 根据id从connections获取connection
func Get(id int64) *Connection {
	//注: id实际为connection的指针地址
	return (*Connection)(unsafe.Pointer(uintptr(id)))
}

func Clear() {
	for i := len(readers) - 1; i > -1; i-- {
		close(readers[i].connChan)
		conns := readers[i].conns
		for conn, _ := range conns {
			conn.Stop()
		}
	}
}
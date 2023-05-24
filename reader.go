package tcpx

import (
	"io"
	"os"
	"time"
)

type Reader struct {
	connChan chan *Connection
	conns map[*Connection]struct{}
}

func (r *Reader) Serve() {
	header := make([]byte, 3)
	for {

		// 读取刚创建完毕的连接
		for i := len(r.connChan); i > -1; i-- {
			select {
			case conn := <- r.connChan:
				r.conns[conn] = struct{}{}
				preHook(conn)
			case <- time.After(Server.Timeout):
			}
		}

		for conn, _ := range r.conns {
			conn.Conn.SetReadDeadline(time.Now().Add(Server.Timeout))
			_, err := io.ReadFull(conn.TcpConn(), header)
			if err != nil {
				if !os.IsTimeout(err) {
					// 读取连接数据失败, 关闭连接资源
					delete(r.conns, conn)
					conn.Stop()
				}
				continue
			}
	
			// 处理附带请求头的消息
			conn.process(header)
		}
	}
}
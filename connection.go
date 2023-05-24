package tcpx

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
)

/*
*
连接模块
*/
type Connection struct {
	//当前连接的socket tcp套接字
	Conn *net.TCPConn
	//连接状态
	isClosed bool
	//防止连接并发写的锁
	writeLock *sync.Mutex
	ExitChan chan struct{}
	//连接属性
	Property property
}

// 连接模块的初始化方法
func NewConnection(conn *net.TCPConn) *Connection {
	return &Connection{
		Conn:     conn,
		isClosed: false,
		writeLock: new(sync.Mutex),
		ExitChan: make(chan struct{}, 1),
		Property: property{
			lock:     &sync.RWMutex{},
			property: make(map[string]interface{}),
		},
	}
}

// 停止连接，让连接停止工作
func (c *Connection) Stop() {
	Info.Println("conn stop, conn id:", c.ID())

	//如果当前连接已关闭，直接退出
	if c.isClosed {
		return
	}

	//调用connection连接销毁之前的hook方法
	if postHook != nil {
		postHook(c)
	}

	//关闭tcp conn连接
	err := c.Conn.Close()
	//通知writer关闭
	c.ExitChan <- struct{}{}
	//从reader中删除 c
	index := c.ID()%int64(len(readers))
	delete(readers[index].conns, c)
	
	close(c.ExitChan)
	c.isClosed = true
	Info.Println("conn close err:", err)
}

// 获取当前连接绑定的 socket conn
func (c *Connection) TcpConn() *net.TCPConn {
	return c.Conn
}

// 获取当前连接模块的id
func (c *Connection) ID() int64 {
	id, err := strconv.ParseInt(fmt.Sprint(&c)[2:], 16, 41)
	if err != nil {
		panic(err)
	}
	return id
}

// 获取远程客户端的 tcp状态 addr
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据给远程客户端
func (c *Connection) Send(msgID uint8, data []byte) error {

	if c.isClosed {
		return errors.New("connection is closed when send message")
	}

	//封装message消息
	message := NewMessage(msgID, data)

	//封包
	response := Pack(message)

	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	_, err := c.Conn.Write(response)
	return err
}

// 处理携带
func (c *Connection) process(header []byte) {
	
	//从header中解析出 datalen 和 id
	msg := UnPack(header)

	//根据msgLen拿到msgData
	_, err := io.ReadFull(c.TcpConn(), msg.Data)
	if err != nil {
		c.Stop()
		Info.Println("read data err:", err)
		return
	}

	//封装Request
	request := &Request{
		conn: c,
		msg:  msg,
	}

	//将当前请求加入请求队列
	DistributeReqest(request)
}
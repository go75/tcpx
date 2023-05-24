package tcpx

import (
	"errors"
	"net"
	"strconv"
	"time"
)

type Engine struct {
	//服务器的服务地址
	Address *net.TCPAddr
	//连接读取消息的超时时间
	Timeout time.Duration
}

// 停止服务器
func (e *Engine) Stop() {
	//将服务器的资源，连接等资源进行销毁或回收
	Clear()
}

// 运行服务器
func (e *Engine) Run() error {
	//监听服务地址
	listener, err := net.ListenTCP(e.Address.Network(), e.Address)
	if err != nil {
		return errors.New("listen tcp err: " + err.Error())
	}

	Info.Println("Engine start success, listenner at addr :" + e.Address.String())
	
	//阻塞等待客户端连接，处理客户端连接业务(读,写)
	for {

		//监听客户端的连接
		conn, err := listener.AcceptTCP()
		if err != nil {
			Info.Println("listener accept err:", err)
			continue
		}

		Info.Println("accept tcp conn,romote addr:", conn.RemoteAddr())

		//封装tcp conn得到connection模块
		connection := NewConnection(conn)

		readers[(connection.ID() % int64(len(readers)))].connChan <- connection
	}
}

// 注册路由
func (s *Engine) Regist(id uint8, handle func(req *Request)) {
	//判断当前msg id是否已经注册过router
	if apis[id] != nil {
		panic("repeat api,msg id:" + strconv.Itoa(int(id)))
	}

	//将处理函数注册到apis中
	apis[id] = handle
}

// 初始化自定义Engine的方法
func New(ipVersion, addr string, dispatcherNum uint16, dispatcherLen uint16, timeout time.Duration, prehook, posthook func(c *Connection), notRegistFn func(r *Request)) (*Engine, error) {
	//获取tcp的addr
	address, err := net.ResolveTCPAddr(ipVersion, addr)
	if err != nil {
		return nil, err
	}

	if prehook == nil {

	}

	Server = &Engine{
		Address:          address,
		Timeout:		  timeout,
	}

	if prehook == nil {
		preHook = func(c *Connection) {}
	} else {
		preHook = prehook
	}

	if prehook == nil {
		postHook = func(c *Connection) {}
	} else {
		postHook = posthook
	}

	if notRegistFn == nil {
		notRegistFn = func(c *Request) {}
	} else {
		notRegistFn = notRegistFn
	}

	//初始化reader池, 请求队列, 调度器
	StartDispatchersAndRequest(dispatcherNum, dispatcherLen)
	return Server, nil
}

// 初始化默认Engine的方法
func Defalut(addr string) (*Engine, error) {
	return New("tcp4", addr, 16, 1024, time.Microsecond << 4, nil, nil, nil)
}
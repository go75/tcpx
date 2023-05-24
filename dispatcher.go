package tcpx

func StartDispatchersAndRequest(dispatcherSize uint16, dispatcherLen uint16) {
	Info.Println("开启请求处理模块")
	readers = make([]Reader, dispatcherSize)
	var i uint16 = 0
	for i < dispatcherSize {
		readers[i] = Reader{
			connChan: make(chan *Connection, dispatcherLen),
			conns: make(map[*Connection]struct{}),
		}
		go readers[i].Serve()
		i++
	}
	//初始化请求队列
	requestQueue = make([]chan *Request, dispatcherSize)
	//根据请求队列的数量开启相应数量的dispatcher
	for i := int(dispatcherSize) - 1; i > -1; i-- {
		//初始化消息队列的对应的dispatcher
		requestQueue[i] = make(chan *Request, dispatcherLen)
		//启动该dispatcher
		go startDispatcher(i)
	}
}

func startDispatcher(id int) {
	Info.Println("dispatcher id=", id, "is started...")
	
	//不断阻塞等待请求队列的请求
	for {
		select {
		//如果接收到请求，就处理该请求
		case req := <-requestQueue[id]:
			DoRequest(req)
		}
	}
}

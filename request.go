package tcpx

// 将连接和数据绑定成一个request
type Request struct {
	//当前连接
	conn *Connection
	//请求的数据
	msg *Message
}

// 从request中得到连接
func (r *Request) Conn() *Connection {
	return r.conn
}

// 从request中得到数据
func (r *Request) Data() []byte {
	return r.msg.Data
}

// 从request中得到消息id
func (r *Request) MsgID() uint8 {
	return r.msg.ID
}

// 请求调度处理
func DoRequest(req *Request) {
	//根据消息id获取绑定的handler方法
	handle := apis[req.MsgID()]
	if handle == nil {
		notRegistFn(req)
		return
	}
	handle(req)
}

// 分配请求
func DistributeReqest(req *Request) {

	//根据连接的id和请求队列的数量取模得到队列号，将请求发送给该请求号对应的队列
	requestQueue[req.conn.ID()%int64(len(requestQueue))] <- req

}

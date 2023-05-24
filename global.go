package tcpx

// 当前服务器对象
var Server *Engine

// 当前消息队列
var requestQueue []chan *Request

//k:msgid , v:*router,将所有消息id和路由方法绑定
var apis [256]func(*Request)


// 所有读协程
var readers []Reader

//连接创建之后的钩子方法
var preHook func(c *Connection)
//连接销毁之前的钩子方法
var postHook func(c *Connection)
//未注册处理函数时的钩子函数
var notRegistFn func(req *Request)
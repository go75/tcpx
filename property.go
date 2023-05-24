package tcpx

import "sync"

// 属性
type property struct {
	//连接属性的锁
	lock *sync.RWMutex
	//连接属性集合
	property map[string]interface{}
}

// 设置属性
func (p *property) Set(key string, value interface{}) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.property[key] = value
}

// 根据key得到属性中相应的value
func (p *property) Get(key string) (interface{}, bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	value, ok := p.property[key]
	return value, ok
}

// 删除属性
func (p *property) Delete(key string) {
	p.lock.Lock()
	p.lock.Unlock()
	delete(p.property, key)
}

package onceFlight

import "sync"

type (
	call struct {
		wg  sync.WaitGroup
		val interface{}
		err error
	}
	onceGroup struct {
		mutex sync.Mutex
		calls map[string]*call
	}

	OnceFlight interface {
		Do(key string, fn func() (interface{}, error)) (interface{}, error)
	}
)

func init() {

}

func (g *onceGroup) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	c, ok := g.createCall(key)
	if ok {
		return c.val, nil
	}

	g.makeCall(c, key, fn)
	return c.val, c.err
}

func (g *onceGroup) createCall(key string) (*call, bool) {
	g.mutex.Lock()
	// 其他相同请求
	if c, ok := g.calls[key]; ok {
		g.mutex.Unlock()
		c.wg.Wait()
		return c, true
	}
	// 第一个请求
	c := new(call)
	c.wg.Add(1)
	g.calls[key] = c
	g.mutex.Unlock()
	return c, false
}

// makeCall .
func (g *onceGroup) makeCall(c *call, key string, fn func() (interface{}, error)) {
	defer func() {
		g.mutex.Lock()
		delete(g.calls, key)
		g.mutex.Unlock()
		c.wg.Done()
	}()
	c.val, c.err = fn()
}

func NewOnceFlight() OnceFlight {
	return &onceGroup{
		calls: make(map[string]*call),
	}
}

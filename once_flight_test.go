package onceFlight_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/syncx"

	"github.com/yino/onceFlight"
)

func TestNewOnceFlight(t *testing.T) {
	round := 10
	var wg sync.WaitGroup
	barrier := onceFlight.NewOnceFlight()
	wg.Add(round)
	for i := 0; i < round; i++ {
		go func() {
			defer wg.Done()
			// 启用10个协程模拟获取缓存操作
			val, err := barrier.Do("get_rand_int", func() (interface{}, error) {
				time.Sleep(time.Second)
				return rand.Int(), nil
			})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(val)
			}
		}()
	}
	wg.Wait()
}

func TestZeroFlight(t *testing.T) {
	round := 10
	var wg sync.WaitGroup
	barrier := syncx.NewSingleFlight()
	wg.Add(round)
	for i := 0; i < round; i++ {
		go func() {
			defer wg.Done()
			// 启用10个协程模拟获取缓存操作
			val, err := barrier.Do("get_rand_int", func() (interface{}, error) {
				time.Sleep(time.Second)
				return rand.Int(), nil
			})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(val)
			}
		}()
	}
	wg.Wait()
}

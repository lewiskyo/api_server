package redis

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	for i := 0; i < 1000; i++ {
		go func() {
			ret := New()
			if ret == nil {
				t.Errorf("new error!")
				return
			}
			ret.Get("a")

			ret1 := New("opposite")
			if ret1 == nil {
				t.Errorf("new error!")
				return
			}
			ret1.Get("a")
		}()
	}
	time.Sleep(1 * time.Second)
	for k, data := range redisConnections {
		t.Logf("k=%#v, v=%#v", k, data)
	}
	time.Sleep(1 * time.Second)
}

func TestRedis_Get(t *testing.T) {
	c := New()
	r, e := c.Get("a").Result()
	if e != nil {
		t.Fatalf("redisCluster get error:%v", e)
	}
	t.Logf("r:%v", r)
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		go func() {
			ret := New()
			if ret == nil {
				b.Error("new error!")
				return
			}
			ret.Get("a").Val()

			ret1 := New("opposite")
			if ret1 == nil {
				b.Error("new error!")
				return
			}
			ret1.Get("a").Val()
		}()
	}
}

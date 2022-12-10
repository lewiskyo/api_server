package rediscluster

import (
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestNew(t *testing.T) {
	for i := 0; i < 100; i++ {
		go func() {
			ret := New()
			if ret == nil {
				t.Errorf("new error!")
				return
			}
			t.Logf("%#v", ret)
			r, err := ret.Get("a").Result()
			if err != nil && err != redis.Nil {
				t.Errorf("Get error:%v", err)
				return
			}
			t.Logf("a:%#v", r)

			ret1 := New("opposite")
			if ret1 == nil {
				t.Errorf("new error!")
				return
			}
			t.Logf("%#v", ret1)
			r1, err := ret1.Get("a").Result()
			if err != nil && err != redis.Nil {
				t.Errorf("Get opposite error:%v", err)
				return
			}
			t.Logf("a:%#v", r1)
		}()
	}
	time.Sleep(6 * time.Second)
	for k, data := range redisConnections {
		t.Logf("k=%#v, v=%#v", k, data)
	}
	time.Sleep(1 * time.Second)
}

func TestRedisCluster_Get(t *testing.T) {
	c := New()
	r, e := c.Get("a").Result()
	if e != nil && e != redis.Nil {
		t.Fatalf("redisCluster get error:%v", e)
	}
	t.Logf("r:%v", r)
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ret := New()
		ret1 := New("opposite")
		ret.Set("a", 1, 0)
		if ret == nil {
			b.Fatal("new error!")
		}
		b.Logf("%v", ret.Get("a").Val())
		b.Logf("ret:%#v", ret)
		b.Logf("ret1:%#v", ret1)

		if ret1 == nil {
			b.Fatal("new error!")
		}
		b.Logf("%v", ret.Get("a").Val())
	}
}

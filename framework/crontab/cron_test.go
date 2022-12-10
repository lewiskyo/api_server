package crontab

import (
	"testing"
	"time"
)

func TestAddFunc(t *testing.T) {
	i, e := AddFunc("*/2 * * * * *", func() {
		t.Logf("I'm 2 seconds.time=%v", time.Now().Unix())
	})
	if e != nil {
		t.Fatalf("AddFunc Error:%v", e)
	}
	t.Logf("ID:%v", i)
	i, e = AddFunc("*/3 * * * * *", func() {
		t.Logf("I'm 3 seconds.time=%v", time.Now().Unix())
	})
	if e != nil {
		t.Fatalf("AddFunc Error:%v", e)
	}
	t.Logf("ID:%v", i)

	i, e = AddFunc("0 */1 * * * *", func() {
		t.Logf("I'm 1 min.time=%v", time.Now().Unix())
	})
	if e != nil {
		t.Fatalf("AddFunc Error:%v", e)
	}
	t.Logf("ID:%v", i)
	StartCronSchedule()
	time.Sleep(30 * time.Second)
}

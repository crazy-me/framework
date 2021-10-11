package redis

import (
	"github.com/crazy-me/framework/cache"
	"strconv"
	"testing"
)

func TestRedis(t *testing.T) {
	bm, err := cache.New("redis", `{"address":":6379","password":"","db":"0","maxIdle":"3","timeout":""}`)
	if err != nil {
		t.Error(err)
	}

	if b := bm.Set("name", "redis"); !b {
		t.Error("set error")
	}

	for i := 1; i <= 1000; i++ {
		s := strconv.Itoa(i)
		if b := bm.Set("test_"+strconv.Itoa(i), "set_test_"+s); !b {
			t.Error("set error")
		}
	}

	if b := bm.SetEx("age", 10, 100); !b {
		t.Error("SetEx error")
	}

	if b := bm.Incr("value"); !b {
		t.Error("Incr error")
	}

	if b := bm.Decr("value"); !b {
		t.Error("Decr error")
	}

	if err = bm.ClearAll(); err != nil {
		t.Error(err)
	}

}

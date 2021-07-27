package rdb

import (
	"errors"
	"strconv"

	"github.com/go-redis/redis/v7"
)

type Z struct {
	Member interface{}
	Score  int64
}

type ZSliceCmd struct {
	err error

	val []Z
}

func (cmd *ZSliceCmd) Val() []Z {
	return cmd.val
}

func (cmd *ZSliceCmd) Result() ([]Z, error) {
	return cmd.val, cmd.err
}

func ZRangeWithScores(r redis.UniversalClient, key string, start, stop int64) *ZSliceCmd {
	cmd := r.Do("zrange", key, start, stop, "withscores")
	if err := cmd.Err(); err != nil {
		return &ZSliceCmd{
			err: err,
			val: make([]Z, 0),
		}
	}

	val, ok := cmd.Val().([]interface{})
	if !ok {
		return &ZSliceCmd{
			err: errors.New("unknown exception"),
			val: make([]Z, 0),
		}
	}

	rs := make([]Z, 0)

	var prev interface{}
	for i, v := range val {
		if i > 0 && i%2 == 1 {
			member := prev.(string)
			score := readZScore(v)
			rs = append(rs, Z{
				Member: member,
				Score:  score,
			})
			continue
		}
		prev = v
	}
	return &ZSliceCmd{
		err: nil,
		val: rs,
	}
}

func readZScore(v interface{}) int64 {
	switch v.(type) {
	// 创梦修改过redis server，返回的score类型是int64
	// 标准redis server，返回的score类型是string
	case int64:
		return v.(int64)
	case string:
		s := v.(string)
		r, e := strconv.ParseInt(s, 10, 64)
		if e != nil {
			return 0
		}
		return r
	default:
		return 0
	}
}

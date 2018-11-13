package db

import (
	"github.com/go-redis/redis"
)

var CfRedis []*redis.Client

func InitRedis(flag bool){
	if len(CfRedis)>0 {
		CfRedis=append(CfRedis[:0],CfRedis[len(CfRedis):]...)
	}
	var s []*redis.Client
	if flag {
		luck := redis.NewClient(
			&redis.Options{
				Addr:     "119.23.219.245:8000",
				Password: "", // no password set
				DB:       4,                  // use default DB)
			})
		sass := redis.NewClient(
			&redis.Options{
				Addr:     "119.23.219.245:8000",
				Password: "", // no password set
				DB:      1,                 // use default DB)
			})
		claw := redis.NewClient(
			&redis.Options{
				Addr:     "119.23.219.245:8000",
				Password: "", // no password set
				DB:       2,                 // use default DB)
			})
		s = append(s, luck)
		s = append(s, sass)
		s = append(s, claw)
	}else {
		luck := redis.NewClient(
			&redis.Options{
				Addr:     "119.23.219.245:8000",
				Password: "", // no password set
				DB:    13   ,                  // use default DB)
			})
		sass := redis.NewClient(
			&redis.Options{
				Addr:     "119.23.219.245:8000",
				Password: "", // no password set
				DB:       14,                 // use default DB)
			})
		claw := redis.NewClient(
			&redis.Options{
				Addr:     "119.23.219.245:8000",
				Password: "", // no password set
				DB:       10,                 // use default DB)
			})
		s = append(s, luck)
		s = append(s, sass)
		s = append(s, claw)
	}
	CfRedis =s
}


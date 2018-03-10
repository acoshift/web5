package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/acoshift/configfile"
	"github.com/acoshift/hime"
	"github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"

	"github.com/acoshift/web5/app"
)

func main() {
	config := configfile.NewReader("config")

	db, err := sql.Open("postgres", config.String("sql"))
	if err != nil {
		log.Fatal(err)
	}

	redisHost := config.String("redis_host")
	redisPool := redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisHost)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	err = hime.New().
		Handler(app.New(app.Config{
			DB:          db,
			RedisPool:   &redisPool,
			RedisPrefix: config.String("redis_prefix"),
		})).
		GracefulShutdown().
		ListenAndServe(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

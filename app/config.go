package app

import (
	"database/sql"

	"github.com/garyburd/redigo/redis"
)

// Config is app's config
type Config struct {
	DB          *sql.DB
	RedisPool   *redis.Pool
	RedisPrefix string
}

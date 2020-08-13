package interfaces

//go:generate moq -out mock/mock_redisclienter.go -pkg mock . RedisClienter
//go:generate moq -out mock/mock_resulter.go -pkg mock . Resulter

import (
	"time"
)

// Resulter - interface for redis.StatusCMD
type Resulter interface {
	Result() (msg string, err error)
}

// RedisClienter - interface for redis
type RedisClienter interface {
	Set(string, interface{}, time.Duration) Resulter
	Ping() Resulter
}

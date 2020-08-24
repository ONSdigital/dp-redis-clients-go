// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"github.com/ONSdigital/dp-redis/interfaces"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

var (
	lockRedisClienterMockPing sync.RWMutex
	lockRedisClienterMockSet  sync.RWMutex
)

// Ensure, that RedisClienterMock does implement RedisClienter.
// If this is not the case, regenerate this file with moq.
var _ interfaces.RedisClienter = &RedisClienterMock{}

// RedisClienterMock is a mock implementation of interfaces.RedisClienter.
//
//     func TestSomethingThatUsesRedisClienter(t *testing.T) {
//
//         // make and configure a mocked interfaces.RedisClienter
//         mockedRedisClienter := &RedisClienterMock{
//             PingFunc: func() *redis.StatusCmd {
// 	               panic("mock out the Ping method")
//             },
//             SetFunc: func(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
// 	               panic("mock out the Set method")
//             },
//         }
//
//         // use mockedRedisClienter in code that requires interfaces.RedisClienter
//         // and then make assertions.
//
//     }
type RedisClienterMock struct {
	// PingFunc mocks the Ping method.
	PingFunc func() *redis.StatusCmd

	// SetFunc mocks the Set method.
	SetFunc func(key string, value interface{}, expiration time.Duration) *redis.StatusCmd

	// calls tracks calls to the methods.
	calls struct {
		// Ping holds details about calls to the Ping method.
		Ping []struct {
		}
		// Set holds details about calls to the Set method.
		Set []struct {
			// Key is the key argument value.
			Key string
			// Value is the value argument value.
			Value interface{}
			// Expiration is the expiration argument value.
			Expiration time.Duration
		}
	}
}

// Ping calls PingFunc.
func (mock *RedisClienterMock) Ping() *redis.StatusCmd {
	if mock.PingFunc == nil {
		panic("RedisClienterMock.PingFunc: method is nil but RedisClienter.Ping was just called")
	}
	callInfo := struct {
	}{}
	lockRedisClienterMockPing.Lock()
	mock.calls.Ping = append(mock.calls.Ping, callInfo)
	lockRedisClienterMockPing.Unlock()
	return mock.PingFunc()
}

// PingCalls gets all the calls that were made to Ping.
// Check the length with:
//     len(mockedRedisClienter.PingCalls())
func (mock *RedisClienterMock) PingCalls() []struct {
} {
	var calls []struct {
	}
	lockRedisClienterMockPing.RLock()
	calls = mock.calls.Ping
	lockRedisClienterMockPing.RUnlock()
	return calls
}

// Set calls SetFunc.
func (mock *RedisClienterMock) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	if mock.SetFunc == nil {
		panic("RedisClienterMock.SetFunc: method is nil but RedisClienter.Set was just called")
	}
	callInfo := struct {
		Key        string
		Value      interface{}
		Expiration time.Duration
	}{
		Key:        key,
		Value:      value,
		Expiration: expiration,
	}
	lockRedisClienterMockSet.Lock()
	mock.calls.Set = append(mock.calls.Set, callInfo)
	lockRedisClienterMockSet.Unlock()
	return mock.SetFunc(key, value, expiration)
}

// SetCalls gets all the calls that were made to Set.
// Check the length with:
//     len(mockedRedisClienter.SetCalls())
func (mock *RedisClienterMock) SetCalls() []struct {
	Key        string
	Value      interface{}
	Expiration time.Duration
} {
	var calls []struct {
		Key        string
		Value      interface{}
		Expiration time.Duration
	}
	lockRedisClienterMockSet.RLock()
	calls = mock.calls.Set
	lockRedisClienterMockSet.RUnlock()
	return calls
}

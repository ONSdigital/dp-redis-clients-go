dp-redis
================

`dp-redis` is a Go client for adding/retrieving user session objects to/from a Redis Cache instance. 

### Getting started
- Add dp-redis to your project using `go get github.com/ONSdigital/dp-redis`

### Dependencies
- No further dependencies other than those defined in go.mod

### Usage

```go
import (
    "crypto/tls"

    dpRedis "github.com/ONSdigital/dp-redis"
)

func main() {
    cfg := dpRedis.Config{
        Addr:     "redis_address",
        Password: "redis_password",
        Database: "database_name",
        TTL:      0, // Time to live config
        TLS: &tls.Config{
            // configure as required
        },
    }

    cli, err := dpredis.NewClient(cfg)
    if err != nil {
        // handle err
    }
    ...
}   
```

Get session by ID:

```go
s, err := cache.GetByID("the_session_id")

if err != nil {
    // handle error
}
```
Get session by email:
```go
s, err := cache.GetByEmail("user_email")

if err != nil {
    // handle error
}
```
Set session:
```go
startTime := time.Now()

s := &Session{
        ID:           "1234",
        Email:        "user@email.com",
        Start:        startTime,
        LastAccessed: startTime, 
    }

if err := cache.Set(s); err != nil {
    // handle error
}
```
Delete all sessions:

```
if err := cache.DeleteAll(); err != nil {
    // handle error
    ...
}
```

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2020, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

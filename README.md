dp-redis
================

### Getting started
- Add dp-redis to your project using `go get github.com/ONSdigital/dp-redis`

### Dependencies
- No further dependencies other than those defined in go.mod

### Usage
Sample use of `GetByID()`:

*s represents a Session in all examples*

```
s, err := cache.GetByID(ID)

if err != nil {
    panic(err)
    return
}
```

Sample use of `SetSession()`:

```
s := &session.Session{
        ID:           "1234",
        Email:        "user@email.com",
        Start:        time.Time{},
        LastAccessed: time.Time{}, 
    }

if err := cache.Set(s); err != nil {
    panic(err)
	return
}
```

Sample use of `DeleteAll()`:

```
if err := cache.DeleteAll(); err != nil {
    panic(err)
    return
}
```

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2020, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

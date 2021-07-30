# go-errors
Errors package for golang

Wrapped on the basis of the "errors" package and added error codes.

Stack printing will be supported in the future.

Features:
- [x] support error codes 
- [x] wrapped message
- [x] wrapped error
- [x] compare error 
- [x] parse from error
- [ ] stack print
## Getting started
```shell
go get -u github.com/lockp111/go-errors
```

## Usage

### New
```golang
// no error code
var err = New("new error")
// with code
var err = New("new error", WithCode(255))
```

### Register
```golang
// register error, duplate error code will panic
var (
    ErrNotFound = Register(101, "not found")
    ErrUnknow = Register(102, "unknow")
)
```

### WithError
```golang
if err := FindUser(1); err != nil {
    return ErrNotFound.WithError(err)
}
```

### WithMessage
```golang
err := FindUser(1)
if err == ErrNotFound {
    return ErrNotFound.WithMessage("find user")
}
```
logdeb
======

logdeb is a Go log library that also supports debug levels, aimed for debugging. It can use many log writers. The library is inspired by `gogits/logs` .

*This is almost a fork of [gogits/logs](https://github.com/gogits/logs).*


### Install the library

To install the library just run

```sh
go get github.com/malix0/logdeb
```

### Supported writers

At the moment the supported log writers are `console` and `file`.

### How to use it?

Import the library

```go
import "github.com/malix0/logdeb"
```

then define configuration, init the log and use it

```go
// Config the console writer to accept messages with severity debug (5)
config := `{"console":{"sev":5}}`
l := logdeb.NewLogDeb(10, config)
defer l.Destroy()
l.Deb("TestConsole", "console - write debug message")
l.Err("TestConsole", "console - write error message")
```

with severity debug, also debug levels can be used

```go
// Config the console writer to accept messages with severity debug (5) and debug level verbose (3)
config := `{"console":{"sev":5, "dlev":3}}`
l := logdeb.NewLogDeb(10, config)
defer l.Destroy()
l.Debl("TestConsole", "console - write debug message with level 3", 3)
// This message will not be written
l.Debl("TestConsole", "console - don't write debug message with level 4", 4)
```

the severity and debug level can be defined by function

```go
// Config the console to use function rules
config := `{"main":{"usefncrules":true},"console":{"sev":5, "fncrules":{"TestConsole.writeme":{"sev":5},"TestConsole.dontwriteme":{"sev":2}}}}`
l := logdeb.NewLogDeb(10, config)
defer l.Destroy()
l.Deb("TestConsole.writeme", "console - write debug message using function rule")
// This message will not be written
l.Deb("TestConsole.dontwriteme", "console - don't write debug message using function rule")
```

### TODO's
- Support more writers
- Add log rotate to file writer
- Add hot reconfiguration for production system debugging purpose

# use fswatch to detect sqlite changes

```
$ go run observer.go &
$ go run writer.go &
$ go run writer.go -slow -interactive
```
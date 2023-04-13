# go-stream
A simple `Stream` api for Go

# Usage

```go
// create
stream := From(1, 2, 3, 4, 5)
// consume
str := stream.Drop(2).Take(2).Join("-")
fmt.Println(str) //3-4
```

Python generator code

```python
def underscore_to_camelcase(s):
    def camelcase():
        yield str.lower
        while True:
            yield str.capitalize

    return ''.join(f(sub) for sub, f in zip(s.split('_'), camelcase()))
```

Equivalent

```go
func underscoreToCamelCase(str string) string {
	capitalize := func(s string) string {
		return strings.ToUpper(s[0:1]) + s[1:]
	}

	stream := func(c Consumer[func(string) string]) {
		c(strings.ToLower)
		for {
			c(capitalize)
		}
	}

	return Zip(stream, strings.Split(str, "_"), func(f func(string) string, str string) string {
		return f(str)
	}).Join("")
}
```

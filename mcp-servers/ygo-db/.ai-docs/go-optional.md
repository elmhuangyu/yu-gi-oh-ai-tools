# Optional Types

When a function may not return a value (e.g., lookup functions that can fail), use [`moznion/go-optional`](https://github.com/moznion/go-optional) instead of returning a zero value to indicate failure.

## Simple usage example

```go
import (
    "github.com/moznion/go-optional"
)

func GetAttributeByName(lang, name string) optional.Option[int] {
    // ... lookup logic
    if found {
        return optional.Some(value)
    }
    return optional.None[int]()
}

// Caller using IsSome() and Unwrap(). Unwrap() will not change the value in the optional.
result := GetAttributeByName(LangEnUS, "Dark")
if result.IsSome() {
    value := result.Unwrap()
    // use value
}

// Or use Take() to get the value and error if it's None. Take() will not change the value in the optional.
result := GetAttributeByName(LangEnUS, "Dark")
if value, ok := result.Take(); ok {
    // use value
}
```

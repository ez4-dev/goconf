This library helps reading configuration from environment variables.
The initial release was inspired by [Spring Boot Relax Binding](https://github.com/spring-projects/spring-boot/wiki/Relaxed-Binding-2.0).
It support primitive types (bool, string, number, map, slice) and complex struct type that contains these types.

# Usage
## Simple struct
You configuration struct
```go
type SimpleConfig struct {
	StrValue   string  `env:"STR"`
	IntValue   int     `env:"INT"`
	UIntValue  uint    `env:"UINT"`
	FloatValue float32 `env:"FLOAT"`
	BoolValue  bool    `env:"BOOL"`
}
```
Your code
```go
	// Given
	os.Setenv("PREFIX_STR", "abc")
	os.Setenv("PREFIX_INT", "-123")
	os.Setenv("PREFIX_UINT", "12")
	os.Setenv("PREFIX_FLOAT", "12.34")
	os.Setenv("PREFIX_BOOL", "true")

	config := SimpleConfig{}
	envconfig.ParseConfigFromEnv("PREFIX", &config)
```

# Change logs
- V1.0
    + Initialize first version

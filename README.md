This library helps manipulate configuration with only std packages in used.

# Usage
```bash
go get https://github.com/trinhpham/goconf
```

# Packages
## envconfig
This helps reading configuration from environment variables.
The initial release was inspired by [Spring Boot Relax Binding](https://github.com/spring-projects/spring-boot/wiki/Relaxed-Binding-2.0).
This supports primitive types (bool, string, number, map, slice) and complex struct type that contains these types.

### Usage
Find more usage examples in unit-test
## Simple struct
Modeling your configuration like this
```go
type SimpleConfig struct {
	StrValue   string  `env:"STR"`
	IntValue   int     `env:"INT"`
	UIntValue  uint    `env:"UINT"`
	FloatValue float32 `env:"FLOAT"`
	BoolValue  bool    `env:"BOOL"`
}
```
You should be able to load your config as below
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
### Nested struct with map and slice
Your configuration struct
```go
type NestedConfig struct {
	StrValue   string                  `env:"STR"`
	MapInt     map[string]int8         `env:"MINT"`
	SliceFloat []float64               `env:"SFL"`
	MapValue   map[string]SimpleConfig `env:"MAP"`
	SliceVal   []SimpleConfig          `env:"SLICE"`
}
```
Your code should be
```go
	// Given
	os.Setenv("PREFIX_STR", "abc")
	os.Setenv("PREFIX_MINT_ITEM1", "1")
	os.Setenv("PREFIX_MINT_ITEM2", "2")
	os.Setenv("PREFIX_SFL_1", "12.34")
	os.Setenv("PREFIX_SFL_4", "56.78")
	os.Setenv("PREFIX_MAP_S1_STR", "abc")
	os.Setenv("PREFIX_SLICE_5_STR", "xyz")
	os.Setenv("PREFIX_SLICE_5_BOOL", "n")

	config := NestedConfig{}
	envconfig.ParseConfigFromEnv("PREFIX", &config)
```
# Change logs
- V1.0
    + Initialize first version

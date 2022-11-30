package envconfig

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

/*
Parse config from environment variable
*/
func Load(envPrefix string, parsedConfig interface{}) {
	// Validate input
	rConfig := reflect.ValueOf(parsedConfig)
	if !rConfig.IsValid() || rConfig.IsZero() || rConfig.IsNil() || rConfig.Kind() != reflect.Pointer {
		panic("Value should be a valid pointer")
	}

	rConfig = reflect.Indirect(rConfig)

	// Load config prefix
	configPrefix := strings.ToUpper(envPrefix)
	prefixLength := len(configPrefix + ENV_DELIMITER)

	for _, env := range os.Environ() {
		// Parse to pair of name:value
		envTokens := strings.SplitN(env, ENV_ASSIGNER, 2)
		envName := strings.ToUpper(envTokens[0])
		if strings.HasPrefix(envName, configPrefix) {
			var envNames []string
			if len(envName) > prefixLength {
				envName = envName[prefixLength:]
				envNames = strings.Split(envName, ENV_DELIMITER)
			}
			setField(rConfig, envNames, envTokens[1])
		}
	}
}

func setField(field reflect.Value, childTokens []string, envValue string) {
	field = reflect.Indirect(field)
	switch field.Type().Kind() {
	case reflect.Map:
		setMapValue(field, childTokens[0], childTokens[1:], envValue)
	case reflect.Struct:
		setStructValue(field, childTokens[0], childTokens[1:], envValue)
	case reflect.Slice:
		setSliceValue(field, childTokens[0], childTokens[1:], envValue)
	case reflect.String:
		setStringValue(field, childTokens, envValue)
	case reflect.Bool:
		setBoolValue(field, childTokens, envValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		setNumberValue(field, childTokens, envValue)
	default:
		panic(fmt.Sprintf("Field type is not supported yet: %v", field.Type()))
	}
}

func setBoolValue(field reflect.Value, childTokens []string, envValue string) {
	assertSettable(field, childTokens)
	val := strings.ToLower(envValue)
	switch val {
	case "true", "t", "1", "yes", "y":
		field.SetBool(true)
	case "false", "f", "0", "no", "n":
		field.SetBool(false)
	default:
		panic("Not support bool value: " + envValue)
	}
}

func setNumberValue(field reflect.Value, childTokens []string, envValue string) {
	assertSettable(field, childTokens)
	switch field.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, e := strconv.ParseInt(envValue, 10, 64)
		if e != nil || field.OverflowInt(i) {
			panic(fmt.Sprintf("Unable to read number [%s]: %v", envValue, e))
		}
		field.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, e := strconv.ParseUint(envValue, 10, 64)
		if e != nil || field.OverflowUint(i) {
			panic(fmt.Sprintf("Unable to read number [%s]: %v", envValue, e))
		}
		field.SetUint(i)
	case reflect.Float32:
		f, e := strconv.ParseFloat(envValue, 32)
		if e != nil || field.OverflowFloat(f) {
			panic(fmt.Sprintf("Unable to read number [%s]: %v", envValue, e))
		}
		field.SetFloat(f)
	case reflect.Float64:
		f, e := strconv.ParseFloat(envValue, 64)
		if e != nil || field.OverflowFloat(f) {
			panic(fmt.Sprintf("Unable to read number [%s]: %v", envValue, e))
		}
		field.SetFloat(f)
	}

}

func setStringValue(field reflect.Value, childTokens []string, envValue string) {
	assertSettable(field, childTokens)
	field.SetString(envValue)
}

func assertSettable(field reflect.Value, childTokens []string) {
	if !(field.CanSet() && len(childTokens) == 0) {
		panic(fmt.Sprintf("Unable to set value %v | %d", field.CanSet(), len(childTokens)))
	}
}

func setStructValue(rStruct reflect.Value, fieldEnvTag string, childTokens []string, envValue string) {
	for i := 0; i < rStruct.Type().NumField(); i++ {
		f := rStruct.Type().Field(i)
		envName := f.Tag.Get("env")
		if envName == fieldEnvTag {
			// Found needed value to set
			field := rStruct.FieldByName(f.Name)
			setField(field, childTokens, envValue)
			break
		}
	}
}

func setSliceValue(field reflect.Value, index string, childTokens []string, envValue string) {
	i, e := strconv.Atoi(index)
	if e != nil || i < 0 {
		panic(fmt.Sprintf("Invalid index: %s | %v", index, e))
	}

	if field.IsNil() {
		// Create array
		field.Set(reflect.MakeSlice(field.Type(), i+1, i+1))
	} else if field.Len() < i {
		// Append item to have i-th item
		missingCount := i - field.Len() + 1
		newSlice := reflect.AppendSlice(field, reflect.MakeSlice(field.Type(), missingCount, missingCount))
		field.Set(newSlice)
	}
	sliceItem := field.Index(i)
	setField(sliceItem, childTokens, envValue)
}

func setMapValue(rMap reflect.Value, mapKey string, childTokens []string, envValue string) {
	if rMap.IsNil() && rMap.CanSet() {
		rMap.Set(reflect.MakeMap(rMap.Type()))
	}

	rMapKey := reflect.ValueOf(strings.ToLower(mapKey))

	// Create new map value with org data
	rMapValue := reflect.New(rMap.Type().Elem()).Elem()
	if orgMapValue := rMap.MapIndex(rMapKey); orgMapValue.IsValid() {
		cloneData(rMapValue, orgMapValue)
	}

	setField(rMapValue, childTokens, envValue)
	rMap.SetMapIndex(rMapKey, rMapValue)
}

func cloneData(dstValue, orgValue reflect.Value) {
	for i := 0; i < dstValue.NumField(); i++ {
		dstField := dstValue.Field(i)
		if dstField.CanSet() {
			dstField.Set(orgValue.Field(i))
		} else {
			panic(fmt.Sprintf("Cannot clone field: %s", dstValue.Type().Field(i).Name))
		}
	}
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

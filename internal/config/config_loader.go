package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
)

const envTagName = "env"

func ReadConfigFile(configPaths []string, configuration interface{}) (err error) {
	configValue := reflect.ValueOf(configuration)
	if typ := configValue.Type(); typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("configuration should be a pointer to a struct type")
	}

	for _, path := range configPaths {
		err = loadConfigFile(path, &configuration)
		if err == nil {
			break
		}
	}

	if err != nil {
		return err
	}

	return enrichConfigFromEnvVariables(configuration)
}

func loadConfigFile(filePath string, configuration interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	ext := filepath.Ext(filePath)

	switch ext {
	case ".json":
		err = json.Unmarshal(data, configuration)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, configuration)
	default:
		return fmt.Errorf("unsupported config file type: %s", ext)
	}

	if err != nil {
		return err
	}

	return nil
}

func enrichConfigFromEnvVariables(configuration interface{}) (err error) {
	val := reflect.ValueOf(configuration)
	typ := reflect.TypeOf(configuration)

	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("enriched type is not Struct it's a: %s", val.Kind())
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		// check if we've got a field name override for the environment
		tagContent := field.Tag.Get(envTagName)
		if tagContent == "" {
			tagContent = field.Name // if tag no tag then field Name
		}
		value := os.Getenv(tagContent)

		if !field.Anonymous && len(value) > 0 {
			f := val.FieldByName(field.Name)
			setValueFromString(f, value)
		}
	}

	return nil
}

func setValueFromString(f reflect.Value, value string) {
	if !f.IsValid() || !f.CanSet() {
		return
	}

	switch f.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		setStringToInt(f, value, f.Type().Bits())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		setStringToUInt(f, value, f.Type().Bits())

	case reflect.Bool:
		setStringToBool(f, value)

	case reflect.Float32, reflect.Float64:
		setStringToFloat(f, value, f.Type().Bits())

	case reflect.String:
		f.SetString(value)

	case reflect.Struct:
		setJSONStringToStruct(f, value)

	case reflect.Slice, reflect.Array:
		setJSONStringToArray(f, value)

	default:
		handleUnsupportedDataTypeError(f.Kind())
	}
}

func setStringToInt(f reflect.Value, value string, bitSize int) {
	var (
		convertedValue interface{}
		err            error
	)

	convertedValue, err = strconv.ParseInt(value, 10, bitSize)

	if err != nil {
		handleParseDataTypeError("number", value, err)
		return
	}

	setIfNoOverflowInt(f, convertedValue.(int64))
}

func setStringToUInt(f reflect.Value, value string, bitSize int) {
	var (
		convertedValue interface{}
		err            error
	)

	convertedValue, err = strconv.ParseUint(value, 10, bitSize)

	if err != nil {
		handleParseDataTypeError("number", value, err)
		return
	}

	setIfNoOverflowUint(f, convertedValue.(uint64))
}

func setIfNoOverflowInt(f reflect.Value, value int64) {
	if !f.OverflowInt(value) {
		f.SetInt(value)
	}
}

func setIfNoOverflowUint(f reflect.Value, value uint64) {
	if !f.OverflowUint(value) {
		f.SetUint(value)
	}
}

func setStringToBool(f reflect.Value, value string) {
	convertedValue, err := strconv.ParseBool(value)
	if err != nil {
		handleParseDataTypeError("boolean", value, err)
		return
	}
	f.SetBool(convertedValue)
}

func setStringToFloat(f reflect.Value, value string, bitSize int) {
	convertedValue, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		handleParseDataTypeError("float", value, err)
		return
	}
	if !f.OverflowFloat(convertedValue) {
		f.SetFloat(convertedValue)
	}
}

func setJSONStringToStruct(f reflect.Value, value string) {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(value), &jsonMap)
	if err != nil {
		handleParseDataTypeError("struct", value, err)
		return
	}

	for k, v := range jsonMap {
		subField := f.FieldByName(k)
		if subField.IsValid() && subField.CanSet() {
			setValueFromString(subField, fmt.Sprintf("%v", v))
		}
	}
}

func setJSONStringToArray(f reflect.Value, value string) {
	var jsonArr []interface{}
	err := json.Unmarshal([]byte(value), &jsonArr)
	if err != nil {
		handleParseDataTypeError("array", value, err)
		return
	}

	f.Set(reflect.MakeSlice(f.Type(), len(jsonArr), len(jsonArr)))
	for i, v := range jsonArr {
		jsonItemVal, err := json.Marshal(v)
		if err != nil {
			handleParseDataTypeError("array element", fmt.Sprintf("%v", v), err)
			return
		}
		setValueFromString(f.Index(i), string(jsonItemVal))
	}
}

func handleParseDataTypeError(dataType, value string, err error) {
	fmt.Printf("Error parsing %s type value \"%s\": %v\n", dataType, value, err)
}

func handleUnsupportedDataTypeError(kind reflect.Kind) {
	fmt.Printf("Unsupported type: %v\n", kind)
}

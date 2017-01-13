package main

import "github.com/davecgh/go-spew/spew"
import "reflect"
import "strings"
import "bufio"
import "strconv"
import "fmt"

type inter struct {
	A string
	B bool
}

type test struct {
	in inter `mparse:"some text"`
	C  int64
	D  string `mparse:"default"`
	A  bool   `mparse:"abc"`
}

func parseString(str string, v interface{}) {
	value := reflect.Indirect(reflect.ValueOf(v)).Elem()
	fields := getStructTypes(value.Type(), true)

	reader := strings.NewReader(str)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "/") {
			f := strings.Fields(strings.TrimPrefix(line, "/"))
			if len(f) > 0 {
				name := strings.ToLower(f[0])

				end := strings.TrimPrefix(line, f[0])
				end = strings.TrimSpace(end)

				if _, ok := fields[name]; ok {
					setField(&value, name, end)
				} else {
					appendDefaultField(&value, end)
				}
			} else {
				appendDefaultField(&value, line)
			}
		}
	}
}

func setValue(v *reflect.Value, value string) {
	switch v.Kind() {
	case reflect.Bool:
		b, _ := strconv.ParseBool(value)
		v.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, _ := strconv.ParseInt(value, 10, 64)
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, _ := strconv.ParseUint(value, 10, 64)
		v.SetUint(i)
	case reflect.Float32, reflect.Float64:
		f, _ := strconv.ParseFloat(value, 64)
		v.SetFloat(f)
	case reflect.String:
		v.SetString(value)
	}
}

func setField(v *reflect.Value, fieldName string, fieldValue string) {
	if _, ok := getStructTypes(v.Type(), false)[fieldName]; ok {
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)

			name := strings.ToLower(field.Name)
			tag := field.Tag.Get("mparse")
			if tag != "" {
				name = strings.ToLower(tag)
			}

			f := v.Field(i)
			if f.CanSet() && (name == fieldName) {
				setValue(&f, fieldValue)
			}
		}
	} else {
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			if field.Type.Kind() == reflect.Struct {
				f := v.Field(i)
				setField(&f, fieldName, fieldValue)
			}
		}
	}
}

func getField(v *reflect.Value, fieldName string) *reflect.Value {
	if _, ok := getStructTypes(v.Type(), false)[fieldName]; ok {
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)

			name := strings.ToLower(field.Name)
			tag := field.Tag.Get("mparse")
			if tag != "" {
				name = strings.ToLower(tag)
			}

			f := v.Field(i)
			if f.CanSet() && (name == fieldName) {
				return &f
			}
		}
	} else {
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			if field.Type.Kind() == reflect.Struct {
				f := v.Field(i)
				return &f
			}
		}
	}
	return nil
}

func appendDefaultField(v *reflect.Value, fieldValue string) {
	value := getField(v, "default").String()
	if value != "" {
		value += "\n"
	}
	value += fieldValue

	setField(v, "default", value)
}

func getStructTypes(v reflect.Type, recursive bool) map[string]reflect.Type {
	types := make(map[string]reflect.Type)
	if v.Kind() != reflect.Struct {
		return types
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Type.Kind() == reflect.Struct {
			if recursive {
				for key, value := range getStructTypes(field.Type, recursive) {
					types[key] = value
				}
			}
		} else {
			name := strings.ToLower(field.Name)
			tag := field.Tag.Get("mparse")
			if tag != "" {
				name = strings.ToLower(tag)
			}
			types[name] = field.Type
		}
	}

	return types
}

func main() {
	msg := `/a true
	/b true
	/c 12312412
	/d test1
	/d test2
	msg`

	a := &test{}

	fmt.Println("Test message:", msg)
	fmt.Println()
	spew.Dump(a)
	fmt.Println()
	parseString(msg, &a)
	spew.Dump(a)
}

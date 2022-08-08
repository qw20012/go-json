package json

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/qw20012/go-basic"
	"github.com/qw20012/go-basic/ref"
	"github.com/qw20012/go-basic/str"
	"github.com/qw20012/go-json/lexer"
	"github.com/qw20012/go-json/parser"
)

func Parse(jsonStr string) *Container {
	lexer := lexer.NewLexer([]byte(jsonStr))
	parser := parser.NewParser(lexer)
	ast := parser.Parse()

	var json = Container{object: ast}

	return &json
}

func Unmarshal[T any](jsonStr string) T {
	lexer := lexer.NewLexer([]byte(jsonStr))
	parser := parser.NewParser(lexer)
	//fmt.Println("ast")
	ast := parser.Parse()
	//fmt.Println(ast)
	var hold T

	ty := reflect.TypeOf(hold)
	switch ty.Kind() {
	case reflect.Map:
		//return buildMap(ty, ast).Interface().(T)
		return ref.GetValue[T](buildMap(ty, ast))
	case reflect.Slice:
		//return buildSlice(ty, ast).Interface().(T)
		return ref.GetValue[T](buildSlice(ty, ast))
	case reflect.Struct:
		//return buildStruct(ty, ast).Interface().(T)
		return ref.GetValue[T](buildStruct(ty, ast))
	case reflect.String:
		//v := reflect.ValueOf(ast.(string))
		//return v.Interface().(T)

		strValue := basic.FromAny[string](ast)
		v := reflect.ValueOf(strValue)
		return ref.GetValue[T](v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		strValue := basic.FromAny[string](ast)
		if i, err := strconv.Atoi(strValue); err == nil {
			v := reflect.ValueOf(i)
			return ref.GetValue[T](v)
			//return v.Interface().(T)
		}
	case reflect.Bool:
		strValue := basic.FromAny[string](ast)
		if realValue, err := strconv.ParseBool(strValue); err == nil {
			v := reflect.ValueOf(realValue)
			//return v.Interface().(T)
			return ref.GetValue[T](v)
		}
	case reflect.Float32:
		strValue := basic.FromAny[string](ast)
		if realValue, err := strconv.ParseFloat(strValue, 32); err == nil {
			v := reflect.ValueOf(realValue)
			//return v.Interface().(T)
			return ref.GetValue[T](v)
		}
	case reflect.Float64:
		strValue := basic.FromAny[string](ast)
		if realValue, err := strconv.ParseFloat(strValue, 64); err == nil {
			v := reflect.ValueOf(realValue)
			//return v.Interface().(T)
			return ref.GetValue[T](v)
		}
	}

	return hold
}

func buildMap(ty reflect.Type, ast any) reflect.Value {
	v := reflect.MakeMap(ty)
	if mapObject, ok := ast.(map[string]any); ok {
		if ty.Elem().Kind() == reflect.Interface {
			v = reflect.ValueOf(mapObject)
			return v
		}

		for key, value := range mapObject {
			strValue := value.(string)
			switch ty.Elem().Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if realValue, err := strconv.Atoi(strValue); err == nil {
					v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(realValue))
				}
			case reflect.String:
				v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(strValue))
			case reflect.Bool:
				if realValue, err := strconv.ParseBool(strValue); err == nil {
					v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(realValue))
				}
			case reflect.Float32:
				if realValue, err := strconv.ParseFloat(strValue, 32); err == nil {
					v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(realValue))
				}
			case reflect.Float64:
				if realValue, err := strconv.ParseFloat(strValue, 64); err == nil {
					v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(realValue))
				}
			}
		}
	}
	return v
}

func buildStruct(ty reflect.Type, ast any) reflect.Value {
	v := reflect.New(ty)

	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}

	if mapObject, ok := ast.(map[string]any); ok {
		for key, value := range mapObject {
			field := v.FieldByName(key)

			if field.Type().Kind() == reflect.Ptr {
				elemType := field.Type().Elem()

				if elemType.Kind() == reflect.Struct {
					field.Set(buildStruct(elemType, value).Addr())
					continue
				} else if elemType.Kind() == reflect.Slice {
					slicePtrVal := reflect.New(elemType)
					sliceVal := reflect.Indirect(slicePtrVal)
					sliceVal.Set(buildSlice(elemType, value))
					field.Set(slicePtrVal)

					continue
				}

				ref.SetStructBasicPtrField(field, value)
				continue
			}

			if field.Type().Kind() == reflect.Struct {
				field.Set(buildStruct(field.Type(), value))
				continue
			} else if field.Type().Kind() == reflect.Slice {
				field.Set(buildSlice(field.Type(), value))
				continue
			}

			ref.SetStructBasicField(field, value)
		}
	}
	return v
}

func buildSlice(ty reflect.Type, ast any) reflect.Value {
	v := reflect.ValueOf(ast.([]any))
	if ty.Elem().Kind() == reflect.Interface {
		return v
	}

	newCap := (v.Cap() + 1) * 2

	aSlice := reflect.MakeSlice(ty, 0, newCap)
	//t = reflect.Indirect(t)
	for i := 0; i < v.Len(); i++ {
		if ty.Elem().Kind() == reflect.Struct {
			aSlice = reflect.Append(aSlice, buildStruct(ty.Elem(), v.Index(i).Interface()))
			continue
		} else if ty.Elem().Kind() == reflect.Ptr {
			aSlice = ref.AppendSliceBasicPtrElem(aSlice, v.Index(i))
			continue
		}

		aSlice = ref.AppendSliceBasicElem(aSlice, v.Index(i))
	}

	return aSlice
}

func Marshal(source any) string {
	json := ""
	ty := reflect.TypeOf(source)
	val := reflect.ValueOf(source)
	switch ty.Kind() {
	case reflect.Struct:
		json = parseStruct(ty, val)
	case reflect.Map:
		json = parseMap(val)
	case reflect.Ptr:
		json = parsePtr(ty, val)
	}

	return json
}

func parsePtr(ty reflect.Type, value reflect.Value) string {
	json := ""

	if value.IsNil() {
		return "nil"
	}

	ty = ty.Elem()
	value = value.Elem()
	val := reflect.Indirect(value)

	switch val.Type().Kind() {
	case reflect.Struct:
		json = parseStruct(ty, value)

	default:

	}

	return json
}

func parseMap(val reflect.Value) string {
	json := ""
	for _, e := range val.MapKeys() {
		v := val.MapIndex(e)
		switch t := v.Interface().(type) {
		case int, int8, int16, int32, int64, float32, float64:

			json = str.Contact(json, e, ":", t, ",")
		case string:

			json = str.Contact(json, e, ":", t, ",")
		case bool:

			json = str.Contact(json, e, ":", t, ",")
		case map[string]any:
			json = str.Contact(json, e, ":", parseMap(reflect.ValueOf(t)), ",")
		default:
			fmt.Println("not found")
		}
	}
	json = removeLastComma(json)
	json = addBrace(json)
	return json
}

func parseStruct(ty reflect.Type, value reflect.Value) string {
	json := ""
	if ty.Kind() == reflect.Ptr {
		if value.IsNil() {
			return "nil"
		}
		//(val.Field(i).Type().Kind() == reflect.Ptr && val.Field(i).Type().Elem().Kind() == reflect.Struct)
		ty = ty.Elem()
		value = value.Elem()
	}

	val := reflect.Indirect(value)

	for i := 0; i < val.NumField(); i++ {
		fieldName := ty.Field(i).Tag.Get("json")
		if str.IsEmpty(fieldName) {
			fieldName = ty.Field(i).Name
		}

		switch val.Field(i).Type().Kind() {
		case reflect.Struct:
			json = str.Contact(json, fieldName, ":", parseStruct(val.Field(i).Type(), val.Field(i)), ",")

		case reflect.Slice, reflect.Array:
			json = str.Contact(json, fieldName, ":", parseSlice(val.Field(i)), ",")
		case reflect.Ptr:

			fieldElem := val.Field(i).Elem()

			if !fieldElem.IsValid() {
				continue
			}

			if fieldElem.Type().Kind() == reflect.Struct {
				json = str.Contact(json, fieldName, ":", parseStruct(fieldElem.Type(), fieldElem), ",")
				continue
			} else if fieldElem.Type().Kind() == reflect.Slice {
				json = str.Contact(json, fieldName, ":", parseSlice(fieldElem), ",")
				continue
			}

			prtValue := parseStructFieldPtr(fieldElem, fieldName)
			json = str.Contact(json, prtValue)
		default:
			json = str.Contact(json, fieldName, ":", val.Field(i), ",")
		}
	}

	json = removeLastComma(json)

	json = addBrace(json)

	return json
}

func parseSlice(value reflect.Value) string {
	sliceStr := ""
	count := value.Len()
	elemTy := value.Type().Elem()

	switch elemTy.Kind() {
	default:
		for sliceIndex := 0; sliceIndex < count; sliceIndex++ {
			child := value.Index(sliceIndex)

			sliceStr = str.Contact(sliceStr, child, ",")
		}
		sliceStr = removeLastComma(sliceStr)
		sliceStr = str.Contact("[", sliceStr, "]")
	case reflect.Struct:
		for sliceIndex := 0; sliceIndex < count; sliceIndex++ {
			child := value.Index(sliceIndex)

			sliceStr = str.Contact(sliceStr, parseStruct(child.Type(), child), ",")
		}
		sliceStr = removeLastComma(sliceStr)
		sliceStr = str.Contact("[", sliceStr, "]")
	case reflect.Ptr:
		for sliceIndex := 0; sliceIndex < count; sliceIndex++ {
			child := value.Index(sliceIndex)
			fieldElem := child.Elem()

			if !fieldElem.IsValid() {
				break
			}

			sliceStr = str.Contact(sliceStr, fieldElem, ",")
		}
		sliceStr = removeLastComma(sliceStr)
		sliceStr = str.Contact("[", sliceStr, "]")
	}
	return sliceStr
}

func addBrace(json string) string {
	json = str.Contact("{", json, "}")
	return json
}

func removeLastComma(json string) string {
	if strings.Contains(json, ",") {
		json = json[:len(json)-1]
	}
	return json
}

func parseStructFieldPtr(value reflect.Value, fieldName string) string {
	json := ""

	switch value.Type().Kind() {
	case reflect.String:
		if realValue, ok := value.Interface().(string); ok {
			json = str.Contact(json, fieldName, ":", realValue, ",")
		}
	case reflect.Int:
		if realValue, ok := value.Interface().(int); ok {
			json = str.Contact(json, fieldName, ":", realValue, ",")
		}

	case reflect.Int8:
		if realValue, ok := value.Interface().(int8); ok {
			json = str.Contact(json, fieldName, ":", realValue, ",")
		}
	case reflect.Int16:
		if realValue, ok := value.Interface().(int16); ok {
			json = str.Contact(json, fieldName, ":", realValue, ",")
		}
	case reflect.Int32:
		if realValue, ok := value.Interface().(int32); ok {
			json = str.Contact(json, fieldName, ":", realValue, ",")
		}
	case reflect.Int64:
		if realValue, ok := value.Interface().(int64); ok {
			json = str.Contact(json, fieldName, ":", realValue, ",")
		}
	case reflect.Bool:
		if realValue, ok := value.Interface().(bool); ok {
			json = str.Contact(json, fieldName, ":", realValue, ",")
		}
	case reflect.Float32:
		if realValue, ok := value.Interface().(float32); ok {
			json = str.Contact(json, fieldName, ":", realValue, ",")
		}
	case reflect.Float64:
		if realValue, ok := value.Interface().(float64); ok {
			json = str.Contact(json, fieldName, ":", realValue, ",")
		}
	}

	return json
}

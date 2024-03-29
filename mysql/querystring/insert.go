package querystring

import (
	"fmt"
	"reflect"
	"strings"
)

type insertSQL struct {
	table      string
	fieldArray []string
	valueArray [][]string
}

func InsertInto(table string) *insertSQL {
	return &insertSQL{
		table:      table,
		valueArray: make([][]string, 0),
	}
}

func (this *insertSQL) GetSQL() string {
	valueArray := make([]string, len(this.valueArray))
	for i, v := range this.valueArray {
		valueArray[i] = fmt.Sprintf("(%s)", strings.Join(v, ","))
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", this.table,
		strings.Join(this.fieldArray, ","),
		strings.Join(valueArray, ","))
}

func (this *insertSQL) SetFieldAndValue(fieldAndValue map[string]interface{}) *insertSQL {
	this.fieldArray = make([]string, len(fieldAndValue))
	valueArray := make([]string, len(fieldAndValue))
	index := 0
	for k, v := range fieldAndValue {
		this.fieldArray[index] = k
		valueArray[index] = formatFieldValue(v)
		index++
	}
	this.valueArray = append(this.valueArray, valueArray)
	return this
}

func (this *insertSQL) SetObject(object interface{}) *insertSQL {
	this.fieldArray = make([]string, 0)
	valueArray := make([]string, 0)

	valueof := reflect.ValueOf(object)
	if valueof.Type().Kind() == reflect.Ptr {
		valueof = valueof.Elem()
	}

	if valueof.Type().Kind() != reflect.Struct {
		panic(fmt.Sprintf("querystring[InsertObject] %s is not a struct", valueof.Type().Name()))
	}

	for i := 0; i < valueof.NumField(); i++ {
		tags := valueof.Type().Field(i).Tag.Get("db")
		if len(tags) == 0 {
			continue
		}
		tagArray := strings.Split(tags, ",")
		if len(tagArray) > 0 && strings.ToLower(tagArray[1]) == "auto_increment" {
			continue
		}
		this.fieldArray = append(this.fieldArray, tagArray[0])
		switch valueof.Type().Field(i).Type.Kind() {
		case reflect.String:
			valueArray = append(valueArray, fmt.Sprintf("'%s'", Escape(valueof.String())))
		case reflect.Bool:
			valueArray = append(valueArray, fmt.Sprintf("%t", valueof.Bool()))
		case reflect.Int, reflect.Int32, reflect.Int64:
			valueArray = append(valueArray, fmt.Sprintf("%d", valueof.Int()))
		case reflect.Uint, reflect.Uint32, reflect.Uint64:
			valueArray = append(valueArray, fmt.Sprintf("%d", valueof.Uint()))
		case reflect.Float32, reflect.Float64:
			valueArray = append(valueArray, fmt.Sprintf("%t", valueof.Bool()))
		}
	}

	this.valueArray = append(this.valueArray, valueArray)

	return this
}

func formatFieldValue(value interface{}) string {
	valueof := reflect.ValueOf(value)
	switch valueof.Type().Kind() {
	case reflect.String:
		return fmt.Sprintf("'%s'", valueof.String())
	case reflect.Bool:
		return fmt.Sprintf("%t", valueof.Bool())
	case reflect.Int, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", valueof.Int())
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", valueof.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%t", valueof.Bool())
	}
	return "''"
}

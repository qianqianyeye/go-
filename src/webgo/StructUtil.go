package webgo

import (
	"reflect"
)

//获取结构体json标签
func GetStructTagJson(stru interface{}) map[string]string {
	t := reflect.TypeOf(stru).Elem()
	resultmap := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		resultmap[t.Field(i).Name] = t.Field(i).Tag.Get("json")
	}
	return resultmap
}
//结构体转Map
func StructToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

//结构体转JSONMap(根据Tag标签)
func StructToJsonMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	resultmap := make(map[string]interface{} )
	for i := 0; i < t.NumField(); i++ {
		resultmap[t.Field(i).Tag.Get("json")] = v.Field(i).Interface()
	}
	return resultmap
}
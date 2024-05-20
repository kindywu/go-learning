package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string `json:"name" comment:"名称"`
	Age  int    `json:"age" comment:"年龄"`
}

func main() {
	// 创建一个Person实例
	person := Person{Name: "Alice", Age: 30}

	// 使用reflect.ValueOf获取Person的值
	personValue := reflect.ValueOf(person)

	// 使用TypeOf获取Person的类型信息
	personType := personValue.Type()

	// 遍历结构体的所有字段
	for i := 0; i < personType.NumField(); i++ {
		// 获取字段信息
		field := personType.Field(i)
		// 获取字段的值
		fieldValue := personValue.Field(i)

		// 打印字段名和值
		fmt.Printf("Field: %s, Value: %v \n", field.Name, fieldValue.Interface())
		fmt.Printf("Tags: %s, json: %v, comment: %s \n", field.Name, field.Tag.Get("json"), field.Tag.Get("comment"))
	}
}

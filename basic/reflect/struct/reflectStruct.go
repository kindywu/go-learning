package basic

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name     string
	Age      int
	isLocked bool
}

func ReflectPerson() {
	// 创建一个Person类型的实例
	p := Person{Name: "Alice", Age: 30, isLocked: true}

	// 使用reflect.ValueOf来获取p的reflect.Value
	// v := reflect.ValueOf(p)

	v := reflect.ValueOf(&p).Elem() // 注意这里我们获取了p的地址，然后使用.Elem()来获取指针指向的值

	/**
	这两段代码都是用来获取 `p` 的 `reflect.Value`，但它们获取的方式和结果有所不同。
	1. `v := reflect.ValueOf(p)`
	   这段代码直接对 `p` 进行反射，获取 `p` 的 `reflect.Value`。
	   这里的 `p` 是一个值类型（在Go中，除了指针、切片、映射、通道、函数、接口之外的所有类型都是值类型）。
	   当你对一个值类型使用 `reflect.ValueOf` 时，你会得到一个不可寻址的 `reflect.Value`。
	   这意味着你不能通过这个 `reflect.Value` 来修改原始对象的值。
	   如果你尝试使用 `.SetInt` 或其他修改值的方法，你会遇到 "panic: reflect: reflect.Value.SetInt using unaddressable value" 的错误。

	2. `v := reflect.ValueOf(&p).Elem()`
	   这段代码首先获取 `p` 的地址，然后通过 `reflect.ValueOf` 获取这个地址的 `reflect.Value`。
	   接着，使用 `Elem()` 方法来获取指针指向的实际值的 `reflect.Value`。这里的 `v` 是一个可寻址的 `reflect.Value`，因为它是通过指针间接引用 `p` 的。
	   这意味着你可以通过这个 `reflect.Value` 来修改 `p` 的值。

	总结来说，`reflect.ValueOf(p)` 获取的是一个不可寻址的 `reflect.Value`，而 `reflect.ValueOf(&p).Elem()` 获取的是一个可寻址的 `reflect.Value`。
	在需要修改值类型变量的值时，你应该使用后者。
		**/

	// 检查类型
	fmt.Println("Type of p:", v.Type())

	// 获取结构体的字段数量
	fmt.Println("Number of fields in Person:", v.NumField())

	// 遍历结构体的字段
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		var value interface{}
		if field.CanInterface() {
			value = field.Interface()
		}
		fmt.Printf("Field %d: Name=%s, Type=%s, CanInterface=%t, Value=%v\n", i, field.Type().Name(), field.Type(), field.CanInterface(), value)
	}

	// 使用反射设置字段的值
	ageField := v.FieldByName("Age")
	if ageField.CanSet() {
		ageField.SetInt(35)
		fmt.Println("After setting Age by reflection:", p)
	}

	// 使用反射获取字段的值
	fmt.Printf("Age of person: %d\n", ageField.Int())
}

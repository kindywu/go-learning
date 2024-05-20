package main

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func main() {
	tryPreLoadModule()

	path := "./lua/plugin2.lua"
	L := lua.NewState()
	defer L.Close()
	fn, err := L.LoadFile(path)
	if err != nil {
		panic(err)
	}
	L.Push(fn)
	err = L.PCall(0, -1, nil)
	if err != nil {
		panic(err)
	}

	m := L.Get(-1)
	println(m.Type().String())
	subtraction := L.GetField(m, "subtraction").(*lua.LFunction)
	println(subtraction)

	err = L.CallByParam(lua.P{
		Fn:      subtraction, // 要调用的函数
		NRet:    1,           // 期望返回值的数量
		Protect: true,        // 是否捕获异常
	}, lua.LNumber(8), lua.LNumber(5)) // 传递参数

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		// 获取并打印结果
		result := L.Get(-1)
		fmt.Println("Result of subtraction:", result)
		L.Pop(1) // 清理栈
	}

	// subtraction := L.GetGlobal("Subtraction").(*lua.LFunction)
	// println(subtraction)

	// err = L.CallByParam(lua.P{
	// 	Fn:      subtraction, // 要调用的函数
	// 	NRet:    1,           // 期望返回值的数量
	// 	Protect: true,        // 是否捕获异常
	// }, lua.LNumber(10), lua.LNumber(5)) // 传递参数

	// if err != nil {
	// 	fmt.Println("Error:", err)
	// } else {
	// 	// 获取并打印结果
	// 	result := L.Get(-1)
	// 	fmt.Println("Result of subtraction:", result)
	// 	L.Pop(1) // 清理栈
	// }

}

func tryPreLoadModule() {
	L := lua.NewState()
	defer L.Close()

	path := "./lua/plugin.lua"

	L.PreloadModule("nakama", loader)
	err := L.DoFile(path)
	if err != nil {
		panic(err)
	}
}

func loader(L *lua.LState) int {
	// 注册函数到表中
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"register_rpc": registerRpc,
	})
	// 注册其他东西
	L.SetField(mod, "name", lua.LString("nakama"))

	// 返回模块
	L.Push(mod)
	return 1
}

func registerRpc(L *lua.LState) int {
	rpc := L.ToString(1)
	fmt.Println("registerRpc=", rpc)
	return 1
}

package main

const (
	MutexLocked = 1 << iota
	MutexWoken
	MutexWaiterShift
)

func main() {
	println(MutexLocked)
	println(MutexWoken)
	println(MutexWaiterShift)

	a := 0b1101 // 13
	b := 0b1011 // 11

	c := a & b // 0b1001, 结果为 9
	println("c =", c)

	c = a | b // 0b1111, 结果为 15
	println("c =", c)

	c = a ^ b // 0b0110, 结果为 6
	println("c =", c)

	c = a ^ 0b1111 // 0b0010, 结果为 2
	println("c =", c)

	c = a << 2 // 0b101100, 结果为 52
	println("c =", c)

	c = a >> 2 // 0b0011, 结果为 3
	println("c =", c)

}

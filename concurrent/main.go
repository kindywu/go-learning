package main

const (
	mutexLocked = 1 << iota // mutex is locked
	mutexWoken
	mutexStarving
	mutexWaiterShift      = iota
	starvationThresholdNs = 1e6
)

func main() {
	println(mutexLocked)
	println(mutexWoken)
	println(mutexStarving)
	println(mutexWaiterShift)
	println(starvationThresholdNs)

}

package algorithm

type Comparator int32

const (
	Comparator_Bigger  Comparator = 1
	Comparator_Smaller Comparator = 0
)

// 单调栈
func MonotonicStack(source []int, comparator Comparator) []int {
	answer := make([]int, len(source))
	stack := []int{}

	for i := len(source) - 1; i >= 0; i-- {
		// fmt.Println("stack =", stack)
		fn := func() bool {
			if comparator == Comparator_Bigger {
				return len(stack) > 0 && stack[len(stack)-1] <= source[i]
			} else {
				return len(stack) > 0 && stack[len(stack)-1] >= source[i]
			}
		}
		for fn() {
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			answer[i] = stack[len(stack)-1]
		} else {
			answer[i] = -1
		}
		stack = append(stack, source[i])
	}
	return answer
}

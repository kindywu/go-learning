package algorithm

// 定义栈结构
type Stack struct {
	data []int
}

// 入栈
func (s *Stack) Push(val int) {
	s.data = append(s.data, val)
}

// 出栈
func (s *Stack) Pop() int {
	if len(s.data) == 0 {
		return -1 // 栈为空
	}
	val := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return val
}

// 获取栈顶元素
func (s *Stack) Top() int {
	if len(s.data) == 0 {
		return -1 // 栈为空
	}
	return s.data[len(s.data)-1]
}

// 判断栈是否为空
func (s *Stack) IsEmpty() bool {
	return len(s.data) == 0
}

func LargestRectangleInHistogram(heights []int) int {
	if len(heights) == 0 {
		return 0
	}

	stack := Stack{}
	heights = append(heights, 0) //尾部追加0，能让栈清空

	maxArea := 0
	for i, height := range heights {
		if stack.IsEmpty() || height > heights[stack.Top()] {
			stack.Push(i)
		} else {
			w := 0
			for !stack.IsEmpty() && height <= heights[stack.Top()] {
				h := heights[stack.Top()]
				stack.Pop()

				if !stack.IsEmpty() {
					w = i - stack.Top() - 1
				}
				maxArea = max(maxArea, w*h)
			}
			stack.Push(i)
		}
	}
	return maxArea
}

func LargestRectangleInHistogram_Scan(heights []int) int {
	if len(heights) == 0 {
		return 0
	}
	maxArea := 0
	for i, height := range heights {
		width := 1
		for k := i - 1; k >= 0; k-- {
			if heights[k] >= height {
				width++
			} else {
				break
			}
		}
		for k := i + 1; k < len(heights); k++ {
			if heights[k] >= height {
				width++
			} else {
				break
			}
		}
		// println(i, height, width)
		maxArea = max(maxArea, height*width)
	}
	return maxArea
}

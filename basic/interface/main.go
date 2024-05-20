package main

import "fmt"

// 定义一个接口
type Shape interface {
	Area() float64
	Perimeter() float64
	Draw() string
}

// 定义一个结构体，并实现 Shape 接口
type Circle struct {
	Shape
	Radius float64
}

// 实现接口的方法
func (c *Circle) Area() float64 {
	return 3.14 * c.Radius * c.Radius
}

func (c *Circle) Perimeter() float64 {
	return 2 * 3.14 * c.Radius
}

func (c *Circle) Draw() string {
	return fmt.Sprintf("Drawing a circle with radius %f", c.Radius)
}

// 另一个结构体实现 Shape 接口
type Rectangle struct {
	Width  float64
	Height float64
}

// 实现接口的方法
func (r *Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r *Rectangle) Perimeter() float64 {
	return 2*r.Width + 2*r.Height
}

func (r *Rectangle) Draw() string {
	return fmt.Sprintf("Drawing a rectangle with width %f and height %f", r.Width, r.Height)
}

// 定义一个函数，它接受一个 Shape 接口作为参数
func printShapeInfo(s Shape) {
	fmt.Println(s.Draw())
	fmt.Printf("Area: %f, Perimeter: %f\n", s.Area(), s.Perimeter())
}

func main() {
	// 创建 Circle 和 Rectangle 的实例
	circle := &Circle{Radius: 5}
	rectangle := &Rectangle{Width: 10, Height: 5}

	// 通过接口变量调用方法
	printShapeInfo(circle)
	printShapeInfo(rectangle)
}

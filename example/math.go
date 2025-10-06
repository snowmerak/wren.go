package main

//go:generate go run ../wrengen/main.go -dir .

import (
	"errors"
	"math"
)

//wren:bind module=main
type Math struct{}

//wren:bind
func (m *Math) Add(a, b int32) int32 {
	return a + b
}

//wren:bind static
func (m *Math) Multiply(a, b float64) float64 {
	return a * b
}

//wren:bind
func (m *Math) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

//wren:bind module=main class=StringUtils name=concat static
func StringConcat(a, b string) string {
	return a + b
}

//wren:bind module=main class=Utils name=greet static
func Greet(name string) string {
	return "Hello, " + name + "!"
}

//wren:bind module=main class=Calculator name=square static
func Square(x float64) float64 {
	return x * x
}

//wren:bind module=main class=Calculator name=sqrt static
func Sqrt(x float64) float64 {
	return math.Sqrt(x)
}

//wren:bind module=main class=Calculator name=power static
func Power(base, exponent float64) float64 {
	return math.Pow(base, exponent)
}

// Geometry module - demonstrates different module usage
//wren:bind module=geometry class=Circle name=area static
func CircleArea(radius float64) float64 {
	return math.Pi * radius * radius
}

//wren:bind module=geometry class=Circle name=circumference static
func CircleCircumference(radius float64) float64 {
	return 2 * math.Pi * radius
}

//wren:bind module=geometry class=Rectangle name=area static
func RectangleArea(width, height float64) float64 {
	return width * height
}

//wren:bind module=geometry class=Rectangle name=perimeter static
func RectanglePerimeter(width, height float64) float64 {
	return 2 * (width + height)
}

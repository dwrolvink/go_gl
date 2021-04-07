package main

import (
	//"fmt"
	
	//"strings"
	"time"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/dwrolvink/gogl"
)

const (
	Width  = 500
	Height = 500

	VertexShaderSource = "shaders/triangle.vert"
	FragmentShaderSource = "shaders/triangle.frag"
	FragmentShaderSource2 = "shaders/triangle2.frag"
)

var (
	WindowTitle = "Test GL Application"
	Triangle = []float32{
		0, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
	}
	triangleVelocity = float32(0.01)
)

func main() {
	window := gogl.Init(WindowTitle, Width, Height)
	
	// Create shaders, and link them together in a program
	program1, err := gogl.MakeProgram("program1", VertexShaderSource, FragmentShaderSource)
	if err != nil {
		panic(err)
	}
	
	program2, err2 := gogl.MakeProgram("program2", VertexShaderSource, FragmentShaderSource2)
	if err2 != nil {
		panic(err2)
	}	
	

	// Main loop
	for !window.ShouldClose() {
		// Update game data (move the triangle)
		UpdateState(&Triangle, &triangleVelocity) 

		// Use different programs for when the triangle is moving in different directions
		if triangleVelocity > 0 {
			gogl.Draw(window, program1, Triangle, gl.TRIANGLES)
		} else {
			gogl.Draw(window, program2, Triangle, gl.TRIANGLES)
		}
		
		// Check if shaders need to be recompiled
		gogl.HotloadShaders()
		
		// Sleep to control the speed
		time.Sleep(0 * time.Millisecond)
	}

	// useless here, but good to keep track of what needs to be deleted
	defer glfw.Terminate()
}

func UpdateState(data *[]float32, d *float32) {
	// Update triangle
	if (*data)[1] >= 1 {
		(*d) *= -1
	} else if (*data)[1] <= 0 {
		(*d) *= -1
	}
	for i := range (*data) {
		(*data)[i] += (*d)
	}
}

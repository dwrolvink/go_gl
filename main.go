package main

import (
	//"fmt"
	//"log"
	
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
)

func main() {
	window := gogl.Init(WindowTitle, Width, Height)
	
	// Create shaders, and link them together in a program
	program1ID, err := gogl.MakeProgram("program1", VertexShaderSource, FragmentShaderSource)
	if err != nil {
		panic(err)
	}
	
	program2ID, err2 := gogl.MakeProgram("program2", VertexShaderSource, FragmentShaderSource2)
	if err2 != nil {
		panic(err2)
	}	
	

	// Main loop
	triangleVelocity := float32(0.01)
	changedShaderFiles := []string{}


	for !window.ShouldClose() {
		// Update game data (move the triangle)
		UpdateState(&Triangle, &triangleVelocity) 

		// Draw and sleep
		
		if triangleVelocity > 0 {
			gogl.Draw(window, program1ID, Triangle, gl.TRIANGLES)
		} else {
			gogl.Draw(window, program2ID, Triangle, gl.TRIANGLES)
		}
		
		//gogl.Draw(window, program1ID, Triangle, gl.TRIANGLES)

		// Check if shaders need to be recompiled
		changedShaderFiles = gogl.GetChangedShaderFiles()
		if len(changedShaderFiles) > 0 {

			// reload prog1 if necessary
			_progID, err := gogl.ReloadProgram("program1", changedShaderFiles)
			if err == nil {
				program1ID = _progID
			}
			// reload prog2 if necessary
			_progID, err = gogl.ReloadProgram("program2", changedShaderFiles)
			if err == nil {
				program2ID = _progID
			}	
		}
		
		
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

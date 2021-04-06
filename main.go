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

	VertexShaderSource = `
		#version 450
		in vec3 vp;
		void main() {
			gl_Position = vec4(vp, 1.0);
		}
	`

	FragmentShaderSource = `
		#version 450
		out vec4 frag_colour;
		void main() {
			frag_colour = vec4(0.5, 1, 1, 1.0);
		}
	`

	FragmentShaderSource2 = `
	#version 450
	out vec4 frag_colour;
	void main() {
		frag_colour = vec4(0.5, 0.5, 1, 1.0);
	}
`
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
	
	// Load shaders
	vertexShaderID := gogl.MakeShader(VertexShaderSource, gl.VERTEX_SHADER)
	fragmentShaderID := gogl.MakeShader(FragmentShaderSource, gl.FRAGMENT_SHADER)

	fragmentShader2ID := gogl.MakeShader(FragmentShaderSource2, gl.FRAGMENT_SHADER)

	// Link everything together in a program
	programID := gogl.MakeProgram(vertexShaderID, fragmentShaderID)
	program2ID := gogl.MakeProgram(vertexShaderID, fragmentShader2ID)

	// Main loop
	triangleVelocity := float32(0.01)
	for !window.ShouldClose() {
		// Update game data (move the triangle)
		UpdateState(&Triangle, &triangleVelocity) 

		// Draw and sleep
		if triangleVelocity > 0 {
			gogl.Draw(window, programID, Triangle, gl.TRIANGLES)
		} else {
			gogl.Draw(window, program2ID, Triangle, gl.TRIANGLES)
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

package main

import (
	"log"
	"time"

	"github.com/dwrolvink/gogl"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	Width  = 500
	Height = 500
)

var (
	WindowTitle = "Test GL Application"
)

func main() {
	window := gogl.Init(WindowTitle, Width, Height)

	// Set Data
	// -----------------------------------------------------------
	var datalist [2]gogl.DataObject

	// Quad
	datalist[0] = gogl.DataObject{
		VAOID:    gogl.GenBindVertexArray(),
		Type:     gogl.GOGL_QUADS,
		Vertices: CreateQuadVertexMatrix(0.8, 0.0, 0.0),
		Indices: []uint32{
			1, 0, 3, // triangle 1
			0, 2, 3, // triangle 2
		},
		VertexShaderSource:   "shaders/quad.vert",
		FragmentShaderSource: "shaders/quadtexture.frag",
	}

	// Triangles
	datalist[1] = gogl.DataObject{
		VAOID: gogl.GenBindVertexArray(),
		Type:  gogl.GOGL_TRIANGLES,
		Vertices: []float32{
			0, 0.5, 0,
			-0.5, -0.5, 0,
			0.5, -0.5, 0,
			0.7, 1.0, 0,
			0.3, 0.3, 0,
			1.1, 0.3, 0,
		},
		VertexShaderSource:   "shaders/triangle.vert",
		FragmentShaderSource: "shaders/triangle.frag",
	}

	// Pick one or the other data set
	data := datalist[0]

	// Link program, and bind vertex data to GL
	// -----------------------------------------------------------
	data.ProcessData()

	// Load image to texture (only used in quad)
	// -----------------------------------------------------------
	textureId := gogl.LoadImageToTexture("assets/img/pepe.png")
	gl.BindTexture(gl.TEXTURE_2D, uint32(textureId))

	// Create oscillating value to animate with
	// -----------------------------------------------------------
	var x float32 = 0.0
	var dir float32 = 1

	// Main loop
	for !window.ShouldClose() {

		// Oscillate
		x += 0.01 * dir
		if x > 0.5 || x < -0.5 {
			dir *= -1.0
		}
		data.Program.SetFloat("x", x)

		// Draw to screen
		Draw(window, data)

		// Check if shaders need to be recompiled
		gogl.HotloadShaders()

		// Sleep to control the speed
		time.Sleep(0 * time.Millisecond)

		if err := gl.GetError(); err != 0 {
			log.Println(err)
		}
	}

	// useless here, but good to keep track of what needs to be deleted
	defer glfw.Terminate()
}

func Draw(window *glfw.Window, data gogl.DataObject) {

	// Clear buffer
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Activate program
	gogl.UseProgram((*data.Program).ID)

	// Compile image
	gl.BindVertexArray(uint32(data.VAOID))

	if data.Type == gogl.GOGL_QUADS {
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	} else if data.Type == gogl.GOGL_TRIANGLES {
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(data.Vertices)/3))
	} else {
		panic("data.Type is unknown")
	}

	// Handle window events
	glfw.PollEvents()

	// Put buffer that we painted on on the foreground
	window.SwapBuffers()
}

// Easy way to create a quad with a certain size and offset
func CreateQuadVertexMatrix(size float32, x_offset float32, y_offset float32) []float32 {
	screen_left := -size + x_offset
	screen_bottom := -size + y_offset
	screen_right := size + x_offset
	screen_top := size + y_offset
	texture_top := float32(1.0)
	texture_bottom := float32(0.0)
	texture_left := float32(0.0)
	texture_right := float32(1.0)
	z := float32(0.0)

	vertices := []float32{
		// x, y, z, texcoordx, texcoordy
		screen_left, screen_top, z, texture_left, texture_top,
		screen_right, screen_top, z, texture_right, texture_top,
		screen_left, screen_bottom, z, texture_left, texture_bottom,
		screen_right, screen_bottom, z, texture_right, texture_bottom,
	}

	return vertices
}

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/dwrolvink/gogl"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"image"
	"image/png"

	"github.com/disintegration/imaging"
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
		},
		VertexShaderSource:   "shaders/triangle.vert",
		FragmentShaderSource: "shaders/triangle.frag",
	}

	// Pick one or the other data set
	data := datalist[0]

	// Apply commandline choice for dataset, if present
	for i := range os.Args {
		if os.Args[i] == "-s" {
			if i+1 < len(os.Args) {
				choice, err := strconv.Atoi(os.Args[i+1])
				if err != nil {
					fmt.Println("ERROR: Dataset not passed in. E.g. '-s 1'. Ignoring.")
				}

				data = datalist[choice]
			} else {
				fmt.Println("ERROR: Dataset not passed in. E.g. '-s 1'. Ignoring.")
			}

		}
	}

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
	var tick float32 = -1.0
	var dir float32 = 1

	// Record
	record := false
	var record_length float32 = 50.0

	for i := range os.Args {
		if os.Args[i] == "-r" {
			record = true
			InitRecording()

			// check if record_length has been passed
			if i+1 < len(os.Args) {
				choice, err := strconv.Atoi(os.Args[i+1])
				// arg can be converted to int
				if err == nil {
					record_length = float32(choice)
				}
			}
		}
	}

	// FPS
	delay_ms := int64(20)

	// Main loop
	for !window.ShouldClose() && (!record || tick < record_length) {

		start := time.Now()

		// Oscillate
		tick += 1.0
		x += 0.01 * dir
		if x > 0.5 || x < -0.5 {
			dir *= -1.0
		}
		data.Program.SetFloat("x", x)
		data.Program.SetFloat("y", x)
		data.Program.SetFloat("t", tick)

		// Clear screen
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Setup program
		gogl.UseProgram((*data.Program).ID)
		gl.BindVertexArray(uint32(data.VAOID))

		// Draw to screen
		if data.Type == gogl.GOGL_QUADS {
			gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

		} else if data.Type == gogl.GOGL_TRIANGLES {
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(data.Vertices)/3))
			data.Program.SetFloat("x", -x)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(data.Vertices)/3))
			data.Program.SetFloat("y", -x)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(data.Vertices)/3))
			data.Program.SetFloat("x", x)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(data.Vertices)/3))

		} else {
			panic("data.Type is unknown")
		}

		// Handle window events
		glfw.PollEvents()

		// Put buffer that we painted on on the foreground
		window.SwapBuffers()

		// Check if shaders need to be recompiled
		gogl.HotloadShaders()

		// Sleep to control the speed
		time.Sleep(0 * time.Millisecond)

		if err := gl.GetError(); err != 0 {
			log.Println(err)
		}

		// Record output
		if record {
			CreateImage(int(tick))
			fmt.Println(tick)
		}

		// FPS
		elapsed := time.Since(start)
		dif_ms := delay_ms - elapsed.Milliseconds()
		time.Sleep(time.Duration(dif_ms * int64(time.Millisecond)))
	}

	// Compile gif
	if record {
		CompileGif()
	}

	// useless here, but good to keep track of what needs to be deleted
	defer glfw.Terminate()
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

// RECORDING
// ----------------------------------------------------------------------

func InitRecording() {
	// Ensure that the recording folder is present
	err := os.Mkdir("recording/temp/", 0755)
	if err != nil {
		panic(err)
	}
}
func CreateImage(number int) {
	filename := fmt.Sprintf("image%03d.png", number)
	width := Width
	height := Height

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewNRGBA(image.Rectangle{upLeft, lowRight})

	gl.ReadBuffer(gl.FRONT)
	gl.ReadPixels(0, 0, Width, Height, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	img = imaging.FlipV(img)

	// Encode as PNG.
	f, _ := os.Create("recording/temp/" + filename)
	png.Encode(f, img)
}

func CompileGif() {
	filename := time.Now().Unix()

	cmd, err := exec.Command("/bin/sh", "scripts/make_gif.sh", fmt.Sprint(filename)).Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	fmt.Println(cmd)
}

func RunBash(scriptPath string) string {
	cmd, err := exec.Command("/bin/sh", "scripts/make_gif.sh", "testt").Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	output := string(cmd)
	return output
}

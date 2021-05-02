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
	Width  = 500 // Width of the main window
	Height = 500 // Height of the main window
)

var (
	WindowTitle = "Test GL Application"

	x        float32 = 0.0  // used to move the sprites around
	dir_x    float32 = 1    // used to change x
	tick     float32 = -1.0 // ticks up every game loop cycle
	delay_ms int64   = 20   // handles frame rate

	DrawMode      string  = "composite"                        // chooses whether to draw one dataset, or all of them ("composite" vs "single_set")
	ChosenDataset int     = 0                                  // Used only when DrawMode = "single_dataset"
	record        bool    = false                              // whether to record the screen.
	record_length float32 = float32(1.0 * (1000.0 / delay_ms)) // After how many ticks to stop recording (and close the program)
)

func main() {
	// Init Window, OpenGL, and Data, get user input from commandline
	// -----------------------------------------------------------
	window := gogl.Init(WindowTitle, Width, Height)
	data, datalist := SetData()
	_ = datalist

	ParseCommandlineArgs()

	// Main loop
	// ===========================================================
	for !window.ShouldClose() && (!record || tick < record_length) {

		// Update game
		// ------------------------------------------------------
		// Naive way to manage FPS. See also bottom of this loop.
		start := time.Now()

		// Increment global clock
		tick += 1.0

		// x is used temporarily to move stuff around, will be removed when
		// there are actors with volition
		x += 0.01 * dir_x
		if x > 1.0 || x < -1.0 {
			dir_x *= -1.0
		}

		// Update DataObjects
		for i := range datalist {
			datalist[i].Update()
		}

		// Draw to screen
		// ------------------------------------------------------
		// Clear screen
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Draw new frame
		if DrawMode == "composite" {
			DrawComposite(datalist)
		}
		if DrawMode == "single_set" {
			DrawDataset(data)
		}

		// Put buffer that we painted on on the foreground
		window.SwapBuffers()

		// Event handling
		// ------------------------------------------------------
		// Handle window events
		glfw.PollEvents()

		// Check if shaders need to be recompiled
		gogl.HotloadShaders()

		if err := gl.GetError(); err != 0 {
			log.Println(err)
		}

		// Record output
		// ------------------------------------------------------
		if record {
			CreateImage(int(tick))
			fmt.Println(tick)
		}

		// FPS management
		// ------------------------------------------------------
		// Sleep for a bit if the loop finished too quickly.
		// A better way would be to update actor positions based on
		// elapsed time, but the neccessary code isn't present yet for
		// that (i.e. volition).
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

// Define the DataObjects that contain our Programs, Shaders, Sprites, etc
func SetData() (gogl.DataObject, []gogl.DataObject) {
	/*
		   Multiple datasets can be defined.
		   Each set contains all that it needs to draw to the screen,
		   think of: Program, VOA, VBO, EBO, Textures, Sprites, etc

		   Below, each dataset is defined, added to datalist, and at
		   the end the commandline args are checked what the choice is.
		   Choices include:
		     - Print either dataset 0, or 1 --> -s 0, -s 1
			 - Print both as a composite --> -s c
	*/

	// List of datasets
	datalist := make([]gogl.DataObject, 2)

	// Fist dataset: Vertex type: Quad, uses Sprites
	// -----------------------------------------------------------
	datalist[0] = gogl.DataObject{
		ProgramName: "DancingPepe",
		Type:        gogl.GOGL_QUADS,
		Vertices:    CreateQuadVertexMatrix(1.0, 0.0, 0.0),
		Indices: []uint32{
			1, 0, 3, // triangle 1
			0, 2, 3, // triangle 2
		},
		VertexShaderSource:   "shaders/quad.vert",
		FragmentShaderSource: "shaders/quadtexture.frag",
	}

	// load sprites and add to sprite list
	datalist[0].AddSprite(gogl.Sprite{
		Name:           "DancingPepe",
		TextureSource:  "assets/img/texture.png",
		Divisions:      4,
		AnimationSpeed: 5,
		AnimationFrames: [][]float32{
			{0, 0},
			{1, 0},
			{2, 0},
			{3, 0},
			{0, 1},
			{1, 1},
		},
		FlipHorizontal: 0.0,
	})

	datalist[0].AddSprite(gogl.Sprite{
		Name:           "Walking Blob",
		TextureSource:  "assets/img/texture.png",
		Divisions:      8,
		AnimationSpeed: 10,
		AnimationFrames: [][]float32{
			{2, 4},
			{3, 4},
			{4, 4},
			{5, 4},
			{6, 4},
			{7, 4},
		},
		FlipHorizontal: 1.0,
		Yn:             1.0,
		Scale:          0.16,
	})

	// Second dataset: Vertex type: Simple triangles
	// -----------------------------------------------------------
	datalist[1] = gogl.DataObject{
		ProgramName: "DiscoTriangles",
		Type:        gogl.GOGL_TRIANGLES,
		Vertices: []float32{
			0, 0.5, 0,
			-0.5, -0.5, 0,
			0.5, -0.5, 0,
		},
		VertexShaderSource:   "shaders/triangle.vert",
		FragmentShaderSource: "shaders/triangle.frag",
	}

	// Link program, and bind vertex data to GL
	// -----------------------------------------------------------
	datalist[0].ProcessData()
	datalist[1].ProcessData()

	// Pick one or the other data set
	// -----------------------------------------------------------
	// When DrawMode = 'composite', this line will be ignored later
	data := datalist[ChosenDataset]

	return data, datalist
}

// DRAWING
// ----------------------------------------------------------------------

// Draw a single dataset.
func DrawDataset(data gogl.DataObject) {
	data.Enable()

	// load uniforms
	data.Program.SetFloat("x", x)
	data.Program.SetFloat("y", x)
	data.Program.SetFloat("scale", 1.0)
	data.Program.SetFloat("t", tick)

	if data.Type == gogl.GOGL_QUADS {

		// load sprite
		sprite := data.SelectSprite(0)
		sprite.SetUniforms(&data)

		// Draw pepe 1
		data.Program.SetFloat("scale", 0.5)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

		// Draw pepe 2
		data.Program.SetFloat("scale", 0.25)
		data.Program.SetFloat("x", -x)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

		// Draw pepe 3
		data.Program.SetFloat("scale", 0.16)
		data.Program.SetFloat("x", -x)
		data.Program.SetFloat("y", -x)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

	} else if data.Type == gogl.GOGL_TRIANGLES {
		// draw triangle 1
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(data.Vertices)/3))

		// draw triangle 2
		data.Program.SetFloat("x", -x)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(data.Vertices)/3))

		// draw triangle 3
		data.Program.SetFloat("y", -x)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(data.Vertices)/3))

		// draw triangle 4
		data.Program.SetFloat("x", x)
		data.Program.SetFloat("z", 2.0)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(data.Vertices)/3))

	} else {
		panic("data.Type is unknown")
	}
}

// Custom: draw using both datasets
func DrawComposite(datalist []gogl.DataObject) {
	// Composite - triangle
	// --------------------------------------------
	data := datalist[1]
	data.Enable()
	data.Program.SetFloat("x", x*0.5)
	data.Program.SetFloat("y", x*0.5)
	data.Program.SetFloat("scale", 1.5)
	data.Program.SetFloat("t", tick)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(data.Vertices)/3))
	data.Program.SetFloat("x", -x*0.5)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(data.Vertices)/3))
	data.Program.SetFloat("y", -x*0.5)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(data.Vertices)/3))
	data.Program.SetFloat("x", x*0.5)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(data.Vertices)/3))

	// Composite - pepe
	// --------------------------------------------
	data = datalist[0]
	data.Enable()
	data.Program.SetFloat("t", tick)
	data.Program.SetFloat("x", x)
	data.Program.SetFloat("y", x)
	data.Program.SetFloat("scale", 1.0)

	// load sprite 0.1 (pepe big)
	sprite := data.SelectSprite(0)
	sprite.Xn = x
	sprite.Yn = x
	sprite.Scale = 0.5
	sprite.SetUniforms(&data)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

	// load sprite 0.2 (pepe small)
	sprite.Scale = 0.25
	sprite.Xn = -x
	sprite.SetUniforms(&data)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

	// load sprite 1 (walking blob)
	sprite = data.SelectSprite(1)
	sprite.Xn = x
	if dir_x < 0.0 {
		sprite.FlipHorizontal = 0.0
	} else {
		sprite.FlipHorizontal = 1.0
	}
	sprite.SetUniforms(&data)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
}

// HELPER FUNCTIONS
// ----------------------------------------------------------------------

// Applies commandline args: --fps <N>, --record <N>, --set <'c', N>
func ParseCommandlineArgs() {

	for i := range os.Args {

		// FPS
		// -----------------------------------------------------------
		// Apply commandline choice for fps, if present.
		// Note that for recording, 50 fps is the max.

		if os.Args[i] == "--fps" {
			// check if fps value has been passed directly after --fps
			if i+1 < len(os.Args) {
				choice, err := strconv.Atoi(os.Args[i+1])
				if err == nil {
					delay_ms = int64(1000 / choice)
				} else {
					fmt.Println("ERROR: Could not parse input after --fps as an int.")
				}
			}
		}

		// Record
		// -----------------------------------------------------------
		// Apply commandline choice for recording settings, if present

		if os.Args[i] == "--record" {
			// Enable recording
			record = true
			InitRecording()

			// Check if record_length has been passed
			if i+1 < len(os.Args) {
				choice, err := strconv.Atoi(os.Args[i+1]) // convert string input to int
				if err == nil {
					record_length = float32(int64(choice) * (1000.0 / delay_ms)) // input is in seconds, convert to ticks
				}
			}
		}

		// DrawMode & ChosenDataset
		// -----------------------------------------------------------
		// Apply commandline choice for dataset

		if os.Args[i] == "--set" {
			if i+1 < len(os.Args) {

				// Print both datasets on top of eachother
				if os.Args[i+1] == "c" {
					DrawMode = "composite"
					continue
				}

				// Print only one dataset
				DrawMode = "single_set"
				choice, err := strconv.Atoi(os.Args[i+1])
				if err != nil {
					fmt.Println("ERROR: Dataset index not passed in. E.g. '-s 1'. Ignoring.")
					continue
				}
				ChosenDataset = choice

			} else {
				fmt.Println("ERROR: Dataset index not passed in. E.g. '-s 1'. Ignoring.")
			}
		}
	}
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

// Ensures that the recording folder (recording/temp) is present.
func InitRecording() {
	err := os.Mkdir("recording/temp/", 0755)
	if err != nil {
		panic(err)
	}
}

// Reads out the pixel data in gl.FRONT, and saves it to recording/temp/image<Tick>.png
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

// Takes all the frame images in recording/temp and makes a palletted gif out of it using ffmpeg.
func CompileGif() {
	filename := time.Now().Unix()

	cmd, err := exec.Command("/bin/sh", "scripts/make_gif.sh", fmt.Sprint(filename)).Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	fmt.Println(cmd)
}

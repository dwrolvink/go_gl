# go_gl
This is an example app for the [go-gl wrapper package](https://github.com/dwrolvink/gogl/tree/main) that I'm writing as an exercise.

- It draws a multiple triangles on a screen, which moves around and changes color. (Dataset 1)
- It draws a picture to the screen (Dataset 0; default). Example of using quads instead of triangle vertices.
- It can record gifs of what happens on the screen (not very performant).

I have built this package on linux with go 1.16, it's not guaranteed to work on any other system, but it probably will. Let me know if you have troubles.
- Recording will not work on non linux systems as of yet!
- Recording requires ffmpeg to be installed.

See https://www.youtube.com/watch?v=EJz71vpNhSU&list=PLDZujg-VgQlZUy1iCqBbe5faZLMkA3g2x&index=42 for the lecture this package is based on.


## Good to know
The gl packages don't have module support, so we need to disable the module system to be able to use them.

We can set this setting permanently by executing:
``` bash
go env -w GO111MODULE=auto
```
This will not apply module mode when the code is located in $GOPATH/src, and no go.mod file is present.

> On Linux, if $GOPATH is empty. Packages are stored in /home/{user}/go/src/. Modules are stored under /home/{user}/go/pkg/mod

If you are working outside of the package folder, you can set GO111MODULE=off:
``` bash 
go env -w GO111MODULE=off
```

Alternatively, you can pass the setting as a first argument:
```bash
GO111MODULE=off go get github.com/go-gl/gl/v4.5-core/gl
GO111MODULE=off go get github.com/go-gl/glfw/v3.2/glfw
```

# Installation
```bash
# Disable Module mode
go env -w GO111MODULE=off

# Download packages
go get github.com/go-gl/gl/v4.1-core/gl
go get github.com/go-gl/glfw/v3.2/glfw
go get github.com/dwrolvink/gogl
```

# Run
```bash
go run ./main.go 
```

# Build
``` bash
# Compile
go build -v

# Run
./go_gl
```

## Args
Choose different dataset (e.g. triangles instead of pepe):
```bash
go run ./main.go -s 1
```


Record gif to `record/output/`. Default length is 1 second (50 ticks @ 20ms per tick).

> Recording requires ffmpeg to be installed.

```bash
go run ./main.go --record
```

Record gif to `record/output/`. Change length to 2 seconds.
```bash
go run ./main.go --record 2
```

Change the speed. Note that fps of 50 is the max that will work with gifs.
```bash
go run ./main.go --fps 20
```

Combinations allowed:
```bash
go run ./main.go -r 100 -s 1
```
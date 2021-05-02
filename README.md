go_gl
====================
This is an example app for the [go-gl wrapper package](https://github.com/dwrolvink/gogl/tree/main) that I'm writing as an exercise.

See also https://www.youtube.com/watch?v=EJz71vpNhSU&list=PLDZujg-VgQlZUy1iCqBbe5faZLMkA3g2x&index=42 for the lecture this package is based on.

Some features:
--------------------
- All that is needed to draw on the screen is contained in `DataObjects`, and the code allows to switch between DataObjects on the fly, and within a draw cycle.
- One DataObject uses plain triangle vertices, the other uses quads.
- One DataObject is really simple by design (one Triangle that is drawn 4 times).
- One DataObject implements a naive system for drawing Sprites, and animating them.
- Frame speed is somewhat configurable.
- It can record a gif of what happens on the screen (not very performant).

Support 
--------------------
I have built this package on linux with go 1.16, it's not guaranteed to work on any other system, but it probably will. Let me know if you have troubles.
- Recording will not work on non linux systems as of yet!
- Recording requires ffmpeg to be installed.




Good to know
--------------------
The gl packages don't have module support (as of 2021/04), so we need to disable the module system to be able to use them.

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

Installation
====================
```bash
# Disable Module mode
go env -w GO111MODULE=off

# Download packages
go get github.com/go-gl/gl/v4.1-core/gl
go get github.com/go-gl/glfw/v3.2/glfw
go get github.com/dwrolvink/gogl
```

Run
====================
```bash
go run ./main.go 
```

Build
====================
``` bash
# Compile
go build -v

# Run
./go_gl
```

Args
====================
--set
-------------------
Choose different dataset (e.g. triangles instead of pepe). Valid values: c, 0, 1, (and further if you add more datasets).
"c" is the default, and this is a custom composite drawmode that uses both dataset 0 and 1.
```bash
go run ./main.go --set 1
```

--record
-------------------
Record gif to `record/output/`. Default length is 1 second (50 ticks @ 20ms per tick).

> Recording requires ffmpeg to be installed.

```bash
go run ./main.go --record
```

Record gif to `record/output/`. Change length to 2 seconds.
```bash
go run ./main.go --record 2
```

--fps
-------------------
Change the speed. Note that fps of 50 is the max that will work with gifs.
```bash
go run ./main.go --fps 20
```

Combinations allowed:
-------------------
```bash
go run ./main.go --record 2 --set 1 --fps 30
```

Dev Notes
====================
VS Code Plugins
--------------------
When working with Go, the following plugins have been tremendously helpful:
- the `Go` plugin, by *Go Team at Google*
- the `Go Doc` plugin, by *Minhaz Ahmed Syrus*

Especially the last one, which allows you to hover over a method and see the description, is one that I use profusely.
If you are staring at the code, not sure what is happening, install the plugin, and hover over some methods.
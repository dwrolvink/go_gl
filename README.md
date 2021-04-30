# go_gl
This is an example app for the [go-gl wrapper package](https://github.com/dwrolvink/gogl/tree/main) that I'm writing as an exercise.
It draws a triangle on a screen, which moves around and changes color.

I have built this package on linux with go 1.16, it's not guaranteed to work on any other system, but it probably will. Let me know if you have troubles.

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
./gl
```

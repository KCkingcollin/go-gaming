package main

import (
	"fmt"
	"time"

	"github.com/KCkingcollin/go-gaming/bin/vars"
	"github.com/KCkingcollin/go-gaming/bin/inputs"
	"github.com/KCkingcollin/go-gaming/bin/spaces"
	"github.com/KCkingcollin/go-help-func/ghf"
	"github.com/KCkingcollin/go-help-func/glf"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/veandco/go-sdl2/sdl"
)

var (
    err error
    window = vars.Window
    shaderProg1 = vars.ShaderProg1
    camera = vars.Camera
)

func main() {
    initWindow()
    initBuffers()
    defer sdl.Quit()
    defer window.Destroy()
    for {
        frameStart := time.Now()
        if inputs.PollEvents() {
            return
        }
        inputs.CameraEvents()
        frameRendering()

        // Frame rate limiter
        for time.Since(frameStart).Nanoseconds() < int64(float64(vars.TimeFactor / vars.FrameRateLimit) * 0.999) {}

        // FPS Counter
        frameTime := time.Since(frameStart).Nanoseconds()
        if time.Since(vars.TimeCount).Nanoseconds() >= vars.TimeFactor {
            frameCount := vars.TimeFactor / frameTime
            println(frameCount)
            vars.TimeCount = time.Now()
        }
        vars.ElapsedTime = time.Since(frameStart) // elapsed time in ns
    }
} 

func initWindow() {
    if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
        panic(err)
    }

    sdl.SetHint(sdl.HINT_VIDEO_WAYLAND_PREFER_LIBDECOR, "1")
    sdl.SetHint(sdl.HINT_VIDEO_X11_NET_WM_BYPASS_COMPOSITOR, "0")

    window, err = sdl.CreateWindow("GO Gaming Engine",
    sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
    vars.WinWidth, vars.WinHeight,
    sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE)
    if err != nil {
        panic(err)
    }

    if _, err = window.GLCreateContext(); err != nil {
        panic(err)
    }

    // Vsync
    sdl.GLSetSwapInterval(1)

    if err := gl.Init();
    err != nil {
        panic(err)
    }
    gl.Enable(gl.DEPTH_TEST)
    glf.PrintVersionGL()

    vars.Vertices, vars.Texture, vars.VAO  = spaces.LocalSpace()
    vars.Positions = spaces.WorldSpace()
}

func initBuffers() {
    shaderProg1, err = glf.NewShaderProgram(vars.VertPath, vars.FragPath)
    if err != nil && ghf.Verbose {
        fmt.Printf("Failed to link shaders: %s \n", err)
    } else if ghf.Verbose {
        println("Program linked successfully")
    }

    vars.UBO1 = glf.GenBindBuffers(gl.UNIFORM_BUFFER)
    glf.BufferData(gl.UNIFORM_BUFFER, make([]float32, 3*16), gl.DYNAMIC_DRAW)
    gl.BindBufferBase(gl.UNIFORM_BUFFER, 1, vars.UBO1)
}

func frameRendering() {
    gl.ClearColor(0.1, 0.1, 0.1, 1.0)
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

    shaderProg1.Use()
    projectionMatrix := mgl64.Perspective(mgl64.DegToRad(45.0), vars.DisplayRatio, 0.1, 100.0)
    viewMatrix := camera.GetViewMatrix()
    vars.Mat4s[2] = projectionMatrix
    vars.Mat4s[1] = viewMatrix

    glf.BindTexture(vars.Texture)
    gl.BindVertexArray(vars.VAO)
    for i, pos := range vars.Positions {
        modelMatrix := mgl64.Ident4()
        angle := 20.0 * float64(i)
        modelMatrix = mgl64.HomogRotate3D(mgl64.DegToRad(angle), mgl64.Vec3{1.0, 0.3, 0.5}).Mul4(modelMatrix)
        modelMatrix = mgl64.Translate3D(pos.X(), pos.Y(), pos.Z()).Mul4(modelMatrix)
        vars.Mat4s[0] = modelMatrix

        glf.BindBufferSubDataMat4(ghf.Mgl64to32Mat4Slice(vars.Mat4s), vars.UBO1)

        gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vars.Vertices)/5*3))
    }
    window.GLSwap()
    glf.CheckShadersforChanges()
}


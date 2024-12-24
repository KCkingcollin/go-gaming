package bin

import (
	"fmt"
	"time"

	"github.com/KCkingcollin/go-help-func/ghf"
	"github.com/KCkingcollin/go-help-func/glf"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/veandco/go-sdl2/sdl"
)

//Init sdl and gl 
func InitWindow() {
    if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
        panic(err)
    }

    sdl.SetHint(sdl.HINT_VIDEO_WAYLAND_PREFER_LIBDECOR, "1")
    sdl.SetHint(sdl.HINT_VIDEO_X11_NET_WM_BYPASS_COMPOSITOR, "0")

    Window, err = sdl.CreateWindow("GO Gaming Engine",
    sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
    WinWidth, WinHeight,
    sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE)
    if err != nil {
        panic(err)
    }

    if _, err = Window.GLCreateContext(); err != nil {
        panic(err)
    }

    sdl.GLSetSwapInterval(1) //Vsync

    if err := gl.Init();
    err != nil {
        panic(err)
    }
    gl.Enable(gl.DEPTH_TEST)
    glf.PrintVersionGL()
}

//Init glsl buffers
func InitBuffers() {
    Vertices, Texture, VAO = LocalSpace()
    Positions = WorldSpace()

    ShaderProg1, err = glf.NewShaderProgram(VertPath, FragPath)
    if err != nil && ghf.Verbose {
        fmt.Printf("Failed to link Shaders: %s \n", err)
    } else if ghf.Verbose {
        println("Program linked successfully")
    }

    UBO1 = glf.GenBindBuffers(gl.UNIFORM_BUFFER)
    glf.BufferData(gl.UNIFORM_BUFFER, make([]float32, 3*16), gl.DYNAMIC_DRAW)
    gl.BindBufferBase(gl.UNIFORM_BUFFER, 1, UBO1)
}

//Render the frame
func frameRendering() {
    gl.ClearColor(0.1, 0.1, 0.1, 1.0)
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

    ShaderProg1.Use()
    projectionMatrix := mgl64.Perspective(mgl64.DegToRad(45.0), DisplayRatio, 0.1, 100.0)
    viewMatrix := Camera.GetViewMatrix()
    Mat4s[2] = projectionMatrix
    Mat4s[1] = viewMatrix

    glf.BindTexture(Texture)
    gl.BindVertexArray(VAO)
    for i, pos := range Positions {
        modelMatrix := mgl64.Ident4()
        angle := 20.0 * float64(i)
        modelMatrix = mgl64.HomogRotate3D(mgl64.DegToRad(angle), mgl64.Vec3{1.0, 0.3, 0.5}).Mul4(modelMatrix)
        modelMatrix = mgl64.Translate3D(pos.X(), pos.Y(), pos.Z()).Mul4(modelMatrix)
        Mat4s[0] = modelMatrix

        glf.BindBufferSubDataMat4(ghf.Mgl64to32Mat4Slice(Mat4s), UBO1)

        gl.DrawArrays(gl.TRIANGLES, 0, int32(len(Vertices)/5*3))
    }
    Window.GLSwap()
    glf.CheckShadersforChanges()
}

//Start main game loop
func MainLoop() {
    for {
        frameStart := time.Now()
        if PollEvents() {
            return
        }
        CameraEvents()
        frameRendering()

        for time.Since(frameStart).Nanoseconds() < int64(float64(TimeFactor / FrameRateLimit) * 0.999) {}

        frameTime := time.Since(frameStart).Nanoseconds()
        if time.Since(TimeCount).Nanoseconds() >= TimeFactor {
            frameCount := TimeFactor / frameTime
            println(frameCount)
            TimeCount = time.Now()
        }
        ElapsedTime = time.Since(frameStart)
    }
}


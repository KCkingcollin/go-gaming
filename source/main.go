package source

import (
	"fmt"
	"time"

	"github.com/KCkingcollin/go-help-func/ghf"
	"github.com/KCkingcollin/go-help-func/glf"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/veandco/go-sdl2/sdl"
)

func Main() {
    LoadConfig(ConfigFile)
    InitWindow()
    InitBuffers()
    go checkConfig()
    defer sdl.Quit()
    defer Window.Destroy()
    // Main game loop
    for {
        frameStart := time.Now()
        if PollEvents() {
            return
        }
        CameraEvents()
        frameRendering()

        for time.Since(frameStart).Nanoseconds() < int64(float64(TimeFactor / FrameRateLimit) * 0.999 * FrameSpeedFactor) {}

        frameTime := time.Since(frameStart).Nanoseconds()
        if time.Since(Time).Nanoseconds() - SecondTime >= TimeFactor {
            frameCount := TimeFactor / frameTime
            println(frameCount)
            SecondTime = time.Since(Time).Nanoseconds()
        }
        ElapsedTime = time.Since(frameStart)
    }
}

// Checks for config updates
func checkConfig() {
    for {
        time.Sleep(1 * time.Second)
        if ghf.FileExists(ConfigFile) {
            configModTimeNow := ghf.GetModifiedTime(ConfigFile)
            if !ConfigModTime.Equal(configModTimeNow) {
                LoadConfig(ConfigFile)
            }
        }
    }
}

// Init sdl and gl 
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

// Init glsl buffers
func InitBuffers() {
    Vertices, Normals, Texture = LocalSpace()
    Positions = WorldSpace()

    ShaderProg1, err = glf.NewShaderProgram(VertPath, FragPath)
    if err != nil && ghf.Verbose {
        fmt.Printf("Failed to link Shaders: %s \n", err)
    } else if ghf.Verbose {
        println("Program linked successfully")
    }

    VAO = glf.GenBindVertexArrays()
    glf.GenBindBuffers(gl.ARRAY_BUFFER)
    glf.BufferData(gl.ARRAY_BUFFER, Vertices, gl.STATIC_DRAW)

    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
    gl.EnableVertexAttribArray(0)
    gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
    gl.EnableVertexAttribArray(1)

    glf.GenBindBuffers(gl.ARRAY_BUFFER)
    glf.BufferData(gl.ARRAY_BUFFER, Normals, gl.STATIC_DRAW)
    gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 3*4, nil)
    gl.EnableVertexAttribArray(2)
    gl.BindVertexArray(0)
    
    UBO0 = glf.CreateNewUniformBuffer(UBMat4s, 0)
    UBO1 = glf.CreateNewUniformBuffer(UBVec3s, 1)

    UBVec3s = []mgl64.Vec3 {
        {2.0, 5.0, 5.0},  // lightPos
        {1.0, 1.0, 1.0},   // lightColor
    }
}

// Render the frame
func frameRendering() {
    gl.ClearColor(0.1, 0.1, 0.1, 1.0)
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

    ShaderProg1.Use()
    projectionMatrix := mgl64.Perspective(mgl64.DegToRad(45.0), DisplayRatio, 0.1, 100.0)
    viewMatrix := Camera.GetViewMatrix()
    UBMat4s[2] = projectionMatrix
    UBMat4s[1] = viewMatrix

    glf.BindTexture(Texture)
    gl.BindVertexArray(VAO)
    for i, pos := range Positions {
        modelMatrix := mgl64.Ident4()
        angle := 20.0 * float64(i)
        modelMatrix = mgl64.HomogRotate3D(mgl64.DegToRad(angle), mgl64.Vec3{1.0, 0.3, 0.5}).Mul4(modelMatrix)
        modelMatrix = mgl64.Translate3D(pos.X(), pos.Y(), pos.Z()).Mul4(modelMatrix)
        UBMat4s[0] = modelMatrix

        glf.SetUBO(UBMat4s, UBO0)

        gl.DrawArrays(gl.TRIANGLES, 0, int32(len(Vertices)/5*3))
    }

    glf.SetUBO(UBVec3s, UBO1)

    Window.GLSwap()
    glf.CheckShadersforChanges()
}


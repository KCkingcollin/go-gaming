package main

import (
    "fmt"
    "time"

    spaces "github.com/KCkingcollin/go-gaming/bin"
    "github.com/KCkingcollin/go-help-func/ghf"
    "github.com/KCkingcollin/go-help-func/glf"
    "github.com/go-gl/gl/v4.6-core/gl"
    "github.com/go-gl/mathgl/mgl64"
    "github.com/veandco/go-sdl2/sdl"
)

var (
    err                 error

    inMenu              bool            = true

    keyboardState       []uint8         = sdl.GetKeyboardState()

    UBO1                uint32
    VAO                 uint32
    texture             uint32

    winWidth, winHeight int32           = 854, 480
    mouseX, mouseY      int32

    vertices            []float32

    timeFactor          int64           = int64(time.Second/time.Nanosecond)
    frameRateLimit      int64           = 500
    prevFrameLimit      int64           = frameRateLimit

    frameCount          int

    displayRatio        float64         = float64(winWidth / winHeight)
    deltaT              float64

    elapsedTime         time.Duration

    limiterStart           time.Time       = time.Now()
    timeCount           time.Time       = time.Now()
    frameStart          time.Time      // the time you started generating the current frame

    camPos              mgl64.Vec3      = mgl64.Vec3{0.0, 0.0, 0.0}
    worldUp             mgl64.Vec3      = mgl64.Vec3{0.0, 1.0, 0.0}
    positions           []mgl64.Vec3                     

    Mat4s               []mgl64.Mat4    = make([]mgl64.Mat4, 3)

    camera              *glf.Camera     = glf.NewCamera(camPos, worldUp, -90.0, 0.0, 0.02, 0.2)

    window              *sdl.Window

    shaderProg1         *glf.ShaderInfo

)

const (
    fragPath = "./shaders/quadtexture.frag.glsl"
    vertPath = "./shaders/main.vert.glsl"
)

func main() {
    initWindow()
    initBuffers()
    defer sdl.Quit()
    defer window.Destroy()
    for {
        frameStart = time.Now()
        if pollEvents() {
            return
        }
        cameraEvents()
        frameRendering()

        elapsedTime = time.Since(frameStart) // elapsed time in ns

        // Frame rate limiter
        for time.Since(frameStart).Nanoseconds() < 1e9 / int64(frameRateLimit) {}

        // FPS Counter
        frameCount++
        if time.Since(timeCount).Nanoseconds() >= timeFactor {
            println(frameCount)
            timeCount = time.Now()
            frameCount = 0
        }
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
    winWidth, winHeight,
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

    vertices, texture, VAO  = spaces.LocalSpace()
    positions = spaces.WorldSpace()
}

func initBuffers() {
    shaderProg1, err = glf.NewShaderProgram(vertPath, fragPath)
    if err != nil && ghf.Verbose {
        fmt.Printf("Failed to link Shaders: %s \n", err)
    } else if ghf.Verbose {
        println("Program linked successfully")
    }

    UBO1 = glf.GenBindBuffers(gl.UNIFORM_BUFFER)
    glf.BufferData(gl.UNIFORM_BUFFER, make([]float32, 3*16), gl.DYNAMIC_DRAW)
    gl.BindBufferBase(gl.UNIFORM_BUFFER, 1, UBO1)
}

func pollEvents() bool {
    mouseX, mouseY = 0, 0
    for event := sdl.PollEvent(); event != nil;  event = sdl.PollEvent() {
        switch event := event.(type) {
        case *sdl.QuitEvent:
            return true
        case *sdl.WindowEvent:
            switch event.Event {
            case sdl.WINDOWEVENT_SIZE_CHANGED:
                winWidth, winHeight = window.GetSize()
                gl.Viewport(0, 0, winWidth, winHeight)
                displayRatio = float64(winWidth) / float64(winHeight)
            case sdl.WINDOWEVENT_FOCUS_GAINED:
                sdl.GLSetSwapInterval(0)
                frameRateLimit = prevFrameLimit
            case sdl.WINDOWEVENT_FOCUS_LOST:
                sdl.GLSetSwapInterval(1)
                prevFrameLimit = frameRateLimit
                frameRateLimit = 5
            }
        case *sdl.KeyboardEvent:
            switch event.Keysym.Sym {
            case sdl.K_ESCAPE:
                keys := make([]bool, 256)
                if event.Type == sdl.KEYDOWN {
                    keys[event.Keysym.Sym] = true
                    inMenu = !inMenu
                    if inMenu {
                        sdl.SetRelativeMouseMode(false)
                        println("mouse disabled")
                    } else {
                        sdl.SetRelativeMouseMode(true)
                        println("mouse enabled")
                    }
                } else if event.Type == sdl.KEYUP {
                    keys[event.Keysym.Sym] = false
                }
            }
        case *sdl.MouseMotionEvent:
            if !inMenu {
                mouseX += event.XRel
                mouseY += event.YRel
            }
        }
    }
    return false
}

func cameraEvents() {
    dir := glf.NoWhere
    deltaT = float64(elapsedTime.Milliseconds())
    if keyboardState[sdl.SCANCODE_LEFT] != 0 || keyboardState[sdl.SCANCODE_A] != 0 {
        dir = glf.Left
        camera.UpdateCamera(dir, deltaT, float64(mouseX), float64(mouseY))
    }
    if keyboardState[sdl.SCANCODE_RIGHT] != 0 || keyboardState[sdl.SCANCODE_D] != 0 {
        dir = glf.Right
        camera.UpdateCamera(dir, deltaT, float64(mouseX), float64(mouseY))
    }
    if keyboardState[sdl.SCANCODE_UP] != 0 || keyboardState[sdl.SCANCODE_W] != 0 {
        dir = glf.Forward
        camera.UpdateCamera(dir, deltaT, float64(mouseX), float64(mouseY))
    }
    if keyboardState[sdl.SCANCODE_DOWN] != 0 || keyboardState[sdl.SCANCODE_S] != 0 {
        dir = glf.Backward
        camera.UpdateCamera(dir, deltaT, float64(mouseX), float64(mouseY))
    }
    camera.UpdateCamera(dir, deltaT, float64(mouseX), float64(mouseY))
}

func frameRendering() {
    gl.ClearColor(0.1, 0.1, 0.1, 1.0)
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

    shaderProg1.Use()
    projectionMatrix := mgl64.Perspective(mgl64.DegToRad(45.0), displayRatio, 0.1, 100.0)
    viewMatrix := camera.GetViewMatrix()
    Mat4s[2] = projectionMatrix
    Mat4s[1] = viewMatrix

    glf.BindTexture(texture)
    gl.BindVertexArray(VAO)
    for i, pos := range positions {
        modelMatrix := mgl64.Ident4()
        angle := 20.0 * float64(i)
        modelMatrix = mgl64.HomogRotate3D(mgl64.DegToRad(angle), mgl64.Vec3{1.0, 0.3, 0.5}).Mul4(modelMatrix)
        modelMatrix = mgl64.Translate3D(pos.X(), pos.Y(), pos.Z()).Mul4(modelMatrix)
        Mat4s[0] = modelMatrix

        glf.BindBufferSubDataMat4(ghf.Mgl64to32Mat4Slice(Mat4s), UBO1)

        gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)/5*3))
    }
    window.GLSwap()
    glf.CheckShadersforChanges()
}


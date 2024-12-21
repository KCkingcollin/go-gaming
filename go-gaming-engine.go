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

var winWidth int32 = 854
var winHeight int32 = 480
const fragPath = "./shaders/quadtexture.frag.glsl"
const vertPath = "./shaders/main.vert.glsl"

func main() {
    if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
        panic(err)
    }
    defer sdl.Quit()

    sdl.SetHint(sdl.HINT_VIDEO_WAYLAND_PREFER_LIBDECOR, "1")
    sdl.SetHint(sdl.HINT_VIDEO_X11_NET_WM_BYPASS_COMPOSITOR, "0")

    window, err := sdl.CreateWindow("GO Gaming Engine",
    sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
    winWidth, winHeight,
    sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE|sdl.WINDOW_BORDERLESS)
    if err != nil {
        panic(err)
    }

    if _, err := window.GLCreateContext(); err != nil {
        panic(err)
    }
    defer window.Destroy()

    if err := gl.Init();
    err != nil {
        panic(err)
    }
    gl.Enable(gl.DEPTH_TEST)
    glf.PrintVersionGL()

    ShaderProg1, err := glf.NewShaderProgram(vertPath, fragPath)
    if err != nil && ghf.Verbose {
        fmt.Printf("Failed to link Shaders: %s \n", err)
    } else if ghf.Verbose {
        println("Program linked successfully")
    }
    
    vertices, texture, VAO := spaces.LocalSpace()

    positions := spaces.WorldSpace()

    Mat4s := make([]mgl64.Mat4, 3)

    UBO1 := glf.GenBindBuffers(gl.UNIFORM_BUFFER)
    glf.BufferData(gl.UNIFORM_BUFFER, make([]float32, 3*16), gl.DYNAMIC_DRAW)
    gl.BindBufferBase(gl.UNIFORM_BUFFER, 1, UBO1)

    keyboardState := sdl.GetKeyboardState()

    camPos := mgl64.Vec3{0.0, 0.0, 3.0}
    worldUp := mgl64.Vec3{0.0, 1.0, 0.0}
    camera := glf.NewCamera(camPos, worldUp, -90.0, -45.0, 0.02, 0.2)

    var elapsedTime float64 = 0.0

    prevMouseX, prevMouseY, _ := sdl.GetMouseState()

    var timeCount float64 = 0.0
    var FramCount int = 0.0

    for {
        frameStart := time.Now()
        for event := sdl.PollEvent(); event != nil;  event = sdl.PollEvent() {
            switch event := event.(type) {
            case *sdl.QuitEvent:
                return
            case *sdl.WindowEvent:
                switch event.Event {
                case sdl.WINDOWEVENT_SIZE_CHANGED:
                    winWidth, winHeight = window.GetSize()
                    gl.Viewport(0, 0, winWidth, winHeight)
                }
            }
        }

        dir := glf.NoWhere
        switch {
        case keyboardState[sdl.SCANCODE_LEFT] != 0 || keyboardState[sdl.SCANCODE_A] != 0:
            dir = glf.Left
        case keyboardState[sdl.SCANCODE_RIGHT] != 0 || keyboardState[sdl.SCANCODE_D] != 0:
            dir = glf.Right
        case keyboardState[sdl.SCANCODE_UP] != 0 || keyboardState[sdl.SCANCODE_W] != 0:
            dir = glf.Forward
        case keyboardState[sdl.SCANCODE_DOWN] != 0 || keyboardState[sdl.SCANCODE_S] != 0:
            dir = glf.Backward
        default:

        }
        mouseX, mouseY, _ := sdl.GetMouseState()

        camera.UpdateCamera(dir, elapsedTime, float64(mouseX-prevMouseX), float64(mouseY-prevMouseY))
        prevMouseX = mouseX
        prevMouseY = mouseY

        gl.ClearColor(0.1, 0.1, 0.1, 1.0)
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

        ShaderProg1.Use()
        projectionMatrix := mgl64.Perspective(mgl64.DegToRad(45.0), float64(winWidth)/float64(winHeight), 0.1, 100.0)
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

        if timeCount < 1000 {
            FramCount++
            timeCount += elapsedTime
        } else {
            println(FramCount)
            timeCount = 0.0
            FramCount = 0.0
        }
        elapsedTime = float64(time.Since(frameStart).Seconds()*1000)
    }
} 


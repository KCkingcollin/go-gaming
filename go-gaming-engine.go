package main

import (
	"fmt"
	"time"

	spaces "github.com/KCkingcollin/go-gaming/bin"
	"github.com/KCkingcollin/go-help-func/ghf"
	"github.com/KCkingcollin/go-help-func/glf"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
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

    Mat4s := make([]mgl32.Mat4, 3)

    UBO1 := glf.GenBindBuffers(gl.UNIFORM_BUFFER)
    glf.BufferData(gl.UNIFORM_BUFFER, make([]float32, 3*16), gl.DYNAMIC_DRAW)
    gl.BindBufferBase(gl.UNIFORM_BUFFER, 1, UBO1)

    keyboardState := sdl.GetKeyboardState()

    camPos := mgl32.Vec3{0.0, 0.0, 3.0}
    worldUp := mgl32.Vec3{0.0, 1.0, 0.0}
    camera := glf.NewCamera(camPos, worldUp, -90.0, -45.0, 0.02, 0.2)

    var elapsedTime float32 = 0.0

    prevMouseX, prevMouseY, _ := sdl.GetMouseState()
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
                    gl.Viewport(0, 0, int32(winWidth), int32(winHeight))
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
        fmt.Print("yaw:", camera.Yaw, " pitch:", camera.Pitch, "\n")

        camera.UpdateCamera(dir, elapsedTime, float32(mouseX-prevMouseX), float32(mouseY-prevMouseY))
        prevMouseX = mouseX
        prevMouseY = mouseY

        gl.ClearColor(0.1, 0.1, 0.1, 1.0)
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

        ShaderProg1.Use()
        projectionMatrix := mgl32.Perspective(mgl32.DegToRad(45.0), float32(winWidth)/float32(winHeight), 0.1, 100.0)
        viewMatrix := camera.GetViewMatrix()
        Mat4s[2] = projectionMatrix
        Mat4s[1] = viewMatrix

        glf.BindTexture(texture)
        gl.BindVertexArray(VAO)
        for i, pos := range positions {
            modelMatrix := mgl32.Ident4()
            angle := 20.0 * float32(i)
            modelMatrix = mgl32.HomogRotate3D(mgl32.DegToRad(angle), mgl32.Vec3{1.0, 0.3, 0.5}).Mul4(modelMatrix)
            modelMatrix = mgl32.Translate3D(pos.X(), pos.Y(), pos.Z()).Mul4(modelMatrix)
            Mat4s[0] = modelMatrix
            
            glf.BindBufferSubDataMat4(Mat4s, UBO1)

            gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)/5*3))
        }

        window.GLSwap()

        glf.CheckShadersforChanges()

        elapsedTime = float32(time.Since(frameStart).Seconds()*1000)
    }
} 


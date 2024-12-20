package main

import (
	"fmt"

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
    glf.PrintVersionGL()

    ShaderProg1, err := glf.NewShaderProgram(vertPath, fragPath)
    if err != nil && ghf.Verbose {
        fmt.Printf("Failed to link Shaders: %s \n", err)
    } else if ghf.Verbose {
        println("Program linked successfully")
    }
    
    vertices, texture, VAO := spaces.LocalSpace()

    positions := spaces.WorldSpace()

    // VirtUniforms := []float32 {
    //     0.0, // float x
    //     0.0, // float y
    // }
    // uniformVars := uniforms

    Mat4s := make([]mgl32.Mat4, 3)

    UBO1 := glf.GenBindBuffers(gl.UNIFORM_BUFFER)
    glf.BufferData(gl.UNIFORM_BUFFER, make([]float32, 3*16), gl.DYNAMIC_DRAW)
    gl.BindBufferBase(gl.UNIFORM_BUFFER, 1, UBO1)

    keyboardState := sdl.GetKeyboardState()

    var x float32
    var z float32

    for {
        for event := sdl.PollEvent(); event != nil;  event = sdl.PollEvent() {
            switch event := event.(type) {
            case *sdl.QuitEvent:
                return
            case *sdl.WindowEvent:
                switch event.Event {
                case sdl.WINDOWEVENT_SIZE_CHANGED:
                    winWidth, winHeight = window.GetSize()
                    gl.Viewport(0, 0, int32(winWidth), int32(winHeight))
                case sdl.WINDOWEVENT_MOVED:
                    // Handle window move events, maybe log or adjust behavior
                case sdl.WINDOWEVENT_FOCUS_GAINED:
                    // Handle focus events
                case sdl.WINDOWEVENT_FOCUS_LOST:
                }
            }
        }

        gl.ClearColor(0.1, 0.1, 0.1, 1.0)
        gl.Clear(gl.COLOR_BUFFER_BIT)

        switch {
        case keyboardState[sdl.SCANCODE_LEFT] != 0 || keyboardState[sdl.SCANCODE_A] != 0:
            x += 0.1
        case keyboardState[sdl.SCANCODE_RIGHT] != 0 || keyboardState[sdl.SCANCODE_D] != 0:
            x -= 0.1
        case keyboardState[sdl.SCANCODE_DOWN] != 0 || keyboardState[sdl.SCANCODE_S] != 0:
            z -= 0.1
        case keyboardState[sdl.SCANCODE_UP] != 0 || keyboardState[sdl.SCANCODE_W] != 0:
            z += 0.1
        }

        ShaderProg1.Use()
        projectionMatrix := mgl32.Perspective(mgl32.DegToRad(45.0), float32(winWidth)/float32(winHeight), 0.1, 100.0)
        viewMatrix := mgl32.Ident4()
        viewMatrix = mgl32.Translate3D(x, 0.0, z)
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

        // uniformVars = []float32 {
        //     uniformVars[0] + 0.01, // x movement 
        //     uniformVars[1] + 0.01, // y movement
        // }
    }
} 


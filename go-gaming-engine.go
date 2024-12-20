package main

import (
	"fmt"
	"unsafe"

	spaces "github.com/KCkingcollin/go-gaming/bin"
	"github.com/KCkingcollin/go-help-func/ghf"
	"github.com/KCkingcollin/go-help-func/glf"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

const winWidth = 854
const winHight = 480
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
    winWidth, winHight,
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

    uniforms0Mat4 := make([]float32, 3*16)
    for i := range uniforms0Mat4 {
        uniforms0Mat4[i] = 0.0
    }

    UBO1 := glf.GenBindBuffers(gl.UNIFORM_BUFFER)
    glf.BufferData(gl.UNIFORM_BUFFER, uniforms0Mat4, gl.DYNAMIC_DRAW)
    gl.BindBufferBase(gl.UNIFORM_BUFFER, 1, UBO1)

    for {
        for event := sdl.PollEvent(); event != nil;  event = sdl.PollEvent() {
            switch event := event.(type) {
            case *sdl.QuitEvent:
                return
            case *sdl.WindowEvent:
                switch event.Event {
                case sdl.WINDOWEVENT_SIZE_CHANGED:
                    width, height := window.GetSize()
                    gl.Viewport(0, 0, int32(width), int32(height))
                case sdl.WINDOWEVENT_MOVED:
                    // Handle window move events, maybe log or adjust behavior
                case sdl.WINDOWEVENT_FOCUS_GAINED:
                    // Handle focus events
                case sdl.WINDOWEVENT_FOCUS_LOST:
                }
            }
        }

        gl.Viewport(0, 0, winWidth, winHight)
        gl.ClearColor(0.1, 0.1, 0.1, 1.0)
        gl.Clear(gl.COLOR_BUFFER_BIT)

        ShaderProg1.Use()
        projectionMatrix := mgl32.Perspective(mgl32.DegToRad(45.0), float32(winWidth)/float32(winHight), 0.1, 100.0)
        viewMatrix := mgl32.Ident4()
        viewMatrix = mgl32.Translate3D(0.0, 0.0, -3.0)
        // var modelMatrix mgl32.Mat4
        // ShaderProg1.SetMat4("projection", projectionMatrix)
        // ShaderProg1.SetMat4("view", viewMatrix)

        glf.BindTexture(texture)
        gl.BindVertexArray(VAO)
        // gl.DrawElementsWithOffset(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, 0)
        for _, pos := range positions {
            modelMatrix := mgl32.Ident4()
            modelMatrix = mgl32.Translate3D(pos.X(), pos.Y(), pos.Z()).Mul4(modelMatrix)
            // angle := 20.0 * float32(i)
            // modelMatrix = mgl32.Rotate3DY(mgl32.DegToRad(angle))
            // ShaderProg1.SetMat4("model", modelMatrix)
            Mat4 := []mgl32.Mat4 {
                modelMatrix, 
                viewMatrix, 
                projectionMatrix, 
            }

            for i := range Mat4 {
                mat := [16]float32(Mat4[i])
                gl.BindBuffer(gl.UNIFORM_BUFFER, UBO1)
                gl.BufferSubData(gl.UNIFORM_BUFFER, i*4*16, 4*16, unsafe.Pointer(&mat))
                gl.BindBuffer(gl.UNIFORM_BUFFER, 0)
            }
            gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)/5*3))
        }

        window.GLSwap()

        glf.CheckShadersforChanges()

        // uniformVars = []float32 {
        //     uniformVars[0] + 0.01, // x movement 
        //     uniformVars[1] + 0.01, // y movement
        // }

        // for i := range uniformVars {
        //     gl.BindBuffer(gl.UNIFORM_BUFFER, UBO0)
        //     gl.BufferSubData(gl.UNIFORM_BUFFER, i*4, 4, unsafe.Pointer(&uniformVars[i]))
        //     gl.BindBuffer(gl.UNIFORM_BUFFER, 0)
        // }
    }
} 


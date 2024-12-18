package main

import (
	"fmt"
	"unsafe"

	"github.com/KCkingcollin/go-help-func/ghf"
	"github.com/KCkingcollin/go-help-func/glf"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

const winWidth = 480
const winHight = 480
const fragPath = "shaders/frag.glsl"
const vertPath = "shaders/vert.glsl"


func main() {
    err := sdl.Init(sdl.INIT_EVERYTHING)
    if err != nil {
        panic(err)
    }
    defer sdl.Quit()

    window, err := sdl.CreateWindow("Hello World", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHight, sdl.WINDOW_OPENGL)
    if err != nil {
        panic(err)
    }

    window.GLCreateContext()
    defer window.Destroy()

    err = gl.Init()
    if err != nil {
        panic(err)
    }
    glf.PrintVersionGL()

    ShaderProg1, err := glf.NewShaderProgram(vertPath, fragPath)
    if err != nil && ghf.Verbose {
        fmt.Printf("Failed to link Shaders: %s \n", err)
    } else if ghf.Verbose {
        println("Program linked successfully")
    }

    vertices := []float32 {
        0.5, 0.5, 0.0, 1.0, 1.0, 
        0.5, -0.5, 0.0, 1.0, 0.0, 
        -0.5, -0.5, 0.0, 0.0, 0.0, 
        -0.5, 0.5, 0.0, 0.0, 1.0, 
    }

    indices := []uint32 {
        0, 1, 3, // triangle 1
        1, 2, 3, // triangle 2
    }
    
    uniforms := []float32 {
        0.0, // float x
        0.0, // float y
    }
    var x float32 = uniforms[0]
    var y float32 = uniforms[1]

    glf.GenBindBuffers(gl.ARRAY_BUFFER)
    VAO := glf.GenBindArrays()
    glf.BufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW)

    glf.GenBindBuffers(gl.ELEMENT_ARRAY_BUFFER)
    glf.BufferData(gl.ELEMENT_ARRAY_BUFFER, indices, gl.STATIC_DRAW)
    
    UBO := glf.GenBindBuffers(gl.UNIFORM_BUFFER)
    glf.BufferData(gl.UNIFORM_BUFFER, uniforms, gl.DYNAMIC_DRAW)
    gl.BindBufferBase(gl.UNIFORM_BUFFER, 0, UBO)

    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
    gl.EnableVertexAttribArray(0)
    gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
    gl.EnableVertexAttribArray(1)
    gl.BindVertexArray(0)

    for {
        for event := sdl.PollEvent(); event != nil;  event = sdl.PollEvent() {
            switch event.(type) {
            case *sdl.QuitEvent:
                return
            }
        }
        gl.ClearColor(0.1, 0.1, 0.1, 1.0)
        gl.Clear(gl.COLOR_BUFFER_BIT)

        ShaderProg1.Use()

        gl.BindVertexArray(VAO)
        gl.DrawElementsWithOffset(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, 0)

        window.GLSwap()

        x += 0.01
        y += 0.01
        vars := []float32 {
            x, 
            y, 
        }
        
        for i := range vars {
            gl.BindBuffer(gl.UNIFORM_BUFFER, UBO)
            gl.BufferSubData(gl.UNIFORM_BUFFER, i*4, 4, unsafe.Pointer(&vars[i]))
            gl.BindBuffer(gl.UNIFORM_BUFFER, 0)
        }

        glf.CheckShadersforChanges()
    }
} 


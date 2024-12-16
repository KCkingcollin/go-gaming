package main

import (
	"fmt"

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

        // 0.5, -0.5, 0.0,
        // -0.5, 0.5, 0.0, 
        // 0.5, 0.5, 0.0, 
    }

    indices := []uint32 {
        0, 1, 3, // triangle 1
        1, 2, 3, // triangle 2
    }
    
    // VBO := glf.GenBindBuffers(gl.ARRAY_BUFFER)
    glf.GenBindBuffers(gl.ARRAY_BUFFER)
    VAO := glf.GenBindArrays()
    glf.BufferDataFloat(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW)
    glf.GenBindBuffers(gl.ELEMENT_ARRAY_BUFFER)
    glf.BufferDataInt(gl.ELEMENT_ARRAY_BUFFER, indices, gl.STATIC_DRAW)

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

        glf.CheckShadersforChanges()
    }
} 


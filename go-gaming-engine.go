package main

import (
	"fmt"

	"github.com/KCkingcollin/go-help-func/ghf"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

const winWidth = 480
const winHight = 480
const fragPath = "shaders/frag.glsl"
const vertPath = "shaders/vert.glsl"

func compileShaders() (uint32, uint32) {
    vertexShaderSource := ghf.LoadFile(vertPath)
    vertexShader, err := ghf.CreateVertexShader(string(vertexShaderSource))
    if err != nil {
        fmt.Printf("Failed to compile vertex shader: %s \n", err)
    } else {
        println("Vertex shader compiled successfully")
    }

    fragmentShaderSource := ghf.LoadFile(fragPath)
    fragmentShader, err := ghf.CreateFragmentShader(string(fragmentShaderSource))
    if err != nil {
        fmt.Printf("Failed to compile fragment shader: %s \n", err)
    } else {
        println("Fragment shader compiled successfully")
    }

    return vertexShader, fragmentShader
}

func makeProgram(vertexShader, fragmentShader uint32) uint32 {
    ProgramID, err := ghf.CreateProgram(vertexShader, fragmentShader)
    if err != nil {
        fmt.Printf("Failed to link Shaders: %s \n", err)
    } else {
        println("Program linked successfully")
    }

    gl.DeleteShader(vertexShader)
    gl.DeleteShader(fragmentShader)

    return ProgramID
}

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
    ghf.PrintVersionGL()

    ProgramID := makeProgram(compileShaders())

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
    
    // VBO := ghf.GenBindBuffers(gl.ARRAY_BUFFER)
    ghf.GenBindBuffers(gl.ARRAY_BUFFER)
    VAO := ghf.GenBindArrays()
    ghf.BufferDataFloat(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW)
    ghf.GenBindBuffers(gl.ELEMENT_ARRAY_BUFFER)
    ghf.BufferDataInt(gl.ELEMENT_ARRAY_BUFFER, indices, gl.STATIC_DRAW)

    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
    gl.EnableVertexAttribArray(0)
    gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 2*4)
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

        gl.UseProgram(ProgramID)
        gl.BindVertexArray(VAO)
        gl.DrawElementsWithOffset(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, 0)

        window.GLSwap()

        if ghf.CheckShadersforChanges() {
            ghf.ResetLoadedShaders()
            ProgramID = makeProgram(compileShaders())
        }
    }
} 


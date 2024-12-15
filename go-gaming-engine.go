package main

import (
	"fmt"

	"github.com/KCkingcollin/go-help-func/ghf"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

const winWidth = 1280
const winHight = 720

type log struct {
    log string
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

    shaderProgram := makeProgram(compileShaders())

    vertices := []float32 {
        -0.5, -0.5, 0.0,
        0.5, -0.5, 0.0,
        0.0, 0.5, 0.0}
    
    var VBO uint32
    gl.GenBuffers(1, &VBO)
    gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
    var VAO uint32
    gl.GenVertexArrays(1, &VAO)
    gl.BindVertexArray(VAO)

    gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
    gl.EnableVertexAttribArray(0)
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

        gl.UseProgram(shaderProgram)
        gl.BindVertexArray(VAO)
        gl.DrawArrays(gl.TRIANGLES, 0, 3)

        window.GLSwap()

    }
} 

func compileShaders() (uint32, uint32) {
    vertexShaderSource := ghf.LoadFile("shaders/vert.glsl")
    vertexShader, err := ghf.CreateVertexShader(string(vertexShaderSource))
    if err != nil {
        fmt.Printf("Failed to compile vertex shader: %s \n", err)
    } else {
        println("Vertex shader compiled successfully")
    }

    fragmentShaderSource := ghf.LoadFile("shaders/frag.glsl")
    fragmentShader, err := ghf.CreateFragmentShader(string(fragmentShaderSource))
    if err != nil {
        fmt.Printf("Failed to compile fragment shader: %s \n", err)
    } else {
        println("Fragment shader compiled successfully")
    }

    return vertexShader, fragmentShader
}

func makeProgram(vertexShader, fragmentShader uint32) uint32 {
    shaderProgram, err := ghf.CreateProgram(vertexShader, fragmentShader)
    if err != nil {
        fmt.Printf("Failed to link Shaders: %s \n", err)
    } else {
        println("Program linked successfully")
    }

    gl.DeleteShader(vertexShader)
    gl.DeleteShader(fragmentShader)

    return shaderProgram
}

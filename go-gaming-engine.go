package main

import (
	"fmt"
	"os"
	"strings"

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

    gl.Init()
    version := ghf.GetVersionGL()
    println(ghf.Test2())
    fmt.Println("OpenGL Version", version)


    vertexShaderSource, err := os.ReadFile("shaders/vert.glsl")
    if err != nil {
        fmt.Println(err)
    }
    
    vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
    csource, free := gl.Strs(string(vertexShaderSource))
    gl.ShaderSource(vertexShader, 1, csource, nil)
    free()
    gl.CompileShader(vertexShader)
    var status int32
    gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetShaderiv(vertexShader, gl.INFO_LOG_LENGTH, &logLength)
        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetShaderInfoLog(vertexShader, logLength, nil, gl.Str(log))
        fmt.Printf("Failed to compile vertex shader: \n" + log)
    }
    println("Vertex shader compiled successfully")
    
    fragmentShaderSource, err := os.ReadFile("shaders/frag.glsl")
    if err != nil {
        fmt.Println(err)
    }
    
    fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
    csource, free = gl.Strs(string(fragmentShaderSource))
    gl.ShaderSource(fragmentShader, 1, csource, nil)
    free()
    gl.CompileShader(fragmentShader)
    gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetShaderiv(fragmentShader, gl.INFO_LOG_LENGTH, &logLength)
        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetShaderInfoLog(fragmentShader, logLength, nil, gl.Str(log))
        fmt.Printf("Failed to compile fragment shader: \n" + log)
    }
    println("Fragment shader compiled successfully")

    shaderProgram := gl.CreateProgram()
    gl.AttachShader(shaderProgram, vertexShader)
    gl.AttachShader(shaderProgram, fragmentShader)
    gl.LinkProgram(shaderProgram)
    var success int32
    gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
    if success == gl.FALSE {
        var logLength int32
        gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)
        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))
        fmt.Printf("Failed to link shaders: \n" + log)
    }
    println("Program linked successfully")

    gl.DeleteShader(vertexShader)
    gl.DeleteShader(fragmentShader)

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

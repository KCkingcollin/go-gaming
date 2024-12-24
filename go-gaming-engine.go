package main

import (
	"github.com/KCkingcollin/go-gaming/bin"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
    bin.InitWindow()
    bin.InitBuffers()
    bin.MainLoop()
    defer sdl.Quit()
    defer bin.Window.Destroy()
} 


package source

import (
	"github.com/KCkingcollin/go-help-func/glf"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

// Begin polling for window events
//
//  - Returns true if sdl.QuitEvent is detected 
//  - Returns false if sdl.QuitEvent is not detected
func PollEvents() bool { 
    MouseX, MouseY = 0, 0
    for event := sdl.PollEvent(); event != nil;  event = sdl.PollEvent() {
        switch event := event.(type) {
        case *sdl.QuitEvent:
            return true
        case *sdl.WindowEvent:
            switch event.Event {
            case sdl.WINDOWEVENT_SIZE_CHANGED:
                WinWidth, WinHeight = Window.GetSize()
                gl.Viewport(0, 0, WinWidth, WinHeight)
                DisplayRatio = float64(WinWidth) / float64(WinHeight)
            case sdl.WINDOWEVENT_FOCUS_GAINED:
                sdl.GLSetSwapInterval(0)
                FrameRateLimit = PrevFrameLimit
            case sdl.WINDOWEVENT_FOCUS_LOST:
                sdl.GLSetSwapInterval(1)
                PrevFrameLimit = FrameRateLimit
                FrameRateLimit = BackgroundLimit
            }
        case *sdl.KeyboardEvent:
            switch event.Keysym.Sym {
            case sdl.K_ESCAPE:
                keys := make([]bool, 256)
                if event.Type == sdl.KEYDOWN {
                    keys[event.Keysym.Sym] = true
                    InMenu = !InMenu
                    if InMenu {
                        sdl.SetRelativeMouseMode(false)
                        println("Input disabled")
                    } else {
                        sdl.SetRelativeMouseMode(true)
                        println("Input enabled")
                    }
                } else if event.Type == sdl.KEYUP {
                    keys[event.Keysym.Sym] = false
                }
            }
        case *sdl.MouseMotionEvent:
            if !InMenu {
                MouseX += event.XRel
                MouseY += event.YRel
            }
        }
    }
    return false
}

// Set camera position and rotation
func CameraEvents() {
    dir := glf.NoWhere
    deltaT := float64(ElapsedTime.Milliseconds()) * 0.1
    if !InMenu {
        if KeyboardState[sdl.SCANCODE_LEFT] != 0 || KeyboardState[sdl.SCANCODE_A] != 0 {
            dir = glf.Left
            Camera.UpdateCamera(dir, deltaT, float64(MouseX), float64(MouseY))
        }
        if KeyboardState[sdl.SCANCODE_RIGHT] != 0 || KeyboardState[sdl.SCANCODE_D] != 0 {
            dir = glf.Right
            Camera.UpdateCamera(dir, deltaT, float64(MouseX), float64(MouseY))
        }
        if KeyboardState[sdl.SCANCODE_UP] != 0 || KeyboardState[sdl.SCANCODE_W] != 0 {
            dir = glf.Forward
            Camera.UpdateCamera(dir, deltaT, float64(MouseX), float64(MouseY))
        }
        if KeyboardState[sdl.SCANCODE_DOWN] != 0 || KeyboardState[sdl.SCANCODE_S] != 0 {
            dir = glf.Backward
            Camera.UpdateCamera(dir, deltaT, float64(MouseX), float64(MouseY))
        }
        Camera.UpdateCamera(dir, deltaT, float64(MouseX), float64(MouseY))
    }
}


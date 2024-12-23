package inputs

import (
	"github.com/KCkingcollin/go-gaming/bin/vars"
	"github.com/KCkingcollin/go-help-func/glf"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

func PollEvents() bool {
    vars.MouseX, vars.MouseY = 0, 0
    for event := sdl.PollEvent(); event != nil;  event = sdl.PollEvent() {
        switch event := event.(type) {
        case *sdl.QuitEvent:
            return true
        case *sdl.WindowEvent:
            switch event.Event {
            case sdl.WINDOWEVENT_SIZE_CHANGED:
                vars.WinWidth, vars.WinHeight = vars.Window.GetSize()
                gl.Viewport(0, 0, vars.WinWidth, vars.WinHeight)
                vars.DisplayRatio = float64(vars.WinWidth) / float64(vars.WinHeight)
            case sdl.WINDOWEVENT_FOCUS_GAINED:
                sdl.GLSetSwapInterval(0)
                vars.FrameRateLimit = vars.PrevFrameLimit
            case sdl.WINDOWEVENT_FOCUS_LOST:
                sdl.GLSetSwapInterval(1)
                vars.PrevFrameLimit = vars.FrameRateLimit
                vars.FrameRateLimit = vars.BackgroundLimit
            }
        case *sdl.KeyboardEvent:
            switch event.Keysym.Sym {
            case sdl.K_ESCAPE:
                keys := make([]bool, 256)
                if event.Type == sdl.KEYDOWN {
                    keys[event.Keysym.Sym] = true
                    vars.InMenu = !vars.InMenu
                    if vars.InMenu {
                        sdl.SetRelativeMouseMode(false)
                        println("Mouse disabled")
                    } else {
                        sdl.SetRelativeMouseMode(true)
                        println("Mouse enabled")
                    }
                } else if event.Type == sdl.KEYUP {
                    keys[event.Keysym.Sym] = false
                }
            }
        case *sdl.MouseMotionEvent:
            if !vars.InMenu {
                vars.MouseX += event.XRel
                vars.MouseY += event.YRel
            }
        }
    }
    return false
}

func CameraEvents() {
    dir := glf.NoWhere
    deltaT := float64(vars.ElapsedTime.Milliseconds()) * 0.1
    if vars.KeyboardState[sdl.SCANCODE_LEFT] != 0 || vars.KeyboardState[sdl.SCANCODE_A] != 0 {
        dir = glf.Left
        vars.Camera.UpdateCamera(dir, deltaT, float64(vars.MouseX), float64(vars.MouseY))
    }
    if vars.KeyboardState[sdl.SCANCODE_RIGHT] != 0 || vars.KeyboardState[sdl.SCANCODE_D] != 0 {
        dir = glf.Right
        vars.Camera.UpdateCamera(dir, deltaT, float64(vars.MouseX), float64(vars.MouseY))
    }
    if vars.KeyboardState[sdl.SCANCODE_UP] != 0 || vars.KeyboardState[sdl.SCANCODE_W] != 0 {
        dir = glf.Forward
        vars.Camera.UpdateCamera(dir, deltaT, float64(vars.MouseX), float64(vars.MouseY))
    }
    if vars.KeyboardState[sdl.SCANCODE_DOWN] != 0 || vars.KeyboardState[sdl.SCANCODE_S] != 0 {
        dir = glf.Backward
        vars.Camera.UpdateCamera(dir, deltaT, float64(vars.MouseX), float64(vars.MouseY))
    }
    vars.Camera.UpdateCamera(dir, deltaT, float64(vars.MouseX), float64(vars.MouseY))
}


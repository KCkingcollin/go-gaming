package source

import (
	"fmt"
	"strconv"

	"github.com/KCkingcollin/go-help-func/ghf"
	"github.com/KCkingcollin/go-help-func/glf"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

// Map of multiple runes (key combinations) to SDL scancodes, including German-specific characters
var CharToScancode = map[string]int{
	// Alphabet keys (lowercase) for QWERTY and German layout
	"a": int(sdl.SCANCODE_A),
	"b": int(sdl.SCANCODE_B),
	"c": int(sdl.SCANCODE_C),
	"d": int(sdl.SCANCODE_D),
	"e": int(sdl.SCANCODE_E),
	"f": int(sdl.SCANCODE_F),
	"g": int(sdl.SCANCODE_G),
	"h": int(sdl.SCANCODE_H),
	"i": int(sdl.SCANCODE_I),
	"j": int(sdl.SCANCODE_J),
	"k": int(sdl.SCANCODE_K),
	"l": int(sdl.SCANCODE_L),
	"m": int(sdl.SCANCODE_M),
	"n": int(sdl.SCANCODE_N),
	"o": int(sdl.SCANCODE_O),
	"p": int(sdl.SCANCODE_P),
	"q": int(sdl.SCANCODE_Q),
	"r": int(sdl.SCANCODE_R),
	"s": int(sdl.SCANCODE_S),
	"t": int(sdl.SCANCODE_T),
	"u": int(sdl.SCANCODE_U),
	"v": int(sdl.SCANCODE_V),
	"w": int(sdl.SCANCODE_W),
	"x": int(sdl.SCANCODE_X),
	"y": int(sdl.SCANCODE_Y),
	"z": int(sdl.SCANCODE_Z),

	// German-specific alphabet keys
	"ä": int(sdl.SCANCODE_APOSTROPHE),
	"ö": int(sdl.SCANCODE_SEMICOLON),
	"ü": int(sdl.SCANCODE_COMMA),
	"ß": int(sdl.SCANCODE_SLASH),
	"Ä": int(sdl.SCANCODE_APOSTROPHE),
	"Ö": int(sdl.SCANCODE_SEMICOLON),
	"Ü": int(sdl.SCANCODE_COMMA),

	// Numbers and symbols
	"1": int(sdl.SCANCODE_1),
	"2": int(sdl.SCANCODE_2),
	"3": int(sdl.SCANCODE_3),
	"4": int(sdl.SCANCODE_4),
	"5": int(sdl.SCANCODE_5),
	"6": int(sdl.SCANCODE_6),
	"7": int(sdl.SCANCODE_7),
	"8": int(sdl.SCANCODE_8),
	"9": int(sdl.SCANCODE_9),
	"0": int(sdl.SCANCODE_0),
	"Space": int(sdl.SCANCODE_SPACE),
	" ": int(sdl.SCANCODE_SPACE),
	"Enter": int(sdl.SCANCODE_RETURN),
	"Esc": int(sdl.SCANCODE_ESCAPE),
	"Tab": int(sdl.SCANCODE_TAB),
	"Backspace": int(sdl.SCANCODE_BACKSPACE),

	// Special symbols
	"-": int(sdl.SCANCODE_MINUS),
	"=": int(sdl.SCANCODE_EQUALS),
	"[": int(sdl.SCANCODE_LEFTBRACKET),
	"]": int(sdl.SCANCODE_RIGHTBRACKET),
	"\\": int(sdl.SCANCODE_BACKSLASH),
	";": int(sdl.SCANCODE_SEMICOLON),
	"'": int(sdl.SCANCODE_APOSTROPHE),
	",": int(sdl.SCANCODE_COMMA),
	".": int(sdl.SCANCODE_PERIOD),
	"/": int(sdl.SCANCODE_SLASH),
    "`": int(sdl.SCANCODE_GRAVE), 

	// Modifier keys
	"LCTRL": int(sdl.SCANCODE_LCTRL),
	"LShift": int(sdl.SCANCODE_LSHIFT),
	"LALT": int(sdl.SCANCODE_LALT),
	"RCTRL": int(sdl.SCANCODE_RCTRL),
	"RShift": int(sdl.SCANCODE_RSHIFT),
	"RALT": int(sdl.SCANCODE_RALT),

	// Additional keys
	"F1": int(sdl.SCANCODE_F1),
	"F2": int(sdl.SCANCODE_F2),
	"F3": int(sdl.SCANCODE_F3),
	"F4": int(sdl.SCANCODE_F4),
	"F5": int(sdl.SCANCODE_F5),
	"F6": int(sdl.SCANCODE_F6),
	"F7": int(sdl.SCANCODE_F7),
	"F8": int(sdl.SCANCODE_F8),
	"F9": int(sdl.SCANCODE_F9),
	"F10": int(sdl.SCANCODE_F10),
	"F11": int(sdl.SCANCODE_F11),
	"F12": int(sdl.SCANCODE_F12),

	// Arrow keys
	"Up": int(sdl.SCANCODE_UP),
	"Down": int(sdl.SCANCODE_DOWN),
	"Left": int(sdl.SCANCODE_LEFT),
	"Right": int(sdl.SCANCODE_RIGHT),

	// Numpad keys
	"Numpad0": int(sdl.SCANCODE_KP_0),
	"Numpad1": int(sdl.SCANCODE_KP_1),
	"Numpad2": int(sdl.SCANCODE_KP_2),
	"Numpad3": int(sdl.SCANCODE_KP_3),
	"Numpad4": int(sdl.SCANCODE_KP_4),
	"Numpad5": int(sdl.SCANCODE_KP_5),
	"Numpad6": int(sdl.SCANCODE_KP_6),
	"Numpad7": int(sdl.SCANCODE_KP_7),
	"Numpad8": int(sdl.SCANCODE_KP_8),
	"Numpad9": int(sdl.SCANCODE_KP_9),
	"NumpadEnter": int(sdl.SCANCODE_KP_ENTER),
	"NumpadPlus": int(sdl.SCANCODE_KP_PLUS),
	"NumpadMinus": int(sdl.SCANCODE_KP_MINUS),
}

var ScancodeToChar = ghf.ReverseMap(CharToScancode)

var Keys = map[string]int{
    "LeftBind":         sdl.SCANCODE_A, 
    "RightBind":        sdl.SCANCODE_D, 
    "ForwardBind":      sdl.SCANCODE_W,  
    "BackBind":         sdl.SCANCODE_S,
    "AltLeftBind":      sdl.SCANCODE_LEFT,
    "AltRightBind":     sdl.SCANCODE_RIGHT, 
    "AltForwardBind":   sdl.SCANCODE_UP, 
    "AltBackBind":      sdl.SCANCODE_DOWN,
    "UpBind":           sdl.SCANCODE_SPACE, 
    "DownBind":         sdl.SCANCODE_LCTRL, 
}

func UpdateKeyBindings() bool {
    needsUpdate := false
    for _, binding := range Config.Binding {
        if _, exists := Keys[binding.Action]; exists {
            if value, err := strconv.Atoi(binding.Key); err == nil {
                Keys[binding.Action] = value
            } else {
                Keys[binding.Action] = CharToScancode[binding.Key]
            }
        } else {
            fmt.Printf("Binding %s is not remapable, or was typed wrong", binding.Action)
        }
    }
    if len(Config.Binding) < len(Keys) {
        for name, value := range Keys {
            exist := false
            for _, binding := range Config.Binding {
                if binding.Action == name {
                    exist = true
                    break
                }
            }
            if !exist {
                Config.Binding = append(Config.Binding, Binding{name, ScancodeToChar[value]})
            }
        }
        needsUpdate = true
    }
    return needsUpdate
}

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
                FrameSpeedFactor = 1
            case sdl.WINDOWEVENT_FOCUS_LOST:
                sdl.GLSetSwapInterval(1)
                FrameSpeedFactor = float64(FrameRateLimit / BackgroundFrameLimit)
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
        if KeyboardState[Keys["LeftBind"]] != 0 || KeyboardState[Keys["AltLeftBind"]] != 0 {
            dir = glf.Left
            Camera.UpdateCamera(dir, deltaT, float64(MouseX), float64(MouseY))
        }
        if KeyboardState[Keys["RightBind"]] != 0 || KeyboardState[Keys["AltRightBind"]] != 0 {
            dir = glf.Right
            Camera.UpdateCamera(dir, deltaT, float64(MouseX), float64(MouseY))
        }
        if KeyboardState[Keys["ForwardBind"]] != 0 || KeyboardState[Keys["AltForwardBind"]] != 0 {
            dir = glf.Forward
            Camera.UpdateCamera(dir, deltaT, float64(MouseX), float64(MouseY))
        }
        if KeyboardState[Keys["BackBind"]] != 0 || KeyboardState[Keys["AltBackBind"]] != 0 {
            dir = glf.Backward
            Camera.UpdateCamera(dir, deltaT, float64(MouseX), float64(MouseY))
        }
        Camera.UpdateCamera(dir, deltaT, float64(MouseX), float64(MouseY))
    }
}


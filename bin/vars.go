package bin

import (
	"time"

	"github.com/KCkingcollin/go-help-func/glf"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/veandco/go-sdl2/sdl"
)

var (
    err                 error 

    InMenu              bool            = true

    KeyboardState       []uint8         = sdl.GetKeyboardState()

    UBO1                uint32
    VAO                 uint32
    Texture             uint32

    WinWidth, WinHeight int32           = 854, 480
    MouseX, MouseY      int32

    Vertices            []float32

    FrameRateLimit      int64           = 500
    PrevFrameLimit      int64           = FrameRateLimit
    BackgroundLimit     int64           = 15 

    DisplayRatio        float64         = float64(WinWidth / WinHeight)

    ElapsedTime         time.Duration

    TimeCount           time.Time       = time.Now()

    CamPos              mgl64.Vec3      = mgl64.Vec3{0.0, 0.0, 0.0}
    WorldUp             mgl64.Vec3      = mgl64.Vec3{0.0, 1.0, 0.0}
    Positions           []mgl64.Vec3                     

    Mat4s               []mgl64.Mat4    = make([]mgl64.Mat4, 3)

    Camera              *glf.Camera     = glf.NewCamera(CamPos, WorldUp, -90.0, 0.0, 0.02, 0.2)

    Window              *sdl.Window

    ShaderProg1         *glf.ShaderInfo
)

const (
    FragPath            string          = "./shaders/quadtexture.frag.glsl"
    VertPath            string          = "./shaders/main.vert.glsl"
    TimeFactor          int64           = int64(time.Second/time.Nanosecond)
)


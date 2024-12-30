package source

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"github.com/KCkingcollin/go-help-func/ghf"
	"github.com/KCkingcollin/go-help-func/glf"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/veandco/go-sdl2/sdl"
)

var (
    err                     error 

    InMenu                  bool            = true

    KeyboardState           []uint8         = sdl.GetKeyboardState()

    UBO1                    uint32
    UBO0                    uint32
    VAO                     uint32
    Texture                 uint32

    WinWidth, WinHeight     int32           = 854, 480
    MouseX, MouseY          int32

    Vertices                []float32
    Normals                 []float32

    FrameRateLimit          int64           = 500
    BackgroundFrameLimit    int64           = 15 
    SecondTime              int64

    FrameSpeedFactor        float64         = 1.0
    DisplayRatio            float64         = float64(WinWidth / WinHeight)

    ElapsedTime             time.Duration

    Time                    time.Time       = time.Now()
    ConfigModTime           time.Time

    CamPos                  mgl64.Vec3      = mgl64.Vec3{0.0, 0.0, 0.0}
    WorldUp                 mgl64.Vec3      = mgl64.Vec3{0.0, 1.0, 0.0}

    Positions               []mgl64.Vec3                     
    UBVec3s                 []mgl64.Vec3    = make([]mgl64.Vec3, 2)

    UBMat4s                 []mgl64.Mat4    = make([]mgl64.Mat4, 3)

    Camera                  *glf.Camera     = glf.NewCamera(CamPos, WorldUp, -90.0, 0.0, 0.02, 0.2)

    Window                  *sdl.Window

    ShaderProg1             *glf.ShaderInfo

    FragPath                string          = "./shaders/quadtexture.frag.glsl"
    VertPath                string          = "./shaders/main.vert.glsl"
)

const (
    TimeFactor              int64           = int64(time.Second/time.Nanosecond)
    ConfigFile              string          = "./config.xml"
)

type Display struct {
    WinWidth                int32           `xml:"WinWidth"`
    WinHeight               int32           `xml:"WinHight"`
    FrameRateLimit          int64           `xml:"FrameRateLimit"`
    BackgroundFrameLimit    int64           `xml:"BackgroundLimit"`
}

type config struct {
    Display                 Display         `xml:"Display"`
    Binding                 []Binding       `xml:"Bindings"`
    Debug                   Debug           `xml:"Debug"`
}

type Binding struct {
	Action                  string          `xml:"Action"`
    Key                     string          `xml:"Key"`
}

type Debug struct {
    FragPath                string          `xml:"FragPath"`
    VertPath                string          `xml:"VertPath"`
}

var Config config

func LoadConfig(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("failed to open config file: ", err)
	}
	defer file.Close()

    Config = config{}

	decoder := xml.NewDecoder(file)
	if err := decoder.Decode(&Config); err != nil {
	    fmt.Println("failed to parse config file: ", err)
	}

    needsUpdate := false
    needsUpdate = setValues()
    if !needsUpdate {
        needsUpdate = UpdateKeyBindings()
    } else {
        UpdateKeyBindings()
    }

    if needsUpdate || !ConfigModTime.Equal(ghf.GetModifiedTime(ConfigFile)) {
        println(len(Config.Binding))
        // Open the file for writing
        writeFile, err := os.Create(filePath) // Creates or truncates the file
        if err != nil {
            fmt.Println("failed to open config file for writing: ", err)
        }
        defer writeFile.Close()

        // Encode the updated Config back to the file
        encoder := xml.NewEncoder(writeFile)
        encoder.Indent("", "  ") // Pretty print the XML
        if err := encoder.Encode(&Config); err != nil {
            fmt.Println("failed to write config to file: ", err)
        }
        ConfigModTime = ghf.GetModifiedTime(ConfigFile)
    }
}

func setValues() bool {
    needsUpdate := false
    if ghf.ContainsString(ConfigFile, "WinWidth") {
        WinWidth                                = Config.Display.WinWidth        
    } else {
        Config.Display.WinWidth                 = WinWidth        
        needsUpdate = true
    }
    if ghf.ContainsString(ConfigFile, "WinHeight") {
        WinHeight                               = Config.Display.WinHeight       
    } else {
        Config.Display.WinHeight                = WinHeight       
        needsUpdate = true
    }
    if ghf.ContainsString(ConfigFile, "FrameRateLimit") {
        FrameRateLimit                          = Config.Display.FrameRateLimit  
    } else {
        Config.Display.FrameRateLimit           = FrameRateLimit  
        needsUpdate = true
    }
    if ghf.ContainsString(ConfigFile, "BackgroundLimit") {
        BackgroundFrameLimit                    = Config.Display.BackgroundFrameLimit 
    } else {
        Config.Display.BackgroundFrameLimit     = BackgroundFrameLimit 
        needsUpdate = true
    }
    if ghf.ContainsString(ConfigFile, "FragPath") {
        FragPath                                = Config.Debug.FragPath        
    } else {
        Config.Debug.FragPath                   = FragPath        
        needsUpdate = true
    }
    if ghf.ContainsString(ConfigFile, "VertPath") {
        VertPath                                = Config.Debug.VertPath        
    } else {
        Config.Debug.VertPath                   = VertPath        
        needsUpdate = true
    }
    return needsUpdate
}

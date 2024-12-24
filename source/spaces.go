// Code for setting up spaces
package source

import (
	"github.com/KCkingcollin/go-help-func/glf"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

// Returns the, vertices + texture locations, normals, texture ID, and VAO ID in that order.
//
// Returns a slice of float32s formated as a vec3 and a vec2 in a slice of float32s in this order:
//  - Vertex positions (3 float32s)
//  - Texture positions (2 float32s)
//
// Returns a slice of float32s formated as a vec3 in a slice of float32s in this order:
//  - Normals (3 float32s)
//
// Returns these after in this order:
//  - Returns the texture ID as a uint32
func LocalSpace() ([]float32, []float32, uint32) {
    texture := glf.LoadTexture("./assets/metalbox_full.png")

	vertices := []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0, 
    }

    normals := make([]float32, len(vertices)/5*3)
    for tri := 0; tri < len(vertices)/5/3; tri++ {
        i := tri * 15
        p1 := mgl32.Vec3{vertices[i], vertices[i+1], vertices[i+2]}
        i += 5
        p2 := mgl32.Vec3{vertices[i], vertices[i+1], vertices[i+2]}
        i += 5
        p3 := mgl32.Vec3{vertices[i], vertices[i+1], vertices[i+2]}
        normal := glf.TriangleNormalCalc(p1, p2, p3)
        normals[tri*9] = normal.X()
        normals[tri*9+1] = normal.Y()
        normals[tri*9+2] = normal.Z()

        normals[tri*9+3] = normal.X()
        normals[tri*9+4] = normal.Y()
        normals[tri*9+5] = normal.Z()

        normals[tri*9+6] = normal.X()
        normals[tri*9+7] = normal.Y()
        normals[tri*9+8] = normal.Z()
    }
    
    return vertices, normals, texture
}

//Returns the positions of objects in order via a slice of Vec3s 
func WorldSpace() []mgl64.Vec3 {
    positions := []mgl64.Vec3{
        {0.0, 0.0, 0.0}, 
        {2.0, 5.0, -10.0}, 
        {1.0, -5.0, 1.0}, 
    }
    return positions
}

package main

import(
	"math"
	"golang.org/x/mobile/exp/f32"
	"encoding/binary"
)

//input semi-major axis length, semi-minor axis length
func makeEllipse(major, minor, incr float32) []byte{
	coordNum := 3*(int(math.Ceil(2*math.Pi/0.1)+2))
	coords := make([]float32, coordNum)
	coords[0] = 0.0 //X Center
	coords[1] = 0.0 //Y Center
	coords[2] = 0.0 //Z Center
	for theta, i := float32(0.0), 1; theta < 2*math.Pi; theta += 0.1 {
		coords[3*i] = major*f32.Cos(theta)
		coords[3*i + 1] = minor*f32.Sin(theta)
		coords[3*i + 2] = 0.0
		i++
	}
	coords[coordNum-3] = 3.0 //completes full circle
	coords[coordNum-2] = 0.0
	coords[coordNum-1] = 0.0
	return f32.Bytes(binary.LittleEndian, coords...)
}

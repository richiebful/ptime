// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin linux windows

// An app that draws a green triangle on a red background.
//
// Note: This demo is an early preview of Go 1.5. In order to build this
// program as an Android APK using the gomobile tool.
//
// See http://godoc.org/golang.org/x/mobile/cmd/gomobile to install gomobile.
//
// Get the basic example and use gomobile to build or install it on your device.
//
//   $ go get -d golang.org/x/mobile/example/basic
//   $ gomobile build golang.org/x/mobile/example/basic # will build an APK
//
//   # plug your Android device to your computer or start an Android emulator.
//   # if you have adb installed on your machine, use gomobile install to
//   # build and deploy the APK to an Android target.
//   $ gomobile install golang.org/x/mobile/example/basic
//
// Switch to your device or emulator to start the Basic application from
// the launcher.
// You can also run the application on your desktop by running the command
// below. (Note: It currently doesn't work on Windows.)
//   $ go install golang.org/x/mobile/example/basic && basic
package main

import (
	"encoding/binary"
	"log"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

var (
	images   *glutil.Images
	fps      *debug.FPS
	program  gl.Program
	position gl.Attrib
	offset   gl.Uniform
	color    gl.Uniform
	buf      gl.Buffer

	green  float32 = 0.8
	touchX float32
	touchY float32
	scrnHeight float32
	scrnWidth float32
)

func main() {
	app.Main(func(a app.App) {
		var glctx gl.Context
		var sz size.Event
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop(glctx)
					glctx = nil
				}
			case size.Event:
				sz = e
				touchX = float32(sz.WidthPx / 2)
				touchY = float32(sz.HeightPx / 2)
				scrnWidth = float32(sz.WidthPx)
				scrnHeight = float32(sz.HeightPx)
			case paint.Event:
				if glctx == nil || e.External {
					// As we are actively painting as fast as
					// we can (usually 60 FPS), skip any paint
					// events sent by the system.
					continue
				}

				onPaint(glctx, sz)
				a.Publish()
				// Drive the animation by preparing to paint the next frame
				// after this one is shown.
				a.Send(paint.Event{})
			case touch.Event:
				touchX = e.X
				touchY = e.Y
			}
		}
	})
}

func onStart(glctx gl.Context) {
	var err error
	program, err = glutil.CreateProgram(glctx, vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	buf = glctx.CreateBuffer()
	glctx.BindBuffer(gl.ARRAY_BUFFER, buf)
	glctx.BufferData(gl.ARRAY_BUFFER, triangleData, gl.STATIC_DRAW)

	position = glctx.GetAttribLocation(program, "position")
	color = glctx.GetUniformLocation(program, "color")
	offset = glctx.GetUniformLocation(program, "offset")

	images = glutil.NewImages(glctx)
	fps = debug.NewFPS(images)
}

func onStop(glctx gl.Context) {
	glctx.DeleteProgram(program)
	glctx.DeleteBuffer(buf)
	fps.Release()
	images.Release()
}

func onPaint(glctx gl.Context, sz size.Event) {
	glctx.ClearColor(1, 0, 0, 1)
	glctx.Clear(gl.COLOR_BUFFER_BIT)

	glctx.UseProgram(program)

	glctx.Uniform4f(color, 0, green, 0, 1)

	glctx.Uniform2f(offset, float32(0.5), float32(0.5))

	glctx.BindBuffer(gl.ARRAY_BUFFER, buf)
	glctx.EnableVertexAttribArray(position)
	glctx.VertexAttribPointer(position, coordsPerVertex, gl.FLOAT, false, 0, 0)
	glctx.DrawArrays(gl.TRIANGLE_FAN, 0, vertexCount)
	glctx.DisableVertexAttribArray(position)

	fps.Draw(sz)
}

var triangleData = f32.Bytes(binary.LittleEndian,
	0.0, 0.0, 0.0, //center of circle
	0.300000, 0.000000, 0.0,
0.298501, 0.029950, 0.0,
0.294020, 0.059601, 0.0,
0.286601, 0.088656, 0.0,
0.276318, 0.116826, 0.0,
0.263275, 0.143828, 0.0,
0.247601, 0.169393, 0.0,
0.229453, 0.193265, 0.0,
0.209012, 0.215207, 0.0,
0.186483, 0.234998, 0.0,
0.162091, 0.252441, 0.0,
0.136079, 0.267362, 0.0,
0.108707, 0.279612, 0.0,
0.080250, 0.289067, 0.0,
0.050990, 0.295635, 0.0,
0.021221, 0.299248, 0.0,
-0.008760, 0.299872, 0.0,
-0.038653, 0.297499, 0.0,
-0.068161, 0.292154, 0.0,
-0.096987, 0.283890, 0.0,
-0.124844, 0.272789, 0.0,
-0.151454, 0.258963, 0.0,
-0.176550, 0.242549, 0.0,
-0.199883, 0.223712, 0.0,
-0.221218, 0.202639, 0.0,
-0.240343, 0.179542, 0.0,
-0.257067, 0.154650, 0.0,
-0.271222, 0.128214, 0.0,
-0.282667, 0.100496, 0.0,
-0.291287, 0.071775, 0.0,
-0.296998, 0.042336, 0.0,
-0.299741, 0.012474, 0.0,
-0.299488, -0.017512, 0.0,
-0.296244, -0.047324, 0.0,
-0.290039, -0.076662, 0.0,
-0.280937, -0.105235, 0.0,
-0.269028, -0.132756, 0.0,
-0.254430, -0.158951, 0.0,
-0.237290, -0.183557, 0.0,
-0.217780, -0.206330, 0.0,
-0.196093, -0.227041, 0.0,
-0.172447, -0.245483, 0.0,
-0.147078, -0.261473, 0.0,
-0.120240, -0.274850, 0.0,
-0.092200, -0.285481, 0.0,
-0.063239, -0.293259, 0.0,
-0.033646, -0.298107, 0.0,
-0.003717, -0.299977, 0.0,
0.026250, -0.298849, 0.0,
0.055954, -0.294736, 0.0,
0.085099, -0.287677, 0.0,
0.113393, -0.277744, 0.0,
0.140555, -0.265036, 0.0,
0.166312, -0.249680, 0.0,
0.190408, -0.231829, 0.0,
0.212601, -0.211662, 0.0,
0.232670, -0.189380, 0.0,
0.250414, -0.165206, 0.0,
0.265656, -0.139381, 0.0,
0.278244, -0.112163, 0.0,
0.288051, -0.083825, 0.0,
0.294981, -0.054649, 0.0,
0.298963, -0.024927, 0.0,

)

const (
	coordsPerVertex = 3
	vertexCount     = 64
)

const vertexShader = `#version 100
uniform vec2 offset;

attribute vec4 position;
void main() {
	// offset comes in with x/y values between 0 and 1.
	// position bounds are -1 to 1.
	vec4 offset4 = vec4(2.0*offset.x-1.0, 1.0-2.0*offset.y, 0, 0);
	gl_Position = position + offset4;
}`

const fragmentShader = `#version 100
precision mediump float;
uniform vec4 color;
void main() {
	gl_FragColor = color;
}`

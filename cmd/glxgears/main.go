package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	pos   = []float32{5.0, 5.0, 10.0, 0.0}
	red   = []float32{0.8, 0.1, 0.0, 1.0}
	green = []float32{0.0, 0.8, 0.2, 1.0}
	blue  = []float32{0.2, 0.2, 1.0, 1.0}

	view_rotx           = float32(20.0)
	view_roty           = float32(30.0)
	view_rotz           = float32(0.0)
	view_zoom           = float32(1.0)
	angle               = 0.0
	gear1, gear2, gear3 uint32
)

func gear(inner_radius, outer_radius, width float64, teeth int, tooth_depth float64) {
	r0 := float64(inner_radius)
	r1 := float64(outer_radius - tooth_depth/2.0)
	r2 := float64(outer_radius + tooth_depth/2.0)

	da := 2.0 * math.Pi / float64(teeth) / 4.0

	gl.ShadeModel(gl.FLAT)

	gl.Normal3f(0.0, 0.0, 1.0)

	/* draw front face */
	gl.Begin(gl.QUAD_STRIP)
	for i := 0; i <= teeth; i++ {
		angle = float64(i) * 2.0 * math.Pi / float64(teeth)
		gl.Vertex3f(float32(r0*math.Cos(angle)), float32(r0*math.Sin(angle)), float32(width*0.5))
		gl.Vertex3f(float32(r1*math.Cos(angle)), float32(r1*math.Sin(angle)), float32(width*0.5))
		if i < teeth {
			gl.Vertex3f(float32(r0*math.Cos(angle)), float32(r0*math.Sin(angle)), float32(width*0.5))
			gl.Vertex3f(float32(r1*math.Cos(angle+3.0*da)), float32(r1*math.Sin(angle+3.0*da)), float32(width*0.5))
		}
	}
	gl.End()

	/* draw front sides of teeth */
	gl.Begin(gl.QUADS)
	da = 2.0 * math.Pi / float64(teeth) / 4.0
	for i := 0; i < teeth; i++ {
		angle = float64(i) * 2.0 * math.Pi / float64(teeth)

		gl.Vertex3f(float32(r1*math.Cos(angle)), float32(r1*math.Sin(angle)), float32(width*0.5))
		gl.Vertex3f(float32(r2*math.Cos(angle+da)), float32(r2*math.Sin(angle+da)), float32(width*0.5))
		gl.Vertex3f(float32(r2*math.Cos(angle+2*da)), float32(r2*math.Sin(angle+2*da)), float32(width*0.5))
		gl.Vertex3f(float32(r1*math.Cos(angle+3*da)), float32(r1*math.Sin(angle+3*da)), float32(width*0.5))
	}
	gl.End()

	gl.Normal3f(0.0, 0.0, -1.0)

	/* draw back face */
	gl.Begin(gl.QUAD_STRIP)
	for i := 0; i <= teeth; i++ {
		angle = float64(i) * 2.0 * math.Pi / float64(teeth)

		gl.Vertex3f(float32(r1*math.Cos(angle)), float32(r1*math.Sin(angle)), float32(-width*0.5))
		gl.Vertex3f(float32(r0*math.Cos(angle)), float32(r0*math.Sin(angle)), float32(-width*0.5))
		if i < teeth {
			gl.Vertex3f(float32(r1*math.Cos(angle+3*da)), float32(r1*math.Sin(angle+3*da)), float32(-width*0.5))
			gl.Vertex3f(float32(r0*math.Cos(angle)), float32(r0*math.Sin(angle)), float32(-width*0.5))
		}
	}
	gl.End()

	gl.Begin(gl.QUADS)
	da = 2.0 * math.Pi / float64(teeth) / 4.0
	for i := 0; i < teeth; i++ {
		angle = float64(i) * 2.0 * math.Pi / float64(teeth)

		gl.Vertex3f(float32(r1*math.Cos(angle+3*da)), float32(r1*math.Sin(angle+3*da)), float32(-width*0.5))
		gl.Vertex3f(float32(r2*math.Cos(angle+2*da)), float32(r2*math.Sin(angle+2*da)), float32(-width*0.5))
		gl.Vertex3f(float32(r2*math.Cos(angle+da)), float32(r2*math.Sin(angle+da)), float32(-width*0.5))
		gl.Vertex3f(float32(r1*math.Cos(angle)), float32(r1*math.Sin(angle)), float32(-width*0.5))
	}
	gl.End()

	gl.Begin(gl.QUAD_STRIP)
	for i := 0; i < teeth; i++ {
		angle = float64(i) * 2.0 * math.Pi / float64(teeth)

		gl.Vertex3f(float32(r1*math.Cos(angle)), float32(r1*math.Sin(angle)), float32(width*0.5))
		gl.Vertex3f(float32(r1*math.Cos(angle)), float32(r1*math.Sin(angle)), float32(-width*0.5))
		u := r2*math.Cos(angle+da) - r1*math.Cos(angle)
		v := r2*math.Sin(angle+da) - r1*math.Sin(angle)
		l := math.Sqrt(u*u + v*v)
		u /= l
		v /= l
		gl.Normal3f(float32(v), float32(-u), 0.0)
		gl.Vertex3f(float32(r2*math.Cos(angle+da)), float32(r2*math.Sin(angle+da)), float32(width*0.5))
		gl.Vertex3f(float32(r2*math.Cos(angle+da)), float32(r2*math.Sin(angle+da)), float32(-width*0.5))
		gl.Normal3f(float32(math.Cos(angle)), float32(math.Sin(angle)), 0.0)
		gl.Vertex3f(float32(r2*math.Cos(angle+2*da)), float32(r2*math.Sin(angle+2*da)), float32(width*0.5))
		gl.Vertex3f(float32(r2*math.Cos(angle+2*da)), float32(r2*math.Sin(angle+2*da)), float32(-width*0.5))
		u = r1*math.Cos(angle+3*da) - r2*math.Cos(angle+2*da)
		v = r1*math.Sin(angle+3*da) - r2*math.Sin(angle+2*da)
		gl.Normal3f(float32(v), float32(-u), 0.0)
		gl.Vertex3f(float32(r1*math.Cos(angle+3*da)), float32(r1*math.Sin(angle+3*da)), float32(width*0.5))
		gl.Vertex3f(float32(r1*math.Cos(angle+3*da)), float32(r1*math.Sin(angle+3*da)), float32(-width*0.5))
		gl.Normal3f(float32(math.Cos(angle)), float32(math.Sin(angle)), 0.0)
	}

	gl.Vertex3f(float32(r1*math.Cos(0)), float32(r1*math.Sin(0)), float32(width*0.5))
	gl.Vertex3f(float32(r1*math.Cos(0)), float32(r1*math.Sin(0)), float32(-width*0.5))

	gl.End()

	gl.ShadeModel(gl.SMOOTH)

	gl.Begin(gl.QUAD_STRIP)
	for i := 0; i <= teeth; i++ {
		angle = float64(i) * 2.0 * math.Pi / float64(teeth)
		gl.Normal3f(float32(-math.Cos(angle)), float32(-math.Sin(angle)), 0.0)
		gl.Vertex3f(float32(r0*math.Cos(angle)), float32(r0*math.Sin(angle)), float32(-width*0.5))
		gl.Vertex3f(float32(r0*math.Cos(angle)), float32(r0*math.Sin(angle)), float32(width*0.5))
	}
	gl.End()
}

func scene() {
	gl.Lightfv(gl.LIGHT0, gl.POSITION, &pos[0])
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.LIGHTING)
	gl.Enable(gl.LIGHT0)
	gl.Enable(gl.DEPTH_TEST)

	/* make the gears */
	gear1 = gl.GenLists(1)
	gl.NewList(gear1, gl.COMPILE)
	gl.Materialfv(gl.FRONT, gl.AMBIENT_AND_DIFFUSE, &red[0])
	gear(1.0, 4.0, 1.0, 20, 0.7)
	gl.EndList()

	gear2 = gl.GenLists(1)
	gl.NewList(gear2, gl.COMPILE)
	gl.Materialfv(gl.FRONT, gl.AMBIENT_AND_DIFFUSE, &green[0])
	gear(0.5, 2.0, 2.0, 10, 0.7)
	gl.EndList()

	gear3 = gl.GenLists(1)
	gl.NewList(gear3, gl.COMPILE)
	gl.Materialfv(gl.FRONT, gl.AMBIENT_AND_DIFFUSE, &blue[0])
	gear(1.3, 2.0, 0.5, 10, 0.7)
	gl.EndList()

	gl.Enable(gl.NORMALIZE)
}

func reshape(width, height int) {
	h := float64(height) / float64(width)

	gl.Viewport(0, 0, int32(width), int32(height))
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Frustum(-1.0, 1.0, -h, h, 5.0, 60.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translatef(0.0, 0.0, -40.0)
}

func draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.PushMatrix()
	gl.Rotatef(view_rotx, 1.0, 0.0, 0.0)
	gl.Rotatef(view_roty, 0.0, 1.0, 0.0)
	gl.Rotatef(view_rotz, 0.0, 0.0, 1.0)
	gl.Scalef(view_zoom, view_zoom, view_zoom)

	gl.PushMatrix()
	gl.Translatef(-3.0, -2.0, 0.0)
	gl.Rotatef(float32(angle), 0.0, 0.0, 1.0)
	gl.CallList(gear1)
	gl.PopMatrix()

	gl.PushMatrix()
	gl.Translatef(3.1, -2.0, 0.0)
	gl.Rotatef(-2.0*float32(angle)-9.0, 0.0, 0.0, 1.0)
	gl.CallList(gear2)
	gl.PopMatrix()

	gl.PushMatrix()
	gl.Translatef(-3.1, 4.2, 0.0)
	gl.Rotatef(-2.0*float32(angle)-25.0, 0.0, 0.0, 1.0)
	gl.CallList(gear3)
	gl.PopMatrix()

	gl.PopMatrix()
}

func keyCb(window *glfw.Window, k glfw.Key, s int, action glfw.Action, mods glfw.ModifierKey) {
	switch glfw.Key(k) {
	case glfw.KeyEscape:
		window.SetShouldClose(true)
	case glfw.KeyLeft:
		view_roty += 5.0
	case glfw.KeyRight:
		view_roty -= 5.0
	case glfw.KeyUp:
		view_rotx += 5.0
	case glfw.KeyDown:
		view_rotx -= 5.0
	case glfw.KeyPageUp:
		view_zoom += 0.1
	case glfw.KeyPageDown:
		view_zoom -= 0.1
	}
}

func main() {
	var information bool
	flag.BoolVar(&information, "info", false, "print OpenGL information")
	flag.Parse()

	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		log.Fatal("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	window, err := glfw.CreateWindow(300, 300, "glxgears", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()

	window.MakeContextCurrent()
	defer glfw.DetachCurrentContext()

	window.SetKeyCallback(keyCb)

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	glfw.SwapInterval(0)

	if information {
		fmt.Printf("GL_RENDERER   = %s\n", gl.GoStr(gl.GetString(gl.RENDERER)))
		fmt.Printf("GL_VERSION    = %s\n", gl.GoStr(gl.GetString(gl.VERSION)))
		fmt.Printf("GL_VENDOR     = %s\n", gl.GoStr(gl.GetString(gl.VENDOR)))
		fmt.Printf("GL_EXTENSIONS = %s\n", gl.GoStr(gl.GetString(gl.EXTENSIONS)))
	}

	reshape(300, 300)
	scene()

	t0 := time.Now()
	frames := 0
	for !window.ShouldClose() {
		angle += 2.0
		draw()
		window.SwapBuffers()

		t := time.Now()
		frames++

		s := t.Sub(t0)
		if s > 5*time.Second {
			fps := float64(frames) / float64(s.Seconds())
			fmt.Printf("%d frames in %3.1f seconds = %6.3f FPS\n", frames, s.Seconds(), fps)
			t0 = t
			frames = 0
		}

		glfw.PollEvents()
	}
}

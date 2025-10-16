package game

import (
	"math"
	"slices"

	"fw/util"
)

type Text struct {
	font   *Font
	text   string
	pos    util.Vector2
	colour util.Pixel
}

func (t *Text) drawTo(buf *util.ScreenBuffer) {
	xp := t.pos.X

	for _, char := range textToChars(t.font, t.text) {
		for y, row := range char.content {
			for x, b := range row {
				// skip empty pixels
				if b != 0 {
					buf.Set(xp+x, t.pos.Y+y, t.colour.Brightness(b))
				}
			}
		}

		xp += int(char.width) - int(t.font.crush)
	}
}

type Grid struct {
	colour1, colour2            util.Pixel
	pos, offset, size, cellSize util.Vector2
}

func (g *Grid) drawTo(buf *util.ScreenBuffer) {
	// no "funny business" around the x and y axes
	for y := range g.size.Y {
		for x := range g.size.X {
			ox, oy := x+g.offset.X, y+g.offset.Y
			cc := ox/g.cellSize.X + oy/g.cellSize.Y
			if ox < 0 {
				cc--
			}
			if oy < 0 {
				cc--
			}

			if cc%2 == 0 {
				buf.Set(g.pos.X+x, g.pos.Y+y, g.colour1)
			} else {
				buf.Set(g.pos.X+x, g.pos.Y+y, g.colour2)
			}
		}
	}
}

type SineWave struct {
	colour                       util.Pixel
	pos                          util.Vector2
	amplitude, wavelength, phase int
}

func (s *SineWave) drawTo(buf *util.ScreenBuffer) {
	for x := s.pos.X; x < util.Width; x++ {
		y := s.pos.Y + int(float64(s.amplitude)*math.Sin(float64(x)/float64(s.wavelength)*2*math.Pi+float64(s.phase)/10))
		buf.Set(x, y, s.colour)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func drawLine(buf *util.ScreenBuffer, p1, p2 util.Vector2, colour util.Pixel) {
	dx := abs(p2.X - p1.X)
	dy := -abs(p2.Y - p1.Y)
	sx := -1
	if p1.X < p2.X {
		sx = 1
	}
	sy := -1
	if p1.Y < p2.Y {
		sy = 1
	}
	e := dx + dy

	for {
		buf.Set(p1.X, p1.Y, colour)
		if p1.X == p2.X && p1.Y == p2.Y {
			break
		}
		e2 := 2 * e
		if e2 >= dy {
			e += dy
			p1.X += sx
		}
		if e2 <= dx {
			e += dx
			p1.Y += sy
		}
	}
}

func floorint(x float32) int {
	if x < 0 {
		return int(x) - 1
	}
	return int(x)
}

func floor(x float32) float32 {
	// return float32(math.Floor(float64(x)))
	if x < 0 {
		return float32(int(x) - 1)
	}
	return float32(int(x))
}

// WARNING: VERY SLOW
func drawLineAntialiased(buf *util.ScreenBuffer, p1, p2 util.Vector2, colour util.Pixel) {
	dx, dy := p2.X-p1.X, p2.Y-p1.Y
	steep := abs(dy) > abs(dx)

	if steep {
		p1, p2 = p1.Swap(), p2.Swap()
		dx, dy = dy, dx
	}

	if p1.X > p2.X {
		p1, p2 = p2, p1
		dx = p2.X - p1.X
		dy = p2.Y - p1.Y
	}

	gradient := float32(dy) / float32(dx)

	// idc about endpoints it works well enough without them

	intery := float32(p1.Y) + gradient

	// main loop
	if steep {
		for x := p1.X; x < p2.X; x++ {
			buf.SetAlpha(floorint(intery), x, colour, uint8((1-(intery-floor(intery)))*0xff))
			buf.SetAlpha(floorint(intery)+1, x, colour, uint8((intery-floor(intery))*0xff))
			intery += gradient
		}
	} else {
		for x := p1.X; x < p2.X; x++ {
			buf.SetAlpha(x, floorint(intery), colour, uint8((1-(intery-floor(intery)))*0xff))
			buf.SetAlpha(x, floorint(intery)+1, colour, uint8((intery-floor(intery))*0xff))
			intery += gradient
		}
	}
}

func f32Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func f32Cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func f32Sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

type Vector3 struct {
	X, Y, Z float32
}

func V3(x, y, z float32) Vector3 {
	return Vector3{X: x, Y: y, Z: z}
}

func (v Vector3) XYZ() (float32, float32, float32) {
	return v.X, v.Y, v.Z
}

func (v Vector3) Add(o Vector3) Vector3 {
	return V3(v.X+o.X, v.Y+o.Y, v.Z+o.Z)
}

func (v Vector3) Div(s float32) Vector3 {
	return V3(v.X/s, v.Y/s, v.Z/s)
}

func (v Vector3) Sin() Vector3 {
	return V3(f32Sin(v.X), f32Sin(v.Y), f32Sin(v.Z))
}

func (v Vector3) Cos() Vector3 {
	return V3(f32Cos(v.X), f32Cos(v.Y), f32Cos(v.Z))
}

func Distance(a, b Vector3) float32 {
	dx, dy, dz := a.X-b.X, a.Y-b.Y, a.Z-b.Z
	return f32Sqrt(dx*dx + dy*dy + dz*dz)
}

type Line3D struct {
	start, end Vector3
	colour     util.Pixel
}

type Object3D struct {
	Lines []Line3D
}

func NewCube3D(colour util.Pixel, position Vector3, size int) *Object3D {
	// 8 corners of the cube
	var points [8]Vector3
	for i := 0; i < 8; i++ {
		x := position.X + float32(size)*(float32((i>>0)&1)-0.5)
		y := position.Y + float32(size)*(float32((i>>1)&1)-0.5)
		z := position.Z + float32(size)*(float32((i>>2)&1)-0.5)
		points[i] = V3(x, y, z)
	}

	// edges
	edges := [12][2]int{
		{0, 1},
		{1, 3},
		{3, 2},
		{2, 0},
		{4, 5},
		{5, 7},
		{7, 6},
		{6, 4},
		{0, 4},
		{1, 5},
		{2, 6},
		{3, 7},
	}

	lines := make([]Line3D, len(edges))
	for i, e := range edges {
		lines[i] = Line3D{start: points[e[0]], end: points[e[1]], colour: colour}
	}

	return &Object3D{Lines: lines}
}

type Scene3D struct {
	objects                []*Object3D
	camera, cameraRotation Vector3
	zoom                   int
}

func (s *Scene3D) drawTo(buf *util.ScreenBuffer) {
	var lines []Line3D
	for _, obj := range s.objects {
		lines = append(lines, obj.Lines...)
	}

	// sort lines by distance from camera (painter's algorithm)
	slices.SortFunc(lines, func(a, b Line3D) int {
		// d1 := Distance(V3((a.start.X+a.end.X)/2, (a.start.Y+a.end.Y)/2, (a.start.Z+a.end.Z)/2), s.camera)
		// d2 := Distance(V3((b.start.X+b.end.X)/2, (b.start.Y+b.end.Y)/2, (b.start.Z+b.end.Z)/2), s.camera)
		d1 := Distance(a.start, s.camera) + Distance(a.end, s.camera)
		d2 := Distance(b.start, s.camera) + Distance(b.end, s.camera)
		if d1 < d2 {
			return 1
		} else if d1 > d2 {
			return -1
		}
		return 0
	})

	fov := float32(s.zoom)

	sin := V3(
		f32Sin(s.cameraRotation.X),
		f32Sin(s.cameraRotation.Y),
		f32Sin(s.cameraRotation.Z),
	)
	cos := V3(
		f32Cos(s.cameraRotation.X),
		f32Cos(s.cameraRotation.Y),
		f32Cos(s.cameraRotation.Z),
	)

	for _, line := range lines {
		// translate line to camera space
		x1, y1, z1 := line.start.X-s.camera.X, line.start.Y-s.camera.Y, line.start.Z-s.camera.Z
		x2, y2, z2 := line.end.X-s.camera.X, line.end.Y-s.camera.Y, line.end.Z-s.camera.Z

		// rotate around Y axis
		x1, z1 = x1*cos.Y-z1*sin.Y, x1*sin.Y+z1*cos.Y
		x2, z2 = x2*cos.Y-z2*sin.Y, x2*sin.Y+z2*cos.Y

		// rotate around X axis
		y1, z1 = y1*cos.X-z1*sin.X, y1*sin.X+z1*cos.X
		y2, z2 = y2*cos.X-z2*sin.X, y2*sin.X+z2*cos.X

		// rotate around Z axis
		x1, y1 = x1*cos.Z-y1*sin.Z, x1*sin.Z+y1*cos.Z
		x2, y2 = x2*cos.Z-y2*sin.Z, x2*sin.Z+y2*cos.Z

		// project to 2D
		if z1 <= 0 || z2 <= 0 {
			continue
		}
		sx1 := int(fov * x1 / z1)
		sy1 := int(fov * y1 / z1)
		sx2 := int(fov * x2 / z2)
		sy2 := int(fov * y2 / z2)

		// convert to screen space
		sx1 += util.Width / 2
		sy1 = util.Height/2 - sy1
		sx2 += util.Width / 2
		sy2 = util.Height/2 - sy2

		// darken depending on distance from camera
		dist := Distance(line.start, s.camera) + Distance(line.end, s.camera)
		fade := max(0xff - dist*2, 0)
		c := line.colour.Brightness(uint8(fade))

		// draw line
		if fade > 0 {
			drawLine(buf, util.V2(sx1, sy1), util.V2(sx2, sy2), c)
		}
	}
}

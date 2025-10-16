package game

import (
	"math"

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

type Cube3D struct {
	colour               util.Pixel
	pos                  util.Vector2
	size, angle1, angle2 int
}

func (c *Cube3D) drawTo(buf *util.ScreenBuffer) {
	// rotating cube in integer 3D
	angle1 := float64(c.angle1) / 10
	angle2 := float64(c.angle2) / 10
	sin1, cos1 := math.Sin(angle1), math.Cos(angle1)
	sin2, cos2 := math.Sin(angle2), math.Cos(angle2)

	// 8 corners of the cube
	points := [8]util.Vector2{}
	for i := range points {
		x := float64((i>>0)&1*2-1) * float64(c.size)
		y := float64((i>>1)&1*2-1) * float64(c.size)
		z := float64((i>>2)&1*2-1) * float64(c.size)

		// rotate around Y axis
		xz := x*cos1 - z*sin1
		z = x*sin1 + z*cos1
		x = xz

		// rotate around X axis
		yz := y*cos2 - z*sin2
		z = y*sin2 + z*cos2
		y = yz

		// project 3D to 2D (simple orthographic projection)
		f := 20.0 / (z + 40) // perspective factor
		sx := int(x*f) + c.pos.X
		sy := int(y*f) + c.pos.Y

		points[i] = util.Vector2{X: sx, Y: sy}
	}

	// draw edges
	edges := [][2]int{
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

	for _, e := range edges {
		drawLineAntialiased(buf, points[e[0]], points[e[1]], c.colour)
	}
}

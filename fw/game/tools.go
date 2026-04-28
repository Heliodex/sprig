package game

import (
	"fw/util"
	"math"
	"slices"
)

type UIElement interface {
	drawTo(*util.ScreenBuffer)
}

type Rect struct {
	pos, size util.Vector2[int]
	colour    util.Pixel
}

func (r *Rect) drawTo(buf *util.ScreenBuffer) {
	for y := r.pos.Y; y < r.pos.Y+r.size.Y; y++ {
		for x := r.pos.X; x < r.pos.X+r.size.X; x++ {
			buf.Set(x, y, r.colour)
		}
	}
}

type Text struct {
	font   *Font
	text   string
	pos    util.Vector2[int]
	colour util.Pixel
}

func (t *Text) drawTo(buf *util.ScreenBuffer) {
	xp := t.pos.X

	for _, char := range textToChars(t.font, t.text) {
		for y, row := range char.content {
			for x, b := range row {
				// skip empty pixels
				if b != 0 {
					buf.SetAlpha(xp+x, t.pos.Y+y, t.colour, b)
				}
			}
		}

		xp += int(char.width) - int(t.font.crush)
	}
}

type Grid struct {
	colour1, colour2            util.Pixel
	pos, offset, size, cellSize util.Vector2[int]
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type DiamondGrid struct {
	colour1, colour2  util.Pixel
	pos, offset, size util.Vector2[int]
	cellSize          int
	heightMulti       float64
}

func (g *DiamondGrid) drawTo(buf *util.ScreenBuffer) {
	// isometric diamond-style grid
	for y := range g.size.Y {
		for x := range g.size.X {
			var c util.Pixel

			xo, yo := abs(x+g.offset.X), abs(y+g.offset.Y)

			diamondWidth := abs(g.cellSize - int(float64(yo)*g.heightMulti)%(g.cellSize*2))
			diamondHeight := abs(g.cellSize - xo%(g.cellSize*2))

			if diamondWidth+diamondHeight <= g.cellSize {
				c = g.colour1
			} else {
				c = g.colour2
			}

			buf.Set(g.pos.X+x, g.pos.Y+y, c)
		}
	}
}

const (
	terrainSize, terrainHeight = 8, 8
	terrainX, terrainY         = terrainSize, terrainSize
	terrainDiagonal            = terrainX + terrainY - 1
)

type TerrainColumn uint8 // 1d, bitpacked

func (col TerrainColumn) Get(z int) bool {
	return (col & (1 << z)) != 0
}

func (col *TerrainColumn) Set(z int, v bool) {
	if v {
		*col |= 1 << z
	} else {
		*col &^= 1 << z
	}
}

type Terrain [terrainX][terrainY]TerrainColumn // 3d

func (t *Terrain) OcclusionCull() {
	const xl = terrainX - 1
	const yl = terrainY - 1
	const zl = terrainHeight - 1

	for x, cols := range *t {
		for y, col := range cols {
			for z := range terrainHeight {
				if !col.Get(z) {
					continue
				}

				// if it's on the front extreme, keep it
				if x == xl || y == yl || z == zl {
					continue
				}

				// if all adjacent blocks in the direction of the "camera" are filled, cull this block
				if (*t)[x+1][y].Get(z) &&
					cols[y+1].Get(z) &&
					col.Get(z+1) {
					(*t)[x][y].Set(z, false)
					continue
				}

				// if all blocks diagonally adjacent in the direction of the "camera" are filled, cull this block
				if (*t)[x+1][y+1].Get(z) &&
					cols[y+1].Get(z+1) &&
					(*t)[x+1][y].Get(z+1) {
					(*t)[x][y].Set(z, false)
					continue
				}

				// if another block is on the diagonal between this block and the "camera", cull it
				for i := 1; i <= min(xl-x, yl-y, zl-z); i++ {
					if (*t)[x+i][y+i].Get(z + i) {
						(*t)[x][y].Set(z, false)
						continue
					}
				}
			}
		}
	}
}

func (t *Terrain) OrderColsDiagonal() (diags [terrainDiagonal][]util.Vector2[int]) {
	// so how this works is the columns will be in the order [(0 0)] [(1 0) (0 1)] [(2 0) (1 1) (0 2)] ...

	for d := range terrainDiagonal {
		for x := range d + 1 {
			y := d - x
			if x < terrainX && y < terrainY {
				diags[d] = append(diags[d], util.V2(x, y))
			}
		}
	}
	return
}

type IsometricProjection struct {
	terrainHeight int
	terrain       Terrain

	topColour, baseColour1, baseColour2 util.Pixel
	minFactor, maxFactor                uint8
	pos, offset, size                   util.Vector2[int]
	cellSize, cellHeight                int

	sprite  UIElement // to be drawn on top of the terrain at the specified Z height
	spriteZ int       // Z up I guess
}

func (g *IsometricProjection) drawTo(buf *util.ScreenBuffer) {
	// isometric diamond-style grid

	var drew int

	// for x, cols := range g.terrain {
	// 	for y, col := range cols {
	// we'll draw diagonally instead
	for i, diag := range g.terrain.OrderColsDiagonal() {
		for _, dcol := range diag {
			x, y := dcol.X, dcol.Y

			extremeX := x == terrainX-1
			extremeY := y == terrainY-1

			sx := g.pos.X + (x-y)*g.cellSize - g.offset.X
			col := g.terrain[x][y]

			for z := range terrainHeight {
				if !col.Get(z) {
					continue
				}

				// calculate brightness factor based on height
				factorf := float64(g.minFactor) + (float64(g.maxFactor-g.minFactor) * float64(z) / float64(g.terrainHeight))
				factor := uint8(min(factorf, 255))

				// project 3d coordinates to 2d isometric
				sy := g.pos.Y + (x+y)*g.cellSize/2 - g.offset.Y - z*g.cellHeight/2
				if sx >= util.Width || sy >= util.Height {
					// cell top is below camera, skip
					continue
				}

				w := g.cellSize * 2
				h := g.cellSize
				if sx+w < 0 || sy+h+g.cellHeight < 0 {
					// cell (including its base) is above camera, skip
					continue
				}

				for dy := range h + g.cellHeight {
					for dx := range w {
						rx, ry := sx+dx, sy+dy
						if rx < 0 || rx >= util.Width || ry < 0 || ry >= util.Height {
							// this particular pixel is out of bounds, skip
							continue
						}

						if dx+dy*2 <= h /* top left */ ||
							dy*2+h <= dx /* top right */ ||
							dy*2-h-g.cellHeight > dx /* bottom left of top */ ||
							dx+dy*2 > h+g.cellHeight+w /* bottom right of top */ {
							// above top or below base, skip
							continue
						}

						extremeZ := z == terrainHeight-1

						if dy*2-h > dx /* bottom left of base */ ||
							dx+dy*2 > h+w /* bottom right of base */ {
							// draw base

							var c util.Pixel
							if dx < w/2 {
								// if there's a block to the front left, no need to draw this
								if !extremeY && g.terrain[x][y+1].Get(z) {
									continue
								}
								// if there's a block above the one to the front left AND one in front, no need to draw this
								if !extremeX && !extremeY && !extremeZ && g.terrain[x][y+1].Get(z+1) && g.terrain[x+1][y+1].Get(z) {
									continue
								}
								c = g.baseColour1 // left base
							} else {
								// if there's a block to the front right, no need to draw this
								if !extremeX && g.terrain[x+1][y].Get(z) {
									continue
								}
								// if there's a block above the one to the front right AND one in front, no need to draw this
								if !extremeX && !extremeY && !extremeZ && g.terrain[x+1][y].Get(z+1) && g.terrain[x+1][y+1].Get(z) {
									continue
								}
								c = g.baseColour2 // right base
							}
							buf.SetAlpha(sx+dx, sy+dy, c.Brightness(factor), 0xff)
							// since we've already checked, we can force set
							// (*buf)[ry][rx] = c.Brightness(factor)
							continue
						}

						// draw top

						// if there's a block on top of this one, no need to draw the top face
						if z < terrainHeight-1 && col.Get(z+1) {
							continue
						}

						buf.SetAlpha(sx+dx, sy+dy, g.topColour.Brightness(factor), 0xff)
						// (*buf)[ry][rx] = g.topColour.Brightness(factor)

						// if rx == util.Width/2 && ry == util.Height/2 {
						// 	println("drawing center pixel of cell at", x, y, z)
						// }
					}
				}
				drew++
			}
		}

		if i == g.spriteZ {
			// g.sprite.drawTo(buf)
		}
	}

	println("drew", drew, "tiles")
}

type SineWave struct {
	colour                       util.Pixel
	pos                          util.Vector2[int]
	amplitude, wavelength, phase int
}

func (s *SineWave) drawTo(buf *util.ScreenBuffer) {
	for x := s.pos.X; x < util.Width; x++ {
		y := s.pos.Y + int(float64(s.amplitude)*math.Sin(float64(x)/float64(s.wavelength)*2*math.Pi+float64(s.phase)/10))
		buf.Set(x, y, s.colour)
	}
}

func drawVLine(buf *util.ScreenBuffer, x, y1, y2 int, colour util.Pixel) {
	for y := y1; y <= y2; y++ {
		buf.Set(x, y, colour)
	}
}

func drawHLine(buf *util.ScreenBuffer, y, x1, x2 int, colour util.Pixel) {
	for x := x1; x <= x2; x++ {
		buf.Set(x, y, colour)
	}
}

func drawLine(buf *util.ScreenBuffer, p1, p2 util.Vector2[int], colour util.Pixel) {
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
func drawLineAntialiased(buf *util.ScreenBuffer, p1, p2 util.Vector2[int], colour util.Pixel) {
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

func clampX(x int) int {
	if x < 0 {
		return -1
	}
	if x > util.Width {
		return util.Width
	}
	return x
}

func clampY(y int) int {
	if y < 0 {
		return -1
	}
	if y > util.Height {
		return util.Height
	}
	return y
}

func drawTriangle2D(buf *util.ScreenBuffer, p1, p2, p3 util.Vector2[int], colour util.Pixel) {
	// sort points by Y
	if p1.Y > p2.Y {
		p1, p2 = p2, p1
	}
	if p1.Y > p3.Y {
		p1, p3 = p3, p1
	}
	if p2.Y > p3.Y {
		p2, p3 = p3, p2
	}

	// compute slopes
	var dx1, dx2, dx3 float32
	if p2.Y-p1.Y > 0 {
		dx1 = float32(p2.X-p1.X) / float32(p2.Y-p1.Y)
	}
	if p3.Y-p1.Y > 0 {
		dx2 = float32(p3.X-p1.X) / float32(p3.Y-p1.Y)
	}
	if p3.Y-p2.Y > 0 {
		dx3 = float32(p3.X-p2.X) / float32(p3.Y-p2.Y)
	}

	// draw upper part
	if p2.Y-p1.Y > 0 {
		for y := clampY(p1.Y); y <= clampY(p2.Y); y++ {
			sx := int(float32(p1.X) + float32(y-p1.Y)*dx1)
			ex := int(float32(p1.X) + float32(y-p1.Y)*dx2)

			if sx > ex {
				sx, ex = ex, sx
			}

			for x := clampX(sx); x <= clampX(ex); x++ {
				buf.Set(x, y, colour)
			}
		}
	}

	// draw lower part
	if p3.Y-p2.Y > 0 {
		for y := clampY(p2.Y); y <= clampY(p3.Y); y++ {
			sx := int(float32(p2.X) + float32(y-p2.Y)*dx3)
			ex := int(float32(p1.X) + float32(y-p1.Y)*dx2)

			if sx > ex {
				sx, ex = ex, sx
			}

			for x := clampX(sx); x <= clampX(ex); x++ {
				buf.Set(x, y, colour)
			}
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

func square(x float64) float64 {
	return x * x
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

func (v Vector3) Mul(s float32) Vector3 {
	return V3(v.X*s, v.Y*s, v.Z*s)
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

type Cube3D struct {
	colour               util.Pixel
	pos                  util.Vector2[int]
	size, angle1, angle2 int
}

func (c *Cube3D) drawTo(buf *util.ScreenBuffer) {
	// rotating cube in integer 3D
	angle1 := float64(c.angle1) / 10
	angle2 := float64(c.angle2) / 10
	sin1, cos1 := math.Sin(angle1), math.Cos(angle1)
	sin2, cos2 := math.Sin(angle2), math.Cos(angle2)

	// 8 corners of the cube
	points := [8]util.Vector2[int]{}
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

		points[i] = util.V2(sx, sy)
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
		drawLine(buf, points[e[0]], points[e[1]], c.colour)
	}
}

type Tri struct {
	a, b, c Vector3
	colour  util.Pixel
}

func tri(a, b, c Vector3, colour util.Pixel) *Tri {
	return &Tri{a: a, b: b, c: c, colour: colour}
}

func drawTriangle3D(tri *Tri, s *Scene3D, buf *util.ScreenBuffer, fov float32, cos, sin Vector3) {
	// Transform vertices to camera space
	p1 := V3(tri.a.X-s.camera.X, tri.a.Y-s.camera.Y, tri.a.Z-s.camera.Z)
	p2 := V3(tri.b.X-s.camera.X, tri.b.Y-s.camera.Y, tri.b.Z-s.camera.Z)
	p3 := V3(tri.c.X-s.camera.X, tri.c.Y-s.camera.Y, tri.c.Z-s.camera.Z)

	// Apply rotation (Y, X, Z order)
	rotate := func(p Vector3) Vector3 {
		// Y axis
		x, z := p.X*cos.Y-p.Z*sin.Y, p.X*sin.Y+p.Z*cos.Y
		// X axis
		y := p.Y*cos.X - z*sin.X
		z = p.Y*sin.X + z*cos.X
		// Z axis
		x, y = x*cos.Z-y*sin.Z, x*sin.Z+y*cos.Z
		return V3(x, y, z)
	}
	p1r, p2r, p3r := rotate(p1), rotate(p2), rotate(p3)

	// Cull triangles behind camera
	if p1r.Z <= 0 || p2r.Z <= 0 || p3r.Z <= 0 {
		return
	}

	// Perspective projection
	project := func(p Vector3) util.Vector2[int] {
		return util.V2(
			int(fov*p.X/p.Z)+util.Width/2,
			util.Height/2-int(fov*p.Y/p.Z),
		)
	}
	s1, s2, s3 := project(p1r), project(p2r), project(p3r)

	// Simple distance-based shading
	dist := Distance(tri.a, s.camera) + Distance(tri.b, s.camera) + Distance(tri.c, s.camera)
	fade := max(0xff-int(dist), 0)
	c := tri.colour.Brightness(uint8(fade))

	if fade > 0 {
		drawTriangle2D(buf, s1, s2, s3, c)
	}
}

type Object3D struct {
	Tris []*Tri
}

func NewCube3D(colour util.Pixel, position Vector3, size float32) *Object3D {
	// 8 vertices
	var points [8]Vector3
	for i := range points {
		x := position.X + size*(float32(i>>0&1)-0.5)
		y := position.Y + size*(float32(i>>1&1)-0.5)
		z := position.Z + size*(float32(i>>2&1)-0.5)
		points[i] = V3(x, y, z)
	}

	// 12 triangles (2 per face)
	tris := []*Tri{
		tri(points[0], points[1], points[2], colour), tri(points[1], points[3], points[2], colour), // front
		tri(points[4], points[5], points[6], colour), tri(points[5], points[7], points[6], colour), // back
		tri(points[0], points[1], points[4], colour), tri(points[1], points[5], points[4], colour), // left
		tri(points[2], points[3], points[6], colour), tri(points[3], points[7], points[6], colour), // right
		tri(points[0], points[2], points[4], colour), tri(points[2], points[6], points[4], colour), // top
		tri(points[1], points[3], points[5], colour), tri(points[3], points[7], points[5], colour), // bottom
	}

	return &Object3D{Tris: tris}
}

type Scene3D struct {
	objects                []*Object3D
	camera, cameraRotation Vector3
	zoom                   int
}

func (s *Scene3D) drawTo(buf *util.ScreenBuffer) {
	var tris []*Tri
	for _, obj := range s.objects {
		tris = append(tris, obj.Tris...)
	}

	// sort tris by distance from camera (painter's algorithm)
	slices.SortFunc(tris, func(a, b *Tri) int {
		d1 := Distance(V3(a.a.X+a.b.X+a.c.X, a.a.Y+a.b.Y+a.c.Y, a.a.Z+a.b.Z+a.c.Z), s.camera)
		d2 := Distance(V3(b.a.X+b.b.X+b.c.X, b.a.Y+b.b.Y+b.c.Y, b.a.Z+b.b.Z+b.c.Z), s.camera)
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

	for _, tri := range tris {
		drawTriangle3D(tri, s, buf, fov, cos, sin)
	}
}

type Circle struct {
	centre util.Vector2[int]
	radius int
	colour util.Pixel
}

func (c *Circle) drawTo(buf *util.ScreenBuffer) {
	x0, y0, r := c.centre.X, c.centre.Y, c.radius

	// draw filled circle
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r {
				buf.Set(x0+x, y0+y, c.colour)
			}
		}
	}
}

type Pendulum struct {
	length, angle, mass float64
}

func (p *Pendulum) position(origin util.Vector2[int]) util.Vector2[int] {
	return util.V2(
		origin.X+int(p.length*math.Sin(p.angle)),
		origin.Y+int(p.length*math.Cos(p.angle)),
	)
}

type DoublePendulum struct {
	origin util.Vector2[int]
	p1, p2 Pendulum
}

func (dp *DoublePendulum) drawTo(buf *util.ScreenBuffer) {
	// calculate positions
	pos1 := dp.p1.position(dp.origin)
	pos2 := dp.p2.position(pos1)

	// draw arms
	drawLine(buf, dp.origin, pos1, util.White)
	drawLine(buf, pos1, pos2, util.White)

	// draw bobs, scaled by mass
	circle1 := &Circle{pos1, int(math.Sqrt(dp.p1.mass/math.Pi) * 5), util.Red}
	circle2 := &Circle{pos2, int(math.Sqrt(dp.p2.mass/math.Pi) * 5), util.Blue}
	circle1.drawTo(buf)
	circle2.drawTo(buf)
}

type Mass struct {
	position, velocity util.Vector2[float64]
	mass, charge       float64
	colour             util.Pixel
}

type Crosshair struct {
	pos    util.Vector2[int]
	colour util.Pixel
	size   int
}

func (c *Crosshair) drawTo(buf *util.ScreenBuffer) {
	drawHLine(buf, c.pos.Y, c.pos.X-c.size, c.pos.X+c.size, c.colour)
	drawVLine(buf, c.pos.X, c.pos.Y-c.size, c.pos.Y+c.size, c.colour)
}

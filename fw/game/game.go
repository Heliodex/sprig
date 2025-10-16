package game

import (
	"runtime"
	"strconv"
	"time"

	"fw/util"
)

type UIElement interface {
	drawTo(*util.ScreenBuffer)
}

// intro to show, in the event that something else is loading or to show information
func Splash(en util.Engine) {
	freq := en.CPUFrequency()

	engineText := &Text{FontDex, "Grips Engine", util.V2(25, 5), util.Red} // I'm calling it this because it's an anagram of Sprig
	buildText := &Text{FontUnifont, strconv.Itoa(runtime.NumCPU()) + " cores", util.V2(5, 35), util.Blue}
	freqText := &Text{FontUnifont, strconv.Itoa(int(freq/1_000_000)) + "MHz", util.V2(5, 50), util.Green}

	ui := []UIElement{engineText, buildText, freqText}
	if freq >= 270_000_000 {
		ui = append(ui, &Text{FontUnifont, "OVERCLOCKED!", util.V2(5, 65), util.Red})
	} else {
		ui = append(ui, &Text{FontUnifont, "(could be better)", util.V2(5, 65), util.White})
	}

	for _, e := range ui {
		e.drawTo(en.ScreenBuffer())
	}
	en.Render()

	time.Sleep(time.Second)
	println("Splash complete")
}

type Line struct {
	start, end util.Vector2
	colour     util.Pixel
}

var (
	lines []*Line
	f     int
	// position   util.Vector2
	position = V3(0, 2, 0)
	rotation Vector3
)

// ran every frame (or, more like this is what makes the frames)
func Update(en util.Engine) {
	btns := en.Buttons()

	if btns[util.W].Pressed() {
		// forwardX, forwardZ := f32Sin(rotation.Y)*f32Cos(rotation.X), f32Cos(rotation.Y)*f32Cos(rotation.X)
		// position.X += forwardX
		// position.Z += forwardZ
		position.X += f32Sin(rotation.Y) * f32Cos(rotation.X)
		position.Z += f32Cos(rotation.Y) * f32Cos(rotation.X)
	}

	if btns[util.S].Pressed() {
		// backwardX, backwardZ := -f32Sin(rotation.Y)*f32Cos(rotation.X), -f32Cos(rotation.Y)*f32Cos(rotation.X)
		// position.X += backwardX
		// position.Z += backwardZ
		position.X -= f32Sin(rotation.Y) * f32Cos(rotation.X)
		position.Z -= f32Cos(rotation.Y) * f32Cos(rotation.X)
	}

	if btns[util.A].Pressed() {
		// leftX, leftZ := -f32Cos(rotation.Y), f32Sin(rotation.Y)
		// position.X += leftX
		// position.Z += leftZ
		position.X -= f32Cos(rotation.Y)
		position.Z += f32Sin(rotation.Y)
	}

	if btns[util.D].Pressed() {
		// rightX, rightZ := f32Cos(rotation.Y), -f32Sin(rotation.Y)
		// position.X += rightX
		// position.Z += rightZ
		position.X += f32Cos(rotation.Y)
		position.Z -= f32Sin(rotation.Y)
	}

	if btns[util.I].Pressed() {
		rotation.X += 0.06
	}

	if btns[util.K].Pressed() {
		rotation.X -= 0.06
	}

	if btns[util.J].Pressed() {
		rotation.Y -= 0.06
	}

	if btns[util.L].Pressed() {
		rotation.Y += 0.06
	}

	x := &Text{FontUnifont, strconv.Itoa(int(rotation.X)), util.V2(70, 2), util.Green}
	y := &Text{FontUnifont, strconv.Itoa(int(rotation.Y)), util.V2(70, 22), util.Green}
	z := &Text{FontUnifont, strconv.Itoa(int(rotation.Z)), util.V2(70, 42), util.Green}

	Texts := [util.ButtonsCount]*Text{
		{FontDex, "W", util.V2(2+20-2, 2), util.Red},
		{FontDex, "A", util.V2(2, 26), util.Red},
		{FontDex, "S", util.V2(2+20, 50), util.Red},
		{FontDex, "D", util.V2(2+40, 26), util.Red},
		{FontDex, "I", util.V2(109+20+1, 2), util.Red},
		{FontDex, "J", util.V2(109, 26), util.Red},
		{FontDex, "K", util.V2(109+20, 50), util.Red},
		{FontDex, "L", util.V2(109+40, 26), util.Red},
	}

	// read button states
	// var leftPressed, rightPressed bool
	for i, b := range btns {
		if b.Pressed() {
			// if i < int(util.ButtonsCount/2) {
			// 	leftPressed = true
			// } else {
			// 	rightPressed = true
			// }
			Texts[i].colour = util.Green
		} else {
			Texts[i].colour = util.Red
		}
	}

	scene := &Scene3D{
		camera:         position,
		cameraRotation: rotation,
		zoom:           110, // 72 or so degrees
		objects: []*Object3D{
			NewCube3D(util.White, V3(20, 5, 5), 10),
			NewCube3D(util.Red, V3(0, 2.5, -10), 5),
			NewCube3D(util.Green, V3(-10, 5, 5), 10),
			NewCube3D(util.Blue, V3(-20, 10, 30), 20),
			NewCube3D(util.Yellow, V3(10, 2.5, 25), 5),
		},
	}

	ui := []UIElement{x, y, z}
	for _, t := range Texts {
		ui = append(ui, t)
	}
	ui = append(ui, scene)

	for _, e := range ui {
		e.drawTo(en.ScreenBuffer())
	}

	en.Render()
	f++
}

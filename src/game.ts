import { type WebEngineAPI, keyMap } from "./webEngine"
import { black, darkgrey, colours, white, lightgrey } from "./colours"

const space = " ".charCodeAt(0)
const allChars = [space, ...colours]

function transformTextures(startTexs: string[]): string[][] {
	const textures: string[][] = []
	for (const tex of startTexs) {
		const texLines = tex
			.replaceAll("\r", "")
			.replaceAll("\t", "")
			.split("\n")
		texLines.shift()

		for (const line of texLines)
			for (const char of line)
				if (!allChars.includes(char.charCodeAt(0)))
					throw new Error(`Unknown colour: ${char.charCodeAt(0)}`)

		textures.push(texLines)
	}
	return textures
}

abstract class Sprite {
	textures: readonly string[][] = []
	#interval = 300
	pos = [0, 0]
	offset = [0, 0]

	get texture() {
		// random texture that changes every [interval] ms
		const index =
			Math.floor(Date.now() / this.#interval) % this.textures.length
		return this.textures[index]
	}

	get size() {
		const width = Math.max(...this.texture.map(l => l.length))
		const height = this.texture.length
		return [width, height]
	}
}

class Ship extends Sprite {
	textures = transformTextures([
		`
			    3
			    2
			   323
			   323
			  33233
			  33233
			 3232323
			 3232323
			3323 3233
			333   333
			 LL   LL
			 96   69
			 6     6`,
		`
			    3
			    2
			   323
			   323
			  33233
			  33233
			 3232323
			 3232323
			3323 3233
			333   333
			 LL   LL
			 99   99
			 96   69
			 6     6`,
	])
	pos = [80, 120]
	offset = [5, 13]
}

class Bullet extends Sprite {
	textures = transformTextures([
		`
			7
			7
			7`,
	])
	offset = [1, 0]
	constructor(x: number) {
		super()
		this.pos = [x, 106]
	}
}

class Enemy extends Sprite {
	// supposed to be a funky planet or asteroid
	textures = transformTextures([
		`
			  LL LLL
			 LH1LHHHL
			1HHHHH1HL 
			LL1HHHHHHL
			1HHHHHLLHL
			1HHH1HHHHL
			1HH1HHHHH1
			1HHLHHLHHL
			 11HHH1H1
			   111 1 `,
	])
	offset = [6, 12]
	constructor(x: number) {
		super()
		this.pos = [x, 0]
	}
}

const ship = new Ship()

export default async (api: WebEngineAPI) => {
	const { onInput, setLegend, setMap } = api

	const spriteLocMap = new Map<number, Sprite>()
	const display = new Int8Array(160 * 128)
	const getPixel = (pos: number) => String.fromCharCode(display[pos])
	function setPixel(pos: number, colour: number) {
		display[pos] = colour
	}

	function starfieldLine(chance: number) {
		const line: boolean[] = []
		for (let i = 0; i < 160; i++) line.push(Math.random() < chance)
		return line
	}

	// darkgrey, lightgrey, white
	const starfields: boolean[][][] = []
	const starcolours = [darkgrey, darkgrey, lightgrey, white]
	const chances = [0.0005, 0.0015, 0.001, 0.0005]
	const updateSpeed = [31, 5, 2, 1]
	const framesSinceLastUpdate = [0, 0, 0, 0]

	for (const field in starcolours) {
		starfields[field] = []
		for (let i = 0; i < 128; i++)
			starfields[field].push(starfieldLine(chances[field]))
	}

	// movement
	let moveDirection = ""
	onInput("a", () => {
		if (moveDirection) return
		moveDirection = "left"
	})
	onInput("d", () => {
		if (moveDirection) return
		moveDirection = "right"
	})
	window.addEventListener("keyup", e => {
		const key = keyMap[e.key]
		if (!key) return
		if (key === "a" || key === "d") moveDirection = ""
	})

	// sprites
	const sprites = new Set<Sprite>([ship])

	// bullets
	let framesSinceLastBullet = 10

	onInput("w", () => {
		if (framesSinceLastBullet < 10) return
		framesSinceLastBullet = 0
		sprites.add(new Bullet(ship.pos[0]))
	})

	// onInput("s", () => {
	sprites.add(new Enemy(ship.pos[0]))
	// })

	function handleCollision(a: Sprite, b: Sprite) {
		if (a instanceof Bullet && b instanceof Enemy) {
			sprites.delete(a)
			sprites.delete(b)
		} else if (a instanceof Enemy && b instanceof Ship) {
			sprites.delete(a)
			sprites.delete(b)
		}
	}

	function render() {
		const startTime = Date.now()

		display.fill(black)
		spriteLocMap.clear()

		// move ship
		if (moveDirection === "left") ship.pos[0] -= 1
		else if (moveDirection === "right") ship.pos[0] += 1

		const o = ship.offset[0]
		if (ship.pos[0] < o) ship.pos[0] = o
		else if (ship.pos[0] > 160 - o + 1) ship.pos[0] = 160 - o + 1

		// move bullets
		const bullets = [...sprites].filter(s => s instanceof Bullet)
		for (const bullet of bullets) {
			bullet.pos[1] -= 3
			if (bullet.pos[1] < 0) sprites.delete(bullet)
		}
		const enemies = [...sprites].filter(s => s instanceof Enemy)
		for (const enemy of enemies) {
			enemy.pos[1] += 0.5
			if (enemy.pos[1] > 142) sprites.delete(enemy)
		}

		// starfields
		for (const field in starcolours) {
			framesSinceLastUpdate[field]++
			if (framesSinceLastUpdate[field] >= updateSpeed[field]) {
				framesSinceLastUpdate[field] = 0

				starfields[field].pop()
				starfields[field].unshift(starfieldLine(chances[field]))
			}

			const colour = starcolours[field]

			for (let y = 0; y < 128; y++)
				for (let x = 0; x < 160; x++)
					if (starfields[field][y][x]) {
						const pos = y * 160 + x
						if (display[pos] === black) display[pos] = colour
						else display[pos] = lightgrey
					}
		}

		framesSinceLastBullet++

		for (const sprite of sprites) {
			const toDraw = sprite.texture

			const [xPos, yPos] = sprite.pos.map(Math.floor)
			const [xOffset, yOffset] = sprite.offset

			const collisionsDetected = new Set<Sprite>()

			for (let y = 0; y < toDraw.length; y++)
				for (let x = 0; x < toDraw[y].length; x++)
					if (toDraw[y].charCodeAt(x) !== space) {
						// it's drawing time
						const pos =
							(y + yPos - yOffset) * 160 + (x + xPos - xOffset)
						if (pos < 0 || pos >= 160 * 128) continue

						const possibleCollision = spriteLocMap.get(pos)
						if (
							possibleCollision &&
							!collisionsDetected.has(possibleCollision)
						) {
							collisionsDetected.add(possibleCollision)
							console.log("collided with", possibleCollision.constructor.name)
							handleCollision(sprite, possibleCollision)
						}
						setPixel(pos, toDraw[y].charCodeAt(x))
						spriteLocMap.set(pos, sprite)
					}
		}

		// screen paint
		const legend: [string, string][] = []
		const screen: string[] = []

		for (let my = 0; my < 8; my++) {
			for (let mx = 0; mx < 10; mx++) {
				const char = String.fromCharCode(mx + my * 10 + 47) // start after .

				const bitmap: string[] = []
				for (let ly = 0; ly < 16; ly++) {
					for (let lx = 0; lx < 16; lx++)
						bitmap.push(
							getPixel((my * 16 + ly) * 160 + mx * 16 + lx)
						)
					bitmap.push("\n")
				}

				legend.push([char, bitmap.join("")])
				screen.push(char)
			}
			screen.push("\n")
		}

		setLegend(...legend)
		setMap(screen.join(""))

		const t = Date.now() - startTime
		if (t > 16) console.log("rendered in", t, "ms")

		requestAnimationFrame(render)
	}

	render()
}

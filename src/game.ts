import { type WebEngineAPI, keyMap } from "./webEngine"
import { black, darkgrey, colours, white, lightgrey } from "./colours"

const space = " ".charCodeAt(0)
const allChars = [space, ...colours]

class Sprite {
	type: "ship" | "bullet" | "enemy" = "ship"
	#textures: string[][] = []
	#interval = 300
	size: [number, number] = [0, 0]
	pos: [number, number] = [0, 0]
	offset: [number, number] = [0, 0]
	constructor(props: {
		type: typeof Sprite.prototype.type
		textures: string[]
		pos: typeof Sprite.prototype.pos
		offset: typeof Sprite.prototype.offset
		interval?: number
	}) {
		for (const i in props.textures) {
			const texture = props.textures[i]
			const texLines = texture
				.replaceAll("\r", "")
				.replaceAll("\t", "")
				.split("\n")
			if (texLines[0] === "") texLines.shift()

			for (const line of texLines)
				for (const char of line)
					if (!allChars.includes(char.charCodeAt(0)))
						throw new Error(`Unknown colour: ${char.charCodeAt(0)}`)

			// get longest line in texture
			const longest = Math.max(...texLines.map(l => l.length))
			this.type = props.type
			this.#textures[i] = texLines
			this.size = [longest, texLines.length]
			this.pos = props.pos
			this.offset = props.offset
		}
	}

	get texture() {
		// random texture that changes every [interval] ms
		const index =
			Math.floor(Date.now() / this.#interval) % this.#textures.length
		return this.#textures[index]
	}
}

const ship = new Sprite({
	type: "ship",
	textures: [
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
	],
	pos: [80, 120],
	offset: [5, 13],
})

const bulletTex = `
	7
	7
	7`

// supposed to be a funky planet or asteroid
const enemyTexs = [
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
]

export default async (api: WebEngineAPI) => {
	const { onInput, setLegend, setMap } = api

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
	const newBullet = (x: number) =>
		new Sprite({
			type: "bullet",
			textures: [bulletTex],
			pos: [x, 108],
			offset: [1, 0],
		})

	onInput("w", () => {
		sprites.add(newBullet(ship.pos[0]))
	})

	// enemies
	const newEnemy = (x: number) =>
		new Sprite({
			type: "enemy",
			textures: enemyTexs,
			pos: [x, 0],
			offset: [6, 12],
		})

	// onInput("s", () => {
	sprites.add(newEnemy(ship.pos[0]))
	// })

	async function render() {
		const startTime = Date.now()

		display.fill(black)

		// move ship
		if (moveDirection === "left") ship.pos[0] -= 1
		else if (moveDirection === "right") ship.pos[0] += 1

		const o = ship.offset[0]
		if (ship.pos[0] < o) ship.pos[0] = o
		else if (ship.pos[0] > 160 - o + 1) ship.pos[0] = 160 - o + 1

		// move bullets
		const bullets = [...sprites].filter(s => s.type === "bullet")
		for (const bullet of bullets) {
			bullet.pos[1] -= 3
			if (bullet.pos[1] < 0) sprites.delete(bullet)
		}
		const enemies = [...sprites].filter(s => s.type === "enemy")
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

		for (const sprite of sprites) {
			const toDraw = sprite.texture

			const [xPos, yPos] = sprite.pos.map(Math.floor)
			const [xOffset, yOffset] = sprite.offset

			for (let y = 0; y < toDraw.length; y++)
				for (let x = 0; x < toDraw[y].length; x++)
					if (toDraw[y].charCodeAt(x) !== space)
						setPixel(
							(y + yPos - yOffset) * 160 + (x + xPos - xOffset),
							toDraw[y].charCodeAt(x)
						)
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

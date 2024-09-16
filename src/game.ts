import type { WebEngineAPI } from "./webEngine"
import { black, darkgrey, colours, white, lightgrey } from "./colours"

const space = " ".charCodeAt(0)
const allChars = [space, ...colours]

class Sprite {
	#texture: string[][] = []
	#interval = 300
	size: [number, number] = [0, 0]
	pos: [number, number] = [0, 0]
	offset: [number, number] = [0, 0]
	constructor(props: {
		textures: string[]
		pos: [number, number]
		offset: [number, number]
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
			this.#texture[i] = texLines
			this.size = [longest, texLines.length]
			this.pos = props.pos
			this.offset = props.offset
		}
	}

	get texture() {
		// random texture that changes every [interval] ms
		const index =
			Math.floor(Date.now() / this.#interval) % this.#texture.length
		return this.#texture[index]
	}
}

const sprite = new Sprite({
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

	async function render() {
		display.fill(black)

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

		const startTime = Date.now()

		// drawing
		const toDraw = sprite.texture

		const [xPos, yPos] = sprite.pos
		const [xOffset, yOffset] = sprite.offset

		for (let y = 0; y < toDraw.length; y++)
			for (let x = 0; x < toDraw[y].length; x++)
				if (toDraw[y].charCodeAt(x) !== space)
					setPixel(
						(y + yPos - yOffset) * 160 + (x + xPos - xOffset),
						toDraw[y].charCodeAt(x)
					)

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

	onInput("w", () => {
		console.log("w")
	})
}

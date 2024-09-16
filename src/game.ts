import type { WebEngineAPI } from "./webEngine"
import { black, darkgrey, transparent, colours, white } from "./colours"

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
	const {
		// afterInput,
		// bitmap,
		// getAll,
		// getFirst,
		// getState,
		// map,
		onInput,
		// playTune,
		setLegend,
		setMap,
		// setPushables,
		// setSolids,
	} = api

	const display = new Int8Array(160 * 128)

	const getPixel = (pos: number) => String.fromCharCode(display[pos])
	function setPixel(pos: number, colour: number) {
		display[pos] = colour
	}

	function update() {
		// sprite.pos[1] -= 1
		setTimeout(update, 100)
	}

	update()

	async function render() {
		display.fill(darkgrey)
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

import type { SpriteType } from "sprig"
import type { WebEngineAPI } from "./webEngine"

const detab = (s: string) => s.replace(/\t/g, "")

// double size both horizontally and vertically
function doublesize(s: string) {
	const lines = s.split("\n")
	const newLines = []
	for (const line of lines) {
		const newLine = []
		for (const char of line) {
			newLine.push(char)
			newLine.push(char)
		}
		newLines.push(newLine.join(""))
		newLines.push(newLine.join(""))
	}
	return newLines.join("\n")
}

const char = detab(`
	.0.0.0.0.0.0.0.0
	0.0.0.0.0.0.0.0.
	.0.0.0.0.0.0.0.0
	0.0.0.0.0.0.0.0.
	.0.0.0.0.0.0.0.0
	0.0.0.0.0.0.0.0.
	.0.0.0.0.0.0.0.0
	0.0.0.0.0.0.0.0.
	.0.0.0.0.0.0.0.0
	0.0.0.0.0.0.0.0.
	.0.0.0.0.0.0.0.0
	0.0.0.0.0.0.0.0.
	.0.0.0.0.0.0.0.0
	0.0.0.0.0.0.0.0.
	.0.0.0.0.0.0.0.0
	0.0.0.0.0.0.0.0.`)

export default async (api: WebEngineAPI) => {
	const {
		afterInput,
		bitmap,
		getAll,
		getFirst,
		getState,
		map,
		onInput,
		playTune,
		setLegend,
		setMap,
		setPushables,
		setSolids,
	} = api

	setSolids([])

	const level = 0
	// 10x8 max size
	const levels = [
		detab(`
1.1.1.1.1.
.1.1.1.1.1
1.1.1.1.1.
.1.1.1.1.1
1.1.1.1.1.
.1.1.1.1.1
1.1.1.1.1.
.1.1.1.1.1
`),
	]
	const legend: [string, string][] = []

	function get(id: string) {
		const got = getFirst(id)
		if (!got) throw new Error(`Sprite ${id} not found`)
		return got
	}

	class Player {
		spriteMap = [
			{
				id: "1",
				sprite: char,
				offset: [0, 0],
			},
			{
				id: "2",
				sprite: char,
				offset: [1, 0],
			},
			{
				id: "3",
				sprite: char,
				offset: [0, 1],
			},
			{
				id: "4",
				sprite: char,
				offset: [1, 1],
			},
		]
		sprites: { [k: string]: SpriteType } = {}
		constructor() {
			for (const { id, sprite } of this.spriteMap) {
				legend.push([id, sprite])
				setLegend(...legend)
			}
		}

		get x() {
			return get("1").x
		}
		get y() {
			return get("1").y
		}
		set x(newX) {
			// set x of all sprites, taking into account offsets
			for (const { id, offset } of this.spriteMap) {
				const sprite = get(id)
				sprite.x = newX + offset[0]
			}
		}
		set y(newY) {
			// set y of all sprites, taking into account offsets
			for (const { id, offset } of this.spriteMap) {
				const sprite = get(id)
				sprite.y = newY + offset[1]
			}
		}
	}
	const player = new Player()

	const charNums: [string, string][] = []

	setMap(levels[level])

	onInput("w", () => {
		player.y -= 1
	})
	onInput("a", () => {
		player.x -= 1
	})
	onInput("s", () => {
		player.y += 1
	})
	onInput("d", () => {
		player.x += 1
	})

	afterInput(() => {})
}

import type { WebEngineAPI } from "./webEngine"

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

	const display = new Int8Array(160 * 128)

	function getPixel(pos: number) {
		return display[pos] ? "1" : "."
	}

	let currentPixelPos = 0
	function update() {
		currentPixelPos++
		display.fill(0)
		display[currentPixelPos] = 1
	}

	function render() {
		console.log("rendering")

		const legend: [string, string][] = []

		let screen = ""
		for (let my = 0; my < 8; my++) {
			for (let mx = 0; mx < 10; mx++) {
				const char = String.fromCharCode(mx + my * 10 + 47) // start after .

				let bitmap = ""
				for (let ly = 0; ly < 16; ly++) {
					for (let lx = 0; lx < 16; lx++)
						bitmap += getPixel((mx * 16 + lx) * 128 + my * 16 + ly)
					bitmap += "\n"
				}

				legend.push([char, bitmap])

				screen += char
			}
			screen += "\n"
		}

		setLegend(...legend)
		setMap(screen)

		update()

		setTimeout(render, 20) // 50fps loll
	}

	render()
}

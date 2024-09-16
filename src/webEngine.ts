import {
	type InputKey,
	type PlayTuneRes,
	VALID_INPUTS,
	type FullSprigAPI,
	type GameState,
} from "../node_modules/sprig/src/api.js"
import {
	type BaseEngineAPI,
	baseEngine,
	textToTune,
} from "../node_modules/sprig/src/base/index.js"
import { bitmapTextToImageData } from "../node_modules/sprig/src/image-data/index.js"
import { getTextImg } from "../node_modules/sprig/src/web/text.js"
import { playTune } from "../node_modules/sprig/src/web/tune.js"
import { makeCanvas } from "../node_modules/sprig/src/web/util.js"

export * from "../node_modules/sprig/src/web/text.js"
export * from "../node_modules/sprig/src/web/tune.js"

export type WebEngineAPI = BaseEngineAPI &
	Pick<FullSprigAPI, "setLegend" | "onInput" | "afterInput" | "playTune"> & {
		getState(): GameState // For weird backwards-compatibility reasons, not part of API
	}

type Ret = {
	api: WebEngineAPI
	state: GameState
	cleanup(): void
}
export default (canvas: HTMLCanvasElement): Ret => {
	const { api, state } = baseEngine()

	const ctx = canvas.getContext("2d")
	const offscreenCanvas = makeCanvas(1, 1)
	const offscreenCtx = offscreenCanvas.getContext("2d")
	if (!ctx || !offscreenCtx) throw new Error("No 2d context")

	const _bitmaps: Record<string, CanvasImageSource> = {}
	let _zOrder: string[] = []

	ctx.imageSmoothingEnabled = false

	const _gameloop = (): void => {
		const width = 10
		const height = 8
		const actualWidth = width * 16
		const actualHeight = height * 16

		ctx.clearRect(0, 0, canvas.width, canvas.height)

		offscreenCanvas.width = actualWidth
		offscreenCanvas.height = actualHeight

		offscreenCtx.fillStyle = "white"
		offscreenCtx.fillRect(0, 0, actualWidth, actualHeight)

		const grid = api.getGrid()

		for (let i = 0; i < width * height; i++) {
			const x = i % width
			const y = Math.floor(i / width)
			const sprites = grid[i]

			if (state.background) {
				const imgData = _bitmaps[state.background]
				offscreenCtx.drawImage(imgData, x * 16, y * 16)
			}

			for (const sprite of sprites.sort(
				(a, b) => _zOrder.indexOf(b.type) - _zOrder.indexOf(a.type)
			)) {
				const imgData = _bitmaps[sprite.type]
				offscreenCtx.drawImage(imgData, x * 16, y * 16)
			}
		}

		ctx.drawImage(
			offscreenCanvas,
			(canvas.width - actualWidth) / 2,
			(canvas.height - actualHeight) / 2,
			actualWidth,
			actualHeight
		)

		const textCanvas = getTextImg(state.texts)
		ctx.drawImage(textCanvas, 0, 0, canvas.width, canvas.height)

		animationId = window.requestAnimationFrame(_gameloop)
	}
	let animationId = window.requestAnimationFrame(_gameloop)

	const setLegend = (...bitmaps: [string, string][]): void => {
		if (bitmaps.length === 0)
			throw new Error(
				"There needs to be at least one sprite in the legend."
			)

		if (!Array.isArray(bitmaps[0]))
			throw new Error(
				"The sprites passed into setLegend each need to be in square brackets, like setLegend with [player, bitmap`...`]."
			)

		for (const [key] of bitmaps) {
			if (key === ".") throw new Error(`Can't reassign "." bitmap`)
			if (key.length !== 1)
				throw new Error("Bitmaps must have one character names")
		}

		state.legend = bitmaps
		_zOrder = bitmaps.map(x => x[0])

		for (let i = 0; i < bitmaps.length; i++) {
			const [key, value] = bitmaps[i]
			const imgData = bitmapTextToImageData(key, value)
			const littleCanvas = makeCanvas(16, 16)
			const ctx = littleCanvas.getContext("2d")
			if (!ctx) throw new Error("No 2d context")

			ctx.putImageData(imgData, 0, 0)
			_bitmaps[key] = littleCanvas
		}
	}

	const tileInputs: Record<InputKey, (() => void)[]> = {
		w: [],
		s: [],
		a: [],
		d: [],
		i: [],
		j: [],
		k: [],
		l: [],
	}
	const afterInputs: (() => void)[] = []

	function keydown(e: KeyboardEvent) {
		// monkeypatch
		const keyMap: { [k: string]: InputKey } = {
			",": "w",
			a: "a",
			o: "s",
			e: "d",
			c: "i",
			h: "j",
			t: "k",
			n: "l",
		}
		const key = keyMap[e.key]
		if (!VALID_INPUTS.includes(key)) return

		for (const validKey of VALID_INPUTS)
			if (key === validKey) for (const fn of tileInputs[key]) fn()
		for (const f of afterInputs) f()
		for (const sprite of state.sprites) {
			const s = sprite as {
				dx: number
				dy: number
			}
			s.dx = 0
			s.dy = 0
		}

		e.preventDefault()
	}
	canvas.addEventListener("keydown", keydown)

	function onInput(key: InputKey, fn: () => void): void {
		if (!VALID_INPUTS.includes(key))
			throw new Error(
				`Unknown input key, "${key}": expected one of ${VALID_INPUTS.join(", ")}`
			)
		tileInputs[key].push(fn)
	}
	const afterInput = (fn: () => void): void => {
		afterInputs.push(fn)
	}

	const tunes: PlayTuneRes[] = []
	return {
		api: {
			...api,
			setLegend,
			onInput,
			afterInput,
			getState: () => state,
			playTune: (text: string, n: number) => {
				const tune = textToTune(text)
				const playTuneRes = playTune(tune, n)
				tunes.push(playTuneRes)
				return playTuneRes
			},
		},
		state,
		cleanup: () => {
			ctx.clearRect(0, 0, canvas.width, canvas.height)
			window.cancelAnimationFrame(animationId)
			canvas.removeEventListener("keydown", keydown)
			for (const tune of tunes) tune.end()
		},
	}
}

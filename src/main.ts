import "./style.css"
import webEngine from "./webEngine.ts"
import gamecode from "./game.ts"

const canvas = document.querySelector<HTMLCanvasElement>("#canvas")
if (!canvas) throw new Error("no canvas")

const game = webEngine(canvas)
gamecode(game.api)


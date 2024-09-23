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
	currentTex = 0
	framesSinceLastTex = 0
	pos = [0, 0]
	offset = [0, 0]
	collidable = true

	get texture() {
		return this.textures[this.currentTex % this.textures.length]
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
	velocity = [0, 0.5]
	constructor(pos: number[], velocity: number[]) {
		super()
		this.pos = pos
		this.velocity = velocity
	}
}

class SmallEnemy extends Sprite {
	textures = transformTextures([
		`

		 	   LLL
		 	  LHHHL
		 	 1HLHHL
		 	 1H1HH1L
		 	 1HHHHHL
		 	  1HLH1
		 	   111 `,
	])
	offset = [6, 12]
	velocity = [0, 0.5]
	immunity = 15
	constructor(pos: number[], velocity: number[]) {
		super()
		this.pos = pos
		this.velocity = velocity
	}
}

class Explosion extends Sprite {
	textures = transformTextures([
		`
		    33 63 3
		    3 6 636
		     66996 3
		    9 96339 3
		    969333969
		    9 9336963
		     36999363
		    33636 6 3
		       996`,
		`
		      96636
		     66999 6
		    3 9939996
		    6993 3939
		    693 6 396
		    6393 39 6
		    399939996
		     6 69993
		      69663`,

		`
		        9
		      66666
		     6999996
		   3 699 996
		    969 C 969
		     699 996
		    96999996
		      66666 3
		        9`,
		`

		        9
		    3  696 9
		      6   6
		     99 C 99
		      6   6
		    9  696  3
		        9`,
		`

		           9
		     3  6
		       9 9
		      6   6
		       9 9 6
		     9  6    3
		         `,
	])
	collidable = false
	offset = [8, 8]
	constructor(pos: number[]) {
		super()
		this.pos = structuredClone(pos)
	}
}

class GameOver extends Sprite {
	textures = transformTextures([
		`
			222222222222222222222222222222222222222
			200000000000000000000000000000000000002
			200000000000000000000000000000000000002
			200000000000000000000000000000000000002
			200002222200000200000200000200222220002
			200020000000000200000220002200200000002
			200020000000002020000202020200200000002
			200020000000002020000202020200222200002
			200020022200020002000200200200200000002
			200020000200022222000200200200200000002
			200020000200200000200200000200200000002
			200002222000200000200200000200222220002
			200000000000000000000000000000000000002
			200000000000000000000000000000000000002
			200000000000000000000000000000000000002
			200000222200020000020022222002222000002
			200002000020020000020020000002000200002
			200002000020002000200020000002000200002
			200002000020002000200022220002222000002
			200002000020000202000020000002200000002
			200002000020000202000020000002020000002
			200002000020000020000020000002002000002
			200000222200000020000022222002000200002
			200000000000000000000000000000000000002
			200000000000000000000000000000000000002
			200000000000000000000000000000000000002
			222222222222222222222222222222222222222
			`,
	])
	collidable = false
	pos = [80, 50]
	offset = [20, 0]
}

const scoreTexs = transformTextures([
	`
		 22
		2  2
		2  2
		2  2
		2  2
		2  2
		 22`,
	`
		  2
		 22
		  2
		  2
		  2
		  2
		 222`,
	`
		 22
		2  2
		   2
		  2
		 2
		2
		2222`,
	`
		2222
		  2
		 2
		222
		   2
		   2
		222`,
	`
		  2
		 22
		2 2
		2222
		  2
		  2
		  2`,
	`
		2222
		2
		222
		   2
		   2
		2  2
		 22`,
	`
		  2
		 2
		2
		222
		2  2
		2  2
		 22`,
	`
		2222
		   2
		   2
		  2
		  2
		 2
		 2`,
	`
		 22
		2  2
		2  2
		 22
		2  2
		2  2
		 22`,
	`
		 22
		2  2
		2  2
		 222
		   2
		  2
		 2`,
])

class ScoreNumber extends Sprite {
	textures = scoreTexs
	offset = [0, 0]
	pos = [0, 0]
	collidable = false
	constructor(pos: number, num: number) {
		super()
		this.pos = [pos * 5, 0]
		this.currentTex = num
	}
}

export default async (api: WebEngineAPI) => {
	const { onInput, setLegend, setMap } = api

	const spriteLocMap = new Map<number, Sprite>()
	const display = new Int8Array(160 * 128)
	const getPixel = (pos: number) => String.fromCharCode(display[pos])

	let gameState = "playing"
	let score = 0
	let stage = 1

	function drawSprite(sprite: Sprite) {
		const toDraw = sprite.texture

		const [xPos, yPos] = sprite.pos.map(Math.floor)
		const [xOffset, yOffset] = sprite.offset

		const collisionsDetected = new Set<Sprite>()

		for (let y = 0; y < toDraw.length; y++)
			for (let x = 0; x < toDraw[y].length; x++)
				if (toDraw[y].charCodeAt(x) !== space) {
					// it's drawing time
					const xp = x + xPos - xOffset
					const yp = y + yPos - yOffset
					if (xp < 0 || xp >= 160) continue

					const pos = yp * 160 + xp
					if (pos < 0 || pos >= 160 * 128) continue // necessary?

					const possibleCollision = spriteLocMap.get(pos)
					if (
						possibleCollision &&
						!collisionsDetected.has(possibleCollision) &&
						possibleCollision !== sprite &&
						possibleCollision.collidable &&
						sprite.collidable
					) {
						collisionsDetected.add(possibleCollision)
						handleCollision(sprite, possibleCollision)
					}
					display[pos] = toDraw[y].charCodeAt(x)
					spriteLocMap.set(pos, sprite)
				}
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
		if (gameState !== "playing") return
		if (moveDirection) return
		moveDirection = "left"
	})
	onInput("d", () => {
		if (gameState !== "playing") return
		if (moveDirection) return
		moveDirection = "right"
	})
	window.addEventListener("keyup", e => {
		if (gameState !== "playing") return
		const key = keyMap[e.key]
		if (!key) return
		if (key === "a" || key === "d") moveDirection = ""
	})

	// sprites
	const ship = new Ship()
	const gameOver = new GameOver()
	const sprites = new Set<Sprite>([ship])

	// bullets
	let framesSinceLastBullet = 10
	let framesSinceLastEnemy = 20
	const bulletCost = 10

	onInput("w", () => {
		if (
			gameState !== "playing" ||
			framesSinceLastBullet < 5 ||
			score < bulletCost
		)
			return

		score -= bulletCost
		framesSinceLastBullet = 0
		sprites.add(new Bullet(ship.pos[0]))
	})

	function handleCollision(...ss: [Sprite, Sprite]) {
		const [a, b] = ss.sort((s1, s2) =>
			s1.constructor.name.localeCompare(s2.constructor.name)
		)
		if (a instanceof Bullet && b instanceof Enemy) {
			sprites.delete(a)
			sprites.delete(b)
			sprites.add(new Explosion(b.pos))
			sprites.add(new SmallEnemy(b.pos, b.velocity))
			score += 200
		} else if (a instanceof Bullet && b instanceof SmallEnemy) {
			sprites.delete(a)
			sprites.delete(b)
			sprites.add(new Explosion(b.pos))
			score += 100
		} else if (
			(a instanceof Enemy && b instanceof Ship) ||
			(a instanceof Ship && b instanceof SmallEnemy)
		) {
			sprites.delete(a)
			sprites.delete(b)
			sprites.add(new Explosion(a.pos))
			sprites.add(new Explosion(b.pos))
			gameState = "over"
		} else if (
			a instanceof SmallEnemy &&
			b instanceof SmallEnemy &&
			a.immunity < 0 &&
			b.immunity < 0
		) {
			sprites.delete(a)
			sprites.delete(b)
			sprites.add(new Explosion(a.pos))
			sprites.add(new Explosion(b.pos))
		} else if (
			a instanceof Enemy &&
			b instanceof SmallEnemy &&
			b.immunity < 0
		) {
			sprites.delete(a)
			sprites.delete(b)
			sprites.add(new SmallEnemy(a.pos, a.velocity))
			sprites.add(new Explosion(b.pos))
		} else if (a instanceof Enemy && b instanceof Enemy) {
			sprites.delete(a)
			sprites.delete(b)
			sprites.add(new SmallEnemy(a.pos, [-a.velocity[0], a.velocity[1]]))
			sprites.add(new SmallEnemy(b.pos, [-b.velocity[0], b.velocity[1]]))
		}
	}

	function render() {
		const startTime = Date.now()

		display.fill(black)
		spriteLocMap.clear()

		if (ship.framesSinceLastTex++ > 15) {
			ship.currentTex++
			ship.framesSinceLastTex = 0
		}

		// move ship
		if (moveDirection === "left") ship.pos[0] -= 1
		else if (moveDirection === "right") ship.pos[0] += 1

		const o = ship.offset[0]
		if (ship.pos[0] < o) ship.pos[0] = o
		else if (ship.pos[0] > 160 - o + 1) ship.pos[0] = 160 - o + 1

		// move bullets
		const ss = [...sprites]
		for (const b of ss.filter(s => s instanceof Bullet)) {
			b.pos[1] -= 3
			if (b.pos[1] < 0) sprites.delete(b)
		}
		for (const e of ss.filter(
			s => s instanceof Enemy || s instanceof SmallEnemy
		)) {
			e.pos[0] += e.velocity[0]
			e.pos[1] += e.velocity[1]
			if (e.pos[1] > 142) sprites.delete(e)
		}

		for (const e of ss.filter(s => s instanceof SmallEnemy)) e.immunity--

		for (const e of ss.filter(s => s instanceof Explosion)) {
			if (e.framesSinceLastTex++ > 5) {
				e.currentTex++
				e.framesSinceLastTex = 0
			}
			if (e.currentTex >= e.textures.length) sprites.delete(e)
		}

		let framesPerEnemy = 100
		if (score > 20000) framesPerEnemy = 15
		else if (score > 15000) framesPerEnemy = 20
		else if (score > 10000) framesPerEnemy = 30
		else if (score > 6000) framesPerEnemy = 40
		else if (score > 4000) framesPerEnemy = 50
		else if (score > 2500) framesPerEnemy = 60
		else if (score > 1000) framesPerEnemy = 80

		let smallEnemyChance = 1
		if (score > 12000) smallEnemyChance = 0
		else if (score > 9000) smallEnemyChance = 0.2
		else if (score > 7500) smallEnemyChance = 0.4
		else if (score > 3000) smallEnemyChance = 0.6
		else if (score > 1500) smallEnemyChance = 0.8

		let velocityMultiplier = 1
		if (score > 30000) velocityMultiplier = 2
		else if (score > 25000) velocityMultiplier = 1.7
		else if (score > 20000) velocityMultiplier = 1.5
		else if (score > 17500) velocityMultiplier = 1.3
		else if (score > 16000) velocityMultiplier = 1.1

		if (stage < 2 && score > 15000) {
			stage++
			// yellow
			for (const tex of scoreTexs)
				for (const i in tex) tex[i] = tex[i].replaceAll("2", "6")
		}
		if (stage < 3 && score > 30000) {
			stage++
			// red
			for (const tex of scoreTexs)
				for (const i in tex) tex[i] = tex[i].replaceAll("6", "3")
		}

		if (
			framesSinceLastEnemy++ > framesPerEnemy &&
			gameState === "playing"
		) {
			const x = Math.random() * 160
			const velocity = [
				Math.random() * 0.5 - 0.25,
				(Math.random() * 0.75 + 0.25) * velocityMultiplier,
			]

			const classToAdd =
				Math.random() < smallEnemyChance ? SmallEnemy : Enemy
			sprites.add(new classToAdd([x, 0], velocity))
			framesSinceLastEnemy = 0
		}

		// starfields
		for (const field in starcolours) {
			if (framesSinceLastUpdate[field]++ >= updateSpeed[field]) {
				framesSinceLastUpdate[field] = 0

				starfields[field].pop()
				starfields[field].unshift(starfieldLine(chances[field]))
			}

			const colour = starcolours[field]

			for (let y = 0; y < 128; y++)
				for (let x = 0; x < 160; x++)
					if (starfields[field][y][x]) display[y * 160 + x] = colour
		}

		framesSinceLastBullet++

		for (const sprite of sprites) drawSprite(sprite)
		if (gameState === "over") drawSprite(gameOver)
		else score += 0.5

		// draw score
		const scoreStr = Math.floor(score).toString()
		for (let i = 0; i < scoreStr.length; i++)
			drawSprite(new ScoreNumber(i, +scoreStr[i]))

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

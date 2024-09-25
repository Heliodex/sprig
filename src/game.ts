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

const suddenDeathify = (s: string) =>
	s.replaceAll("H", "3").replaceAll("L", "9").replaceAll("1", "6")

const enemyTexs = transformTextures([
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

class Enemy extends Sprite {
	// supposed to be a funky planet or asteroid
	textures = enemyTexs
	offset = [6, 12]
	velocity = [0, 0.5]
	suddenDeath = false
	constructor(pos: number[], velocity: number[], suddenDeath = false) {
		super()

		if (pos[0] < 0) velocity[0] = Math.abs(velocity[0]) + 0.2
		else if (pos[0] > 160) velocity[0] = -Math.abs(velocity[0]) - 0.2

		this.pos = pos
		this.velocity = velocity
		this.suddenDeath = suddenDeath
	}
}

const smallEnemyTexs = transformTextures([
	`

	 	   LLL
	 	  LH1HL
	 	 1HHHHL
	 	 1H1HLHL
	 	 1LHHHL
	 	  1HLH1
	 	   111 `,
])

class SmallEnemy extends Sprite {
	textures = smallEnemyTexs
	offset = [6, 12]
	velocity = [0, 0.5]
	immunity = 15
	suddenDeath = false
	constructor(pos: number[], velocity: number[], suddenDeath = false) {
		super()

		if (pos[0] < 0) velocity[0] = Math.abs(velocity[0]) + 0.2
		else if (pos[0] > 160) velocity[0] = -Math.abs(velocity[0]) - 0.2

		this.pos = pos
		this.velocity = velocity
		this.suddenDeath = suddenDeath
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
		    3 9939696
		    6963 3939
		    693 6 396
		    6363 39 6
		    399939696
		     6 69993
		      69663`,

		`
		        9
		      66666
		     6999936
		   3 696 996
		    969 C 669
		     639 996
		    96969936
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
			200020000000002020000200200200222200002
			200020022200020002000200000200200000002
			200020000200022222000200000200200000002
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
			222222222222222222222222222222222222222`,
	])
	collidable = false
	pos = [80, 50]
	offset = [20, 0]
}

class Intro extends Sprite {
	textures = transformTextures([
		`	
			HHHLLHHHHL000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001HHHHHHHH1100000000000000000000000000000000000000
			LHHHHH1HHL000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001LHHLLHHHL000000000000000000000000000000000000000
			LHHHHHHH100000000000000000000L0000000000000000000000000000000000000000000000000000000000000000000000000000000001HHHHLHHHL000000000000000000000000000000000000000
			0LHHHHH1L000000000000000000000000000000000000000000000000000000LLLLL0000000000000000000000000000000000000000000LHH1HH1HH0000000000000000000000000000000000000000
			00LL11L000000000000000000000000000000000000000000000000000000LLHHHHH1L0000000000000000000000000000000000000000001HHHHHHL0000000000000000000000000000000000000000
			000000000000000000000000000000000000000000000000000000000000LLH1HHHHH1L0000000000000C000000000000000000000000000L1HHHHL00000000000000000000000000000000000000000
			00000000000000000000000000000000000000000000000000000000000LLH11HHHLHH10000000000000000000000000000000000000000000LLL0000000000000000000000000000000000000000000
			00000000000000000000000000000000000000000000000000000000000LHHHHHHHHHHHL0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
			0000000000000000000000000000000000000000000000000000000000LLHHHHHHHHHHH1000000000000000000000000000000000000000000000000000000000LL00000000000000000000000000000
			0000000000000000000000000000000000000000000000000000000000LHLLHHLLHHH1H1000000000000LCC00000C0000000000000000000000000000000000001L00000000000000000000000000000
			0000000000000000000000000000000000000000000000000000000000LHHHHHLLHHH1110000000000000C9C00000000000000000000000000000000000000000000000000000000000000000LLLLL00
			0000000000000000000000000000000000000000000000000000L00000LHHHHHHHHHHHH10000000000000C99C00000000000000000000000000000000000000000000000000000000000000L1HHHHHL0
			00000000000000000000000000000000000000000000000000000000000LHHHHHHHLHHH10000000000C0CC9990000000000000000000000000000000000000000000000000000000000000LHHHHLLHHL
			00000000000000000000000000000000000000000000000000000000000LHHHHHHHHHH1L000000000000CC66990000C90000000000000000000000000000000000000000000000000000001HHHHHHHHH
			00000000000000000000000000000000000001L000000000000000000000LHHH11HHH1L00000000L9C0C3C6669900C3300000000000000000000000000000000000000000000000000000LHHHHHHHHHH
			00000000000000L0000000000000000000000000000000000000000000000L1HHHHH1L000000000L99CC336666999933CCC00C00000000000000000000000000000000000000000000000LHHHHHLH11H
			00000000000LL11111L0000000000000000000000000000000000000000000LL111L000000000000999933662669933339F000000000000000000000000000LL000000000000000000000LHHHHLLLHHH
			0000000000LHHHHHHH110000000000000000L1HH1L00000000000000000000000000000000000000C99993669969933366000000000000000000000000LL0000000000000000000000000LHHHHHLHHHH
			00000000LLHH11HHHHHH10000000000000LLHHHHHH100000000000000000000000000000C00000L668666896966888966F00000000000000000000000000000000000000000000000000001H1HHHHLLH
			00000000LHHH11HHHHHH1L000000000000LHHHHLHHHL0000L100000000000000000000000000000F6666669666688866600000000000000000000000000000000000000000000000000000LHHHHHHHHH
			0000000LHHHHHHHLLHHHH100000000000LHHHHHHHHHH0000L100000000000000000000000000000066966666666966669C00000000000000000000000000000000000000000000000000000L1HHHHHHL
			0000000LHHHHHHHHHHHHH1L0000000000LHHHHHHH1HHL000000000000000000000000000000000009622666626666666999C0000000000000000000000000000000000000000000000000000LL111L10
			0000000LHHLHHHHHHHH1HHL0000000000HHHHHLLH1HHL000000000000000000L00000000000000C9982266662666669969C0000000000000000000000000000000000000000000000000000000000000
			0000000LHHHHHLLLHHHHHHL0000000000HHHHHLLHHHHL00000000000000000000000000000000C9999289F16666666666F00000000000000000000000000000000000000000000000000000000000000
			0000000LHHHHHLLLHHHHHHL0000000000LHH1HHHHHHHL0000000000000000000000000000000C33333391555966626666F0000LL00000000000000000000000000000000000000000000L00000000000
			0000000LHHHHHHHHH1HHH1L0000000000LHH1HHHHLH10000000000000000000000000000000C9993396615559666639F0000001L00000000000000000000000000000000000000000000L00000000000
			0000000LHHH1HHHHHHHHH1000000000000LHHHHHHH1L0000000000000000000000000000000CCCCC9666777F96699333C000000000000000000000000000000000000000000000000000000000000000
			00000000LH11HHHHHHHH1L0000000000000LHHHHH1L0000000000000000000000000000000000000666777663699993333C0000000000000000000000000000000000000000000000000000000000000
			00000000LLHHHHHLHHHH1000L000000000000LLLL00000000000000000000000000000000000000C6F177669369999983330000000000000000000000000000000000000000000000000000000000000
			000000000LLHHHHHHHH1000L2000000000000000000000000000000000000L00000000000000C00L0L77F66339828CCC3C0000L00000000000000000000000000000000000000000LLLLL00000000000
			00000000000LLHHH11L00000L0000000000000000000000000000000000001L0000000000L100000L7719993398220003C00000000000000000000000000000000000000000000L11HHH1LL000000000
			0000000000000LLLLL00000000000000000000000000000000010000000000000000000000000000771999333C98F0000C0000000000000000000000000000000000000000000L1HHHHH1HLL00000000
			00000000000000000000000000000000000000000000000000L1000000000000000000000000000777999C333099C000L00C0000000000000000000000000000000000000000L1HHLHHH11HL00000000
			00000000000000000000000000000000000000000000000000000000000000000000000000000077709C03CC30C9C00L100000000000000000000000000000000000000000001HHHHHHHHHHHL0000000
			0000000000000000000000000000000000000000000000000000000000000000000000000000077700000000000900000000000000000000000000000000000000000000000L1HHHHHHHHHHHL0000000
			0000000000000000000000000000000000000000000000000000000000000000000000000000L77000000000000C0000000000000000000000000000000000000000000L000LH11HHLLHHLHHL0000000
			0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000L00000000000000000000000000L0000000000L000LHH1HHLLHHHHHL0000000
			000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000L100000000000000000000000000000000000000000L1HHHHHHHHHHHL0000000
			000000000000000000000000000000000000000000000000000000000000000000000000000000000000000C00000000000000000000000000000000000000000000000000001HHHHHHHH1HHL0000000
			0000000000000000000000000000000000000000000L000000000000000000000LL0000000000000000C0000000000000000000000000000000000000000000000000L100000LHHLHHHHHHHLL0000000
			0000000000000000000000000000000000000000000L00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000L10000001HHHHHHHHHL00000000
			000000000000000000000000000000000000000000000000000000000L00000000000000000000000000000000000000000000000000000000000000000000000000000000000011HHHHHLL000000000
			00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000L00000000000000000000000000000LL1LLLL0000000000
			00000000001L00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000LLLLLL00000000000000000000000000000000000000000
			0000000000000000000000000000LL0000000000000000000000000000000000000000000000000000000000000000000000000000000000L1HHHHHL0000000000000000000000000000000000000000
			0000000000000000000000000000LL000000000000000000000000000000000000000000000000000000000000000000000000000000000LHHHHHHHHL000000000000000LL0000000000000000000000
			000000000000000000000000000000000000000000000000L000000000000000000000000000000000000000000000000000000000000001HH11HHLHHL000000000000001L0000000000000000000000
			0000000000000000000000000000000000000L100000000010000000000000000000000000000000000000000000000000000000000000LHHHHHHHHHHL00000001L00000000000000000000000000000
			0000000000000000000000000000000000000LL00000000000000000000000000000000000000000000000000000000000000000000000LHLHHHLLHHHH00000000000000000000000000000000000000
			00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000LHHHHHLLHHHH00000000000000000000000000000000000000
			00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000LHHHHHHHHHHL00000000000000000000000000000000000000
			00000000000000000000000000000LLLLL0000000000000000000L000000000000000000000000000000000000000000000000000000000HHH1HHH1HHL000000000000000000000000000000000L1000
			LLLL000000000000000000000000LHHHHHL000000000000000000000000000000000000000000000000000000000L100000000000000000LHHHHHHHHL00000000000LL00000000000000000000000000
			HHHH1L000000000000000000000LHHHHHHHL000000001L0000000000000000000000000000000000LL0000000000LL000000000000000000LHHHHHHL000000000000LL00000000000000000000000000
			HHHHHH100000000000000000000HHHHHHHH100000000100000000000000000000000000000000000000000000000000000000000000000000LLHLLL000000000000000000000000L0000000000000000
			HLHHHHHL000000000000000000LHHHHLHHHHL000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
			HHHHHHHH000000000000000000LH11HLHHHHL00000000000000000000000777000L000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
			HHHHHHHHL000000000000000000HHHHHHHHH0000000000000000000000077700011000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
			HLLLHHHH1000000000000000000LHHHHHHH1000000000000000000000077700000000000000000000000000000000LL00000000000000000000000000000000000000000000100000000000000000000
			HLLLH11H10000000000000000000LHHHHH10000000000000000000000L77000000000000000000000000000000000LL0000000000000000000000000000000000000000000L100000000000000000000
			HHHHHHHH100000000000000000000LHH1L0000000000000000000000L77L000000LL00000000L0000000000000000000000000000000L000000000000000000000000000000000000000000000000000
			HHHHHHHHL00000000000000000000000000000000000000000LLL00L77L0000000LL000000000000000000000000000000000000000000000000000000L0000000000000000000000000000000000000
			HH1HHHH10000000000000000000000000000000000000000777557777L000000000000000000000000000000000000000000000001000000000000000010000000000000000000000000000000000000
			HH1HHHHL00000000000000000LL000000000000000000007555555777000000000000000000000000000000000000000000000000L000000000000000000000000000000000000000000000000000000
			HHHHHHL000000000000000000000000000000000000000L755555777L0000000000000222222222222222222222222222222222222222222222222222222222222222222222222222222220000000000
			HHHH1L000000000000000000000000000000000000000075555557757000000000000020000000000000000000000000000000000000000000000000000000000000000000000000000002LL11111L00
			0LL0000000000000000000000000000000000000000000755555555570000000000000200000000000000000000000000000000000000000000000000000000000000000000000000000021HHHHHH1L0
			000000000000000000000000000000000000000000000075555555557000000000000120000000000000000000000000000000000000000000000000000000000000000000000000000002HHHHHHHHHL
			00000000L0000000000000000000000000L00000000000755H5555557000000000000020000222202222002222002220022220000002222020002002220002220022222022220222200002HHLHHH11HH
			00000000110000000000L00000000000001L0000000000L533555555L000000000000020002000002000202000200200200000000020000020002020002020002000200020000200020002HHHHHHH1HH
			0000000000000000000000000000000000000000000000C33355555700000000L0000020002000002000202000200200200000000020000020002020002020002000200020000200020002HHHHHHHHHH
			000000000000000000000000000000000000000000000333335557700000000000000020000222002222002222000200202220000002220022222020002020002000200022200222200002HHHLHHHHHH
			0000000000000000000000000000000000000000000C333833L77L000000000000000020000000202000002020000200200020000000002020002020002020002000200020000202000002HHLLLHLLHH
			000000000000000000000000000000000000000000C33828330000000000000000000020000000202000002002000200200020000000002020002020002020002000200020000200200002HHHLLHHHHH
			0000000000000000000000000000000000000000C3332228330000000000000000000020002222002000002000202220022200000022220020002002220002220000200022220200020002HHHHHHHHHH
			000000000000000000000000000000000000000C33322228330000000000000000000020000000000000000000000000000000000000000000000000000000000000000000000000000002LHHHHHH1HH
			000000000000LL00000000000000000000000C33338222233300000000000000L0000020000000000000000000000000000000000000000000000000000000000000000000000000000002LHHH1HH1HH
			000000000000000000000000000000000000C33338222233330000000000000000000020000000000000000000000000000000000000000000000000000000000000000000000000000002HHHHHHHHHL
			000000000000000000000000000000000003333332222833330000000000000000000L20000000000000000000000000000000000000000000000000000000000000000000000000000002HHHHHHHHLL
			0000000000000000000L0000000000000C333333822223333C0000000000000000000020000000000000000000000000000000000000000000000000000000000000000000000000000002LLHHHHLL00
			0000000000000000000000000000000033333338222233333C000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000200LLLL0000
			000000000000000000000000000000C338833332222833333C0002L000000000000000200022200000000000000000002220020000220020002000000000000000000000000000000000020000000000
			0000000000000000000000000000033382333382222333333C000L0000000000000000200002000000000000000000002002020002002002020000000000000000000000000000000000020000000000
			000000000000000000000000000C333223333822223333333C00000000000000000000200002000000000000001110002220020002222000200000000000000000000000000000000000020000000000
			0000000000000000000000000033382283333222283333833C00000000000000000000200002000000000000000000002000020002002000200000000000000000000000000000000000020000000000
			000000000000000000000000C333222233338222233338833C00000000000000000000200022200000000000000000002000022202002000200000000000000000000000000000000000020000000000
			00000000000000000000000C3382222333382222333382833000000000000000000000200000000000000000000000000000000000000000000000000000000000000000000000000000020000000000
			000000000000000000000C333322228333322228333322833000000000000000000000200000000000000000000000000000000000000000000000000000000000000000000000000000020000000000
			00000000000000000000C3333822223333822223333822833000000000000000000000200020002000000000000000002220222022200222000000000000000000000000000000000000020000000000
			0000000000000000000333338222233338222233338222833000000000000000000000200020002000000000000000002000020020020200000000000000000000000000000000000000020000000000
			00000000000000000C333333222283333288283333222233300000000L000000000000200020202000000000001110002220020022200220000000000000000000000000000000000000020000000000
			0000000000000000C33333382222333333332333382222333000000000000000000000200022022000000000000000002000020020020200000000000000000000000000000000000000020000000000
			0000000000000000333333822223333333333333822223333000000000000000000000200020002000000000000000002000222020020222000000000000000000000000000000000000020000000000
			000000000000000333333322228333333C33333322228333C000000000000000000000200000000000000000000000000000000000000000000000000000000000000000000000000000020000000000
			00000000000000C333333822283333CCCC33333822223333C00000000L0000000000L1200000000000000000000000000000000000000000000000000000000000000000000000000000020000000000
			0000000000000033333333333333CCCCCC33338222233333C00000000000000000000L200002200001022200000000002000200220020002022200000000000000000000000000000000L20000000000
			00000000000000CCCC33333333CCCCCCCC33332222833333C00000000000000000000020002002000102002000000000220220200202000202000000000000000000000000000000000LL20000000000
			00000000000000CCLLLLLC33CCCCCCCCCC33332222333333C00000000000000000000020002222001002002000111000202020200200202002200000000000000000000000000000000LL20000000000
			00000000000000CLLLLLLLLCCCCCCCCCC0C3332223333333C0000000000000000000002000200201000200200000000020002020020020200200000000000000000000000000000000LLL20000000000
			00000000000000F111LLLLLCCCCCCC0000C333333333333CL0000000000000000000002000200201000222000000000020002002200002000222000000000000000000000000000000LLL20000000000
			000000000000006822111LLCCCCCL00000C333333333333CLL00L10000000000000000200000000000000000000000000000000000000000000000000000000000000000000000000LLL020000000000
			0000000000000F9333886LLCCC0LL00000LLLC33333333CC000000000000000000000020000000000000000000000000000000000000000000000000000000000000000000000000LLL0020000000000
			000000000000C69999996FLL0000000000LLLLLLLC333CCC000000000000000000000020002002000000000000000000222000220020020022202220000000000000000000000000LLL0020000000000
			000000000000699999666FL0000000000LLLLLLLLLCCCCCC00000000000000000000002000202000000000000000000020020200202002020000200000000000000000000000000LLL00020000000000
			00001000000F999966666F0000000000C11111LLLLCCCCC000000000000000000000002000220000000000000011100022200222202002002200220000000000000000000000000LLL00L20000000000
			0000L00000C6966666666L000000000063382211LLCCCCC00000000000000000000000200020200000000000000000002000020020200200002020000000000000000000000000LLL000L20000000000
			0000000000666666666660000000000L69933366LLCCCC00000000000000000000000020002002000000000000000000200002002002200222002220000000000000000000000LLL000LL20000000000
			000000000F6666666666F0000000000L69999966LLCCC000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000000000LLL000L020000000000
			000000000666666666FL00000000000L66999666L0000000000000000000000000L0002000000000000000000000000000000000000000000000000000000000000000000000LLL000LL020000000000
			0000000006666666FL0000000000000L66996666L00LL0000000000000000000001L002000200000000000000000000020002022202002020020000000000000000000000000LLL00LL002L000000000
			00000000066666F000000000000000006699666F0002L0000000000000000000000000200020000000000000000000002202202000220202002000000000000000000000000LLL000L0002L000000000
			00000000066FC00000000000000000006696666C00000000000000000000000000000020002000000000000000111000202020220020220200200000000000000000000000LLL000LL0002L000000000
			000000000FL000000000000000000000F66666F000000000000000000000000000000020002000000000000000000000200020200020020200200000000000000000000000LLL000L00002L000000000
			00000000000000000000000000000000C66666000000000000000000000000000L00002000222000000000000000000020002022202002002200000000000000000000000LLL000LL00002L000000000
			0000000000000000000000000000000006666C00000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000LLL000LL0000020000000000
			000000000000000000000000000000000666F000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000LLL000L00000020000000000
			000000000000000000000000000000000F6600000000000000000000000000000000LL20000000000000000000000000000000000000000000000000000000000000000LLL000LL00000020000000000
			00000000000000000000000000L000000L6C00000000000000000000000000000000L1222222222222222222222222222222222222222222222222222222222222222222222222222222220000000000
			00000000000000000000000000L0000000F0000000000000000000000000000000001HHHHHHHHHHHHHHH1000000000000000000000000000000000000000000000000000000000LLLL00000000000000
			0000000000000000000000000000000000000000000000000000000000000000000LHHHHHHHHHHHHH11HHL0000000000000000000000000000000000000000000000000000000000LL00000000000000
			0000000000000000000000000000000000000LL00000000000000000000000000001HHHHH1HHHHHHH11HH1000000000000000000000000000000000000000L0000000000000000000000000000000000
			00000000000000000000000000000000000000000000000000000000000000000001HHHHHHHHHHHHH1HHH10000000000000L000000000000000000000000000000000000000000000000000000000000
			000000000000000000000000000000000000000000000000000000000000000000L1HHHHHHHLLLHHHHHHH1L0000000000000000000000000000000000000000000000000000000000000000000000000
			000000000000000000000000000000000000000000000000000000000000000000LHH11HHHHLLLLHHHHHHHL1L00000000000000000000000000000LL0000000000000000000000000000000000000000
			000000000000000000000000000000000000000000000000000000000000000000LHHHHHHHHHLLHHHHHHH1L0000000000000000000000000000000LL0000000000000000000000000000000000000000
			0000000000000000000000000000000000000000000000000000000000LL00000001HHHHHHHHHHHHHHHHH100000000000000000000000000000000000000000000000000000000000000000000000000
			00000000000000000000000000000000000000000000000000000000002L0000000LHHHHHLLHHHHHHLHHHL00000000000000000000000000000000000000000000000000000000000000000000000000
			0000000000000000000000000000000000000000000000000000000000000000000LHHHHHHHHHHHHHHHHHL00000000000000000000000000000000000000000000000000000000000000000000000000
			
		`,
	])
	collidable = false
	pos = [0, 0]
	offset = [0, 0]
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

function framesPerEnemy(s: number) {
	if (s > 65000) return 2
	if (s > 40000) return 3
	if (s > 30000) return 7
	if (s > 23000) return 10
	if (s > 15000) return 15
	if (s > 10000) return 20
	if (s > 6000) return 25
	if (s > 4000) return 33
	if (s > 2500) return 40
	if (s > 1000) return 50
	return 60
}

function smallEnemyChance(s: number) {
	if (s > 12000) return 0
	if (s > 9000) return 0.2
	if (s > 7500) return 0.4
	if (s > 3000) return 0.6
	if (s > 1500) return 0.8
	return 1
}

function velocityMultiplier(s: number) {
	if (s > 60000) return 3
	if (s > 50000) return 2.5
	if (s > 30000) return 2
	if (s > 27000) return 1.7
	if (s > 20000) return 1.5
	if (s > 18000) return 1.3
	if (s > 16000) return 1.1
	return 1
}

export default async (api: WebEngineAPI) => {
	const { onInput, setLegend, setMap } = api

	const spriteLocMap = new Map<number, Sprite>()
	const display = new Int8Array(160 * 128)
	const getPixel = (pos: number) => String.fromCharCode(display[pos])

	function drawSprite(sprite: Sprite) {
		const toDraw = sprite.texture

		const [xPos, yPos] = sprite.pos.map(Math.round)
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

	for (const field in starcolours) {
		starfields[field] = []
		for (let i = 0; i < 128; i++)
			starfields[field].push(starfieldLine(chances[field]))
	}

	// sprites
	let ship = new Ship()
	const gameOver = new GameOver()
	const sprites = new Set<Sprite>()

	let gameState = "menu"
	let score = 0
	let stage = 1
	let framesSinceLastUpdate = [0, 0, 0, 0]
	let framesSinceLastBullet = 0
	let framesSinceLastEnemy = 60

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

	onInput("i", () => {
		if (!["paused", "over", "menu"].includes(gameState)) return
		if (gameState !== "paused") {
			score = 0
			stage = 1
			framesSinceLastUpdate = [0, 0, 0, 0]
			framesSinceLastBullet = 0
			framesSinceLastEnemy = 60
			sprites.clear()
			ship = new Ship()
			sprites.add(ship)
		}
		gameState = "playing"
	})
	onInput("k", () => {
		if (gameState !== "playing") return 
		gameState = "paused"
		moveDirection = ""
	})
	onInput("l", () => {
		if (["paused", "over"].includes(gameState)) gameState = "menu"
	})

	onInput("w", () => {
		if (gameState !== "playing" || framesSinceLastBullet < 10) return
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
			sprites.add(new Explosion(a.pos))
			sprites.delete(b)
			sprites.add(new Explosion(b.pos))
			gameState = "over"
		} else if (
			a instanceof SmallEnemy &&
			b instanceof SmallEnemy &&
			a.immunity < 0 &&
			b.immunity < 0
		) {
			if (stage >= 4) return
			sprites.delete(a)
			sprites.add(new Explosion(a.pos))
			sprites.delete(b)
			sprites.add(new Explosion(b.pos))
		} else if (
			a instanceof Enemy &&
			b instanceof SmallEnemy &&
			b.immunity < 0
		) {
			sprites.delete(a)
			sprites.add(new SmallEnemy(a.pos, a.velocity, stage >= 4))
			if (stage >= 4) return
			sprites.delete(b)
			sprites.add(new Explosion(b.pos))
		} else if (a instanceof Enemy && b instanceof Enemy) {
			sprites.delete(a)
			sprites.delete(b)
			sprites.add(
				new SmallEnemy(
					a.pos,
					[-a.velocity[0], a.velocity[1]],
					stage >= 4
				)
			)
			sprites.add(
				new SmallEnemy(
					b.pos,
					[-b.velocity[0], b.velocity[1]],
					stage >= 4
				)
			)
		}
	}

	let lastFrame = 0

	function paintAndFinish() {
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
		requestAnimationFrame(render)
	}

	function render() {
		if (Date.now() - lastFrame < 1000 / 30) {
			// should we render? nah can't be botherd yet
			requestAnimationFrame(render)
			return
		}
		lastFrame = Date.now()

		display.fill(black)
		spriteLocMap.clear()

		if (gameState === "playing") {
			if (ship.framesSinceLastTex++ > 7) {
				ship.currentTex++
				ship.framesSinceLastTex = 0
			}

			if (moveDirection === "left")
				// move ship
				ship.pos[0] -= 1.5
			else if (moveDirection === "right") ship.pos[0] += 1.5

			const o = ship.offset[0]
			if (ship.pos[0] < o) ship.pos[0] = o
			else if (ship.pos[0] > 160 - o + 1) ship.pos[0] = 160 - o + 1

			// score/difficulty calculation

			if (stage < 2 && score > 15000) {
				stage++
				// yellow score
				for (const tex of scoreTexs)
					for (const i in tex) tex[i] = tex[i].replaceAll("2", "6")
			}
			if (stage < 3 && score > 40000) {
				stage++
				// red score
				for (const tex of scoreTexs)
					for (const i in tex) tex[i] = tex[i].replaceAll("6", "3")
			}
			if (stage < 4 && score > 70000) {
				stage++
				// small sudden death
				for (const tex of smallEnemyTexs)
					for (const i in tex) tex[i] = suddenDeathify(tex[i])
			}
			if (stage < 5 && score > 75000) {
				stage++
				// sudden death
				for (const tex of enemyTexs)
					for (const i in tex) tex[i] = suddenDeathify(tex[i])
			}

			const framesPerE = framesPerEnemy(score)
			const velMult = velocityMultiplier(score)
			const smallEChance = smallEnemyChance(score)

			if (framesSinceLastEnemy++ > framesPerE) {
				const x = Math.random() * 240 - 40
				const velocity = [
					(Math.random() * 0.75 - 0.3) * Math.min(1, velMult - 1),
					(Math.random() * 0.75 + 0.3) * velMult,
				]

				if (Math.random() < smallEChance)
					sprites.add(new SmallEnemy([x, 0], velocity, stage >= 4))
				else sprites.add(new Enemy([x, 0], velocity, stage >= 5))
				framesSinceLastEnemy = 0
			}
		}

		if (["playing", "over"].includes(gameState)) {
			// move bullets
			const ss = [...sprites]
			for (const b of ss.filter(s => s instanceof Bullet)) {
				b.pos[1] -= 3
				if (b.pos[1] < 0) sprites.delete(b)
			}
			for (const e of ss.filter(
				s => s instanceof Enemy || s instanceof SmallEnemy
			)) {
				if (e.suddenDeath) {
					// accelerate towards the ship, muahahahha
					e.velocity[0] += e.pos[0] < ship.pos[0] ? 0.02 : -0.02
					e.velocity[1] += 0.03
				}
				e.pos[0] += e.velocity[0]
				e.pos[1] += e.velocity[1]
				if (e.pos[1] > 142) sprites.delete(e)
			}
			for (const e of ss.filter(s => s instanceof SmallEnemy))
				e.immunity--
			for (const e of ss.filter(s => s instanceof Explosion)) {
				if (e.framesSinceLastTex++ > 2) {
					e.currentTex++
					e.framesSinceLastTex = 0
				}
				if (e.currentTex >= e.textures.length) sprites.delete(e)
			}

			// starfields
			for (const field in starcolours) {
				if (framesSinceLastUpdate[field]++ >= updateSpeed[field]) {
					framesSinceLastUpdate[field] = 0

					starfields[field].pop()
					starfields[field].unshift(starfieldLine(chances[field]))
				}

				for (let y = 0; y < 128; y++)
					for (let x = 0; x < 160; x++)
						if (starfields[field][y][x])
							display[y * 160 + x] = starcolours[field]
			}

			framesSinceLastBullet++
		} else if (gameState === "menu") {
			drawSprite(new Intro())
			paintAndFinish()
			return
		}

		for (const sprite of sprites) drawSprite(sprite)

		if (gameState === "over") drawSprite(gameOver)
		else if (gameState !== "paused") score += 0.5

		// draw score
		const scoreStr = Math.floor(score).toString()
		for (let i = 0; i < scoreStr.length; i++)
			drawSprite(new ScoreNumber(i, +scoreStr[i]))

		paintAndFinish()
	}

	render()
}

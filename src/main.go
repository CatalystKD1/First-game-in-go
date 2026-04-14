package main

import (
	//"crypto/x509"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Window Size
const (
	windowX = 640
	windowY = 480
	margin = 16
)

// wall info
var (
	gapSize = 120
	wallGap = 200
)

const DEBUG = false
const ZERO = 0

// Bullet info
const (
	fireRate = 10
	bulletOffsetX = 50
	bulletOffsetY = 5
	bulletSpeed = 5
)
// information about the ship
const (
	startX = 25
	startY = 100

	shipRadius = 50
)

// set image variables
var (
	default_ship *ebiten.Image
	firing_ship *ebiten.Image

	currentShip *ebiten.Image
	bullet *ebiten.Image

	wallImg *ebiten.Image

	debugCircle *ebiten.Image
	debugPixel *ebiten.Image
)


func init() {
	var err error
	default_ship, _, err = ebitenutil.NewImageFromFile("../assets/ship_guy.png")
	if err != nil {
		log.Fatal(err)
	}

	firing_ship, _, err = ebitenutil.NewImageFromFile("../assets/ship_shoot.png")
	if err != nil {
		log.Fatal(err)
	}

	bullet, _, err = ebitenutil.NewImageFromFile("../assets/bullet.png")
	if err != nil {
		log.Fatal(err)
	}

	wallImg = ebiten.NewImage(1, 1)
	wallImg.Fill(color.RGBA{120, 0, 0, 255})

	debugCircle = ebiten.NewImage(64, 64)

	// draw a circle using pixels
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			dx := float64(x - 32)
			dy := float64(y - 32)
			if dx*dx+dy*dy <= 32*32 {
				debugCircle.Set(x, y, color.RGBA{255, 255, 255, 80})
			}
		}
	}

	debugPixel = ebiten.NewImage(2, 2)
	debugPixel.Fill(color.RGBA{255, 0, 0, 200}) // red semi-visible dot

	currentShip = default_ship

}

type Game struct{
	x float64
	y float64

	walls []Wall

	velocity float64 // speed the walls will go
	score int

	bullets []Bullet
	fireTimer int

	spawnTimer int
	nextSpawnTick int
	speed float64

	debugHitX float64
	debugHitY float64
}

type Bullet struct {
	x float64
	y float64
	speed float64
}

type Wall struct {
	x float64

	gapY float64
	gapH float64

	hp int // 0 = gap, >0 = brekable window

	passed bool // for scoring
}

func getColor(hp int) color.RGBA {
	switch hp {
	case 4:
		return color.RGBA{0, 0, 255, 255} // blue
	case 3:
		return color.RGBA{0, 255, 0, 255} // green
	case 2:
		return color.RGBA{255, 255, 0, 255} // yellow
	case 1:
		return color.RGBA{255, 255, 255, 255} // white
	default:
		return color.RGBA{120, 0, 0, 255}
	}
}

func (g *Game) SpawnWall() {
	gapY := rand.Float64() * float64(windowY - gapSize)

	w := Wall{
		x: windowX,
		gapY: gapY,
		gapH: float64(gapSize),
		hp: 0, // making it a solid wall for the moment
	}

	if rand.Float64() < 0.5 {
		w.hp = 4
	}

	g.walls = append(g.walls, w)
}

func (g *Game) Update() error {
	var err error
	


	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.fireTimer == ZERO {
		currentShip = firing_ship
		b := Bullet{
			x: g.x + bulletOffsetX, // adjust based on ship width
			y: g.y + bulletOffsetY, // adjust based on ship height
			speed: bulletSpeed,
		}
		g.bullets = append(g.bullets, b)

		g.fireTimer = fireRate
	} else if g.fireTimer > ZERO {
		g.fireTimer--
	} else {
		currentShip = default_ship
	}

	// spawn timer for walls
	// 🔥 spawn timer
	g.spawnTimer++

	if len(g.walls) < 3 && g.spawnTimer >= g.nextSpawnTick {

		// spacing check
		if len(g.walls) == 0 || g.walls[len(g.walls)-1].x < float64(windowX) - float64(wallGap) {

			g.SpawnWall()

			g.spawnTimer = 0
			g.nextSpawnTick = rand.Intn(91) + 30 // 30–120 ticks
		}
	}

	// collision 
	shipCX, shipCY := g.shipCenter()

	for _, w := range g.walls {

		wallLeft := w.x
		wallRight := w.x + 10
	
		if shipCX+shipRadius > wallLeft && shipCX-shipRadius < wallRight {
	
			inGap := shipCY > w.gapY && shipCY < w.gapY+w.gapH
	
			if !inGap {
	
				// 💥 compute collision point (IMPORTANT)
				collideX := shipCX
	
				// clamp Y into wall region (top or bottom block)
				var collideY float64
				if shipCY < w.gapY {
					collideY = w.gapY // hit top block
				} else {
					collideY = w.gapY + w.gapH // hit bottom block
				}
	
				if DEBUG {
					g.debugHitX = collideX
					g.debugHitY = collideY
				}
				
				if DEBUG {
					log.Println("HIT WALL")
				} else {
					return ebiten.Termination
				}
			}
		}
	}


	// bullet collision
	for bi := len(g.bullets) - 1; bi >= 0; bi-- {
		for wi := 0; wi < len(g.walls); wi++ {
	
			w := &g.walls[wi]
	
			if w.hp <= 0 {
				continue
			}
	
			if g.bullets[bi].x > w.x && g.bullets[bi].x < w.x+20 {
				if g.bullets[bi].y > w.gapY && g.bullets[bi].y < w.gapY+w.gapH {
	
					w.hp--
	
					// safe removal
					g.bullets = append(g.bullets[:bi], g.bullets[bi+1:]...)
					break
				}
			}
		}
	}

	//player movement

	if (ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW)) && g.y != ZERO {
		g.y -= g.speed
	}

	if (ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS)) && g.y != windowY {
		g.y += g.speed
	}

	// bullet logic

	//move bullets
	for i := range g.bullets {
		g.bullets[i].x += g.bullets[i].speed
	}

	// remove bullet that are off screen
	var aliveBullets []Bullet
	for _, b := range g.bullets {
		if b.x <= windowX + margin {
			aliveBullets = append(aliveBullets, b)
		}
	}
	g.bullets = aliveBullets

	// wall scoring
	for i := range g.walls {
		if !g.walls[i].passed && g.walls[i].x < g.x {
			g.walls[i].passed = true
			g.score++
			// could add a velocity increse here
		}
	}
	
	// walls moving
	for i := range g.walls {
		g.walls[i].x -= g.velocity
	}

	var aliveWalls []Wall
	for _, w := range g.walls {
		if w.x > -margin {
			aliveWalls = append(aliveWalls, w)
		}
	}
	g.walls = aliveWalls

	if err != nil {
		return err
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// fill background with (R, G, B, A)
	screen.Fill(color.RGBA{50, 51, 59, 255}) // light blue

	// set position of the ship
	op := &ebiten.DrawImageOptions{}

	//op.GeoM.Scale(shipW, shipH)
	op.GeoM.Translate(g.x, g.y)

	// draw bullets
	for _, b := range g.bullets {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(b.x, b.y)

		screen.DrawImage(bullet, op)
	}

	for _, w := range g.walls {
		// TOP of the wall
		opTop := &ebiten.DrawImageOptions{}
		opTop.GeoM.Scale(20, w.gapY)
		opTop.GeoM.Translate(w.x, 0)
		screen.DrawImage(wallImg, opTop)

		if w.hp > 0 {
			col := getColor(w.hp)

			windowImg := ebiten.NewImage(1, 1)
			windowImg.Fill(col)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(20, w.gapH)
			op.GeoM.Translate(w.x, w.gapY)

			screen.DrawImage(windowImg, op)
		}
		
		// BOTTOM of the wall
		opBot := &ebiten.DrawImageOptions{}
		opBot.GeoM.Scale(20, float64(windowY)-(w.gapY+w.gapH))
		opBot.GeoM.Translate(w.x, w.gapY+w.gapH)
		screen.DrawImage(wallImg, opBot)
	}

	screen.DrawImage(currentShip, op)
	if DEBUG {
		shipCX, shipCY := g.shipCenter()
	
		op := &ebiten.DrawImageOptions{}
	
		// debug circle image is assumed 64x64 centered
		const debugSize = 64.0
	
		scale := (shipRadius * 2) / debugSize
	
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(shipCX-shipRadius, shipCY-shipRadius)
	
		screen.DrawImage(debugCircle, op)
	}

	if DEBUG && g.debugHitX != 0 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(g.debugHitX, g.debugHitY)
	
		screen.DrawImage(debugPixel, op)
	}
}

func (g *Game) shipCenter() (float64, float64) {
	return g.x + float64(currentShip.Bounds().Dx())/2,
		   g.y + float64(currentShip.Bounds().Dy())/2
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowX, windowY
}

func main() {
	game := &Game {
		x: startX,
		y: startY,
		fireTimer: ZERO,

		velocity : 5,
		nextSpawnTick: rand.Intn(91) + 30,

		speed: 4,
	}

	ebiten.SetWindowSize(windowX, windowY)
	ebiten.SetWindowTitle("Ship Guy")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
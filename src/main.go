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
	margin  = 16
)

// wall info
var (
	gapSize = 120
	wallGap = 200
)

const DEBUG = true
const ZERO = 0

// Bullet info
const (
	fireRate      = 10
	bulletOffsetX = 50
	bulletOffsetY = 5
	bulletSpeed   = 5
)

// information about the ship
const (
	startX = 25
	startY = 100

	shipRadius     = 20  // <- tightened: this is now the TRUE circle radius used for collision
	shipHitboxBias = 0.9 // <- optional: scale radius down slightly for forgiveness
)

// set image variables
var (
	default_ship *ebiten.Image
	firing_ship  *ebiten.Image

	currentShip *ebiten.Image
	bullet      *ebiten.Image

	wallImg *ebiten.Image

	debugCircle *ebiten.Image
	debugPixel  *ebiten.Image
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

	debugPixel = ebiten.NewImage(2, 2)
	debugPixel.Fill(color.RGBA{255, 0, 0, 200})

	currentShip = default_ship
}

type Game struct {
	x float64
	y float64

	walls []Wall

	velocity float64
	score    int

	bullets   []Bullet
	fireTimer int

	spawnTimer    int
	nextSpawnTick int
	speed         float64

	debugHitX float64
	debugHitY float64
}

type Bullet struct {
	x     float64
	y     float64
	speed float64
}

type Wall struct {
	x    float64
	gapY float64
	gapH float64
	hp   int
	passed bool
}

func getColor(hp int) color.RGBA {
	switch hp {
	case 4:
		return color.RGBA{0, 0, 255, 255}
	case 3:
		return color.RGBA{0, 255, 0, 255}
	case 2:
		return color.RGBA{255, 255, 0, 255}
	case 1:
		return color.RGBA{255, 255, 255, 255}
	default:
		return color.RGBA{120, 0, 0, 255}
	}
}

func (g *Game) SpawnWall() {
	gapY := rand.Float64() * float64(windowY-gapSize)

	w := Wall{
		x:    windowX,
		gapY: gapY,
		gapH: float64(gapSize),
		hp:   0,
	}

	if rand.Float64() < 0.5 {
		w.hp = 4
	}

	g.walls = append(g.walls, w)
}

// circleCollidesWall checks if a circle (cx, cy, r) overlaps a wall's solid regions.
// The wall has a solid top block [0, gapY] and solid bottom block [gapY+gapH, windowY].
// Returns true if the circle is touching either solid block.
func circleCollidesWall(cx, cy, r float64, w Wall) bool {
	wallLeft := w.x
	wallRight := w.x + 20

	// Broad phase: is the circle anywhere near the wall column horizontally?
	if cx+r < wallLeft || cx-r > wallRight {
		return false
	}

	// The gap is the SAFE zone. The two solid rects are:
	//   Top block:    x=[wallLeft,wallRight], y=[0, gapY]
	//   Bottom block: x=[wallLeft,wallRight], y=[gapY+gapH, windowY]
	//
	// For circle-vs-AABB we clamp the circle centre to the rect and check distance.
	topRect    := [4]float64{wallLeft, 0,           wallRight, w.gapY}
	bottomRect := [4]float64{wallLeft, w.gapY + w.gapH, wallRight, float64(windowY)}

	return circleHitsRect(cx, cy, r, topRect) || circleHitsRect(cx, cy, r, bottomRect)
}

// circleHitsRect does a proper circle-vs-AABB test.
// rect is [left, top, right, bottom].
func circleHitsRect(cx, cy, r float64, rect [4]float64) bool {
	// Clamp circle centre to the rectangle
	clampX := clamp(cx, rect[0], rect[2])
	clampY := clamp(cy, rect[1], rect[3])

	// Distance from clamped point to circle centre
	dx := cx - clampX
	dy := cy - clampY

	return dx*dx+dy*dy <= r*r
}

func clamp(val, lo, hi float64) float64 {
	if val < lo {
		return lo
	}
	if val > hi {
		return hi
	}
	return val
}

// checkCollisions runs ship-vs-wall circle collision for every wall.
// Returns true if the ship has hit something and the game should end.
func (g *Game) checkCollisions() bool {
	cx, cy := g.shipCenter()
	r := float64(shipRadius) * shipHitboxBias

	for _, w := range g.walls {
		if circleCollidesWall(cx, cy, r, w) {
			if DEBUG {
				log.Println("HIT WALL at", cx, cy)
				g.debugHitX = cx
				g.debugHitY = cy
				return false // in debug mode keep running so we can see the hit
			}
			return true
		}
	}
	return false
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.fireTimer == ZERO {
		currentShip = firing_ship
		b := Bullet{
			x:     g.x + bulletOffsetX,
			y:     g.y + bulletOffsetY,
			speed: bulletSpeed,
		}
		g.bullets = append(g.bullets, b)
		g.fireTimer = fireRate
	} else if g.fireTimer > ZERO {
		g.fireTimer--
	} else {
		currentShip = default_ship
	}

	// Spawn walls
	g.spawnTimer++
	if len(g.walls) < 3 && g.spawnTimer >= g.nextSpawnTick {
		if len(g.walls) == 0 || g.walls[len(g.walls)-1].x < float64(windowX)-float64(wallGap) {
			g.SpawnWall()
			g.spawnTimer = 0
			g.nextSpawnTick = rand.Intn(91) + 30
		}
	}

	// ── Collision ──────────────────────────────────────────────────────────────
	if g.checkCollisions() {
		return ebiten.Termination
	}

	// Bullet vs wall
	for bi := len(g.bullets) - 1; bi >= 0; bi-- {
		for wi := 0; wi < len(g.walls); wi++ {
			w := &g.walls[wi]
			if w.hp <= 0 {
				continue
			}
			if g.bullets[bi].x > w.x && g.bullets[bi].x < w.x+20 {
				if g.bullets[bi].y > w.gapY && g.bullets[bi].y < w.gapY+w.gapH {
					w.hp--
					g.bullets = append(g.bullets[:bi], g.bullets[bi+1:]...)
					break
				}
			}
		}
	}

	// Player movement
	if (ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW)) && g.y > ZERO {
		g.y -= g.speed
	}
	if (ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS)) && g.y < windowY {
		g.y += g.speed
	}

	// Move bullets, cull off-screen
	for i := range g.bullets {
		g.bullets[i].x += g.bullets[i].speed
	}
	var aliveBullets []Bullet
	for _, b := range g.bullets {
		if b.x <= windowX+margin {
			aliveBullets = append(aliveBullets, b)
		}
	}
	g.bullets = aliveBullets

	// Scoring
	for i := range g.walls {
		if !g.walls[i].passed && g.walls[i].x < g.x {
			g.walls[i].passed = true
			g.score++
		}
	}

	// Move walls, cull off-screen
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

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{50, 51, 59, 255})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.x, g.y)

	for _, b := range g.bullets {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(b.x, b.y)
		screen.DrawImage(bullet, op)
	}

	for _, w := range g.walls {
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

		opBot := &ebiten.DrawImageOptions{}
		opBot.GeoM.Scale(20, float64(windowY)-(w.gapY+w.gapH))
		opBot.GeoM.Translate(w.x, w.gapY+w.gapH)
		screen.DrawImage(wallImg, opBot)
	}

	screen.DrawImage(currentShip, op)

	if DEBUG {
		cx, cy := g.shipCenter()
		r := float64(shipRadius) * shipHitboxBias
		diameter := int(r * 2)

		debugCircle = ebiten.NewImage(diameter, diameter)
		for y := 0; y < diameter; y++ {
			for x := 0; x < diameter; x++ {
				dx := float64(x) - r
				dy := float64(y) - r
				if dx*dx+dy*dy <= r*r {
					debugCircle.Set(x, y, color.RGBA{255, 255, 255, 80})
				}
			}
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(cx-r, cy-r)
		screen.DrawImage(debugCircle, op)
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
	game := &Game{
		x:             startX,
		y:             startY,
		fireTimer:     ZERO,
		velocity:      2,
		nextSpawnTick: rand.Intn(91) + 30,
		speed:         2,
	}

	ebiten.SetWindowSize(windowX, windowY)
	ebiten.SetWindowTitle("Ship Guy")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
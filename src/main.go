package main

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// Window Size
const (
	windowX = 640
	windowY = 480
	margin  = 16
)

// default settings
var (
	shipSpeed = 4.0
	wallVelocity = 4.0
	fireTime = 0
	score = 0
	spawnTimer = 0
	nextSpawnTick = int(rand.Intn(60) + 60)
	debugHitX = 0.0
	debugHitY = 0.0
	state = StateMenu
)

// Wall info
var (
	gapSize = 120
	wallGap = 200
)

const DEBUG = false
const ZERO = 0

// Bullet info
const (
	fireRate      = 10
	bulletOffsetX = 50
	bulletOffsetY = 5
	bulletSpeed   = 5
)

// Ship info
const (
	startX         = 25
	startY         = 100
	shipRadius     = 20
	shipHitboxBias = 0.9
)

// Images
var (
	default_ship *ebiten.Image
	firing_ship  *ebiten.Image
	currentShip  *ebiten.Image
	bullet       *ebiten.Image
	wallImg      *ebiten.Image
	debugCircle  *ebiten.Image
	debugPixel   *ebiten.Image
)

type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
)

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

	state GameState
}

type Bullet struct {
	x     float64
	y     float64
	speed float64
}

type Wall struct {
	x      float64
	gapY   float64
	gapH   float64
	hp     int
	passed bool
}


func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowX, windowY
}

func main() {
	game := &Game{
		x: startX,
		y: startY,
		fireTimer: fireTime,

		// changes the speed of the walls
		velocity: wallVelocity,

		nextSpawnTick: nextSpawnTick,

		// changes the speed of the player
		speed: shipSpeed,

		state: StateMenu,
	}

	ebiten.SetWindowSize(windowX, windowY)
	ebiten.SetWindowTitle("Ship Guy")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
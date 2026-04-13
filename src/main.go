package main

import (
	//"crypto/x509"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Window Size
const windowX = 640
const windowY = 480
const margin = 16

const ZERO = 0

// Bullet info
const fireRate = 15
const bulletOffsetX = 50
const bulletOffsetY = 5
const bulletSpeed = 5

// start positions
const startX = 25
const startY = 100


// set image variables
var default_ship *ebiten.Image
var firing_ship *ebiten.Image

var currentShip *ebiten.Image
var bullet *ebiten.Image


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

	currentShip = default_ship

}

type Game struct{
	x float64
	y float64

	bullets []Bullet
	fireTimer int
}

type Bullet struct {
	x float64
	y float64
	speed float64
}

func (g *Game) Update() error {
	var err error

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.fireTimer == 0 {
		currentShip = firing_ship
		b := Bullet{
			x: g.x + bulletOffsetX, // adjust based on ship width
			y: g.y + bulletOffsetY, // adjust based on ship height
			speed: bulletSpeed,
		}
		g.bullets = append(g.bullets, b)

		g.fireTimer = fireRate
	} else if g.fireTimer > 0 {
		g.fireTimer--
	} else {
		currentShip = default_ship
	}

	if (ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW)) && g.y != 0 {
		g.y -= 2
	}

	if (ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS)) && g.y != windowY {
		g.y += 2
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
	op.GeoM.Translate(g.x, g.y) // (x, y)

	// draw bullets
	for _, b := range g.bullets {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(b.x, b.y)

		screen.DrawImage(bullet, op)
	}

	screen.DrawImage(currentShip, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowX, windowY
}

func main() {
	game := &Game {
		x: startX,
		y: startY,
		fireTimer: ZERO,
	}

	ebiten.SetWindowSize(windowX, windowY)
	ebiten.SetWindowTitle("Render an image")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
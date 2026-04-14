package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	ebitentext "github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"
)

const (
	menuButtonX = windowX/2 - 60
	menuButtonY = windowY/2 - 20
	menuButtonW = 120
	menuButtonH = 40
)

func (g *Game) resetGame() {
	g.x = startX
	g.y = startY
	g.fireTimer = fireTime
	g.velocity = wallVelocity
	g.speed = shipSpeed
	g.score = score
	g.spawnTimer = spawnTimer
	g.nextSpawnTick = nextSpawnTick
	g.walls = nil
	g.bullets = nil
	g.debugHitX = debugHitX
	g.debugHitY = debugHitY
	g.state = state
}

func (g *Game) updateMenu() {
	mx, my := ebiten.CursorPosition()

	inButton := mx >= menuButtonX && mx <= menuButtonX+menuButtonW &&
		my >= menuButtonY && my <= menuButtonY+menuButtonH

	if inButton && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.state = StatePlaying
	}
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	screen.Fill(color.RGBA{50, 51, 59, 255})

	face := ebitentext.NewGoXFace(basicfont.Face7x13)

	// Title
	titleText := "Ship Guy"
	op := &ebitentext.DrawOptions{}
	op.GeoM.Translate(float64(windowX/2-len(titleText)*7/2), float64(windowY/2-60))
	op.ColorScale.ScaleWithColor(color.RGBA{220, 220, 220, 255})
	ebitentext.Draw(screen, titleText, face, op)

	// Button background
	buttonImg := ebiten.NewImage(menuButtonW, menuButtonH)

	mx, my := ebiten.CursorPosition()
	inButton := mx >= menuButtonX && mx <= menuButtonX+menuButtonW &&
		my >= menuButtonY && my <= menuButtonY+menuButtonH

	if inButton {
		buttonImg.Fill(color.RGBA{100, 100, 120, 255})
	} else {
		buttonImg.Fill(color.RGBA{70, 70, 90, 255})
	}

	imgOp := &ebiten.DrawImageOptions{}
	imgOp.GeoM.Translate(float64(menuButtonX), float64(menuButtonY))
	screen.DrawImage(buttonImg, imgOp)

	// Button label
	btnText := "Start Game"
	btnOp := &ebitentext.DrawOptions{}
	btnOp.GeoM.Translate(float64(menuButtonX+menuButtonW/2-len(btnText)*7/2), float64(menuButtonY+menuButtonH/2+5))
	btnOp.ColorScale.ScaleWithColor(color.RGBA{220, 220, 220, 255})
	ebitentext.Draw(screen, btnText, face, btnOp)
}
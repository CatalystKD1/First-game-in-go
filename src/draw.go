package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

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

func (g *Game) Draw(screen *ebiten.Image) {
	if g.state == StateMenu {
		g.drawMenu(screen)
		return
	}
	
	screen.Fill(color.RGBA{50, 51, 59, 255})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.x, g.y)

	// Bullets
	for _, b := range g.bullets {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(b.x, b.y)
		screen.DrawImage(bullet, op)
	}

	// Walls
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

	// Ship
	screen.DrawImage(currentShip, op)

	// Debug hitbox circle
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
package main

import (
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
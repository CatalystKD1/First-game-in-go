package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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

func (g *Game) Update() error {
	if g.state == StateMenu {
		g.updateMenu()
		return nil
	}

	// Firing
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

	// Collision
	if g.checkCollisions() {
		g.resetGame()  // ← was: return ebiten.Termination
		return nil
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

	// Move & cull bullets
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

	// Move & cull walls
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
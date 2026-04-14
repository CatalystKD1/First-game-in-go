package main

import "log"

// circleCollidesWall checks if a circle (cx, cy, r) overlaps a wall's solid regions.
func circleCollidesWall(cx, cy, r float64, w Wall) bool {
	wallLeft := w.x
	wallRight := w.x + 20

	if cx+r < wallLeft || cx-r > wallRight {
		return false
	}

	topRect    := [4]float64{wallLeft, 0, wallRight, w.gapY}
	bottomRect := [4]float64{wallLeft, w.gapY + w.gapH, wallRight, float64(windowY)}

	if circleHitsRect(cx, cy, r, topRect) || circleHitsRect(cx, cy, r, bottomRect) {
		return true
	}

	// If the gap has HP remaining it is a breakable window — solid until destroyed
	if w.hp > 0 {
		windowRect := [4]float64{wallLeft, w.gapY, wallRight, w.gapY + w.gapH}
		return circleHitsRect(cx, cy, r, windowRect)
	}

	return false
}

// circleHitsRect does a proper circle-vs-AABB test.
// rect is [left, top, right, bottom].
func circleHitsRect(cx, cy, r float64, rect [4]float64) bool {
	clampX := clamp(cx, rect[0], rect[2])
	clampY := clamp(cy, rect[1], rect[3])

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
				return false
			}
			return true
		}
	}
	return false
}
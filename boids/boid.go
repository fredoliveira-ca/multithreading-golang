package main

import (
	"math"
	"math/rand"
	"time"
)

// Boid is the object which represents a point in the screen
type Boid struct {
	position Vector2D
	velocity Vector2D
	id       int
}

func (b *Boid) moveOne() {
	b.velocity = b.velocity.Add(b.calculateAcceleration()).limit(-1, 1)
	boidMap[int(b.position.x)][int(b.position.y)] = -1
	b.position = b.position.Add(b.velocity)
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	next := b.position.Add(b.velocity)
	if next.x > screenWidth || next.x < 0 {
		b.velocity = Vector2D{-b.velocity.x, b.position.y}
	}
	if next.y > screenHeight || next.y < 0 {
		b.velocity = Vector2D{b.velocity.x, b.position.y}
	}
}

func (b *Boid) start() {
	for {
		b.moveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

func (b *Boid) calculateAcceleration() Vector2D {
	upper, lower := b.position.AddV(viewRadius), b.position.AddV(-viewRadius)
	averageVelocity := Vector2D{0,0}
	count := 0.0

	for i := math.Max(lower.x, 0); i <= math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j <= math.Min(upper.y, screenHeight); j++ {
			if otherBoidId := boidMap[int(i)][int(j)]; otherBoidId != -1 && otherBoidId != b.id {
				if dist := boids[otherBoidId].position.Distance(b.position); dist < viewRadius {
					count++
					averageVelocity = averageVelocity.Add(boids[otherBoidId].velocity)
				}
			}
		}
	}
	accel := Vector2D{0, 0}
	if count > 0 {
		averageVelocity = averageVelocity.DivisionV(count)
		accel = averageVelocity.Subtraction(b.velocity).MultiplyV(adjRate)
	}
	return accel
}

func createBoid(bid int) {
	b := Boid{
		position: Vector2D{rand.Float64() * screenWidth, rand.Float64() * screenHeight},
		velocity: Vector2D{(rand.Float64() * 2) - 1.0, (rand.Float64() * 2) - 1.0},
		id:       bid,
	}
	boids[bid] = &b
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	go b.start()
}

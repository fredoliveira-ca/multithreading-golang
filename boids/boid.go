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
func (b *Boid) start() {
	for {
		b.moveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

func (b *Boid) calculateAcceleration() Vector2D {
	upper, lower := b.position.AddV(viewRadius), b.position.AddV(-viewRadius)
	averagePosition, averageVelocity, separation := Vector2D{0,0}, Vector2D{0,0}, Vector2D{0,0}
	count := 0.0

	rWlock.RLock()
	for i := math.Max(lower.x, 0); i <= math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j <= math.Min(upper.y, screenHeight); j++ {
			if otherBoidId := boidMap[int(i)][int(j)]; otherBoidId != -1 && otherBoidId != b.id {
				if dist := boids[otherBoidId].position.Distance(b.position); dist < viewRadius {
					count++
					averageVelocity = averageVelocity.Add(boids[otherBoidId].velocity)
					averagePosition = averageVelocity.Add(boids[otherBoidId].position)
					separation = separation.Add(b.position.Subtraction(boids[otherBoidId].position).DivisionV(dist))
				}
			}
		}
	}
	rWlock.RUnlock()

	accel := Vector2D{0, 0}
	if count > 0 {
		averagePosition, averageVelocity = averagePosition.DivisionV(count), averageVelocity.DivisionV(count)
		accelAlignment := averageVelocity.Subtraction(b.velocity).MultiplyV(adjRate)
		accelCohesion := averagePosition.Subtraction(b.position).MultiplyV(adjRate)
		accelSeparation := separation.MultiplyV(adjRate)
		accel = accel.Add(accelAlignment).Add(accelCohesion).Add(accelSeparation)
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


func (b *Boid) moveOne() {
	acceleration := b.calculateAcceleration()
	rWlock.Lock()
	b.velocity = b.velocity.Add(acceleration).limit(-1, 1)
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
	rWlock.Unlock()
}

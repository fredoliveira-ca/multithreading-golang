package main

import (
	"math"
	"math/rand"
	"time"
)

type Boid struct {
	position Vector2D
	velocity Vector2D
	id       int
}

func (b *Boid) calculateAcceleration() Vector2D {
	upper, lower := b.position.AddV(viewRadius), b.position.AddV(-viewRadius)
	averagePosition, averageVelocity, separation := Vector2D{0, 0}, Vector2D{0, 0}, Vector2D{0, 0}
	count := 0.0

	rWlock.RLock()
	for i := math.Max(lower.x, 0); i <= math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j <= math.Min(upper.y, screenHeight); j++ {
			if otherBoidId := boidMap[int(i)][int(j)]; otherBoidId != -1 && otherBoidId != b.id {
				if dist := boids[otherBoidId].position.Distance(b.position); dist < viewRadius {
					count++
					averageVelocity = averageVelocity.Add(boids[otherBoidId].velocity)
					averagePosition = averagePosition.Add(boids[otherBoidId].position)
					separation = separation.Add(b.position.Subtract(boids[otherBoidId].position).DivideV(dist))
				}
			}
		}
	}
	rWlock.RUnlock()

	accel := Vector2D{b.borderBounce(b.position.x, screenWidth), b.borderBounce(b.position.y, screenHeight)}
	if count > 0 {
		averagePosition, averageVelocity = averagePosition.DivideV(count), averageVelocity.DivideV(count)
		accelAlignment := averageVelocity.Subtract(b.velocity).MultiplyV(adjRate)
		accelCohesion := averagePosition.Subtract(b.position).MultiplyV(adjRate)
		accelSeparation := separation.MultiplyV(adjRate)
		accel = accel.Add(accelAlignment).Add(accelCohesion).Add(accelSeparation)
	}
	return accel
}

func (b *Boid) borderBounce(position, maxBorderPos float64) float64 {
	if position < viewRadius {
		return 1 / position
	} else if position > maxBorderPos - viewRadius {
		return 1 / (position - maxBorderPos)
	}
	return 0
}

func (b *Boid) moveOne() {
	acceleration := b.calculateAcceleration()
	rWlock.Lock()
	b.velocity = b.velocity.Add(acceleration).limit(-1, 1)
	boidMap[int(b.position.x)][int(b.position.y)] = -1
	b.position = b.position.Add(b.velocity)
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	rWlock.Unlock()
}

func (b *Boid) start() {
	for {
		b.moveOne()
		time.Sleep(5 * time.Millisecond)
	}
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
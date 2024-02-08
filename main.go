package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Vector struct {
	x, y float32
}

type Entity struct {
	pos, vel Vector
}

const size = 2

type Game struct {
	world_dims Vector
	entities   []Entity
}

func (v *Vector) Add(other Vector) Vector {
	return Vector{
		x: v.x + other.x,
		y: v.y + other.y,
	}
}

func (v *Vector) Sub(other Vector) Vector {
	return Vector{
		x: v.x - other.x,
		y: v.y - other.y,
	}
}

func (v *Vector) Mult(scalar float32) Vector {
	return Vector{
		x: v.x * scalar,
		y: v.y * scalar,
	}
}

func (v *Vector) Div(scalar float32) Vector {
	return Vector{
		x: v.x / scalar,
		y: v.y / scalar,
	}
}

func Zero() Vector {
	return Vector{
		x: 0.0,
		y: 0.0,
	}
}

func RandomInRange(min, max float32) float32 {
	return min + (max-min)*float32(rand.Float64())
}

func RandomVector(min, max float32) Vector {
	return Vector{
		x: RandomInRange(min, max),
		y: RandomInRange(min, max),
	}
}

func (g *Game) StepEntities() {
	for i := range g.entities {
		entity := &g.entities[i]
		entity.pos = entity.pos.Add(entity.vel)
	}
}

func (g *Game) Bounce() {
	for i := range g.entities {
		entity := &g.entities[i]
		if entity.pos.x < 0 {
			entity.vel.x = -entity.vel.x
			entity.pos.x = 0
		}

		if (entity.pos.x + size) > g.world_dims.x {
			entity.vel.x = -entity.vel.x
			entity.pos.x = g.world_dims.x - size
		}

		if entity.pos.y < 0 {
			entity.vel.y = -entity.vel.y
			entity.pos.y = 0
		}

		if (entity.pos.y + size) > g.world_dims.y {
			entity.vel.y = -entity.vel.y
			entity.pos.y = g.world_dims.y - size
		}

	}
}

func (g *Game) HandleInputs() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) || ebiten.IsKeyPressed(ebiten.KeyEscape) {
		fmt.Println("Game closed")
		return ebiten.Termination
	}
	return nil
}

func (g *Game) Init() {
	num_entities := 30000
	for i := 0; i < num_entities; i++ {
		g.entities = append(g.entities, Entity{
			pos: Vector{
				x: RandomInRange(0, g.world_dims.x),
				y: RandomInRange(0, g.world_dims.y),
			},
			vel: RandomVector(-1, 1),
		})
	}
	fmt.Println(g.entities[0])
}

func (g *Game) Update() error {
	if err := g.HandleInputs(); err != nil {
		return err
	}

	g.StepEntities()
	g.Bounce()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawEntities(screen)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()), 51, 51)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(g.world_dims.x), int(g.world_dims.y)
}

func main() {
	world_dims := Vector{640, 480}
	ebiten.SetWindowSize(int(world_dims.x), int(world_dims.y))
	ebiten.SetWindowTitle("Hello, World!")
	game := &Game{
		world_dims: world_dims,
	}
	game.Init()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) DrawEntities(screen *ebiten.Image) {
	color := color.RGBA{255, 255, 255, 255}
	for _, entity := range g.entities {
		vector.DrawFilledRect(screen, entity.pos.x, entity.pos.y, size, size, color, false)
	}
}

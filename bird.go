package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const gravaty = 1

type bird struct {
	time     int
	textures []*sdl.Texture

	x     float64
	y     float64
	speed float64
}

func newBird(r *sdl.Renderer) (*bird, error) {
	var textures []*sdl.Texture
	path := "res/imgs/bird_frame_%d.png"
	for i := 1; i <= 4; i++ {
		texture, err := img.LoadTexture(r, fmt.Sprintf(path, i))
		if err != nil {
			return nil, fmt.Errorf("could not load bird: %v", err)
		}
		textures = append(textures, texture)
	}
	return &bird{time: 0, textures: textures, x: 10, y: 250 - 43/2}, nil
}

func (b *bird) paint(r *sdl.Renderer) error {
	b.time++
	b.speed += gravaty
	b.y += b.speed
	if b.y >= 500-43 {
		b.speed = -10
	}
	rect := &sdl.Rect{X: int32(b.x), Y: int32(b.y), W: 50, H: 43}
	if err := r.Copy(b.textures[b.time%len(b.textures)], nil, rect); err != nil {
		return fmt.Errorf("could not copy bird: %v", err)
	}
	return nil
}

func (b *bird) destroy() {
	for _, t := range b.textures {
		t.Destroy()
	}
}

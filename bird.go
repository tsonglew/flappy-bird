package main

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	bw        = 50
	bh        = 43
	startBX   = 20
	gravaty   = 0.1
	jumpSpeed = 5
)

type bird struct {
	sync.RWMutex
	dead     bool
	time     int
	x, y     int
	speed    float64
	textures []*sdl.Texture
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
	return &bird{time: 0, textures: textures, x: startBX, y: (winHeigth - bh) / 2}, nil
}

func (b *bird) update() {
	b.Lock()
	defer b.Unlock()
	b.y += int(b.speed)
	if b.y >= winHeigth-bh {
		b.dead = true
	}
	b.time++
	b.speed += gravaty
}

func (b *bird) jump() {
	b.Lock()
	defer b.Unlock()
	b.speed = -jumpSpeed
}

func (b *bird) isDead() bool {
	b.Lock()
	defer b.Unlock()
	return b.dead
}

func (b *bird) restart() {
	b.Lock()
	defer b.Unlock()
	b.y = (winWidth - bh) / 2
	b.dead = false
	b.speed = 0
}

func (b *bird) touch(ps *pipes) {
	ps.RLock()
	defer ps.RUnlock()
	for _, p := range ps.pipes {
		p.RLock()
		if p.x-pw/2 <= b.x+bw/2 && b.x-bw/2 <= p.x+pw/2 {
			if p.inverted {
				if b.y <= p.h {
					b.dead = true
				}
			} else {
				if b.y >= winHeigth-p.h {
					b.dead = true
				}
			}
		}
		p.RUnlock()
	}
}

func (b *bird) paint(r *sdl.Renderer) error {
	b.RLock()
	defer b.RUnlock()
	rect := &sdl.Rect{X: int32(b.x), Y: int32(b.y), W: bw, H: bh}
	if err := r.Copy(b.textures[b.time/10%len(b.textures)], nil, rect); err != nil {
		return fmt.Errorf("could not copy bird: %v", err)
	}
	return nil
}

func (b *bird) destroy() {
	b.Lock()
	defer b.Unlock()
	for _, t := range b.textures {
		t.Destroy()
	}
}

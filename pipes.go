package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	pw     = 52
	phmin  = 100
	pSpeed = 2
)

type pipes struct {
	sync.RWMutex
	speed   float64
	texture *sdl.Texture
	pipes   []*pipe
}

type pipe struct {
	sync.RWMutex
	inverted bool
	x        int
	h        int
	w        int
}

func newPipes(r *sdl.Renderer) (*pipes, error) {
	texture, err := img.LoadTexture(r, "res/imgs/pipe.png")
	if err != nil {
		return nil, fmt.Errorf("could not load pipe: %v", err)
	}

	ps := &pipes{
		texture: texture,
		speed:   pSpeed,
	}

	go func() {
		for {
			ps.Lock()
			ps.pipes = append(ps.pipes, newPipe())
			ps.Unlock()
			time.Sleep(600 * time.Millisecond)
		}
	}()
	return ps, nil
}

func newPipe() *pipe {
	return &pipe{
		inverted: rand.Float64() > 0.5,
		x:        winWidth,
		h:        phmin + rand.Intn(220),
		w:        pw,
	}
}

func (ps *pipes) paint(r *sdl.Renderer) error {
	ps.RLock()
	defer ps.RUnlock()

	for _, p := range ps.pipes {
		if err := p.paint(r, ps.texture); err != nil {
			return fmt.Errorf("could not paint pipes: %v", err)
		}
	}
	return nil
}

func (p *pipe) paint(r *sdl.Renderer, texture *sdl.Texture) error {
	p.RLock()
	defer p.RUnlock()
	rect := &sdl.Rect{X: int32(p.x), Y: int32(winHeigth - p.h), W: int32(p.w), H: int32(p.h)}
	flip := sdl.FLIP_NONE
	if p.inverted {
		rect.Y = 0
		flip = sdl.FLIP_VERTICAL
	}
	if err := r.CopyEx(texture, nil, rect, 0, nil, flip); err != nil {
		return fmt.Errorf("could not copy pipe: %v", err)
	}
	return nil
}

func (ps *pipes) update() {
	ps.Lock()
	defer ps.Unlock()

	var rem []*pipe
	for _, p := range ps.pipes {
		p.Lock()
		defer p.Unlock()
		p.x -= int(ps.speed)
		if p.x >= -pw {
			rem = append(rem, p)
		}
	}
	ps.pipes = rem
}

func (ps *pipes) restart() {
	ps.Lock()
	defer ps.Unlock()
	ps.pipes = nil
}

func (ps *pipes) destroy() {
	ps.Lock()
	defer ps.Unlock()
	ps.texture.Destroy()
}

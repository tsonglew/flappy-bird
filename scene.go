package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const frequency = 10

type scene struct {
	bg    *sdl.Texture
	bird  *bird
	pipes *pipes
}

func newScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background: %v", err)
	}
	bird, err := newBird(r)
	if err != nil {
		return nil, fmt.Errorf("could not new bird: %v", err)
	}

	pipes, err := newPipes(r)
	if err != nil {
		return nil, fmt.Errorf("could not create pipes: %v", err)
	}

	return &scene{bg: bg, bird: bird, pipes: pipes}, nil
}

func (s *scene) run(events chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)
		tick := time.Tick(frequency * time.Millisecond)
		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
			case <-tick:
				s.update()
				if s.bird.dead {
					time.Sleep(time.Second)
					if err := drawTitle(r, "Game Over"); err != nil {
						errc <- err
					}
					time.Sleep(time.Second)
					s.restart()
				}
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()
	return errc
}

func (s *scene) handleEvent(e sdl.Event) bool {
	switch e.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseButtonEvent:
		s.bird.jump()
	default:
	}
	return false
}

func (s *scene) update() {
	s.bird.update()
	s.pipes.update()
	s.bird.touch(s.pipes)
}

func (s *scene) restart() {
	s.bird.restart()
	s.pipes.restart()
}

func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()

	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	if err := s.bird.paint(r); err != nil {
		return fmt.Errorf("could not paint bird: %v", err)
	}

	if err := s.pipes.paint(r); err != nil {
		return fmt.Errorf("could not paint pipes: %v", err)
	}

	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.destroy()
	s.pipes.destroy()
}

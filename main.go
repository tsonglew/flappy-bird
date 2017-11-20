package main

import (
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
}

func run() error {
	// Initialize SDL
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("could not init sdl: %v", err)
	}
	defer sdl.Quit()

	// Initialize TTF
	if err := ttf.Init(); err != nil {
		return fmt.Errorf("could not init ttf: %v", err)
	}
	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer w.Destroy()

	if err := drawTitle(r); err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}
	time.Sleep(5 * time.Second)

	s, err := newScene(r)
	if err != nil {
		return fmt.Errorf("could not create new scene: %v", err)
	}

	if err := s.paint(r); err != nil {
		return fmt.Errorf("could not paint: %v", err)
	}

	time.Sleep(5 * time.Second)
	return nil
}

func drawTitle(r *sdl.Renderer) error {
	r.Clear()

	// Open fonts from ttf
	f, err := ttf.OpenFont("res/fonts/Flappy.ttf", 15)
	if err != nil {
		return fmt.Errorf("could not open font: %v", err)
	}
	defer f.Close()

	s, err := f.RenderUTF8_Solid("Flappy Go", sdl.Color{R: 255, G: 0, B: 0, A: 255})
	if err != nil {
		return fmt.Errorf("could not render UTF8: %v", err)
	}
	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not create texture from surface: %v", err)
	}
	defer t.Destroy()

	if err := r.Copy(t, nil, nil); err != nil {
		return fmt.Errorf("could not copy texutre: %v", err)
	}

	r.Present()
	return nil
}

func drawBackground(r *sdl.Renderer) error {
	r.Clear()

	t, err := img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return fmt.Errorf("could not laod texture: %v", err)
	}
	defer t.Destroy()

	if r.Copy(t, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}
	r.Present()
	return nil
}

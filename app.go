package main

import (
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

type App struct {
	univ Universe
	univMtx sync.Mutex

	screen tcell.Screen
	quit   chan struct{}
}

func NewApp() App {
	app := App{
		quit: make(chan struct{}),
		univMtx: sync.Mutex{},
	}
	return app
}

func (app *App) initUniverse() {
	TermW, TermH := app.screen.Size()
	app.univ = NewUniverse(TermW, TermH)
	app.univ.Randomize()
}

func (app *App) Run() error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}

	app.screen = screen

	err = app.screen.Init()
	if err != nil {
		return err
	}
	app.initUniverse()

	defer app.screen.Fini()

	app.screen.Clear()
	screen.EnableMouse()

	err = app.loop()
	if err != nil {
		return err
	}

	return nil
}

func (app *App) putRune(i, j int, p rune, s tcell.Style) {
	app.screen.SetContent(j, i, p, nil, s)
}

func (app *App) putUniverse() {
	for i := 0; i < app.univ.h; i++ {
		for j := 0; j < app.univ.w; j++ {
			if app.univ.Alive(i, j) {
				app.putRune(i, j, tcell.RuneBlock, tcell.StyleDefault)
			} else {
				app.putRune(i, j, ' ', tcell.StyleDefault)
			}
		}
	}
}

func (app *App) loop() (err error) {
	go app.termevents()

	for {
		select {
		case <-app.quit:
			return

		case <-time.After(time.Second/10):
			app.univMtx.Lock()
			app.univ.Tick()
			app.univMtx.Unlock()
			app.putUniverse()
			app.screen.Show()
		}
	}
}

func (app *App) termevents() {
	for {
		ev := app.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC || ev.Key() == tcell.KeyESC {
				app.quit <- struct{}{}
			}
		case *tcell.EventMouse:
			but := ev.Buttons()
			if but == tcell.Button1 {
				j, i := ev.Position()
				if !app.univ.Exists(i, j) {
					continue
				}
				app.univMtx.Lock()
				app.univ.Set(i, j, true)
				app.univMtx.Unlock()

				app.putRune(i, j, tcell.RuneBlock, tcell.StyleDefault.Foreground(tcell.ColorAqua))
				app.screen.Show()
			}
		}
	}
}

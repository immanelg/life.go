package main

import (
	"time"

	"github.com/gdamore/tcell"
)

type App struct {
	univ Universe

	screen tcell.Screen
	events chan interface{}
	ticker time.Ticker
}

func NewApp() App {
	app := App{
		events: make(chan interface{}),
		ticker: *time.NewTicker(time.Second / 3),
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
	app.screen.EnableMouse()

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
	go app.listenTermeEvents()

	for {
		select {

		case ev := <-app.events:
			switch ev := ev.(type) {
			case QuitEvent:
				return nil
			case ResetEvent:
				app.univ.Randomize()
				app.screen.Show()
			case ClickEvent:
				i, j := ev.y, ev.x
				if !app.univ.Exists(i, j) {
					continue
				}
				app.univ.Set(i, j, true)
				app.putRune(i, j, tcell.RuneBlock, tcell.StyleDefault.Foreground(tcell.ColorAqua))
				app.screen.Show()
			}

		case <-app.ticker.C:
			app.univ.Tick()
			app.putUniverse()
			app.screen.Show()
		}
	}
}

func (app *App) listenTermeEvents() {
	for {
		ev := app.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC || ev.Key() == tcell.KeyESC {
				app.events <- QuitEvent{}
			} else if ev.Rune() == 'r' {
				app.events <- ResetEvent{}
			}
		case *tcell.EventMouse:
			but := ev.Buttons()
			if but == tcell.Button1 {
				x, y := ev.Position()
				app.events <- ClickEvent{x, y}
			}
		}
	}
}

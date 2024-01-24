package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"os"
	"time"
)

func main() {
	vacacionesFlag := flag.String("vacaciones", "", "El exacto momento que vas salir de vacaciones in RFC3339 format (e.g. 2024-02-01T18:30:00-03:00)")
	flag.Parse()
	if *vacacionesFlag == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	vacaciones, err := parseVacacionesFlag(vacacionesFlag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	g := &game{
		vacaciones: vacaciones,
		width:      600,
		height:     480,
	}

	ebiten.SetWindowTitle(fmt.Sprintf("vacaciones countdown to %s", vacaciones.Format(time.DateTime)))
	if err = ebiten.RunGame(g); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseVacacionesFlag(flag *string) (time.Time, error) {
	vacaciones, err := time.Parse(time.RFC3339, *flag)
	if err != nil {
		return time.Time{}, err
	}

	now := time.Now()
	if vacaciones.Before(now) {
		return time.Time{}, errors.New("you can not salir de vacaciones en el pasado!")
	}

	return vacaciones, nil
}

type game struct {
	vacaciones        time.Time
	timeRemaining     time.Duration
	vacacionesStarted bool
	message           string
	width, height     int
}

func (g *game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("too much to wait, right?")
	}

	g.loop()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, g.message)
}

func (g *game) Layout(w, h int) (int, int) {
	return g.width, g.height
}

func (g *game) loop() {
	if g.vacacionesStarted {
		return
	}

	g.timeRemaining = g.vacaciones.Sub(time.Now())
	if int(g.timeRemaining.Seconds()) < 0 {
		g.vacacionesStarted = true
		g.message = `tu momento llegó, dale! \o/`
		return
	}

	g.message = g.fmtDuration()
}

func (g *game) fmtDuration() string {
	total := int(g.timeRemaining.Seconds())
	days := total / (60 * 60 * 24)
	hours := total / (60 * 60) % 24
	minutes := int(total/60) % 60
	seconds := total % 60
	return fmt.Sprintf("días: %d horas: %d minutos: %d segundos: %d", days, hours, minutes, seconds)
}

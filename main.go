package main

import (
	"errors"
	"flag"
	"fmt"
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

	for range time.Tick(1 * time.Second) {
		timeRemaining := calculateTimeRemaining(vacaciones)

		if int(timeRemaining.Seconds()) <= 0 {
			fmt.Println("tu momento llegó, dale!")
			os.Exit(0)
		}

		fmt.Println(fmtDuration(timeRemaining))
	}
}

func parseVacacionesFlag(flag *string) (time.Time, error) {
	vacaciones, err := time.Parse(time.RFC3339, *flag)
	if err != nil {
		return time.Time{}, err
	}

	now := time.Now()
	if vacaciones.Before(now) {
		return time.Time{}, errors.New("you can not salir de vacaciones en el pasado")
	}

	return vacaciones, nil
}

func calculateTimeRemaining(vacaciones time.Time) time.Duration {
	now := time.Now()
	return vacaciones.Sub(now)
}

func fmtDuration(d time.Duration) string {
	total := int(d.Seconds())
	days := total / (60 * 60 * 24)
	hours := total / (60 * 60) % 24
	minutes := int(total/60) % 60
	seconds := total % 60
	return fmt.Sprintf("días: %d horas: %d minutos: %d segundos: %d", days, hours, minutes, seconds)
}

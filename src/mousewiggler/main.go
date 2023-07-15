package main

import (
	"flag"
	"math/rand"

	"time"

	"github.com/go-vgo/robotgo"
)

func doTask(sec int) {

	for {
		x, y := robotgo.GetMousePos()
		high := 1.0
		mouseDelay := 0.05
		for i := 0; i < 5; i++ {
			x1 := x + rand.Intn(100)
			time.Sleep(10 * time.Millisecond)
			robotgo.MoveSmooth(x1, y, high, mouseDelay)
			time.Sleep(10 * time.Millisecond)
			robotgo.MoveSmooth(x, y, high, mouseDelay)
		}
		time.Sleep(time.Duration(sec) * time.Second)

	}

}

func main() {
	secPtr := flag.Int("sec", 300, "-sec xxx")
	flag.Parse()
	doTask(*secPtr)
}

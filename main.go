package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"golang_test2/subpack"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 800
	screenHeight = 600
	space        = 10
	plotWidth    = screenWidth - space*2
	plotHeight   = screenHeight - space*2
)

type Game struct {
	count int
}

type PlotQso struct {
	WorkingNo int
	UnixStart int64
	UnixEnd   int64
	Band      string
}

var gZlogqso []*subpack.ZlogQso
var gPlotqso []*PlotQso
var gLatestWorkingNo int
var timestarted time.Time
var timelinesDrawed bool
var accImage *ebiten.Image

func init() {
	gZlogqso = subpack.Readfile()
	timestarted = time.Now()
	ebiten.SetTPS(1)
	gLatestWorkingNo = 1
	timelinesDrawed = false
}

func (g *Game) Update() error {
	g.count++
	g.count %= plotHeight

	qsonow := checkElapsedTime()
	if qsonow.Callsign != "NA" {
		fmt.Println(qsonow.TimeQSO, qsonow.Callsign, qsonow.Band)
		Add(qsonow)
	}
	return nil
}

func Add(a subpack.ZlogQso) error {
	ut := ToUnixTime(a.DateQSO + a.TimeQSO)
	for _, gp := range gPlotqso {
		if gp.Band == a.Band {
			if (ut - gp.UnixEnd) < 10*60 {
				gp.UnixEnd = ut
				return nil
			}
		}
	}
	gLatestWorkingNo++
	plot := PlotQso{
		WorkingNo: gLatestWorkingNo,
		UnixStart: ut,
		UnixEnd:   ut,
		Band:      a.Band,
	}
	gPlotqso = append(gPlotqso, &plot)
	return nil
}

func ToUnixTime(timeString string) int64 {
	layout := "2006/01/0215:04"
	t, err := time.Parse(layout, timeString)
	if err != nil {
		fmt.Println(err)
	}
	return t.Unix()
}
func DrawTimelines(screen *ebiten.Image) {
	ypos := space * 3
	step_y := plotHeight / (1 * 60)

	for i := 0; i <= 1*60; i++ {
		xs := float32(space * 3)
		if i%10 == 0 {
			xs = space * 2
		} else if i%60 == 0 {
			xs = space
		}
		vector.StrokeLine(accImage, xs, float32(ypos), space+plotWidth, float32(ypos), 1, color.White, true)
		ypos += step_y
	}
	timelinesDrawed = true
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !timelinesDrawed {
		DrawTimelines(screen)
	}

	if g.count%10 == 0 {

		//fmt.Println("Draw: ", g.count, ", ")
		//fmt.Print("ZlogQso: ", gZlogqso[g.count], ", ")
	}

	//vector.StrokeLine(screen, space, space, space+plotWidth, space, 1, color.RGBA{0xff, 0xff, 0xff, 0xff}, true)
	//vector.StrokeLine(screen, space, space+plotHeight, space+plotWidth, space+plotHeight, 1, color.RGBA{0xff, 0xff, 0x00, 0xff}, true)

	for i := 0; i < g.count; i += 20 {
		//cf := float32(i)
		//vector.StrokeLine(screen, space, space+cf, space+plotWidth, space+cf, 1, color.White, true)
	}

	//for i := 0; i < g.count; i += 10 {
	cf := float32(g.count)
	//vector.DrawFilledRect(screen, 50+cf, 50+cf, 100+cf, 100+cf, color.RGBA{0x80, 0x80, 0x80, 0xc0}, true)

	vector.StrokeRect(screen, space+plotWidth/2, space, plotWidth/3, cf, 2, color.RGBA{0x00, 0xff, 0x00, 0xff}, false)
	//}
	screen.DrawImage(accImage, &ebiten.DrawImageOptions{})
	//vector.DrawFilledCircle(screen, 400, 400, 100, color.RGBA{0x80, 0x00, 0x80, 0x80}, true)
	//vector.StrokeCircle(screen, 400, 400, 10+cf, 10+cf/2, color.RGBA{0xff, 0x80, 0xff, 0xff}, true)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("ActualTPS: %0.2f", ebiten.ActualTPS()))
	ebitenutil.DebugPrintAt(screen, fmt.Sprint(g.count), 0, 20)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	accImage = ebiten.NewImage(screenWidth, screenHeight)
	for _, qso := range gZlogqso {
		fmt.Print(qso.Callsign, ", ")
	}
	fmt.Println("ZlogQso: ", len(gZlogqso), "records")

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Shapes (Ebitengine Demo)")
	ebiten.SetWindowPosition(100, 100)

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func checkElapsedTime() subpack.ZlogQso {
	min1 := 0.0
	min2 := 0.0
	for _, record := range gZlogqso {
		dur := record.TimeQSO[0:2] + "m" + record.TimeQSO[3:5] + "s"
		qsoelapsedSeconds2, _ := time.ParseDuration(dur)
		min1 = qsoelapsedSeconds2.Seconds()
		nowelapsedSeconds := time.Since(timestarted)
		min2 = math.Floor(nowelapsedSeconds.Seconds())

		if min1 == min2 {
			//fmt.Printf("The elapsed time for record %f %s has passed.\n", min1, record.Callsign)
			return *record
		} else {
			//fmt.Printf("The elapsed time for record %s has not passed yet.\n", record.CallSign)
			//return *gZlogqso[0]
		}
		//fmt.Println(min1, min2)
	}
	//fmt.Println(min1, min2)
	na := *gZlogqso[0]
	na.Callsign = "NA"
	return na
}

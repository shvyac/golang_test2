package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"showQsoTX2/subpack"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth      = 800
	screenHeight     = 1200
	space            = 10
	plotWidth        = screenWidth - space*2
	plotHeight       = screenHeight - space*2
	contestTimeStart = "2015/04/2601:00"
	contestTimeEnd   = "2015/04/2621:00"
	plotRangeMinutes = 60 * 4
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
var bands []string
var bandPlotWidth int

func init() {
	gZlogqso = subpack.Readfile()
	timestarted = time.Now()
	ebiten.SetTPS(1)
	gLatestWorkingNo = 0
	timelinesDrawed = false
	bands = []string{"3.5", "7", "14", "21", "28", "50"}
	bandPlotWidth = (plotWidth - space*2) / len(bands)
}

func (g *Game) Update() error {
	g.count++

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
			if (ut - gp.UnixEnd) <= 10*60 {
				gp.UnixEnd = ut
				//fmt.Println("gPlotqso: ", len(gPlotqso),gp.WorkingNo, gp.UnixStart, gp.UnixEnd, gp.Band)
				return nil
			}
		}
	}
	gLatestWorkingNo++
	plot := PlotQso{
		WorkingNo: gLatestWorkingNo,
		UnixStart: ut,
		UnixEnd:   ut + 60,
		Band:      a.Band,
	}
	gPlotqso = append(gPlotqso, &plot)
	//fmt.Println("gPlotqso: ", len(gPlotqso),gPlotqso[len(gPlotqso)-1].WorkingNo, gPlotqso[len(gPlotqso)-1].UnixStart, gPlotqso[len(gPlotqso)-1].UnixEnd, gPlotqso[len(gPlotqso)-1].Band)
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
	strokeWidth := float32(.5)
	loc, _ := time.LoadLocation("Asia/Tokyo")
	layout := "2006/01/0215:04"
	ts, err := time.ParseInLocation(layout, contestTimeStart, loc)
	fmt.Print("ts: ", ts.Hour(), ts.Minute(), ", ")
	if err != nil {
		fmt.Println(err)
	}
	ypos := space * 3
	minutes := plotRangeMinutes
	step_y := plotHeight / minutes

	for i := 0; i <= minutes; i++ {
		xs := float32(space * 3)
		vector.StrokeLine(accImage, xs, float32(ypos), space+plotWidth, float32(ypos), strokeWidth, color.White, true)
		if i%60 == 0 {
			xs = space
			vector.StrokeLine(accImage, xs, float32(ypos), space+plotWidth, float32(ypos), strokeWidth, color.White, true)
			ebitenutil.DebugPrintAt(accImage, fmt.Sprintf("%02d:00", ts.Hour()), int(5), ypos-10)
		} else if i%10 == 0 {
			xs = space * 2
			vector.StrokeLine(accImage, xs, float32(ypos), space+plotWidth, float32(ypos), strokeWidth, color.White, true)
			ebitenutil.DebugPrintAt(accImage, fmt.Sprintf(":%02d", ts.Minute()), int(10), ypos-10)
		}
		ts = ts.Add(time.Minute)
		ypos += step_y
	}
	xpos := space * 3
	xinc := bandPlotWidth
	for _, ba := range bands {
		vector.StrokeLine(accImage, float32(xpos), space*3, float32(xpos), space+plotHeight, strokeWidth, color.White, true)
		ebitenutil.DebugPrintAt(accImage, fmt.Sprintf("%s", ba), xpos+xinc/2, space)
		xpos += xinc
	}
	vector.StrokeLine(accImage, float32(xpos), space*3, float32(xpos), space+plotHeight, strokeWidth, color.White, true)
	timelinesDrawed = true
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !timelinesDrawed {
		DrawTimelines(screen)
	}

	for _, gp := range gPlotqso {
		xpos := space * 3
		ypos := space * 3
		width := bandPlotWidth
		widthPad := width / 10

		te := time.Unix(gp.UnixEnd, 0)
		ts := time.Unix(gp.UnixStart, 0)
		boxHeight := te.Sub(ts)
		cs := time.Unix(ToUnixTime(contestTimeStart), 0)
		boxStart := ts.Sub(cs)
		height := boxHeight.Minutes() * (plotHeight / plotRangeMinutes)
		//fmt.Println(te.Format("15:04"), ts.Format("15:04"), tt.Minutes())
		//fmt.Printf("gp.UnixEnd: %d, gp.UnixStart: %d, height: %d\n", gp.UnixEnd, gp.UnixStart, height)
		for i, ba := range bands {
			if gp.Band == ba {
				xpos += width*i + widthPad
				ypos += int(boxStart.Minutes() * (plotHeight / plotRangeMinutes))
				//fmt.Println("ts.Minute() ", xpos, ypos, ts.Minute(), height)
				vector.StrokeRect(screen, float32(xpos), float32(ypos),
					float32(width-2*widthPad), float32(height), 1, color.White, false)
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", gp.WorkingNo), xpos, ypos)
			}
		}
	}
	screen.DrawImage(accImage, &ebiten.DrawImageOptions{})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("ActualTPS: %0.2f", ebiten.ActualTPS()))
	ebitenutil.DebugPrintAt(screen, fmt.Sprint(g.count), 200, 00)
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
	for _, qso := range gZlogqso {
		contestStart, _ := time.Parse("2006/01/0215:04", contestTimeStart)
		clockElapsedSeconds := time.Since(timestarted).Seconds()
		clockElapsedMinutes := clockElapsedSeconds
		contestElapsed := contestStart.Add(time.Duration(clockElapsedMinutes) * time.Minute)

		//fmt.Println(contestElapsed.Format("2006/01/02 15:04"))
		elaDate := contestElapsed.Format("2006/01/02")
		elaTime := contestElapsed.Format("15:04")
		// dur := qso.TimeQSO[0:2] + "m" + qso.TimeQSO[3:5] + "s"
		// qsoelapsedSeconds2, _ := time.ParseDuration(dur)
		// min1 = qsoelapsedSeconds2.Seconds()
		// nowelapsedSeconds := time.Since(timestarted)
		// min2 = math.Floor(nowelapsedSeconds.Seconds())
		conDate := qso.DateQSO
		conTime := qso.TimeQSO

		if conDate == elaDate && conTime == elaTime {
			return *qso
		} else {
		}
	}
	na := *gZlogqso[0]
	na.Callsign = "NA"
	return na
}

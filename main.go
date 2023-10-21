package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"showQsoTX2/subpack"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth      = 2400
	screenHeight     = 1200
	space            = 10
	plotWidth        = screenWidth - space*2
	plotHeight       = screenHeight - space*2
	contestTimeStart = "2015/04/2521:00"
	contestTimeEnd   = "2015/04/2621:00"
	plotRangeMinutes = 60 * 3
)

type Game struct {
	count int
	keys  []ebiten.Key
}

type PlotQsoBox struct {
	WorkingNo int
	UnixStart int64
	UnixEnd   int64
	Band      string
	//NumberQso int
}
type PlotQso struct {
	UnixStart int64
	Callsign  string
}

var (
	gZlogqso         []*subpack.ZlogQso
	gPlotQsoBox      []*PlotQsoBox
	gPlotCall        map[int][]*PlotQso
	gLatestWorkingNo int
	timeAppStarted   time.Time
	timelinesDrawed  bool
	accImage         *ebiten.Image
	bands            []string
	bandPlotWidth    int
	contestShowStart time.Time
	lastInputTime    time.Time
	acceptInput      bool
)

func init() {
	gZlogqso = subpack.Readfile()
	timeAppStarted = time.Now()
	ebiten.SetTPS(4)
	gLatestWorkingNo = 0
	timelinesDrawed = false
	bands = []string{"3.5", "7", "14", "21", "28", "50", "se1", "se2"}
	bandPlotWidth = (plotWidth - space*2) / len(bands)
	contestShowStart = ToJstTimeFromString(contestTimeStart)
	gPlotCall = make(map[int][]*PlotQso)
}

func (g *Game) Update() error {
	g.count++

	qsonows := checkElapsedTime()
	for _, qsonow := range qsonows {
		if qsonow.Callsign != "NA" {
			fmt.Println(qsonow.TimeQSO, qsonow.Callsign, qsonow.Band)
			Add(qsonow)
		}
	}
	//g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	// KeyArrowDown
	if acceptInput && inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		acceptInput = false
		lastInputTime = time.Now()
		fmt.Println("KeyArrowDown", lastInputTime.Format("2006/01/02 15:04:05"))
		contestShowStart = contestShowStart.Add(time.Hour)
		//fmt.Println("new time---", contestShowStart.Format("2006/01/02 15:04"))

		if contestShowStart.Before(ToJstTimeFromString(contestTimeStart)) {
			//contestShowStart = ToTime(contestTimeStart)
		} else if contestShowStart.After(ToJstTimeFromString(contestTimeEnd)) {
			//contestShowStart = ToTime(contestTimeEnd)
		}
		//fmt.Println(contestShowStart.Format("2006/01/02 15:04"))
		timelinesDrawed = false
	}
	// KeyArrowUp
	if acceptInput && inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		acceptInput = false
		lastInputTime = time.Now()
		//fmt.Println("KeyArrowUp", lastInputTime.Format("2006/01/02 15:04:05"))
		contestShowStart = contestShowStart.Add(-time.Hour)
		//fmt.Println("new time---", contestShowStart.Format("2006/01/02 15:04"))
		timelinesDrawed = false
	}

	// if !acceptInput && time.Since(lastInputTime).Seconds() > 1 {
	// 	acceptInput = true
	// 	//fmt.Println("acceptInput = true", time.Now().Format("2006/01/02 15:04:05"))
	// }

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !timelinesDrawed {
		DrawTimelines(screen)
	}
	screen.DrawImage(accImage, &ebiten.DrawImageOptions{})

	for _, gp := range gPlotQsoBox {
		xpos := space * 3
		ypos := space * 3
		width := bandPlotWidth
		widthPad := width / 10

		te := time.Unix(gp.UnixEnd, 0)
		ts := time.Unix(gp.UnixStart, 0)
		boxHeight := te.Sub(ts)
		cs := ToJstTimeFromString(contestTimeStart) //time.Unix(ToUnixTime(contestTimeStart), 0)
		cs = time.Unix(ToUnixTimeFromString(contestShowStart.Format("2006/01/0215:04")), 0)
		//cs2 := contestShowStart
		//fmt.Println("cs: ", cs.Format("2006/01/02 15:04 "), cs2.Format("2006/01/02 15:04"))
		boxStart := ts.Sub(cs)
		height := boxHeight.Minutes() * (plotHeight / plotRangeMinutes)
		//fmt.Println(te.Format("15:04"), ts.Format("15:04"), cs.Format("15:04"))
		//fmt.Printf("gp.UnixEnd: %d, gp.UnixStart: %d, height: %d\n", gp.UnixEnd, gp.UnixStart, height)
		for i, ba := range bands {
			if gp.Band == ba {
				xpos += width*i + widthPad
				ypos += int(boxStart.Minutes() * (plotHeight / plotRangeMinutes))
				//fmt.Println("ts.Minute() ", xpos, ypos, ts.Minute(), height)
				vector.StrokeRect(screen, float32(xpos), float32(ypos),
					float32(width-2*widthPad), float32(height), 1, color.RGBA{0xff, 0xff, 0xff, 0xff}, true)
				startjst := ToJstTimeFromUnix(gp.UnixStart)
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s", startjst.Format("15:04")), xpos, ypos)
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", len(gPlotCall[gp.WorkingNo])), xpos+50, ypos)
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", (gp.UnixEnd-gp.UnixStart)/60), xpos+80, ypos)
			}
		}
		for _, pq := range gPlotCall[gp.WorkingNo] {
			call := pq.Callsign
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s", call), xpos+110, ypos)
			ypos += 10
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("ActualTPS: %0.2f", ebiten.ActualTPS()))
	ebitenutil.DebugPrintAt(screen, fmt.Sprint(g.count), 200, 00)
	ebitenutil.DebugPrintAt(screen, fmt.Sprint(contestShowStart.Format("2006/01/02 15:04")), 300, 00)
	acceptInput = true
}

func Add(a subpack.ZlogQso) error {
	ut := ToUnixTimeFromString(a.DateQSO + a.TimeQSO)
	for _, gp := range gPlotQsoBox {
		if gp.Band == a.Band {
			if (ut - gp.UnixEnd) <= 10*60 {
				gp.UnixEnd = ut + 60
				//gp.NumberQso++
				//fmt.Println("append gPlotqso: ", len(gPlotqso), gp.WorkingNo, gp.UnixStart, gp.UnixEnd, gp.Band)
				if !IsCall(gPlotCall, a.Callsign) {
					gPlotCall[gp.WorkingNo] = append(gPlotCall[gp.WorkingNo],
						&PlotQso{UnixStart: gp.UnixStart, Callsign: a.Callsign})
				}
				return nil
			}
		}
	}
	plot := PlotQsoBox{
		WorkingNo: gLatestWorkingNo,
		UnixStart: ut,
		UnixEnd:   ut + 60,
		Band:      a.Band,
		//NumberQso: 1,
	}

	gPlotQsoBox = append(gPlotQsoBox, &plot)
	//fmt.Println("new gPlotqso: ", len(gPlotqso), gPlotqso[len(gPlotqso)-1].WorkingNo, gPlotqso[len(gPlotqso)-1].UnixStart, gPlotqso[len(gPlotqso)-1].UnixEnd, gPlotqso[len(gPlotqso)-1].Band)
	//if IsCall(gPlotCall, gLatestWorkingNo-1, a.Callsign) {
	gPlotCall[gLatestWorkingNo] = append(gPlotCall[gLatestWorkingNo],
		&PlotQso{UnixStart: ut, Callsign: a.Callsign})
	gLatestWorkingNo++
	return nil
}

func IsCall(MapCalls map[int][]*PlotQso, Callsign string) bool {
	for _, slice := range MapCalls {
		for _, s := range slice {
			if s.Callsign == Callsign {
				//fmt.Printf("Found struct at key %d index %d\n", key, i)
				return true
			}
		}
	}
	return false
}

func ToUnixTimeFromString(timeString string) int64 {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	layout := "2006/01/0215:04"
	ts, err := time.ParseInLocation(layout, timeString, loc)
	if err != nil {
		fmt.Println(err)
	}
	return ts.Unix()
}

func ToJstTimeFromString(timeString string) time.Time {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	layout := "2006/01/0215:04"
	ts, err := time.ParseInLocation(layout, timeString, loc)
	if err != nil {
		fmt.Println(err)
	}
	return ts
}

func ToJstTimeFromUnix(unixtime int64) time.Time {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	t := time.Unix(unixtime, 0)
	return t.In(loc)
}

func DrawTimelines(screen *ebiten.Image) {
	accImage.Clear()
	//accImage.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff}) // white
	//screen.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})   // white
	strokeWidth := float32(.5)
	loc, _ := time.LoadLocation("Asia/Tokyo")
	layout := "2006/01/0215:04"
	ts, err := time.ParseInLocation(layout, contestTimeStart, loc)
	ts = contestShowStart
	fmt.Print("ts: ", ts.Hour(), ts.Minute(), ", ")
	if err != nil {
		fmt.Println(err)
	}
	ypos := space * 3
	minutes := plotRangeMinutes
	step_y := plotHeight / minutes
	// horizontal lines, hours and minutes
	for i := 0; i <= minutes; i++ {
		xs := float32(space * 3)
		vector.StrokeLine(accImage, xs, float32(ypos), space+plotWidth, float32(ypos), strokeWidth,
			color.RGBA{0x80, 0x80, 0x80, 0xff}, true)
		if i%60 == 0 {
			xs = space
			vector.StrokeLine(accImage, xs, float32(ypos), space+plotWidth, float32(ypos), 1,
				color.RGBA{0x00, 0xff, 0x00, 0xff}, true)
			ebitenutil.DebugPrintAt(accImage, fmt.Sprintf("%02d:00", ts.Hour()), int(5), ypos-10)
		} else if i%10 == 0 {
			xs = space * 2
			vector.StrokeLine(accImage, xs, float32(ypos), space+plotWidth, float32(ypos), strokeWidth,
				color.RGBA{0x80, 0x80, 0x00, 0xff}, true)
			ebitenutil.DebugPrintAt(accImage, fmt.Sprintf(":%02d", ts.Minute()), int(10), ypos-10)
		}
		ts = ts.Add(time.Minute)
		ypos += step_y
	}
	// vertical lines, bands
	xpos := space * 3
	xinc := bandPlotWidth
	for _, ba := range bands {
		vector.StrokeLine(accImage, float32(xpos), space*3, float32(xpos), space+plotHeight, strokeWidth,
			color.RGBA{0x80, 0x80, 0x00, 0xff}, true)
		ebitenutil.DebugPrintAt(accImage, fmt.Sprintf("%s", ba), xpos+xinc/2, space)
		xpos += xinc
	}
	vector.StrokeLine(accImage, float32(xpos), space*3, float32(xpos), space+plotHeight, strokeWidth,
		color.RGBA{0x80, 0x80, 0x00, 0xff}, true)
	timelinesDrawed = true
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	acceptInput = true
	lastInputTime = time.Now()

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

func checkElapsedTime() []subpack.ZlogQso {
	var QSOs []subpack.ZlogQso
	contestStart, _ := time.Parse("2006/01/0215:04", contestTimeStart)
	contestStart = ToJstTimeFromString(contestTimeStart)
	clockElapsedSeconds := time.Since(timeAppStarted).Seconds()
	clockElapsedMinutes := clockElapsedSeconds
	contestElapsed := contestStart.Add(time.Duration(clockElapsedMinutes) * time.Minute)

	for _, qso := range gZlogqso {

		//fmt.Println(contestElapsed.Format("2006/01/02 15:04"))
		//elaDate := contestElapsed.Format("2006/01/02")
		//elaTime := contestElapsed.Format("15:04")
		// dur := qso.TimeQSO[0:2] + "m" + qso.TimeQSO[3:5] + "s"
		// qsoelapsedSeconds2, _ := time.ParseDuration(dur)
		// min1 = qsoelapsedSeconds2.Seconds()
		// nowelapsedSeconds := time.Since(timestarted)
		// min2 = math.Floor(nowelapsedSeconds.Seconds())
		conDateTime := ToJstTimeFromString(qso.DateQSO + qso.TimeQSO)
		//conTime :=ToJstTimeFromString( qso.TimeQSO)
		con := conDateTime.Format("15:04")
		ela := contestElapsed.Format("15:04")
		if con == ela {
			QSOs = append(QSOs, *qso)
		} else if conDateTime.After(contestElapsed) {
			return QSOs
		}
	}

	na := *gZlogqso[0]
	na.Callsign = "NA"
	QSOs = append(QSOs, na)
	return QSOs
}

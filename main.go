package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"showQsoTX2/subpack"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth      = 1200
	screenHeight     = 1200
	space            = 10
	plotWidth        = screenWidth - space*2
	plotHeight       = screenHeight - space*2
	plotRangeMinutes = 60 * 3
)

type Game struct {
	count int
	keys  []ebiten.Key
}

type BoxQso struct {
	BoxNo    int
	BoxStart time.Time
	BoxEnd   time.Time
	BoxBand  string
}
type CallQso struct {
	CallStart time.Time
	Callsign  string
}

var (
	gZlogqso         []*subpack.ZlogQso
	gPlotQsoBox      []*BoxQso
	gPlotCall        map[int][]*CallQso
	gLatestWorkingNo int
	ReadTime         string
	timeAppStarted   time.Time
	timelinesDrawed  bool
	accImage         *ebiten.Image
	bands            []string
	bandPlotWidth    float32
	contestShowStart time.Time
	lastKeyinTime    time.Time
	acceptKeyin      bool
	QsoBoxNo         int
	gfile            *os.File
	contestTimeStart = time.Date(2015, 4, 26, 02, 0, 0, 0, time.Local) //"2015/04/2521:00"
	contestTimeEnd   = time.Date(2015, 4, 26, 21, 0, 0, 0, time.Local) //"2015/04/2621:00"
)

func init() {
	gZlogqso = subpack.Readfile()
	timeAppStarted = time.Now()
	ebiten.SetTPS(1)
	gLatestWorkingNo = 0
	timelinesDrawed = false
	bands = []string{"3.5", "7", "14", "21", "28", "50", "se1", "se2"}
	bandPlotWidth = float32((plotWidth - space*2) / len(bands))
	contestShowStart = contestTimeStart
	gPlotCall = make(map[int][]*CallQso)
	ReadTime = ""
	QsoBoxNo = 0
}

func (g *Game) Update() error {
	g.count++
	qsonows := checkElapsedTime()
	for _, qsonow := range qsonows {
		if qsonow.Callsign != "NA" {
			fmt.Println("Update--> ", len(qsonows), qsonow.DateTime, qsonow.Callsign, qsonow.Band)
			AddQso(qsonow)
		}
	}
	//g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	// KeyArrowDown
	if acceptKeyin && inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		acceptKeyin = false
		lastKeyinTime = time.Now()
		fmt.Println("KeyArrowDown", lastKeyinTime.Format("2006/01/02 15:04:05"))
		contestShowStart = contestShowStart.Add(time.Hour)
		//fmt.Println("new time---", contestShowStart.Format("2006/01/02 15:04"))

		if contestShowStart.Before(contestTimeStart) {
			//contestShowStart = ToTime(contestTimeStart)
		} else if contestShowStart.After(contestTimeEnd) {
			//contestShowStart = ToTime(contestTimeEnd)
		}
		//fmt.Println(contestShowStart.Format("2006/01/02 15:04"))
		timelinesDrawed = false
	}
	// KeyArrowUp
	if acceptKeyin && inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		acceptKeyin = false
		lastKeyinTime = time.Now()
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
		xpos := float32(space * 3)
		ypos := float32(space * 3)
		width := bandPlotWidth
		widthPad := width / 10
		te := gp.BoxEnd
		ts := gp.BoxStart
		boxHeight := te.Sub(ts)
		cs := contestTimeStart //time.Unix(ToUnixTime(contestTimeStart), 0)
		cs = time.Unix(ToUnixTimeFromString(contestShowStart.Format("2006/01/0215:04")), 0)
		//cs2 := contestShowStart
		//fmt.Println("cs: ", cs.Format("2006/01/02 15:04 "), cs2.Format("2006/01/02 15:04"))
		boxStart := ts.Sub(cs)
		height := boxHeight.Minutes() * (plotHeight / plotRangeMinutes)
		//fmt.Println(te.Format("15:04"), ts.Format("15:04"), cs.Format("15:04"))
		//fmt.Printf("gp.UnixEnd: %d, gp.UnixStart: %d, height: %d\n", gp.UnixEnd, gp.UnixStart, height)
		for i, ba := range bands {
			if gp.BoxBand == ba {
				xpos += width * float32(i) // + widthPad
				ypos += float32(boxStart.Minutes()) * float32(plotHeight/plotRangeMinutes)
				//fmt.Println("ts.Minute() ", xpos, ypos, ts.Minute(), height)
				vector.StrokeRect(screen, float32(xpos), float32(ypos),
					float32(width-2*widthPad), float32(height), 1, color.RGBA{0xff, 0xff, 0xff, 0xff}, true)
				startjst := gp.BoxStart
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s", startjst.Format("15:04")), int(xpos), int(ypos))
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", len(gPlotCall[gp.BoxNo])), int(xpos+40), int(ypos))
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%v", gp.BoxEnd.Sub(gp.BoxStart)), int(xpos+0), int(ypos+10))
			}
		}
		for _, pq := range gPlotCall[gp.BoxNo] {
			call := pq.Callsign
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s", call), int(xpos+60), int(ypos))
			ypos += 10
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("ActualTPS: %0.2f", ebiten.ActualTPS()))
	ebitenutil.DebugPrintAt(screen, fmt.Sprint(g.count), 200, 00)
	ebitenutil.DebugPrintAt(screen, fmt.Sprint(contestShowStart.Format("2006/01/02 15:04")), 300, 00)
	acceptKeyin = true
}

func AddQso(a subpack.ZlogQso) error {

	if a.Callsign == "JH4WBY" {
		fmt.Print("JH4WBY: ", a.Callsign, a.DateTime, a.Band)
	}
	ut := a.DateTime
	for _, gp := range gPlotQsoBox {
		if gp.BoxBand == a.Band {
			if ut.Sub(gp.BoxEnd) <= 10*time.Minute {
				gp.BoxEnd = ut.Add(time.Minute)
				if !IsCall(gPlotCall[gp.BoxNo], a.Callsign) {
					gPlotCall[gp.BoxNo] = append(
						gPlotCall[gp.BoxNo], &CallQso{CallStart: ut, Callsign: a.Callsign})
				}
				return nil
			}
		}
	}
	plot := BoxQso{
		BoxNo:    gLatestWorkingNo,
		BoxStart: ut,
		BoxEnd:   ut.Add(time.Minute),
		BoxBand:  a.Band,
	}
	QsoBoxNo = -1
	for i, qb := range gPlotQsoBox {
		if qb.BoxBand == plot.BoxBand {
			t1 := qb.BoxStart.Format("15:04")
			t2 := qb.BoxEnd.Format("15:04")
			fmt.Println("\t\t\t\tBox: ", i, t1, t2, qb.BoxBand)
			QsoBoxNo = qb.BoxNo
		}
	}
	gPlotQsoBox = append(gPlotQsoBox, &plot)
	//fmt.Println("new gPlotqso: ", len(gPlotqso), gPlotqso[len(gPlotqso)-1].WorkingNo, gPlotqso[len(gPlotqso)-1].UnixStart, gPlotqso[len(gPlotqso)-1].UnixEnd, gPlotqso[len(gPlotqso)-1].Band)
	//if IsCall(gPlotCall, gLatestWorkingNo-1, a.Callsign) {
	gPlotCall[gLatestWorkingNo] = append(gPlotCall[gLatestWorkingNo],
		&CallQso{CallStart: ut, Callsign: a.Callsign})
	gLatestWorkingNo++

	if QsoBoxNo > -1 {
		for i, pq := range gPlotCall[QsoBoxNo] {
			s1 := pq.CallStart.Format("15:04")
			fmt.Println("\t\t\t\tCall: ", QsoBoxNo, i, s1, pq.Callsign, plot.BoxBand)
			_, err := fmt.Fprintln(gfile, QsoBoxNo, s1, pq.Callsign, plot.BoxBand)
			if err != nil {
				fmt.Println(err)
			}
		}
		for i, box := range gPlotQsoBox {
			s1 := box.BoxStart
			s2 := box.BoxEnd
			if box.BoxNo == QsoBoxNo && s2.Sub(s1).Minutes() < 10 {
				box.BoxEnd = box.BoxStart.Add(11 * time.Minute)
				fmt.Println("\t\t\t\t10Min: ", i, s1, s2, box.BoxEnd.Format("15:04"), box.BoxBand)
			}
		}
	}
	return nil
}

func IsCall(MapCalls []*CallQso, Callsign string) bool {
	for _, slice := range MapCalls {
		//for _, s := range slice {
		if slice.Callsign == Callsign {
			//fmt.Printf("Found struct at key %d index %d\n", key, i)
			return true
		}
		//}
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
	// loc, _ := time.LoadLocation("Asia/Tokyo")
	// layout := "2006/01/0215:04"
	// ts, err := time.ParseInLocation(layout, contestTimeStart, loc)
	ts := contestShowStart
	//fmt.Print("ts: ", ts.Hour(), ts.Minute(), ", ")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	ypos := float32(space * 3)
	//minutes := plotRangeMinutes
	step_y := float32(plotHeight / plotRangeMinutes)
	// horizontal lines, hours and minutes
	xpos := float32(space * 3)
	for i := 0; i <= plotRangeMinutes; i++ {
		vector.StrokeLine(accImage, xpos, (ypos), space+plotWidth, (ypos), strokeWidth,
			color.RGBA{0x80, 0x80, 0x80, 0xff}, true)
		if i%60 == 0 {
			xpos = float32(space)
			vector.StrokeLine(accImage, xpos, (ypos), space+plotWidth, (ypos), 1,
				color.RGBA{0x00, 0xff, 0x00, 0xff}, true)
			ebitenutil.DebugPrintAt(accImage, fmt.Sprintf("%02d:00", ts.Hour()), int(5), int(ypos-10))
		} else if i%10 == 0 {
			xpos = float32(space * 2)
			vector.StrokeLine(accImage, xpos, (ypos), space+plotWidth, (ypos), strokeWidth,
				color.RGBA{0x80, 0x80, 0x00, 0xff}, true)
			ebitenutil.DebugPrintAt(accImage, fmt.Sprintf(":%02d", ts.Minute()), int(10), int(ypos-10))
		}
		ts = ts.Add(time.Minute)
		ypos += step_y
	}
	// vertical lines, bands
	//xpos = space * 3
	xinc := bandPlotWidth
	for _, ba := range bands {
		vector.StrokeLine(accImage, float32(xpos), space*3, float32(xpos), space+plotHeight, strokeWidth,
			color.RGBA{0x80, 0x80, 0x00, 0xff}, true)
		ebitenutil.DebugPrintAt(accImage, fmt.Sprintf("%s", ba), int(xpos+xinc/2), space)
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
	acceptKeyin = true
	lastKeyinTime = time.Now()

	accImage = ebiten.NewImage(screenWidth, screenHeight)
	for _, qso := range gZlogqso {
		fmt.Print(qso.Callsign, ", ")
	}
	fmt.Println("ZlogQso: ", len(gZlogqso), "records")

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Shapes (Ebitengine Demo)")
	ebiten.SetWindowPosition(100, 100)

	f, err := os.Create("log.txt")
	if err != nil {
		fmt.Println(err)
	}
	gfile = f
	defer f.Close()

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func checkElapsedTime() []subpack.ZlogQso {
	var QSOs []subpack.ZlogQso
	//contestStart, _ := time.Parse("2006/01/0215:04", contestTimeStart)
	contestStart := contestTimeStart
	clockElapsedSeconds := time.Since(timeAppStarted).Seconds()
	clockElapsedMinutes := clockElapsedSeconds
	contestElapsed := contestStart.Add(time.Duration(clockElapsedMinutes) * time.Minute)
	if contestElapsed.After(contestTimeEnd) {
		os.Exit(0)
	}
	for _, qso := range gZlogqso {
		conDateTime := qso.DateTime // ToJstTimeFromString(qso.DateQSO + qso.TimeQSO)
		//conTime :=ToJstTimeFromString( qso.TimeQSO)
		con := conDateTime.Format("15:04")
		ela := contestElapsed.Format("15:04")
		if con == ela && ReadTime != ela {
			QSOs = append(QSOs, *qso)
		}
		if conDateTime.After(contestElapsed) {
			ReadTime = ela
			return QSOs
		}
	}
	na := *gZlogqso[0]
	na.Callsign = "NA"
	QSOs = append(QSOs, na)
	return QSOs
}

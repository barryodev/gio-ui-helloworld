package main

import (
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"
)

func main() {
	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			drawProgressBar(&ops, time.Now())

			drawHeader(gtx)

			drawRoundedSquareWithTriangle(&ops)

			drawSecondSquare(&ops)

			drawThirdCircle(&ops)

			drawFiveRectangles(&ops)

			drawImage(&ops)

			doButton(gtx)

			e.Frame(gtx.Ops)
		}
	}
}

func addColorOperation(colorToApply color.NRGBA, ops *op.Ops) {
	paint.ColorOp{Color: colorToApply}.Add(ops)
}

func drawRect(x, y int, ops *op.Ops) {
	clip.Rect{Max: image.Pt(x, y)}.Add(ops)
	paint.PaintOp{}.Add(ops)
}

func moveDrawing(x, y float32, ops *op.Ops) {
	op.Offset(f32.Pt(x, y)).Add(ops)
}

func drawTriangle(ops *op.Ops) {
	var path clip.Path
	path.Begin(ops)
	path.Move(f32.Pt(50, 0))
	path.Line(f32.Pt(50, 100))
	path.Line(f32.Pt(-100, 0))
	path.Line(f32.Pt(50, -100))
	clip.Outline{Path: path.End()}.Op().Add(ops)
	drawRect(100,100, ops)
}


func drawRoundedSquareWithTriangle(ops *op.Ops) {
	defer op.Save(ops).Load()

	const r = 15 // roundness
	bounds := f32.Rect(100, 100, 200, 200)
	clip.RRect{Rect: bounds, SE: r, SW: r, NW: r, NE: r}.Add(ops)
	moveDrawing(100, 100, ops)
	addColorOperation(color.NRGBA{R: 0x80, A: 0xFF}, ops)
	drawRect(100, 100, ops)

	drawTriangle(ops)
	addColorOperation(color.NRGBA{G: 0x80, A: 0xFF}, ops)
	paint.PaintOp{}.Add(ops)
}

func drawSecondSquare(ops *op.Ops) {
	defer op.Save(ops).Load()

	moveDrawing(250, 100, ops)
	addColorOperation(color.NRGBA{B: 0x80, A: 0xFF}, ops)
	drawRect(100,100, ops)
}

func drawThirdCircle(ops *op.Ops) {
	defer op.Save(ops).Load()

	moveDrawing(350, 100, ops)
	addColorOperation(color.NRGBA{R: 190, G: 82, B: 209, A: 0xFF}, ops)
	drawCircle(100, 50, 50, ops)
}

func drawCircle(x, y, radius float32, ops *op.Ops) {
	clip.Circle{Center: f32.Pt(x, y), Radius: radius}.Add(ops)
	paint.PaintOp{}.Add(ops)
}

func drawFiveRectangles(ops *op.Ops) {
	defer op.Save(ops).Load()

	moveDrawing(50, 225, ops)
	addColorOperation(color.NRGBA{R: 252, G: 245, B: 53, A: 0xFF}, ops)

	for i := 1; i < 6; i++ {
		drawRectAndMove(float32(i * 75), 0, ops)
	}
}

func drawRectAndMove(x, y float32, ops *op.Ops) {
	defer op.Save(ops).Load()

	moveDrawing(x, y, ops)
	drawRect(50, 50, ops)
}

var startTime = time.Now()
func drawProgressBar(ops *op.Ops, now time.Time) {
	// Calculate how much of the progress bar to draw,
	// based on the current time.
	elapsed := now.Sub(startTime)

	var duration = 10 * time.Second
	progress := elapsed.Seconds() / duration.Seconds()
	if progress < 1 {
		op.InvalidateOp{}.Add(ops)
	} else {
		progress = 1
	}

	defer op.Save(ops).Load()
	width := 200 * float32(progress)
	clip.Rect{Max: image.Pt(int(width), 20)}.Add(ops)
	paint.ColorOp{Color: color.NRGBA{R: 0x80, A: 0xFF}}.Add(ops)
	paint.ColorOp{Color: color.NRGBA{G: 0x80, A: 0xFF}}.Add(ops)
	paint.PaintOp{}.Add(ops)
}

func drawImage(ops *op.Ops) {
	defer op.Save(ops).Load()

	existingImageFile, err := os.Open("img/gopher.png")
	if err != nil {
		panic(err)
	}
	defer existingImageFile.Close()

	loadedImage, err := png.Decode(existingImageFile)
	if err != nil {
		panic(err)
	}

	moveDrawing(100, 300, ops)
	imageOp := paint.NewImageOp(loadedImage)
	imageOp.Add(ops)
	op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(4, 4)))
	paint.PaintOp{}.Add(ops)
}

var tag = new(bool) // We could use &pressed for this instead.
var pressed = false

func doButton(gtx layout.Context) {
	// Make sure we don’t pollute the graphics context.
	defer op.Save(gtx.Ops).Load()

	moveDrawing(350, 300, gtx.Ops)

	// Process events that arrived between the last frame and this one.
	for _, ev := range gtx.Queue.Events(tag) {
		if x, ok := ev.(pointer.Event); ok {
			switch x.Type {
			case pointer.Press:
				pressed = true
			case pointer.Release:
				pressed = false
			}
		}
	}

	// Confine the area of interest to a 100x100 rectangle.
	pointer.Rect(image.Rect(0, 0, 100, 100)).Add(gtx.Ops)
	// Declare the tag.
	pointer.InputOp{
		Tag:   tag,
		Types: pointer.Press | pointer.Release,
	}.Add(gtx.Ops)

	writeLabel(gtx, "Click me")

	moveDrawing(0, 20, gtx.Ops)

	clip.Rect{Max: image.Pt(100, 100)}.Add(gtx.Ops)
	var c color.NRGBA
	if pressed {
		c = color.NRGBA{R: 0xFF, A: 0xFF}
	} else {
		c = color.NRGBA{G: 0xFF, A: 0xFF}
	}
	paint.ColorOp{Color: c}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func writeLabel(gtx layout.Context, labelText string) {
	l := material.Body2(material.NewTheme(gofont.Collection()), labelText )
	l.Color = color.NRGBA{R: 145, G: 50, B: 168, A: 255}
	l.Alignment = text.Start
	l.Layout(gtx)
}

func drawHeader(gtx layout.Context) {
	defer op.Save(gtx.Ops).Load()
	moveDrawing(250, 0, gtx.Ops)

	writeLabel(gtx,"Trying out the gio gui")
}
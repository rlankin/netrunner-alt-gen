package basic

import (
	"fmt"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

func (fb FrameBasic) Back() art.Drawer {

	return art.DrawerFunc(func(ctx *canvas.Context, card *nrdb.Printing) error {

		canvasWidth, canvasHeight := ctx.Size()

		attributionBoxTop := canvasHeight * 0.12
		// attributionBoxBottom := canvasHeight * 0.032
		attributionBoxBottom := canvasHeight * 0.0591
		attributionBoxHeight := attributionBoxTop - attributionBoxBottom
		attributionBoxLeft := canvasWidth * 0.25
		attributionBoxRight := canvasWidth * 0.75
		attributionBoxRadius := canvasWidth * 0.01

		fb.drawRoundedBox(ctx, attributionBoxTop, attributionBoxBottom, attributionBoxLeft, attributionBoxRight, attributionBoxRadius)

		attributionFontSize := attributionBoxHeight * 0.6
		attributionTextMaxWidth := (attributionBoxRight - attributionBoxLeft) * 0.9

		attributionString := fmt.Sprintf("%s<BR>Generated using \"%s\" by %s<BR>netrunner-alt-gen %s", card.Attributes.Title, fb.Algorithm, fb.Designer, fb.Version)

		if fb.Algorithm == "" && fb.Designer == "" {
			attributionString = fmt.Sprintf("%s<BR>Layout by netrunner-alt-gen %s", card.Attributes.Title, fb.Version)
		}
		if fb.Algorithm == "" && fb.Designer != "" {
			attributionString = fmt.Sprintf("%s<BR>Design by %s<BR>Layout by netrunner-alt-gen %s", card.Attributes.Title, fb.Designer, fb.Version)
		}

		attributionText := fb.getFittedText(ctx, attributionString, attributionFontSize, attributionTextMaxWidth, 0, canvas.Center)

		attributionTextX := attributionBoxLeft + ((attributionBoxRight - attributionBoxLeft) * 0.05)
		attributionTextY := (attributionBoxTop - (attributionBoxHeight-attributionText.Bounds().H)*0.5)

		ctx.DrawText(attributionTextX, attributionTextY, attributionText)

		return nil
	})
}

func (fb FrameBasic) drawRoundedBox(ctx *canvas.Context, top, bottom, left, right, radius float64) {

	strokeWidth := getStrokeWidth(ctx)

	ctx.Push()
	ctx.SetFillColor(fb.getColorBG())
	ctx.SetStrokeColor(fb.getColorBorder())
	ctx.SetStrokeWidth(strokeWidth)

	path := &canvas.Path{}
	path.MoveTo(left, top-radius)
	path.QuadTo(left, top, left+radius, top)
	path.LineTo(right-radius, top)
	path.QuadTo(right, top, right, top-radius)
	path.LineTo(right, bottom+radius)
	path.QuadTo(right, bottom, right-radius, bottom)
	path.LineTo(left+radius, bottom)
	path.QuadTo(left, bottom, left, bottom+radius)
	path.Close()

	ctx.DrawPath(0, 0, path)
	ctx.Pop()
}
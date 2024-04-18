package netrunner

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/mangofeet/netrunner-alt-gen/art"
	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

type FrameResource struct{}

func (FrameResource) Draw(ctx *canvas.Context, card *nrdb.Printing) error {

	canvasWidth, canvasHeight := ctx.Size()

	strokeWidth := canvasHeight * 0.0023

	log.Printf("strokeWidth: %f", strokeWidth)

	factionBaseColor := art.GetFactionBaseColor(card.Attributes.FactionID)
	factionColor := color.RGBA{
		R: uint8(math.Max(0, math.Min(float64(int64(factionBaseColor.R)-48), 255))),
		G: uint8(math.Max(0, math.Min(float64(int64(factionBaseColor.G)-48), 255))),
		B: uint8(math.Max(0, math.Min(float64(int64(factionBaseColor.B)-48), 255))),
		A: 0xff,
	}

	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	titleBoxHeight := (canvasHeight / 16)

	titleBoxTop := canvasHeight - (canvasHeight / 12)
	titleBoxBottom := titleBoxTop - titleBoxHeight
	titleBoxArcStart := canvasWidth - (canvasWidth / 3)
	titleBoxRight := canvasWidth - (canvasWidth / 16)
	titleBoxArcCP1 := titleBoxRight - (canvasWidth / 48)

	costContainerR := titleBoxHeight * 0.667
	costContainerStart := costContainerR

	titlePath := &canvas.Path{}
	titlePath.MoveTo(0, titleBoxTop)

	// background for cost, top
	titlePath.LineTo(costContainerStart, titleBoxTop)
	titlePath.QuadTo(costContainerStart+costContainerR, titleBoxTop+(costContainerR*0.8), costContainerStart+(costContainerR*2), titleBoxTop)

	// arc down on right side
	titlePath.LineTo(titleBoxArcStart, titleBoxTop)
	titlePath.QuadTo(titleBoxArcCP1, titleBoxTop, titleBoxRight, titleBoxBottom)

	// background for cost, bottom
	titlePath.LineTo(costContainerStart+(costContainerR*2), titleBoxBottom)
	titlePath.QuadTo(costContainerStart+costContainerR, titleBoxBottom-(costContainerR*0.8), costContainerStart, titleBoxBottom)

	// finish title box
	titlePath.LineTo(0, titleBoxBottom)
	titlePath.Close()

	ctx.DrawPath(0, 0, titlePath)
	ctx.Pop()

	// outline for cost circle
	ctx.Push()
	ctx.SetFillColor(transparent)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)
	costOutline := canvas.Circle(costContainerR)
	ctx.DrawPath(costContainerStart+(costContainerR), titleBoxTop-(titleBoxHeight*0.5), costOutline)
	ctx.Pop()

	// bottom text box
	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	textBoxHeight := canvasHeight / 3
	textBoxLeft := canvasWidth / 8
	textBoxRight := canvasWidth - (canvasWidth / 12)
	textBoxArcRadius := (canvasHeight / 32)
	// textBoxArc1StartY := textBoxHeight - textBoxArcRadius
	// textBoxArc1EndX := textBoxLeft + textBoxArcRadius

	textBoxArc2StartX := textBoxRight - textBoxArcRadius
	textBoxArc2EndY := textBoxHeight - textBoxArcRadius

	ctx.MoveTo(textBoxLeft, 0)
	// ctx.LineTo(textBoxLeft, textBoxArc1StartY)
	// ctx.QuadTo(textBoxLeft, textBoxHeight, textBoxArc1EndX, textBoxHeight)
	ctx.LineTo(textBoxLeft, textBoxHeight)

	ctx.LineTo(textBoxArc2StartX, textBoxHeight)
	ctx.QuadTo(textBoxRight, textBoxHeight, textBoxRight, textBoxArc2EndY)

	ctx.LineTo(textBoxRight, 0)

	ctx.FillStroke()
	ctx.Pop()

	ctx.Push()
	ctx.SetFillColor(factionColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	influenceHeight := textBoxHeight * 0.55
	influenceWidth := canvasHeight / 48

	influenceCost := 0
	if card.Attributes.InfluenceCost != nil {
		influenceCost = *card.Attributes.InfluenceCost
	}
	ctx.DrawPath(textBoxRight-(influenceWidth/2), 0, influence(influenceHeight, influenceWidth, influenceCost))

	ctx.Pop()

	// type box
	ctx.Push()
	ctx.SetFillColor(bgColor)
	ctx.SetStrokeColor(textColor)
	ctx.SetStrokeWidth(strokeWidth)

	typeBoxHeight := textBoxHeight * 0.17
	typeBoxBottom := textBoxHeight + strokeWidth*0.5
	typeBoxLeft := textBoxLeft
	typeBoxRight := canvasWidth - (canvasWidth / 6)

	typeBoxArcRadius := (canvasHeight / 32)
	typeBoxArc1StartY := typeBoxBottom + typeBoxHeight - typeBoxArcRadius
	typeBoxArc1EndX := typeBoxLeft + typeBoxArcRadius

	typeBoxArc2StartX := typeBoxRight - typeBoxArcRadius
	typeBoxArc2EndY := typeBoxBottom + typeBoxHeight - typeBoxArcRadius

	ctx.MoveTo(typeBoxLeft, typeBoxBottom)
	ctx.LineTo(typeBoxLeft, typeBoxArc1StartY)
	ctx.QuadTo(typeBoxLeft, typeBoxHeight+typeBoxBottom, typeBoxArc1EndX, typeBoxHeight+typeBoxBottom)

	ctx.LineTo(typeBoxArc2StartX, typeBoxHeight+typeBoxBottom)
	ctx.QuadTo(typeBoxRight, typeBoxHeight+typeBoxBottom, typeBoxRight, typeBoxArc2EndY)

	ctx.LineTo(typeBoxRight, typeBoxBottom)

	ctx.FillStroke()

	ctx.Pop()

	// render card text

	// not sure how these sizes actually correlate to the weird
	// pixel/mm setup I'm using, but these work
	fontSizeTitle := titleBoxHeight * 2
	fontSizeCost := titleBoxHeight * 3
	fontSizeCard := titleBoxHeight * 1.2

	titleTextX := costContainerStart + (costContainerR * 2) + (costContainerR / 3)
	titleTextY := titleBoxTop - titleBoxHeight*0.1
	ctx.DrawText(titleTextX, titleTextY, getCardText(getTitleText(card), fontSizeTitle, titleBoxRight, titleBoxHeight))
	// ctx.DrawText(titleTextX, titleTextY, canvas.NewTextLine(getFont(fontSizeTitle, canvas.FontRegular), getTitleText(card), canvas.Left))

	if card.Attributes.Cost != nil {
		costTextX := costContainerStart
		costTextY := titleBoxBottom + titleBoxHeight/2
		ctx.DrawText(costTextX, costTextY, canvas.NewTextBox(
			getFont(fontSizeCost, canvas.FontBlack), fmt.Sprint(*card.Attributes.Cost),
			costContainerR*2, 0,
			canvas.Center, canvas.Center, 0, 0))
	}

	cardTextPadding := canvasWidth * 0.02
	cardTextX := textBoxLeft + cardTextPadding
	cardTextY := textBoxHeight - cardTextPadding
	typeTextX := cardTextX
	typeTextY := typeBoxBottom + typeBoxHeight - cardTextPadding
	cardTextBoxW := textBoxRight - textBoxLeft - (cardTextPadding * 2)
	cardTextBoxH := textBoxHeight
	typeTextBoxW := typeBoxRight - typeBoxLeft - (cardTextPadding * 2)
	typeTextBoxH := typeBoxHeight

	var tText *canvas.Text

	typeName := getTypeName(card.Attributes.CardTypeID)

	if card.Attributes.DisplaySubtypes != nil {
		tText = getCardText(fmt.Sprintf("<strong>%s</strong> - %s", typeName, *card.Attributes.DisplaySubtypes), fontSizeCard, typeTextBoxW, typeTextBoxH)
	} else {
		tText = getCardText(fmt.Sprintf("<strong>%s</strong>", typeName), fontSizeCard, typeTextBoxW, typeTextBoxH)
	}

	ctx.DrawText(typeTextX, typeTextY, tText)

	cText := getCardText(card.Attributes.Text, fontSizeCard, cardTextBoxW, cardTextBoxH)

	_, lastLineH := cText.Heights()

	for lastLineH > cardTextBoxH*0.75 {
		fontSizeCard -= strokeWidth
		cText = getCardText(card.Attributes.Text, fontSizeCard, cardTextBoxW, cardTextBoxH)
		_, lastLineH = cText.Heights()
	}

	ctx.DrawText(cardTextX, cardTextY, cText)

	return nil
}

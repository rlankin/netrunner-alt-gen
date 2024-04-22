package basic

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mangofeet/nrdb-go"
	"github.com/tdewolff/canvas"
)

func getTitle(card *nrdb.Printing) string {

	if !card.Attributes.IsUnique {
		return card.Attributes.Title
	}

	return fmt.Sprintf("♦ %s", card.Attributes.Title)
}

func getTitleText(ctx *canvas.Context, card *nrdb.Printing, fontSize, maxWidth, height float64) *canvas.Text {
	title := getTitle(card)

	text := getCardText(title, fontSize, maxWidth*2, height, canvas.Left)

	strokeWidth := getStrokeWidth(ctx)

	for text.Bounds().W > maxWidth {
		fontSize -= strokeWidth
		text = getCardText(title, fontSize, maxWidth*2, height, canvas.Left)
	}

	return text
}

func getCardText(text string, fontSize, cardTextBoxW, cardTextBoxH float64, align canvas.TextAlign) *canvas.Text {

	regFace := getFont(fontSize, canvas.FontRegular)
	boldFace := getFont(fontSize, canvas.FontBold)
	italicFace := getFont(fontSize, canvas.FontItalic)
	arrowFace := getFont(fontSize, canvas.FontExtraBold)
	uniqueFace := getFont(fontSize, canvas.FontExtraBold)

	rt := canvas.NewRichText(regFace)

	var parts []string
	strongParts := strings.Split(text, "<strong>")
	for _, p := range strongParts {
		emParts := strings.Split(p, "<em>")
		parts = append(parts, emParts...)
	}

	for _, part := range parts {

		if strings.Contains(part, "→") {
			subParts := strings.Split(part, "→")
			writeTextPart(rt, subParts[0], regFace, boldFace, italicFace)
			rt.WriteFace(arrowFace, "→")
			part = subParts[1]
		}

		if strings.Contains(part, "♦") {
			subParts := strings.Split(part, "♦")
			writeTextPart(rt, subParts[0], regFace, boldFace, italicFace)
			rt.WriteFace(uniqueFace, "♦")
			part = subParts[1]
		}

		writeTextPart(rt, part, regFace, boldFace, italicFace)
	}

	return rt.ToText(
		cardTextBoxW, cardTextBoxH,
		align, canvas.Top,
		0, 0)

}

func writeTextPart(rt *canvas.RichText, text string, regFace, boldFace, italicFace *canvas.FontFace) {

	text = strings.ReplaceAll(text, "\n", "\n\n")

	if strings.Contains(text, "</strong>") {
		subParts := strings.Split(text, "</strong>")
		writeChunk(rt, subParts[0], boldFace)
		text = subParts[1]
	}

	if strings.Contains(text, "</em>") {
		subParts := strings.Split(text, "</em>")
		writeChunk(rt, subParts[0], italicFace)
		text = subParts[1]
	}

	writeChunk(rt, text, regFace)

}

var replacementCheck = regexp.MustCompile(`\[[a-z-]+\]`)

func replaceSymbol(rt *canvas.RichText, symbol, svgName, text string, face *canvas.FontFace, scaleFactor, translateFactor float64) string {
	if strings.Contains(text, symbol) {
		subParts := strings.Split(text, symbol)
		for _, chunk := range subParts[:len(subParts)-1] {
			writeChunk(rt, chunk, face)
			path := mustLoadGameAsset(svgName).Scale(face.Size*scaleFactor, face.Size*scaleFactor).Transform(canvas.Identity.ReflectY().Translate(0, face.Size*-1*translateFactor))
			rt.WritePath(path, textColor, canvas.FontMiddle)
		}
		text = subParts[len(subParts)-1]
		if len(text) == 0 || (text[0] != ' ' && text[0] != ',' && text[0] != '.') {
			text = " " + text
		}
	}

	return text

}

func writeChunk(rt *canvas.RichText, text string, face *canvas.FontFace) {

	text = replaceSymbol(rt, "[mu]", "Mu", text, face, 0.0002, 0.8)
	text = replaceSymbol(rt, "[credit]", "CREDIT", text, face, 0.000025, 0.8)
	text = replaceSymbol(rt, "[recurring-credit]", "RECURRING_CREDIT", text, face, 0.00014, 0.8)
	text = replaceSymbol(rt, "[click]", "CLICK", text, face, 0.0002, 1)
	text = replaceSymbol(rt, "[subroutine]", "SUBROUTINE", text, face, 0.0002, 1)
	text = replaceSymbol(rt, "[trash]", "TRASH_ABILITY", text, face, 0.0002, 1)

	if replacementCheck.MatchString(text) {
		writeChunk(rt, text, face)
		return
	}

	rt.WriteFace(face, text)

}

func getTypeName(typeID string) string {
	switch typeID {
	case "program":
		return "Program"
	case "resource":
		return "Resource"
	case "hardware":
		return "Hardware"
	case "event":
		return "Event"
	case "runner_identity", "corp_identity":
		return "Identity"
	case "ice":
		return "Ice"
	case "asset":
		return "Asset"
	case "upgrade":
		return "Upgrade"
	}

	return typeID
}

type textBoxDimensions struct {
	left, right, bottom, top float64
	width, height            float64
	align                    canvas.TextAlign
}

func drawCardText(ctx *canvas.Context, card *nrdb.Printing, fontSize, indentCutoff, indent float64, box textBoxDimensions) {

	if box.align == 0 {
		box.align = canvas.Left
	}

	strokeWidth := getStrokeWidth(ctx)

	paddingLR, paddingTB := getCardTextPadding(ctx)
	x := box.left + paddingLR
	y := box.height - paddingTB
	if box.top != 0 {
		y = box.top - paddingTB
	}
	w := box.right - box.left - (paddingLR * 2.5)
	h := box.height

	cText := getCardText(card.Attributes.Text, fontSize, w, h, box.align)

	var leftoverText string

	_, lastLineH := cText.Heights()

	for lastLineH > h*0.75 {
		fontSize -= strokeWidth
		cText = getCardText(card.Attributes.Text, fontSize, w, h, box.align)
		_, lastLineH = cText.Heights()
	}

	i := 0
	_, lastLineH = cText.Heights()
	for lastLineH > indentCutoff {

		i++

		lines := strings.Split(card.Attributes.Text, "\n")

		leftoverText = strings.Join(lines[len(lines)-i:], "\n")
		newText := strings.Join(lines[:len(lines)-i], "\n")

		cText = getCardText(newText, fontSize, w, h, box.align)

		_, lastLineH = cText.Heights()

	}

	ctx.DrawText(x, y, cText)

	if leftoverText != "" {
		newCardTextX := x + indent
		if !cText.Empty() {
			y = y - (lastLineH + fontSize*0.4)
		}

		cText := getCardText(leftoverText, fontSize, w-(newCardTextX-x)-w*0.03, h, box.align)
		ctx.DrawText(newCardTextX, y, cText)
	}

}

func getCardTextPadding(ctx *canvas.Context) (lr, tb float64) {
	canvasWidth, _ := ctx.Size()

	lr = canvasWidth * 0.03
	tb = canvasWidth * 0.02

	return lr, tb

}

func drawTypeText(ctx *canvas.Context, card *nrdb.Printing, fontSize float64, box textBoxDimensions) {

	if box.align == 0 {
		box.align = canvas.Left
	}

	paddingLR, PaddingTB := getCardTextPadding(ctx)

	x := box.left + paddingLR
	y := box.bottom + box.height - PaddingTB
	w := box.right - box.left - (paddingLR * 2)
	h := box.height

	typeText := getTypeText(card, fontSize, w, h, box.align)

	ctx.DrawText(x, y, typeText)

}

func getTypeText(card *nrdb.Printing, fontSize, w, h float64, align canvas.TextAlign) *canvas.Text {
	var tText *canvas.Text
	typeName := getTypeName(card.Attributes.CardTypeID)

	if card.Attributes.DisplaySubtypes != nil {
		tText = getCardText(fmt.Sprintf("<strong>%s</strong> - %s", typeName, *card.Attributes.DisplaySubtypes), fontSize, w, h, align)
	} else {
		tText = getCardText(fmt.Sprintf("<strong>%s</strong>", typeName), fontSize, w, h, align)
	}

	return tText
}

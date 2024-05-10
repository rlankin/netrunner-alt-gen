package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/mangofeet/netrunner-alt-gen/art/entangler"
	"github.com/spf13/cobra"
)

var entanglerCmd = &cobra.Command{
	Use:   "entangler [card name or printing ID]",
	Args:  cobra.MinimumNArgs(1),
	Short: `Generate a card using the "entangler" algorithm`,
	Run: func(cmd *cobra.Command, args []string) {

		cardName := strings.Join(args, " ")

		if err := generateCardEntangler(cardName); err != nil {
			log.Println("error:", err)
			os.Exit(1)
		}

	},
}

func generateCardEntangler(cardName string) error {
	printing, err := getCardData(cardName)
	if err != nil {
		return err
	}
	log.Printf("generating %s", printing.Attributes.StrippedTitle)

	var nGridP *float64
	if gridPercent >= 0 {
		nGridP = &gridPercent
	}

	ns := entangler.Entangler{
		MinWalkers:   netspaceWalkersMin,
		MaxWalkers:   netspaceWalkersMax,
		GridPercent:  nGridP,
		Color:        parseColor(baseColor),
		ColorBG:      parseColor(netspaceColorBG),
		WalkerColor1: parseColor(walkerColor1),
		WalkerColor2: parseColor(walkerColor2),
		WalkerColor3: parseColor(walkerColor3),
		WalkerColor4: parseColor(walkerColor4),
		GridColor1:   parseColor(gridColor1),
		GridColor2:   parseColor(gridColor2),
		GridColor3:   parseColor(gridColor3),
		GridColor4:   parseColor(gridColor4),
		RingColor1:   parseColor(altColor1),
		RingColor2:   parseColor(altColor2),
		RingColor3:   parseColor(altColor3),
		RingColor4:   parseColor(altColor4),
	}

	return generateCard(ns, printing, "entangler", "mangofeet")
}

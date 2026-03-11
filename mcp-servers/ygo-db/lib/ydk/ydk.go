package ydk

import (
	"strconv"
	"strings"
)

// ParseYDKFile returns main, extra, side, each map is code of card -> number of card.
func ParseYDKFile(s string) (map[int]int, map[int]int, map[int]int) {
	mainDeck := make(map[int]int)
	extraDeck := make(map[int]int)
	sideDeck := make(map[int]int)

	currentSection := "" // "" | "main" | "extra" | "side"

	lines := strings.Split(s, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for section markers
		if strings.HasPrefix(line, "#main") {
			currentSection = "main"
			continue
		}
		if strings.HasPrefix(line, "#extra") {
			currentSection = "extra"
			continue
		}
		if strings.HasPrefix(line, "!side") {
			currentSection = "side"
			continue
		}

		// Skip other comment lines (lines starting with #)
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Try to parse as a card code
		code, err := strconv.Atoi(line)
		if err != nil {
			// Skip lines that are not valid card codes
			continue
		}

		// Add to the appropriate deck
		switch currentSection {
		case "main":
			mainDeck[code]++
		case "extra":
			extraDeck[code]++
		case "side":
			sideDeck[code]++
		}
	}

	return mainDeck, extraDeck, sideDeck
}

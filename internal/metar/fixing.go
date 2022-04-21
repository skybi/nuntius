package metar

import "strings"

var replacements = map[rune]rune{
	'–': '-',
	' ': ' ',
}

func normalize(r rune) rune {
	rep, ok := replacements[r]
	if ok {
		return rep
	}
	return r
}

// fix applies fixes on a raw METAR that solve common problems experienced over time
func fix(raw string) string {
	runes := []rune(raw)
	var builder strings.Builder

	// Clean up the raw string
	space := false
	for _, r := range runes {
		// Normalize the rune (i.e. replacing weird characters with ASCII ones)
		r = normalize(r)

		// Remove stacked spaces ('   ' -> ' ')
		if !space && r == ' ' {
			space = true
			builder.WriteRune(r)
			continue
		} else if space && r != ' ' {
			space = false
		}
		if !space {
			builder.WriteRune(r)
		}
	}
	clean := []rune(builder.String())

	// Sometimes the time is not led by a trailing 'Z'
	if len(clean) >= 12 && clean[11] == ' ' {
		appendix := append([]rune{'Z'}, clean[11:]...)
		clean = append(clean[:11], appendix...)
	}

	return string(clean)
}

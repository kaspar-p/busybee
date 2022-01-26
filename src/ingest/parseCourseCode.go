package ingest

import "strings"

func ParseEventTitle(summary string) string {
	potentialMarkers := []string{"H1", "Y1"}

	// Default to the entire string if no marker found
	var index int
	for _, marker := range potentialMarkers {
		index = strings.Index(summary, marker)
		if index != -1 {
			break
		}
	}

	if index == -1 {
		index = len(summary)
	}

	return summary[:index]
}

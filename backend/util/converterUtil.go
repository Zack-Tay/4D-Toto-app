package util

import (
	"strconv"
	"strings"
)

// parseCurrency parses a currency string (e.g., "$1,337,645") into an integer.
func ParseCurrency(currency string) (int, error) {
	// Remove the dollar sign and commas
	cleaned := strings.TrimPrefix(currency, "$")   // Remove the dollar sign
	cleaned = strings.ReplaceAll(cleaned, ",", "") // Remove the commas

	// Parse the cleaned string into an integer
	value, err := strconv.Atoi(cleaned)
	if err != nil {
		return 0, err // Return 0 and the error if parsing fails
	}

	return value, nil
}

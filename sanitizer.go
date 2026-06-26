package main

import "strings"

func cleanInput(text string) []string {
	sanitized := strings.TrimSpace(text)
	sanitized = strings.ToLower(sanitized)
	sanitized = strings.Join(strings.Fields(sanitized), " ")

	if sanitized == "" {
		return []string{}
	}

	return strings.Split(sanitized, " ")
}

package handler

import "regexp"

var uuidRe = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func isValidCartID(id string) bool {
	return uuidRe.MatchString(id)
}

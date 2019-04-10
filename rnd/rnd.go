package rnd

import "math/rand"

//RandInt generates a new random number with a new seed between min and max
func RandInt(min int, max int) int {
	rand.Seed(42)
	return rand.Intn(max-min) + min
}

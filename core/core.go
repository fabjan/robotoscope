// Package core contains code and types shared between all robotoscope packages.
package core

// RobotInfo says how many times a given robot was seen.
type RobotInfo struct {
	Seen      int
	UserAgent string
}

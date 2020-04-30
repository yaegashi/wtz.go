package wtz

import (
	"fmt"
	"time"
)

// LoadLocation is an alias of NameToLocation
func LoadLocation(name string) (*time.Location, error) {
	return NameToLocation(name)
}

// NameToLocation converts Windows zone name to Location
func NameToLocation(name string) (*time.Location, error) {
	if unixName, ok := WindowsToUnixMap[name]; ok {
		return time.LoadLocation(unixName)
	}
	return nil, fmt.Errorf("Unknown Windows time zone: %s", name)
}

// LocationToName converts Location to Windows zone name
func LocationToName(loc *time.Location) (string, error) {
	if name, ok := UnixToWindowsMap[loc.String()]; ok {
		return name, nil
	}
	_, offset := time.Now().In(loc).Zone() // XXX: affected by dst
	if name, ok := OffsetToWindowsMap[offset]; ok {
		return name, nil
	}
	return "", fmt.Errorf("No Windows time zone for: %s", loc)
}

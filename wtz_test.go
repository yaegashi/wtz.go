package wtz

import (
	"testing"
	"time"
)

func TestNameToLocation(t *testing.T) {
	testZoneList := []struct {
		Windows, Unix string
	}{
		{"UTC", "Etc/GMT"},
		{"Tokyo Standard Time", "Asia/Tokyo"},
		{"Dateline Standard Time", "Etc/GMT+12"},
	}
	for _, testZone := range testZoneList {
		t.Run(testZone.Windows, func(t *testing.T) {
			loc, err := NameToLocation(testZone.Windows)
			if err != nil {
				t.Error(err)
			}
			if loc.String() != testZone.Unix {
				t.Errorf("got %q want %q", loc.String(), testZone.Unix)
			}
		})
	}
}

func TestLocationToName(t *testing.T) {
	testZoneList := []struct {
		Unix, Windows string
	}{
		{"Etc/GMT", "UTC"},
		{"Asia/Macau", "China Standard Time"},     // by UnixToWindowsMap
		{"Asia/Macao", "Singapore Standard Time"}, // by OffsetToWindowsMap
		{"Asia/Tokyo", "Tokyo Standard Time"},     // by UnixToWindowsMap
		{"Japan", "Tokyo Standard Time"},          // by OffsetToWindowsMap
		{"Etc/GMT+12", "Dateline Standard Time"},
	}
	for _, testZone := range testZoneList {
		t.Run(testZone.Unix, func(t *testing.T) {
			loc, err := time.LoadLocation(testZone.Unix)
			if err != nil {
				t.Error(err)
			}
			name, err := LocationToName(loc)
			if err != nil {
				t.Error(err)
			}
			if name != testZone.Windows {
				t.Errorf("got %q want %q", name, testZone.Windows)
			}
		})
	}
}

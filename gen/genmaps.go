// +build ignore

package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"go/format"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"text/template"
	"time"
)

var filename = flag.String("output", "wtz_maps.go", "output file name")

type zone struct {
	WindowsName string
	UnixName    string
	Offset      int
}

const windowsZonesURL = "https://raw.githubusercontent.com/unicode-org/cldr/master/common/supplemental/windowsZones.xml"

type MapZone struct {
	Other     string `xml:"other,attr"`
	Territory string `xml:"territory,attr"`
	Type      string `xml:"type,attr"`
}

type SupplementalData struct {
	Zones []MapZone `xml:"windowsZones>mapTimezones>mapZone"`
}

func main() {
	flag.Parse()

	r, err := http.Get(windowsZonesURL)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var sd SupplementalData
	err = xml.Unmarshal(data, &sd)
	if err != nil {
		log.Fatal(err)
	}

	windowsZones := []*zone{}
	unixZones := []*zone{}
	offsetZones := []*zone{}
	offsetZoneMap := map[int]*zone{}

	for _, z := range sd.Zones {
		if z.Territory == "001" {
			windowsZones = append(windowsZones, &zone{WindowsName: z.Other, UnixName: z.Type})
		} else {
			for _, unixName := range strings.Split(z.Type, " ") {
				if unixName == "" {
					continue
				}
				unixZones = append(unixZones, &zone{WindowsName: z.Other, UnixName: unixName})
				loc, err := time.LoadLocation(unixName)
				if err != nil {
					continue
				}
				_, off := time.Now().In(loc).Zone()
				if oz, ok := offsetZoneMap[off]; !ok || !strings.HasPrefix(oz.UnixName, "Etc/GMT") {
					offsetZoneMap[off] = &zone{WindowsName: z.Other, UnixName: unixName, Offset: off}
				}
			}
		}
	}

	for _, zone := range offsetZoneMap {
		offsetZones = append(offsetZones, zone)
	}

	sort.Slice(windowsZones, func(i, j int) bool { return windowsZones[i].WindowsName < windowsZones[j].WindowsName })
	sort.Slice(unixZones, func(i, j int) bool { return unixZones[i].UnixName < unixZones[j].UnixName })
	sort.Slice(offsetZones, func(i, j int) bool { return offsetZones[i].Offset < offsetZones[j].Offset })

	var v = struct {
		URL          string
		WindowsZones []*zone
		UnixZones    []*zone
		OffsetZones  []*zone
	}{
		windowsZonesURL,
		windowsZones,
		unixZones,
		offsetZones,
	}

	var buf bytes.Buffer
	err = template.Must(template.New("prog").Parse(prog)).Execute(&buf, v)
	if err != nil {
		log.Fatal(err)
	}

	data, err = format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(*filename, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

const prog = `
// Code generated by genmaps.go; DO NOT EDIT.
// Based on information from {{.URL}}

package wtz

// WindowsToUnixMap is the time zone translation map from Windows to Unix.
var WindowsToUnixMap = map[string]string{
{{range .WindowsZones}}	"{{.WindowsName}}": "{{.UnixName}}",
{{end}}}

// UnixToWindowsMap is the time zone translation map from Unix to Windows.
var UnixToWindowsMap = map[string]string{
{{range .UnixZones}} "{{.UnixName}}": "{{.WindowsName}}",
{{end}}}

// OffsetToWindowsMap translates from time zone offset seconds to Windows zone name.
var OffsetToWindowsMap = map[int]string{
{{range .OffsetZones}} {{.Offset}}: "{{.WindowsName}}",
{{end}}}
`
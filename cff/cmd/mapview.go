package cmd

import (
	"context"
	"fmt"
	"math"

	mapscii "github.com/audergonv/go-mapscii"
	"github.com/audergonv/go-mapscii/canvas"
	"github.com/audergonv/go-mapscii/style"
	"github.com/audergonv/go-mapscii/tileprovider"
)

const tileURL = "https://vectortiles.geo.admin.ch/tiles/ch.swisstopo.base.vt/v1.0.0/{z}/{x}/{y}.pbf"

type mapPoint struct {
	lat, lon float64
	label    string
}

func collectStations(connection Connection) []mapPoint {
	var points []mapPoint
	seen := map[string]bool{}

	add := func(station Station) {
		if seen[station.ID] {
			return
		}
		seen[station.ID] = true
		points = append(points, mapPoint{lat: station.Coordinate.X, lon: station.Coordinate.Y, label: station.Name})
	}

	for i, section := range connection.Sections {
		if i == 0 {
			add(section.Departure.Station)
		}
		add(section.Arrival.Station)
	}

	return points
}

func renderTripMap(points []mapPoint, width, height int) (string, error) {
	if len(points) == 0 {
		return "", fmt.Errorf("aucune gare à afficher")
	}

	minLat, maxLat := points[0].lat, points[0].lat
	minLon, maxLon := points[0].lon, points[0].lon
	for _, p := range points[1:] {
		minLat = math.Min(minLat, p.lat)
		maxLat = math.Max(maxLat, p.lat)
		minLon = math.Min(minLon, p.lon)
		maxLon = math.Max(maxLon, p.lon)
	}

	centerLat := (minLat + maxLat) / 2
	centerLon := (minLon + maxLon) / 2
	spanKm := haversineKm(minLat, minLon, maxLat, maxLon)

	sty := style.Default()
	sty.RoadClasses["rail"] = canvas.Color{R: 100, G: 100, B: 100}
	sty.RoadWidths["rail"] = 1.5
	sty.RoadClasses["motorway"] = sty.Background
	sty.RoadWidths["motorway"] = 0
	sty.RoadClasses["trunk"] = sty.Background
	sty.RoadWidths["trunk"] = 0

	provider := tileprovider.NewHTTPProvider(tileURL)

	m, err := mapscii.New(mapscii.Options{
		Provider: provider,
		Width:    width,
		Height:   height,
		Center:   mapscii.LatLon{Lat: centerLat, Lon: centerLon},
		Zoom:     zoomForDistance(spanKm),
		Style:    sty,
	})
	if err != nil {
		return "", err
	}

	//Draw lines
	if len(points) >= 2 {
		line := make([]mapscii.LatLon, len(points))
		for i, p := range points {
			line[i] = mapscii.LatLon{Lat: p.lat, Lon: p.lon}
		}
		m.DrawLine(line, mapscii.WithLineColor(canvas.Color{R: 255, G: 0, B: 0}), mapscii.WithLineWidth(1))
	}

	for _, p := range points {
		m.AddPin(p.lat, p.lon, p.label)
	}

	return m.Render(context.Background())
}

//Distance between two gps points
func haversineKm(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

//Raw estimation
func zoomForDistance(distanceKm float64) float64 {
	switch {
	case distanceKm < 5:
		return 12
	case distanceKm < 10:
		return 11
	case distanceKm < 20:
		return 10
	case distanceKm < 35:
		return 9
	case distanceKm < 75:
		return 8
	case distanceKm < 100:
		return 7
	case distanceKm < 200:
		return 6.5
	case distanceKm < 400:
		return 6
	case distanceKm < 500:
		return 5.5
	case distanceKm < 600:
		return 4
	case distanceKm < 700:
		return 3
	default:
		return 2
	}
}
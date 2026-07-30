/*
Copyright © 2026 Aurélien Bulliard
*/
package cmd

import (
	"fmt"
	"time"
)

func printConnection(id int, connection Connection) {
	fmt.Printf("Trajet #%d\n", id)
	fmt.Println("--------------------------------------------------")
	for _, section := range connection.Sections {
		if section.Journey == nil {
			walkDuration := time.Duration(*section.Arrival.ArrivalTimestamp-*section.Departure.DepartureTimestamp) * time.Second
			fmt.Printf("\n%dmin Correspondance / Marche de %s à %s\n\n", int(walkDuration.Minutes()+0.5), section.Departure.Station.Name, section.Arrival.Station.Name)
			continue
		}
		var color string
		var category = section.Journey.Category
		if category == "TGV" || category == "EC" || category == "TER" || category == "RJX" {
			//Correct the category by excluding train number
			section.Journey.Number = ""
		}
		switch category {
		case "IC":
			color = "97;41"
		case "IR":
			color = "97;101"
		case "RE":
			color = "31;107"
		case "TGV":
			color = "3;97;41"
		case "EC":
			color = "3;97;41"
		case "B":
			color = "97;44"
		default:
			color = "30;107"
		}
		dt, err1 := time.Parse("2006-01-02T15:04:05-0700", section.Departure.Departure)
		at, err2 := time.Parse("2006-01-02T15:04:05-0700", section.Arrival.Arrival)
		var departureTime string
		var arrivalTime string
		if err1 != nil || err2 != nil {
			departureTime = section.Departure.Departure // fallback
			arrivalTime = section.Arrival.Arrival
		} else {
			departureTime = dt.Format("15:04")
			arrivalTime = at.Format("15:04")
		}
		if section.Departure.Platform != "" {
			fmt.Printf("%s * %-25s Voie %s", departureTime, section.Departure.Station.Name, section.Departure.Platform)
		} else {
			fmt.Printf("%s * %s", departureTime, section.Departure.Station.Name)
		}
		var depdelay = section.Departure.Delay
		if depdelay != nil && *depdelay != 0 {
			fmt.Printf("   \033[93m+%d min\033[0m\n", *depdelay)
		} else {
			fmt.Printf("\n")
		}
		fmt.Println("      |")
		fmt.Printf("      | \033[%sm%s%s\033[0m Direction %s\n", color, section.Journey.Category, section.Journey.Number, section.Arrival.Station.Name)
		fmt.Println("      |")
		if section.Arrival.Platform != "" {
			fmt.Printf("%s * %-25s Voie %s", arrivalTime, section.Arrival.Station.Name, section.Arrival.Platform)
		} else {
			fmt.Printf("%s * %s", arrivalTime, section.Arrival.Station.Name)
		}
		var arrdelay = section.Arrival.Delay
		if arrdelay != nil && *arrdelay != 0 {
			fmt.Printf("   \033[93m+%d min\033[0m\n", *arrdelay)
		} else {
			fmt.Printf("\n")
		}

	}
	fmt.Println("--------------------------------------------------")
	var d, h, m, s int
	fmt.Sscanf(connection.Duration, "%dd%d:%d:%d", &d, &h, &m, &s)
	fmt.Printf("Temps total: %s\n", formatDuration(d, h, m, s))
	fmt.Println("==================================================")
}
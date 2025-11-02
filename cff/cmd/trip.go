/*
Copyright © 2025 Aurélien Bulliard
*/
package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

func ftrip(cmd *cobra.Command, args []string) {
	url := fmt.Sprintf("http://transport.opendata.ch/v1/connections?from=%s&to=%s", args[0], args[1])
	//Query API
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	//Read response:
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Parse
	var conn Connections
	errJson := json.Unmarshal(body, &conn)
	if errJson != nil {
		log.Fatalln(errJson)
	}
	//Check if connexions found
	if len(conn.Connections) == 0 {
		fmt.Println("Aucune connexion trouvée pour les paramètres donnés...")
		return
	}
	fmt.Printf("Trajet de %s à %s:\n", args[0], args[1])
	for _, connection := range conn.Connections {
		fmt.Println("--------------------------------------------------")
		for _, section := range connection.Sections {
			if section.Journey == nil {
				walkDuration := time.Duration(*section.Arrival.ArrivalTimestamp-*section.Departure.DepartureTimestamp) * time.Second
				fmt.Printf("\n%dmin Correspondance / Marche de %s à %s\n\n", int(walkDuration.Minutes()+0.5), section.Departure.Station.Name, section.Arrival.Station.Name)
				continue
			}
			var color string
			switch section.Journey.Category {
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
				fmt.Printf("%s * %-25s Voie %s\n", departureTime, section.Departure.Station.Name, section.Departure.Platform)
			} else {
				fmt.Printf("%s * %s\n", departureTime, section.Departure.Station.Name)
			}
			fmt.Println("      |")
			fmt.Printf("      | \033[%sm%s%s\033[0m Direction %s\n", color, section.Journey.Category, section.Journey.Number, section.Arrival.Station.Name)
			fmt.Println("      |")
			if section.Arrival.Platform != "" {
				fmt.Printf("%s * %-25s Voie %s\n", arrivalTime, section.Arrival.Station.Name, section.Arrival.Platform)
			} else {
				fmt.Printf("%s * %s\n", arrivalTime, section.Arrival.Station.Name)
			}

		}
		fmt.Println("--------------------------------------------------")
		break
	}

}

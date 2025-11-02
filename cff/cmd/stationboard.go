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
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

func fstationboard(cmd *cobra.Command, args []string) {
	//Prepare URL depending on the flags
	nb, err := cmd.Flags().GetInt("nConnexions")
	if err != nil {
		log.Fatalln(err)
	}
	date, err := cmd.Flags().GetString("date")
	//encode
	encodedDate := url.QueryEscape(date)
	var url string
	if err != nil {
		log.Fatalln(err)
	} else if date != "" {
		url = fmt.Sprintf("http://transport.opendata.ch/v1/stationboard?station=%s&limit=%d&datetime=%s", args[0], nb, encodedDate)
	} else {
		url = fmt.Sprintf("http://transport.opendata.ch/v1/stationboard?station=%s&limit=%d", args[0], nb)
	}
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
	var sb Stationboard
	errJson := json.Unmarshal(body, &sb)
	if errJson != nil {
		log.Fatalln(errJson)
	}
	if date != "" {
		fmt.Printf("%d départs dès le %s depuis %s:\n", len(sb.Stationboard), date, sb.Stationboard[0].Stop.Station.Name)
	} else {
		fmt.Printf("%d prochains départs depuis %s:\n", len(sb.Stationboard), sb.Stationboard[0].Stop.Station.Name)
	}

	fmt.Println("---------------------------------------------")
	for _, entry := range sb.Stationboard {
		var color string
		switch entry.Category {
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
		default:
			color = "30;107"
		}
		fmt.Printf("\033[%sm%s%s\033[0m --> \033[1m%s\033[0m \n", color, entry.Category, entry.Number, entry.To)
		t, err := time.Parse("2006-01-02T15:04:05-0700", entry.Stop.Departure)
		if err != nil {
			fmt.Println("Départ:", entry.Stop.Departure) // fallback
		} else {
			fmt.Println("Départ:", t.Format("15:04"))
		}
		fmt.Println("Voie:", entry.Stop.Platform)
		fmt.Println("---------------------------------------------")
	}
}

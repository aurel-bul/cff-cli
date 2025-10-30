/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cff",
	Short: "Liste les prochains départs depuis la gare CFF en donnée.",
	Long:  `Un outil CLI qui permet de lister les départs de train de n'importe quelle gare Suisse.`,

	Run: func(cmd *cobra.Command, args []string) {
		//Query API
		url := fmt.Sprintf("http://transport.opendata.ch/v1/stationboard?station=%s&limit=5", args[0])
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

		fmt.Printf("5 prochains départs depuis %s:\n", sb.Stationboard[0].Stop.Station.Name)
		fmt.Println("---------------------------------------------")
		for _, entry := range sb.Stationboard {
			var color string
			switch entry.Category {
			case "IC":
				color = "97;41"
			case "IR":
				color = "97;101"
			case "S":
				color = "30;107"
			case "RE":
				color = "31;107"
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
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cff.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//Structs

type Stationboard struct {
	Stationboard []Entry `json:"stationboard"`
}

type Entry struct {
	Stop     Stop   `json:"stop"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Number   string `json:"number"`
	Operator string `json:"operator"`
	To       string `json:"to"`
}

type Stop struct {
	Station            Station   `json:"station"`
	Arrival            *string   `json:"arrival"`          // nullable
	ArrivalTimestamp   *int64    `json:"arrivalTimestamp"` // nullable
	Departure          string    `json:"departure"`
	DepartureTimestamp int64     `json:"departureTimestamp"`
	Platform           string    `json:"platform"`
	Prognosis          Prognosis `json:"prognosis"`
}

type Station struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Score      *float64   `json:"score"` // nullable
	Coordinate Coordinate `json:"coordinate"`
}

type Coordinate struct {
	Type string  `json:"type"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}

type Prognosis struct {
	Platform  *string `json:"platform"`  // nullable
	Arrival   *string `json:"arrival"`   // nullable
	Departure *string `json:"departure"` // nullable
	Capacity1 string  `json:"capacity1st"`
	Capacity2 string  `json:"capacity2nd"`
}

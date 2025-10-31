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
	"os"
	"time"

	"github.com/spf13/cobra"
)

//go:embed ascii_train.txt
var asciiArt string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cff <gare>",
	Short: "Liste les prochains départs depuis la gare CFF donnée.",
	Long:  asciiArt + "\nUn outil CLI qui permet de lister les départs de train de n'importe quelle gare Suisse.",
	Example: `# Afficher les 5 prochains départs depuis Fribourg
cff Fribourg

# Afficher les 10 prochains départs depuis Romont
cff Romont -n 10

# Afficher les départs à une date/heure donnée
cff Fribourg -d "2025-01-11 13:30"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
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
	rootCmd.Flags().StringP("date", "d", "", "date/heure, format: \"2025-10-31 17:30\"")
	rootCmd.Flags().IntP("nConnexions", "n", 5, "nombre de connexions à afficher")
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

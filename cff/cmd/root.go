/*
Copyright © 2025 Aurélien Bulliard
*/
package cmd

import (
	_ "embed"
	"os"

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
	Run:  fstationboard,
}

var trip = &cobra.Command{
	Use:   "trip <depuis> <vers>",
	Short: "Décrit le prochain trajet le plus court entre les deux gares données.",
	Long:  asciiArt + "\nDécrit le prochain trajet le plus court entre les deux gares données en paramètres.",
	Example: `# Afficher le prochain trajet entre Romont et Fribourg
cff trip Romont Fribourg

# Afficher le premier trajet pour aller de Fribourg à Bern le 03.11.2025 à 13h30
cff trip Fribourg Bern -d 2025-03-11 -t 13:30`,
	Args: cobra.ExactArgs(2),
	Run:  ftrip,
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
	rootCmd.Flags().StringP("date", "d", "", "date/heure, format: \"2025-10-31 17:30\"")
	rootCmd.Flags().IntP("nConnexions", "n", 5, "nombre de connexions à afficher")
	rootCmd.AddCommand(trip)
	trip.Flags().StringP("date", "d", "", "date du départ, format: \"2025-10-31\"")
	trip.Flags().StringP("heure", "t", "", "heure du départ, format: \"17:30\"")
}

//Structs

type Stationboard struct {
	Stationboard []Entry `json:"stationboard"`
}

type Connection struct {
	From        Stop      `json:"from"`
	To          Stop      `json:"to"`
	Duration    string    `json:"duration"`
	Transfers   int       `json:"transfers"`
	Service     *string   `json:"service"` // nullable
	Products    []string  `json:"products"`
	Capacity1st *int      `json:"capacity1st"`
	Capacity2nd *int      `json:"capacity2nd"`
	Sections    []Section `json:"sections"`
}

type Connections struct {
	Connections []Connection `json:"connections"`
	From        LocationInfo `json:"from"`
	To          LocationInfo `json:"to"`
	Stations    StationsInfo `json:"stations"`
}

type Entry struct {
	Stop     Stop   `json:"stop"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Number   string `json:"number"`
	Operator string `json:"operator"`
	To       string `json:"to"`
}

type Section struct {
	Journey   *Journey `json:"journey"` // nullable
	Walk      *any     `json:"walk"`    // nullable
	Departure Stop     `json:"departure"`
	Arrival   Stop     `json:"arrival"`
}

type Journey struct {
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	Subcategory  *string `json:"subcategory"`
	CategoryCode *string `json:"categoryCode"`
	Number       string  `json:"number"`
	Operator     string  `json:"operator"`
	To           string  `json:"to"`
	PassList     []Stop  `json:"passList"`
	Capacity1st  *int    `json:"capacity1st"`
	Capacity2nd  *int    `json:"capacity2nd"`
}

type Stop struct {
	Station              Station   `json:"station"`
	Arrival              string    `json:"arrival"`
	ArrivalTimestamp     *int64    `json:"arrivalTimestamp"`
	Departure            string    `json:"departure"`
	DepartureTimestamp   *int64    `json:"departureTimestamp"`
	Delay                *int      `json:"delay"`
	Platform             string    `json:"platform"`
	Prognosis            Prognosis `json:"prognosis"`
	RealtimeAvailability *any      `json:"realtimeAvailability"`
	Location             Station   `json:"location"`
}

type StationsInfo struct {
	From []LocationInfo `json:"from"`
	To   []LocationInfo `json:"to"`
}

type LocationInfo struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Score      *float64   `json:"score"`
	Coordinate Coordinate `json:"coordinate"`
	Distance   *float64   `json:"distance"`
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

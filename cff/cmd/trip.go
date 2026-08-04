/*
Copyright © 2026 Aurélien Bulliard
*/
package cmd

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)


func ftrip(cmd *cobra.Command, args []string) {
	date, err := cmd.Flags().GetString("date")
	heure, err2 := cmd.Flags().GetString("heure")
	nConnexions, err3 := cmd.Flags().GetInt("nConnexions")
	showMap, err4 := cmd.Flags().GetBool("carte")

	if err != nil || err2 != nil || err3 != nil || err4 != nil {
		log.Fatalln(err)
	}

	conn, err := fetchConnections(args[0], args[1], date, heure)
	if err != nil {
		log.Fatalln(err)
	}
	//Check if connexions found
	if len(conn.Connections) == 0 {
		fmt.Println("Aucune connexion trouvée pour les paramètres donnés...")
		return
	}

	//Limit to the number of connections actually shown, so IDs match what's cached
	shown := conn.Connections
	if len(shown) > nConnexions {
		shown = shown[:nConnexions]
	}

	if err := saveTripCache(args[0], args[1], shown); err != nil {
		fmt.Fprintln(os.Stderr, "Impossible de sauvegarder le cache:", err)
	}

	fmt.Printf("Prochain(s) trajet(s) de %s à %s:\n", args[0], args[1])
	for i, connection := range shown {
		printConnection(i+1, connection, showMap)
	}
}

func formatDuration(d, h, m, s int) string {
	result := ""
	if d > 0 {
		result += fmt.Sprintf("%dj ", d)
	}

	if h > 0 {
		result += fmt.Sprintf("%dh ", h)
	}
	if m > 0 {
		result += fmt.Sprintf("%dmin ", m)
	}
	if s > 0 && d == 0 {
		result += fmt.Sprintf("%ds ", s)
	}

	if result == "" {
		return "0min"
	}

	return strings.TrimSpace(result)
}

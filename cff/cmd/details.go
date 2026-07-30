/*
Copyright © 2026 Aurélien Bulliard
*/
package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

var details = &cobra.Command{
	Use:   "details <id> | <depuis> <vers>",
	Short: "Affiche le détail d'un trajet.",
	Long:  asciiArt + "\nAffiche le détail d'un trajet: soit celui d'un cff trip précédent via son ID, soit le prochain trajet direct entre deux gares.",
	Example: `# Détail du trajet 1 affiché par le dernier "cff trip"
cff details 1

# Détail du prochain trajet direct entre Romont et Fribourg
cff details Romont Fribourg

# Détail du prochain trajet direct à une heure donnée
cff details Fribourg Bern -t 13:30`,
	Args: cobra.RangeArgs(1, 2),
	Run:  fdetails,
}

func fdetails(cmd *cobra.Command, args []string) {
	if len(args) == 1 {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln("cff details attend soit un ID de trajet (cff details 1), soit deux gares (cff details Fribourg Bern)")
		}

		cache, err := loadTripCache()
		if err != nil {
			log.Fatalln("Aucun trajet en cache, lancer d'abord \"cff trip <depuis> <vers>\":", err)
		}

		if id < 1 || id > len(cache.Connections) {
			log.Fatalln("ID de trajet invalide, doit être entre 1 et", len(cache.Connections))
		}

		printConnection(id, cache.Connections[id-1])
		return
	}

	//read flags
	date, err := cmd.Flags().GetString("date")
	heure, err2 := cmd.Flags().GetString("heure")
	if err != nil || err2 != nil {
		log.Fatalln(err)
	}

	conn, err := fetchConnections(args[0], args[1], date, heure)
	if err != nil {
		log.Fatalln(err)
	}

	if len(conn.Connections) == 0 {
		fmt.Println("Aucune connexion trouvée pour les paramètres donnés...")
		return
	}

	printConnection(1, conn.Connections[0])

}
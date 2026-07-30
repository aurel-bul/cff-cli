/*
Copyright © 2026 Aurélien Bulliard
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


func fetchConnections(from, to, date, heure string) (Connections, error){
	var url string
	switch {
	case heure != "" && date != "":
		url = fmt.Sprintf("http://transport.opendata.ch/v1/connections?from=%s&to=%s&date=%s&time=%s", from, to, date, heure)
	case heure != "" && date == "":
		url = fmt.Sprintf("http://transport.opendata.ch/v1/connections?from=%s&to=%s&time=%s", from, to, heure)
	default:
		url = fmt.Sprintf("http://transport.opendata.ch/v1/connections?from=%s&to=%s", from, to)
	}

	//Query API
	resp, err := http.Get(url)
	if err != nil {
		return Connections{}, err
	}
	defer resp.Body.Close()
	//Read response:
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Connections{}, err
	}
	//Parse
	var conn Connections
	errJson := json.Unmarshal(body, &conn)
	if errJson != nil {
		return Connections{}, errJson
	}
	return conn, nil
}
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func getJSON() ([]byte, error) {
	resp, err := http.Get("http://stats.nba.com/js/data/ptsd/stats_ptsd.js")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buff := bytes.NewBuffer([]byte{})

	if _, err := io.Copy(buff, resp.Body); err != nil {
		return nil, err
	}

	if buff.Len() < 17 {
		return nil, fmt.Errorf("json response too short")
	}
	// strip javascript from result hopefully leaving only valid json.
	return buff.Bytes()[17 : buff.Len()-1], nil
}

func main() {
	start := time.Now()
	defer func() {
		now := time.Now()
		log.Printf("took %s", now.Sub(start))
	}()

	b, err := getJSON()

	if err != nil {
		log.Fatal(err)
	}

	nba := struct {
		Data struct {
			Players []interface{} `json:"players"`
		} `json:"data"`
	}{}

	if err := json.Unmarshal(b, &nba); err != nil {
		log.Fatal(err)
	}

	for _, player := range nba.Data.Players {
		row := player.([]interface{})
		if len(row) < 2 {
			continue
		}
		fmt.Printf("%s\n", row[1])
	}
}

package main

import (
	"encoding/json"
	"escape/algo"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func setCitiesFromJson(codeToCity *map[string]algo.City) {
	jsonFile, err := os.Open("cities.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, codeToCity)
}

func main() {
	// Load json to memory
	var codeToCity map[string]algo.City
	setCitiesFromJson(&codeToCity)

	// Map cities to continents
	continentToCities := make(map[string][]algo.City)
	for _, city := range codeToCity {
		continentToCities[city.Continent] = append(continentToCities[city.Continent], city)
	}

	// Take input from command line
	id := os.Args[1]
	fmt.Println("Running...")
	algo.PrintBestPath(codeToCity[id], &continentToCities)
}

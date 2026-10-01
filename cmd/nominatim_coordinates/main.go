package main

import (
	"flag"
	"fmt"
	"log"

	"darbelis.eu/taxi/internal/di"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, using default values")
	}

	nominatimClient := di.GetNominatimClient()

	address := flag.String("address", "Gedimino 1, Kaunas", "Real address in earth")
	help := flag.Bool("help", false, "Print help")
	flag.Parse()

	if *help {
		log.Fatal(`Parameters: -address "Gedimino 1, Kaunas"`)
		return
	}

	if *address == "" {
		log.Fatal("-address is required")
	}

	fmt.Println("Geocoding address:", *address)

	lat, lon, err := nominatimClient.Geocode(*address)

	if err != nil {
		log.Fatal("error from nominatim:", err.Error())
		return
	}

	fmt.Printf("Lat: %f ; Lon: %f\n", lat, lon)
}

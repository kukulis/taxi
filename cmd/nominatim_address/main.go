package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"darbelis.eu/taxi/internal/di"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, using default values")
	}

	nominatimClient := di.GetNominatimClient()

	lat := flag.String("lat", "55.0", "lat Coordinate in Earth")
	lon := flag.String("lon", "24.0", "lon Coordinate in Earth")
	help := flag.Bool("help", false, "Print help")
	flag.Parse()

	if *help {
		log.Fatal("Parameters: -lat 55.0  -lon 24.0")
		return
	}

	if *lat == "" {
		log.Fatal("-lat is required")
	}

	if *lon == "" {
		log.Fatal("-lon is required")
	}

	fLat, err := strconv.ParseFloat(*lat, 64)

	if err != nil {
		log.Fatal("Error parsing latitude", err.Error())
		return
	}
	fLon, err := strconv.ParseFloat(*lon, 64)

	if err != nil {
		log.Fatal("Error parsing longitued", err.Error())
		return
	}

	fmt.Printf("Geocoding coordinates: Lat: %f, Lon: %f\n", fLat, fLon)

	address, err := nominatimClient.ReverseGeocode(fLat, fLon)

	if err != nil {
		fmt.Println("error from nominatim:", err.Error())
		return
	}

	fmt.Printf("Address: %s \n", address)

}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var Version = "0.0.0"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	openWeatherMapAPIKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	if openWeatherMapAPIKey == "" {
		log.Fatal("OPENWEATHERMAP_API_KEY environment variable is not set")
	}

	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s", r.Method, r.URL)
		getWeather(w, r, openWeatherMapAPIKey)
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("Server is running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getWeather(w http.ResponseWriter, r *http.Request, apiKey string) {
	location := r.URL.Query().Get("location")
	if location == "" {
		http.Error(w, "Location parameter is required", http.StatusBadRequest)
		return
	}

	apiURL := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", location, apiKey)
	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Error fetching weather data", http.StatusInternalServerError)
		log.Printf("Error fetching weather data: %s", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error fetching weather data", http.StatusInternalServerError)
		log.Printf("Error fetching weather data. Status code: %d", resp.StatusCode)
		return
	}

	var weatherData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		http.Error(w, "Error decoding weather data", http.StatusInternalServerError)
		log.Printf("Error decoding weather data: %s", err)
		return
	}

	mainInfo, ok := weatherData["main"].(map[string]interface{})
	if !ok {
		http.Error(w, "Error parsing weather data", http.StatusInternalServerError)
		log.Printf("Error parsing weather data: main info not found")
		return
	}

	temperature, ok := mainInfo["temp"].(float64)
	if !ok {
		http.Error(w, "Error parsing weather data", http.StatusInternalServerError)
		log.Printf("Error parsing weather data: temperature not found")
		return
	}

	descriptionArray, ok := weatherData["weather"].([]interface{})
	if !ok || len(descriptionArray) == 0 {
		http.Error(w, "Error parsing weather data", http.StatusInternalServerError)
		log.Printf("Error parsing weather data: weather description not found")
		return
	}

	descriptionInfo, ok := descriptionArray[0].(map[string]interface{})
	if !ok {
		http.Error(w, "Error parsing weather data", http.StatusInternalServerError)
		log.Printf("Error parsing weather data: weather description not found")
		return
	}

	description, ok := descriptionInfo["description"].(string)
	if !ok {
		http.Error(w, "Error parsing weather data", http.StatusInternalServerError)
		log.Printf("Error parsing weather data: weather description not found")
		return
	}

	city, ok := weatherData["name"].(string)
	if !ok {
		http.Error(w, "Error parsing weather data", http.StatusInternalServerError)
		log.Printf("Error parsing weather data: city name not found")
		return
	}

	sysInfo, ok := weatherData["sys"].(map[string]interface{})
	if !ok {
		http.Error(w, "Error parsing weather data", http.StatusInternalServerError)
		log.Printf("Error parsing weather data: sys info not found")
		return
	}

	country, ok := sysInfo["country"].(string)
	if !ok {
		http.Error(w, "Error parsing weather data", http.StatusInternalServerError)
		log.Printf("Error parsing weather data: country code not found")
		return
	}

	caser := cases.Title(language.English)
	responseData := map[string]interface{}{
		"city":          city,
		"country":       country,
		"temperature":   temperature,
		"description":   caser.String(description),
		"webApiVersion": fmt.Sprintf("%s - cool feature", Version),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Error encoding weather data", http.StatusInternalServerError)
		log.Printf("Error encoding weather data: %s", err)
	}
}

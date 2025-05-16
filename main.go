package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"weather-dashboard/weather"

	"github.com/pkg/browser"
)

func main() {
	// Define the port for the server
	const port = "8080"

	// Serve static files (index.html, CSS, JS) from the "static" directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Redirect root (/) to index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/index.html", http.StatusFound)
	})

	// API endpoint to fetch current weather data
	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		// Get the city from query parameters
		city := r.URL.Query().Get("city")
		if city == "" {
			http.Error(w, "City name is required", http.StatusBadRequest)
			return
		}

		// Use the provided API key
		apiKey := "03b93e35ee7b70117f8c483b6fefc665"

		// Fetch weather data
		weatherData, err := weather.GetWeather(city, apiKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send the weather data as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(weatherData)
	})

	// API endpoint to fetch forecast data
	http.HandleFunc("/forecast", func(w http.ResponseWriter, r *http.Request) {
		// Get the city and datetime from query parameters
		city := r.URL.Query().Get("city")
		datetime := r.URL.Query().Get("datetime")
		if city == "" || datetime == "" {
			http.Error(w, "City name and datetime are required", http.StatusBadRequest)
			return
		}

		// Use the provided API key
		apiKey := "03b93e35ee7b70117f8c483b6fefc665"

		// Fetch forecast data (supports up to 10 days if API allows)
		forecastData, err := weather.GetForecast(city, datetime, apiKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send the forecast data as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(forecastData)
	})

	// Start the server
	serverAddr := fmt.Sprintf(":%s", port)
	fmt.Printf("Server running at http://localhost:%s\n", port)

	// Automatically open the browser
	go func() {
		// Give the server a moment to start
		time.Sleep(500 * time.Millisecond)
		err := browser.OpenURL(fmt.Sprintf("http://localhost:%s", port))
		if err != nil {
			fmt.Println("Failed to open browser:", err)
		}
	}()

	// Keep the server running
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}

package weather

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"
)

// GetWeather fetches current weather data for a given city
func GetWeather(city string, apiKey string) (*WeatherData, error) {
	// Encode the city name to handle spaces, commas, etc.
	encodedCity := url.QueryEscape(city)

	// Create the API URL for current weather
	apiURL := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", encodedCity, apiKey)

	// Send the request to OpenWeatherMap API
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the API: %v", err)
	}
	defer resp.Body.Close()

	// Handle non-200 status codes
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("city not found or API error (Status Code: %d)", resp.StatusCode)
	}

	// Decode the response into the WeatherResponse struct
	var weatherData WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weatherData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode weather data: %v", err)
	}

	// Convert the timestamp to local time using the timezone offset
	utcTime := time.Unix(int64(weatherData.Dt), 0)
	localTime := utcTime.Add(time.Second * time.Duration(weatherData.Timezone))

	// Set Sri Lanka's time zone offset (UTC+5:30)
	sriLankaOffset := time.Hour*5 + time.Minute*30
	sriLankaTime := utcTime.Add(sriLankaOffset)

	// Calculate the time difference (in hours) relative to Sri Lanka
	timeDifference := localTime.Sub(sriLankaTime)
	timeDifferenceHours := timeDifference.Hours()

	// Get the current time in the city's local timezone
	currentTime := time.Now().UTC().Add(time.Second * time.Duration(weatherData.Timezone))

	// Prepare the response data
	data := &WeatherData{
		Name:           weatherData.Name,
		Temperature:    weatherData.Main.Temp,
		Humidity:       weatherData.Main.Humidity,
		WindSpeed:      weatherData.Wind.Speed,
		Condition:      weatherData.Weather[0].Description,
		TimeDifference: timeDifferenceHours,
		CurrentTime:    currentTime.Format("2006-01-02 15:04:05"),
	}

	return data, nil
}

// GetForecast fetches forecast weather data for a given city and datetime (aiming for a 10-day range)
func GetForecast(city string, datetime string, apiKey string) (*ForecastData, error) {
	// Parse the requested datetime
	requestedTime, err := time.Parse("2006-01-02 15:04:05", datetime)
	if err != nil {
		return nil, fmt.Errorf("invalid datetime format: %v", err)
	}

	// Round the requested time to the nearest 3-hour interval (frontend provides 1-hour increments)
	roundedTime := roundToNearestHour(requestedTime, 3)

	// Encode the city name
	encodedCity := url.QueryEscape(city)

	// Create the API URL for forecast (using 5-day forecast API; for 10 days, a paid API may be required)
	apiURL := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?q=%s&appid=%s&units=metric", encodedCity, apiKey)

	// Send the request to OpenWeatherMap API
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the API: %v", err)
	}
	defer resp.Body.Close()

	// Handle non-200 status codes
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("city not found or API error (Status Code: %d)", resp.StatusCode)
	}

	// Decode the response into the ForecastResponse struct
	var forecastData ForecastResponse
	err = json.NewDecoder(resp.Body).Decode(&forecastData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode forecast data: %v", err)
	}

	// Find the forecast closest to the rounded requested time
	var closestForecast struct {
		Dt   int64 `json:"dt"`
		Main struct {
			Temp     float64 `json:"temp"`
			Humidity int     `json:"humidity"`
		} `json:"main"`
		Wind struct {
			Speed float64 `json:"speed"`
		} `json:"wind"`
		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`
	}
	minTimeDiff := math.MaxInt64
	requestedUnix := roundedTime.Unix()

	for _, item := range forecastData.List {
		timeDiff := int(math.Abs(float64(item.Dt - requestedUnix)))
		if timeDiff < minTimeDiff {
			minTimeDiff = timeDiff
			closestForecast = item
		}
	}

	if minTimeDiff == math.MaxInt64 {
		return nil, fmt.Errorf("no forecast data available for the requested time")
	}

	// Convert the forecast timestamp to local time
	utcTime := time.Unix(int64(closestForecast.Dt), 0)
	localTime := utcTime.Add(time.Second * time.Duration(forecastData.City.Timezone))

	// Set Sri Lanka's time zone offset (UTC+5:30)
	sriLankaOffset := time.Hour*5 + time.Minute*30
	sriLankaTime := utcTime.Add(sriLankaOffset)

	// Calculate the time difference (in hours) relative to Sri Lanka
	timeDifference := localTime.Sub(sriLankaTime)
	timeDifferenceHours := timeDifference.Hours()

	// Prepare the response data
	data := &ForecastData{
		Name:           forecastData.City.Name,
		Temperature:    closestForecast.Main.Temp,
		Humidity:       closestForecast.Main.Humidity,
		WindSpeed:      closestForecast.Wind.Speed,
		Condition:      closestForecast.Weather[0].Description,
		TimeDifference: timeDifferenceHours,
		ForecastTime:   localTime.Format("2006-01-02 15:04:05"),
	}

	return data, nil
}

// roundToNearestHour rounds a time to the nearest specified hour interval
func roundToNearestHour(t time.Time, interval int) time.Time {
	hour := t.Hour()
	minute := t.Minute()
	second := t.Second()

	// Calculate the nearest hour within the interval
	roundedHour := int(math.Round(float64(hour)/float64(interval)) * float64(interval))
	if roundedHour >= 24 {
		roundedHour -= 24
		t = t.Add(24 * time.Hour)
	}

	// Use minute and second to adjust rounding if needed
	if minute >= 30 || (minute == 29 && second >= 30) {
		roundedHour += interval
		if roundedHour >= 24 {
			roundedHour -= 24
			t = t.Add(24 * time.Hour)
		}
	}

	// Create a new time with the rounded hour, resetting minutes and seconds
	return time.Date(t.Year(), t.Month(), t.Day(), roundedHour, 0, 0, 0, t.Location())
}

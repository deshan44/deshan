package weather

// WeatherResponse is the structure for the OpenWeatherMap current weather API response
type WeatherResponse struct {
	Name string `json:"name"`
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
	Dt       int64 `json:"dt"`       // Unix timestamp of the data
	Timezone int   `json:"timezone"` // Timezone offset in seconds
}

// ForecastResponse is the structure for the OpenWeatherMap 5-day forecast API response
type ForecastResponse struct {
	City struct {
		Name     string `json:"name"`
		Timezone int    `json:"timezone"`
	} `json:"city"`
	List []struct {
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
	} `json:"list"`
}

// WeatherData is the structure for the JSON response to the frontend for current weather
type WeatherData struct {
	Name           string  `json:"name"`
	Temperature    float64 `json:"temperature"`
	Humidity       int     `json:"humidity"`
	WindSpeed      float64 `json:"windSpeed"`
	Condition      string  `json:"condition"`
	TimeDifference float64 `json:"timeDifference"`
	CurrentTime    string  `json:"currentTime"`
}

// ForecastData is the structure for the JSON response to the frontend for forecast weather
type ForecastData struct {
	Name           string  `json:"name"`
	Temperature    float64 `json:"temperature"`
	Humidity       int     `json:"humidity"`
	WindSpeed      float64 `json:"windSpeed"`
	Condition      string  `json:"condition"`
	TimeDifference float64 `json:"timeDifference"`
	ForecastTime   string  `json:"forecastTime"`
}

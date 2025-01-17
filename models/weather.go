// models/weather.go
package models

import (
	"encoding/json"
	"fmt"
	"net/http"
) // Function to retrieve weather data from OpenWeatherMap API
type WeatherData struct {
	Weather     string  `json:"weather"`
	Temperature float64 `json:"temperature"`
}

func GetWeather(latitude, longitude string) (WeatherData, error) {
	// Construct the API URL
	apiKey := "65bd6689ba60f41871774e40059c6129"
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s", latitude, longitude, apiKey)

	// Send HTTP GET request to the API
	resp, err := http.Get(url)
	if err != nil {
		return WeatherData{}, err
	}
	defer resp.Body.Close()

	// Check if the API response was successful
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		return WeatherData{}, err
	}

	// Decode the JSON response into a WeatherData struct
	var weatherData struct {
		Weather []struct {
			Main string `json:"main"`
		} `json:"weather"`
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return WeatherData{}, err
	}

	// Extract relevant weather information
	weather := WeatherData{
		Weather:     weatherData.Weather[0].Main,
		Temperature: weatherData.Main.Temp,
	}

	return weather, nil
}

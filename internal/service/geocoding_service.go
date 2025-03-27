package service

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type GeocodingService struct {
    ApiKey string
}

type LocationInfo struct {
    City    string
    Country string
}

func (s *GeocodingService) GetLocationInfo(lat, long float64) (*LocationInfo, error) {
    // Using OpenStreetMap Nominatim API (free, no API key required)
    url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f", lat, long)
    
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    // Required by Nominatim's terms of use
    req.Header.Set("User-Agent", "NostosTrips/1.0")
    
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result struct {
        Address struct {
            City    string `json:"city"`
            Country string `json:"country"`
        } `json:"address"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &LocationInfo{
        City:    result.Address.City,
        Country: result.Address.Country,
    }, nil
}
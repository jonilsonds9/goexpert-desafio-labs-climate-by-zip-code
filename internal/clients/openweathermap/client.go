package openweathermap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client interface {
	GetCoordinates(ctx context.Context, cityName, countryCode string) (*GeoLocation, int, error)
}

type client struct {
	httpClient *http.Client
	apiKey     string
}

func NewClient(httpClient *http.Client, apiKey string) Client {
	return &client{httpClient: httpClient, apiKey: apiKey}
}

type GeoLocation struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
	State   string  `json:"state"`
}

func (c *client) GetCoordinates(ctx context.Context, cityName, countryCode string) (*GeoLocation, int, error) {
	query := cityName
	if countryCode != "" {
		query = fmt.Sprintf("%s,%s", cityName, countryCode)
	}

	endpoint := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s",
		url.QueryEscape(query),
		url.QueryEscape(c.apiKey))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var locations []GeoLocation
	if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil {
		return nil, resp.StatusCode, err
	}

	if len(locations) == 0 {
		return nil, resp.StatusCode, fmt.Errorf("location not found")
	}

	return &locations[0], resp.StatusCode, nil
}

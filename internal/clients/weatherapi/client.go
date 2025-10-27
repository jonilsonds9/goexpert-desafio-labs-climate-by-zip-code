package weatherapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client interface {
	CurrentTempCByCoords(ctx context.Context, lat, lon float64) (float64, int, error)
}

type client struct {
	httpClient *http.Client
	apiKey     string
}

func NewClient(httpClient *http.Client, apiKey string) Client {
	return &client{httpClient: httpClient, apiKey: apiKey}
}

type currentResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func (c *client) CurrentTempCByCoords(ctx context.Context, lat, lon float64) (float64, int, error) {
	endpoint := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%f,%f",
		url.QueryEscape(c.apiKey), lat, lon)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var cr currentResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return 0, resp.StatusCode, err
	}

	return cr.Current.TempC, resp.StatusCode, nil
}

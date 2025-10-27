package viacep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jonilsonds9/goexpert-desafio-labs-climate-by-zip-code/internal/domain"
)

type Client interface {
	ConsultCEP(ctx context.Context, cep string) (*domain.ViaCEPAddress, int, error)
}

type client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) Client {
	return &client{httpClient: httpClient}
}

func (c *client) ConsultCEP(ctx context.Context, cep string) (*domain.ViaCEPAddress, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var addr domain.ViaCEPAddress
	if err := json.NewDecoder(resp.Body).Decode(&addr); err != nil {
		return nil, resp.StatusCode, err
	}

	return &addr, resp.StatusCode, nil
}

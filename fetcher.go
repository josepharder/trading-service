package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	requestMethod = "GET"
	requestHost   = "api.coinbase.com"
	requestPath   = "/api/v3/brokerage/orders/historical/batch"
)

func FetchOrders() ([]Order, error) {
	uri := fmt.Sprintf("%s %s%s", requestMethod, requestHost, requestPath)

	jwt, err := buildJWT(uri)
	if err != nil {
		return nil, fmt.Errorf("building JWT: %w", err)
	}

	url := fmt.Sprintf("https://%s%s?product_type=FUTURE&order_placement_source=RETAIL_ADVANCED&contract_expiry_type=UNKNOWN_CONTRACT_EXPIRY_TYPE&sort_by=LAST_FILL_TIME&use_simplified_total_value_calculation=false&order_status=FILLED", requestHost, requestPath)

	req, err := http.NewRequest(requestMethod, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+jwt)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var ordersResp OrdersResponse
	if err := json.Unmarshal(body, &ordersResp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return ordersResp.Orders, nil
}

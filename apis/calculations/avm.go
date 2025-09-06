package apicalculation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rentieranalytics/rentier-core/httpx"
)

type ApiCalculationConfig struct {
	Addr string
}

type ApiCalculation struct {
	client httpx.HTTPClient
	config ApiCalculationConfig
}

func NewApiCalculation(
	client httpx.HTTPClient,
	config ApiCalculationConfig,
) ApiCalculation {
	return ApiCalculation{
		client: client,
		config: config,
	}
}

func (c *ApiCalculation) AVM(
	ctx context.Context,
	data *AVMCalculationRequest,
) (AVMCalculationResponse, error) {
	v, e := json.Marshal(data)
	if e != nil {
		return AVMCalculationResponse{}, e
	}
	path := fmt.Sprintf("%s/v1.0/avm", c.config.Addr)
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		path,
		bytes.NewBuffer(v),
	)
	if err != nil {
		return AVMCalculationResponse{}, err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(request)
	if err != nil {
		return AVMCalculationResponse{}, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return AVMCalculationResponse{}, NewHTTPError(
			resp.StatusCode,
			fmt.Sprintf("external call failed: %s", body),
		)
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return AVMCalculationResponse{}, err
	}
	var avmR AVMCalculationResponse
	err = json.Unmarshal(b, &avmR)
	if err != nil {
		return AVMCalculationResponse{}, err
	}

	return avmR, nil
}

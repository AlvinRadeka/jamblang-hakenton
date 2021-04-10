package repository

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/alvinradeka/jamblang-hakenton/warehouse/internal/domain"
	"github.com/gojektech/heimdall/v6/httpclient"
	"github.com/sirupsen/logrus"
)

type barcodeRepository struct {
	logger     *logrus.Logger
	config     domain.BarcodeRepositoryConfig
	httpClient *httpclient.Client
}

func New(logger *logrus.Logger, cfg domain.BarcodeRepositoryConfig, httpClient *httpclient.Client) domain.BarcodeRepository {
	return &barcodeRepository{
		logger:     logger,
		config:     cfg,
		httpClient: httpClient,
	}
}

func (b *barcodeRepository) ParseToLambda(file64 string) (domain.BarcodeLambdaResponse, error) {
	var (
		lambdaResponse domain.BarcodeLambdaResponse
	)

	payload := map[string]interface{}{
		"img": file64,
	}

	jsonOut, err := json.Marshal(&payload)
	if err != nil {
		return lambdaResponse, err
	}

	resp, err := b.httpClient.Post(b.config.LambdaURL, bytes.NewBuffer(jsonOut), http.Header{
		"content-type": []string{"application/json"},
	})
	if err != nil {
		return lambdaResponse, err
	}

	// Read JSON
	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return lambdaResponse, err
	}

	if err := json.Unmarshal(output, &lambdaResponse); err != nil {
		return lambdaResponse, err
	}

	return lambdaResponse, nil
}

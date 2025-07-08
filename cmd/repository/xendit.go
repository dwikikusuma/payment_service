package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"payment_service/infra/log"
	"payment_service/models"
)

type XenditClient interface {
	CrateInvoice(ctx context.Context, param models.XenditInvoiceRequest) (models.XenditInvoiceResponse, error)
	CheckInvoiceStatus(ctx context.Context, externalID string) (string, error)
}

type xenditClient struct {
	APISecret string
}

func NewXenditClient() XenditClient {
	return &xenditClient{}
}

func (c *xenditClient) CrateInvoice(ctx context.Context, param models.XenditInvoiceRequest) (models.XenditInvoiceResponse, error) {
	var response models.XenditInvoiceResponse
	payload, err := json.Marshal(&param)
	if err != nil {
		return models.XenditInvoiceResponse{}, err
	}

	uri := "https://api.xendit.co/v2/invoices"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, bytes.NewBuffer(payload))
	if err != nil {
		return models.XenditInvoiceResponse{}, err
	}

	httpReq.SetBasicAuth(c.APISecret, "")
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return models.XenditInvoiceResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"message": "failed to close body content",
				"err":     err.Error(),
			}).Error("error occurred on xendit.CreateInvoice() defer func(Body io.ReadCloser)")
		}
	}(resp.Body)

	if resp.StatusCode > 300 {
		body, _ := io.ReadAll(resp.Body)
		return models.XenditInvoiceResponse{}, errors.New(fmt.Sprintf("xendit.CreateInvoice() got error %s", string(body)))
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return models.XenditInvoiceResponse{}, err
	}

	return response, nil
}

func (c *xenditClient) CheckInvoiceStatus(ctx context.Context, externalID string) (string, error) {
	uri := fmt.Sprintf("https://api.xendit.co/v2/invoices?external_id=%s", externalID)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return "", nil
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	var invoice []models.XenditInvoiceResponse
	err = json.NewDecoder(resp.Body).Decode(&invoice)
	if err != nil {
		return "", nil
	}

	return invoice[0].Status, nil
}

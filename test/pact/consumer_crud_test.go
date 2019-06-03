package pact_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"log"
	"net/http"
	"testing"
)

func TestConsumer(t *testing.T) {
	type Payment struct {
		Type           string   `json:"type"`
		ID             string   `json:"id"`
		Version        int16    `json:"version"`
		OrganisationID string   `json:"organisation_id"`
		Attributes     struct{} `json:"attributes"`
	}

	type PaymentCollection struct {
		Payments []Payment `json:"payments"`
		Meta     struct {
			More bool `json:"more"`
		} `json:"meta"`
	}

	pact := &dsl.Pact{
		Consumer: "Consumer",
		Provider: "Provider",
		Host:     "localhost",
	}
	defer pact.Teardown()

	payment := Payment{
		Type:           "Payment",
		ID:             "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec45",
		Version:        0,
		OrganisationID: "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec45",
	}

	updatePayload := Payment{
		Type:           "Payment",
		Version:        0,
		OrganisationID: "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec46",
	}

	updatedPayment := Payment{
		Type:           "Payment",
		ID:             "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec45",
		Version:        0,
		OrganisationID: "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec46",
	}

	var test = func() (err error) {
		// Create
		body, err := json.Marshal(payment)
		if err != nil {
			return err
		}
		u := fmt.Sprintf(fmt.Sprintf("http://localhost:%d/payments", pact.Server.Port))
		req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		if _, err = http.DefaultClient.Do(req); err != nil {
			return err
		}

		// Read All
		req, err = http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			return err
		}
		if _, err = http.DefaultClient.Do(req); err != nil {
			return err
		}

		// Update
		body, err = json.Marshal(updatePayload)
		if err != nil {
			return err
		}
		u = fmt.Sprintf(fmt.Sprintf("http://localhost:%d/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec45", pact.Server.Port))
		req, err = http.NewRequest(http.MethodPut, u, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		if _, err = http.DefaultClient.Do(req); err != nil {
			return err
		}

		// Read One
		req, err = http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			return err
		}
		if _, err = http.DefaultClient.Do(req); err != nil {
			return err
		}

		// Delete
		req, err = http.NewRequest(http.MethodDelete, u, nil)
		if err != nil {
			return err
		}
		if _, err = http.DefaultClient.Do(req); err != nil {
			return err
		}

		return nil
	}

	pact.
		AddInteraction().
		Given("system is empty").
		UponReceiving("a request to create a payment").
		WithRequest(dsl.Request{
			Method:  http.MethodPost,
			Path:    dsl.String("/payments"),
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
			Body:    payment,
		}).
		WillRespondWith(dsl.Response{
			Status:  http.StatusCreated,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json; charset=utf-8")},
			Body:    dsl.Match(payment),
		})

	pact.
		AddInteraction().
		Given("payment exists").
		UponReceiving("a request to fetch all payments").
		WithRequest(dsl.Request{
			Method: http.MethodGet,
			Path:   dsl.String("/payments"),
		}).
		WillRespondWith(dsl.Response{
			Status:  http.StatusOK,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json; charset=utf-8")},
			Body:    dsl.Match(PaymentCollection{Payments: []Payment{updatedPayment}}),
		})

	pact.
		AddInteraction().
		Given("payment exists").
		UponReceiving("a request to update a payment").
		WithRequest(dsl.Request{
			Method: http.MethodPut,
			Path:   dsl.String("/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec45"),
			Body:   updatePayload,
		}).
		WillRespondWith(dsl.Response{
			Status:  http.StatusOK,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json; charset=utf-8")},
			Body:    dsl.Match(updatedPayment),
		})

	pact.
		AddInteraction().
		Given("payment exists").
		UponReceiving("a request to fetch a payment").
		WithRequest(dsl.Request{
			Method: http.MethodGet,
			Path:   dsl.String("/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec45"),
		}).
		WillRespondWith(dsl.Response{
			Status:  http.StatusOK,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json; charset=utf-8")},
			Body:    dsl.Match(updatedPayment),
		})

	pact.
		AddInteraction().
		Given("payment exists").
		UponReceiving("a request to delete a payment").
		WithRequest(dsl.Request{
			Method: http.MethodDelete,
			Path:   dsl.String("/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec45"),
		}).
		WillRespondWith(dsl.Response{
			Status: http.StatusOK,
		})

	if err := pact.Verify(test); err != nil {
		log.Fatalf("Error on Verify: %v", err)
	}
}

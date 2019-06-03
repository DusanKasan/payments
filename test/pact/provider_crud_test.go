package pact_test

import (
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"path/filepath"
	"testing"
)

// TODO: This should be ran inside a container, connecting to a disposable instance of the server.
func TestProvider(t *testing.T) {
	pact := &dsl.Pact{
		Consumer: "MyConsumer",
		Provider: "MyProvider",
		Host: "localhost",
	}

	if _, err := pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL:        "http://localhost:8080",
		PactURLs:               []string{filepath.ToSlash(fmt.Sprintf("%s/consumer-provider.json", "./pacts"))},
		StateHandlers: types.StateHandlers{
			"system is empty": func() error {
				return nil
			},
			"payment exists": func() error {
				return nil
			},
		},
	}); err != nil {
		t.Error(err)
	}
}

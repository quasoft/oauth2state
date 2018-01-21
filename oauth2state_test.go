package oauth2state

import (
	"testing"
)

func TestCryptoValueGenerator_String(t *testing.T) {
	generator := CryptoValueGenerator{}
	got1 := generator.String()
	if len(got1) == 0 {
		t.Error(`CryptoValueGenerator.String() returned empty string`)
	}

	got2 := generator.String()
	if len(got2) == 0 {
		t.Error(`CryptoValueGenerator.String() returned empty string`)
	}

	if got1 == got2 {
		t.Error(`CryptoValueGenerator.String() returns the same string`)
	}
}

// Package oauth2state declares a simple storage interface for OAuth2 state values
// and provides a sample in-memory implementation for it.
package oauth2state

import (
	"crypto/rand"
	"encoding/base64"
)

// ValueGenerator generates random state values
type ValueGenerator interface {
	// Generate a random string suitable for use as state value
	String() string
}

// CryptoValueGenerator generates random state values using crypto/rand
type CryptoValueGenerator struct {
}

func (c CryptoValueGenerator) String() string {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

// TODO: Use some mechanism to expire state values

// StateStorer is the interface of a simple storage for state values exchanged with OAuth2
type StateStorer interface {
	// NewState creates a new random state value and associates the given URL with that value
	NewState(url string) (string, error)

	// Add a new state/url combination to the store
	Add(state, url string) error

	// Contains checks if the given state value exists in the store.
	Contains(state string) (bool, error)

	// URL retrieves the URL that is associated with a given state value
	URL(state string) (string, error)

	// Delete the given state value from the store
	Delete(state string) error
}

package ledgerapi

import (
	"strings"
)

// SplitKey - Splits a key on colon.
func SplitKey(key string) []string {
	return strings.Split(key, ":")
}

// MakeKey - Joins key parts using colon.
func MakeKey(keyParts ...string) string {
	return strings.Join(keyParts, ":")
}

// StateInterface - Interface states must implement for use in a list.
type StateInterface interface {
	// GetSplitKey return components that combine to form the key
	GetSplitKey() []string
	Serialize() ([]byte, error)
}

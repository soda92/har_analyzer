package har

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchStatus(t *testing.T) {
	assert.True(t, MatchStatus(200, "200"))
	assert.True(t, MatchStatus(204, "2xx"))
	assert.True(t, MatchStatus(502, "5xx"))
	assert.False(t, MatchStatus(200, "500"))
	assert.False(t, MatchStatus(404, "5xx"))
	assert.True(t, MatchStatus(200, ""))
}

func TestMatchURL(t *testing.T) {
	assert.True(t, MatchURL("https://example.com/api/users", "api"))
	assert.True(t, MatchURL("https://example.com/api/users", `api/v?[a-z]+`)) // Regex pattern
	assert.True(t, MatchURL("https://example.com/api/users", `users$`))
	assert.False(t, MatchURL("https://example.com/api/users", `^api`))
}

func TestMatchEntry(t *testing.T) {
	entry := Entry{
		"_resourceType": "fetch",
		"request": map[string]interface{}{
			"method": "GET",
			"url":    "https://example.com/api/v1/data",
		},
		"response": map[string]interface{}{
			"status": float64(200),
		},
	}

	// 1. Matches all criteria
	optsAllMatch := FilterOpts{
		Methods:    []string{"GET"},
		URLPattern: "v1",
		Status:     "200",
		Types:      []string{"fetch"},
	}
	assert.True(t, MatchEntry(entry, optsAllMatch))

	// 2. Fails method match
	optsFailMethod := FilterOpts{
		Methods: []string{"POST"},
	}
	assert.False(t, MatchEntry(entry, optsFailMethod))

	// 3. Fails status match
	optsFailStatus := FilterOpts{
		Status: "5xx",
	}
	assert.False(t, MatchEntry(entry, optsFailStatus))

	// 4. Fails type match
	optsFailType := FilterOpts{
		Types: []string{"xhr"},
	}
	assert.False(t, MatchEntry(entry, optsFailType))
}

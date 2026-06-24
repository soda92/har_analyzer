package har

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntryGetters(t *testing.T) {
	rawEntry := Entry{
		"time": float64(150.5),
		"_resourceType": "xhr",
		"request": map[string]interface{}{
			"method": "POST",
			"url":    "https://example.com/api",
			"headers": []interface{}{
				map[string]interface{}{"name": "X-Requested-With", "value": "XMLHttpRequest"},
				map[string]interface{}{"name": "Content-Type", "value": "application/json"},
			},
		},
		"response": map[string]interface{}{
			"status":     float64(200),
			"statusText": "OK",
		},
	}

	assert.Equal(t, 150.5, rawEntry.GetTime())
	assert.Equal(t, "xhr", rawEntry.GetResourceType())
	assert.Equal(t, "POST", rawEntry.GetMethod())
	assert.Equal(t, "https://example.com/api", rawEntry.GetURL())
	assert.Equal(t, 200, rawEntry.GetStatus())
	assert.Equal(t, "OK", rawEntry.GetStatusText())
	assert.True(t, rawEntry.HasXHRHeader())

	// Test fallback header check when ResourceType is empty
	rawEntryNoType := Entry{
		"request": map[string]interface{}{
			"headers": []interface{}{
				map[string]interface{}{"name": "X-Requested-With", "value": "XMLHttpRequest"},
			},
		},
	}
	assert.True(t, rawEntryNoType.HasXHRHeader())
}

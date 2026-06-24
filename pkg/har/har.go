package har

import (
	"strings"
)

// Entry wraps map[string]interface{} representing a single HAR request/response entry
type Entry map[string]interface{}

// HAR represents the structure of the HTTP Archive file
type HAR struct {
	Log Log `json:"log"`
}

// Log represents the log field inside HAR
type Log struct {
	Version string      `json:"version"`
	Creator interface{} `json:"creator,omitempty"`
	Browser interface{} `json:"browser,omitempty"`
	Pages   interface{} `json:"pages,omitempty"`
	Entries []Entry     `json:"entries"`
	Comment string      `json:"comment,omitempty"`
}

// GetRequest extracts request map from entry
func (e Entry) GetRequest() map[string]interface{} {
	if req, ok := e["request"].(map[string]interface{}); ok {
		return req
	}
	return nil
}

// GetResponse extracts response map from entry
func (e Entry) GetResponse() map[string]interface{} {
	if resp, ok := e["response"].(map[string]interface{}); ok {
		return resp
	}
	return nil
}

// GetResourceType returns resource type specified in HAR (usually Chrome/Firefox extension field)
func (e Entry) GetResourceType() string {
	if t, ok := e["_resourceType"].(string); ok {
		return t
	}
	return ""
}

// GetURL returns the requested URL
func (e Entry) GetURL() string {
	req := e.GetRequest()
	if req == nil {
		return ""
	}
	if url, ok := req["url"].(string); ok {
		return url
	}
	return ""
}

// GetMethod returns the HTTP method used in the request
func (e Entry) GetMethod() string {
	req := e.GetRequest()
	if req == nil {
		return ""
	}
	if method, ok := req["method"].(string); ok {
		return method
	}
	return ""
}

// GetStatus returns the HTTP response status code
func (e Entry) GetStatus() int {
	resp := e.GetResponse()
	if resp == nil {
		return 0
	}
	if statusVal, ok := resp["status"].(float64); ok {
		return int(statusVal)
	}
	return 0
}

// GetStatusText returns HTTP status text
func (e Entry) GetStatusText() string {
	resp := e.GetResponse()
	if resp == nil {
		return ""
	}
	if text, ok := resp["statusText"].(string); ok {
		return text
	}
	return ""
}

// GetTime returns elapsed time in ms
func (e Entry) GetTime() float64 {
	if tVal, ok := e["time"].(float64); ok {
		return tVal
	}
	return 0.0
}

// GetRequestHeaders extracts request headers as a slice of maps
func (e Entry) GetRequestHeaders() []map[string]interface{} {
	req := e.GetRequest()
	if req == nil {
		return nil
	}
	if headers, ok := req["headers"].([]interface{}); ok {
		res := make([]map[string]interface{}, 0, len(headers))
		for _, h := range headers {
			if hMap, ok := h.(map[string]interface{}); ok {
				res = append(res, hMap)
			}
		}
		return res
	}
	return nil
}

// HasXHRHeader returns true if X-Requested-With header is XMLHttpRequest
func (e Entry) HasXHRHeader() bool {
	headers := e.GetRequestHeaders()
	for _, h := range headers {
		name, _ := h["name"].(string)
		val, _ := h["value"].(string)
		if strings.ToLower(name) == "x-requested-with" && strings.ToLower(val) == "xmlhttprequest" {
			return true
		}
	}
	return false
}

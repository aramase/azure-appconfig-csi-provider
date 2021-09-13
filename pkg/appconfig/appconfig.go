package appconfig

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aramase/azure-appconfig-csi-provider/pkg/types"
)

const (
	endpointKey = "endpoint"
	idKey       = "id"
	secretKey   = "secret"

	apiVersion = "1.0"
)

type AppConfig interface {
	GetKV(key, label string) ([]types.KV, error)
}

// client is the client for the AppConfig service
type client struct {
	endpoint string
	id       string
	secret   string
}

// New creates a new AppConfig client
func New(connectionString string) AppConfig {
	parts := parseConnectionString(connectionString)

	return &client{
		endpoint: parts[endpointKey],
		id:       parts[idKey],
		secret:   parts[secretKey],
	}
}

// GetKV returns the value of the keys for the key prefix and label
func (c *client) GetKV(key, label string) ([]types.KV, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/kv?key=%s&label=%s&api-version=%s", c.endpoint, key, label, apiVersion), nil)
	if err != nil {
		return nil, err
	}
	err = SignRequest(req, c.id, c.secret)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var kvs types.KVItems
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, &kvs)
	if err != nil {
		return nil, err
	}
	return kvs.Items, nil
}

// parseConnectionString parses the connection string into a map of key/value pairs
func parseConnectionString(cs string) map[string]string {
	parts := make(map[string]string)

	// the connection string for the appconfig service is in the form:
	// Endpoint=<endpoint>;Id=<id>;Secret=<secret>
	for _, pair := range strings.Split(cs, ";") {
		if pair == "" {
			continue
		}
		equalDex := strings.IndexByte(pair, '=')
		if equalDex <= 0 {
			continue
		}

		key := strings.TrimSpace(strings.ToLower(pair[:equalDex]))
		value := strings.TrimSpace(pair[equalDex+1:])
		parts[key] = value
	}

	return parts
}

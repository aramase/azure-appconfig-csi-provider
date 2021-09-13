package types

// KV is a key/value pair with the label.
type KV struct {
	// Key refers to the key of the KV pair.
	// Supported filters: https://docs.microsoft.com/en-us/azure/azure-app-configuration/rest-api-key-value#supported-filters
	Key string `json:"key" yaml:"key"`
	// Label is the label of the KV pair.
	// Supported filters: https://docs.microsoft.com/en-us/azure/azure-app-configuration/rest-api-key-value#supported-filters
	Label string `json:"label" yaml:"value"`
	// Value is the value of the KV pair.
	Value string `json:"value"`
	// LastModified is the last time the KV pair was modified.
	LastModified string `json:"lastModified"`
	// ETag is the ETag of the KV pair.
	ETag string `json:"eTag"`
}

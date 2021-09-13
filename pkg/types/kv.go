package types

// KV is a key/value pair with the label.
type KV struct {
	// Key refers to the key of the KV pair.
	// Supported filters: https://docs.microsoft.com/en-us/azure/azure-app-configuration/rest-api-key-value#supported-filters
	Key string `json:"key" yaml:"key"`
	// Label is the label of the KV pair.
	// Supported filters: https://docs.microsoft.com/en-us/azure/azure-app-configuration/rest-api-key-value#supported-filters
	Label string `json:"label" yaml:"label"`
	// Value is the value of the KV pair.
	Value string `json:"value"`
	// LastModified is the last time the KV pair was modified.
	LastModified string `json:"last_modified"`
	// ETag is the ETag of the KV pair.
	ETag string `json:"eTag"`
}

// KVItems is a list of KV pairs.
type KVItems struct {
	Items []KV `json:"items"`
}

// StringArray is an array of strings.
type StringArray struct {
	Array []string `json:"array" yaml:"array"`
}

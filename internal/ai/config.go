package ai

// Config configures an AI client.
type Config struct {
	Provider Provider
	APIKey   string
	Model    string
}

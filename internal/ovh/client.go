package ovh

import (
	"fmt"

	"github.com/golgoth31/aliasme/internal/utils"
	"github.com/ovh/go-ovh/ovh"
	"github.com/rs/zerolog/log"
)

// Alias represents an email alias
type Alias struct {
	Source      string
	Destination string
}

// Config holds the OVH API configuration
type Config struct {
	Endpoint          string
	ApplicationKey    string
	ApplicationSecret string
	ConsumerKey       string
}

// Client wraps the OVH client
type Client struct {
	client *ovh.Client
}

// NewClient creates a new OVH client
func NewClient(endpoint, applicationKey, applicationSecret, consumerKey string) (*Client, error) {
	client, err := ovh.NewClient(
		endpoint,
		applicationKey,
		applicationSecret,
		consumerKey,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create OVH client")
		return nil, fmt.Errorf("failed to create OVH client: %w", err)
	}

	return &Client{client: client}, nil
}

// CreateAlias creates a new email alias
func (c *Client) CreateAlias(domain, prefix, destination string) (*Alias, error) {
	suffix := utils.GenerateRandomString(10)

	aliasFrom := prefix + "." + suffix + "@" + domain

	newAlias := &aliasesPost{
		From:      aliasFrom,
		LocalCopy: false,
		To:        "david.sabatie@notrenet.com",
	}

	err := c.client.Post("/email/domain/aliasme.ovh/redirection", newAlias, nil)

	return &Alias{
		Source:      aliasFrom,
		Destination: destination,
	}, err
}

// CreateEmailAlias creates a new email alias
func (c *Client) CreateEmailAlias(domain, alias, target string) error {
	// Implementation will depend on the specific OVH API endpoints
	// This is a placeholder for the actual implementation
	return nil
}

// DeleteEmailAlias deletes an email alias
func (c *Client) DeleteEmailAlias(domain, alias string) error {
	// Implementation will depend on the specific OVH API endpoints
	// This is a placeholder for the actual implementation
	return nil
}

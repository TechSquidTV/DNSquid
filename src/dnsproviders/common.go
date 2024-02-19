package dnsproviders

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/go-resty/resty/v2"
	"time"
)

// All DNS providers must implement this interface
type DNSProvider interface {
	Name() string
	PromptCredentials() error                     // Prompt user for credentials and set in keychain
	LoadCredentials() error                       // Load credentials from keyring
	SaveCredentials() error                       // Save the current credentials to the keychain
	New(options interface{}) (DNSProvider, error) // Create a new instance of the provider
	GetDomains() ([]string, error)                // Get a list of domains
}

type DNSProviderContext struct {
	RegisteredProviders map[string]DNSProvider
	AllProviders        map[string]DNSProvider
}

func NewDNSProviderContext() *DNSProviderContext {
	return &DNSProviderContext{
		RegisteredProviders: make(map[string]DNSProvider),
		AllProviders:        make(map[string]DNSProvider),
	}
}

// Initialize the provider context with all available providers
func (ctx *DNSProviderContext) Initialize() {
	// Init a provider
	p := &Porkbun{}

	// Add the provider to the map of all providers
	ctx.AllProviders[p.Name()] = p
}

// Register a provider by name
// This will create a new instance of the provider and store it in the map of registered providers
func (ctx *DNSProviderContext) RegisterProvider(providerName string) {

	// Check if the provider is already registered
	if _, exists := ctx.RegisteredProviders[providerName]; exists {
		log.Warnf("Provider %s is already registered", providerName)
		return
	}

	// Check if the provider is available
	provider, exists := ctx.AllProviders[providerName]
	if !exists {
		log.Fatalf("Unable to find provider %s", providerName)
	}

	// Register the provider
	ctx.RegisteredProviders[providerName] = provider
}

// Check to see if a provider exists
func (ctx *DNSProviderContext) ProviderExists(name string) bool {
	_, exists := ctx.RegisteredProviders[name]
	return exists
}

// Get a registered provider by name
func (ctx *DNSProviderContext) GetProvider(name string) (DNSProvider, error) {

	provider, exists := ctx.RegisteredProviders[name]
	if !exists {
		return nil, fmt.Errorf("provider %s is not registered", name)
	}
	return provider, nil
}

// Get a list of all available providers
func (ctx *DNSProviderContext) ListProviders() []string {

	var providerNames []string
	for name := range ctx.AllProviders {
		providerNames = append(providerNames, name)
	}
	return providerNames
}

// Return a preconfigured http client with retry logic
func NewClient() *resty.Client {
	return resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second).
		AddRetryCondition(
			func(response *resty.Response, err error) bool {
				statusCode := response.StatusCode()
				// Retry on specific status codes
				return statusCode == 503 || statusCode == 504 || statusCode == 429 || statusCode == 500 || statusCode == 502
			},
		)
}

package dnsproviders

import (
	"fmt"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/zalando/go-keyring"
)

type Porkbun struct {
	apiKey    string
	apiSecret string
}

// ListDomainsResponse represents the top-level structure of the JSON response.
type ListDomainsResponse struct {
	Status  string   `json:"status"`
	Domains []Domain `json:"domains"`
}

// Domain represents each domain's details in the response.
type Domain struct {
	Domain       string    `json:"domain"`
	Status       string    `json:"status"` // Assuming status could be a string, including empty/null.
	TLD          string    `json:"tld"`
	CreateDate   time.Time `json:"createDate"` // Using time.Time for date fields, assuming RFC3339 formatting.
	ExpireDate   time.Time `json:"expireDate"`
	SecurityLock string    `json:"securityLock"`
	WhoisPrivacy string    `json:"whoisPrivacy"`
	AutoRenew    string    `json:"autoRenew"`
	NotLocal     int       `json:"notLocal"` // Assuming int based on the provided values.
}

func (p *Porkbun) fetchDomains() ([]Domain, error) {
	client := NewClient()
	postData := fmt.Sprintf(`{"secretapikey":"%s","apikey":"%s"}`, p.apiSecret, p.apiKey)
	var response ListDomainsResponse

	_, err := client.R().
		SetBody(postData).
		SetResult(&response).
		Post("https://porkbun.com/api/json/v3/domain/listAll")

	if err != nil {
		log.Warnf("Unable to get domains: %s", err)
		return nil, err
	}

	return response.Domains, nil
}

func (p *Porkbun) GetDomains() ([]string, error) {
	domains, err := p.fetchDomains()
	if err != nil {
		return nil, err
	}

	var domainNames []string
	for _, domain := range domains {
		domainNames = append(domainNames, domain.Domain)
	}

	return domainNames, nil
}

func (p *Porkbun) Name() string {
	return "porkbun"
}

func (p *Porkbun) LoadCredentials() error {
	var err error
	p.apiKey, err = keyring.Get(p.Name(), "apikey")
	if err != nil {
		log.Warnf("Unable to load credentials for %s: %s", p.Name(), err)
		return err
	}
	p.apiSecret, err = keyring.Get(p.Name(), "secretapikey")
	if err != nil {
		log.Warnf("Unable to load credentials for %s: %s", p.Name(), err)
		return err
	}
	return nil
}

func (p *Porkbun) New(options interface{}) (DNSProvider, error) {
	return &Porkbun{}, nil
}

// Prompt the user to enter their credentials to be stored in the keychain.
func (p *Porkbun) PromptCredentials() error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("API Key").
				Placeholder("<apikey>").
				Value(&p.apiKey),
			huh.NewInput().
				Title("Secret API Key").
				Placeholder("<secretapikey>").
				Value(&p.apiSecret),
		),
	)

	err := form.Run()
	if err != nil {
		log.Warnf("Unable to collect credentials for %s: %s", p.Name(), err)
		return err
	}
	return nil
}

// Save the user's credentials to the keychain.
func (p *Porkbun) SaveCredentials() error {
	err := keyring.Set(p.Name(), "apikey", p.apiKey)
	if err != nil {
		log.Warnf("Unable to save apikey for %s: %s", p.Name(), err)
		return err
	}
	err = keyring.Set(p.Name(), "secretapikey", p.apiSecret)
	if err != nil {
		log.Warnf("Unable to save secretapikey for %s: %s", p.Name(), err)
		return err
	}
	return nil
}

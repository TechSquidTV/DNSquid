package config

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

type Providers struct {
	Porkbun bool `ini:"porkbun"`
}

type Config struct {
	Providers Providers `ini:"providers"`
}

func EnableProvider(providerName string) {
	// TODO: Ensure the provider exists and matches the available names
	viper.GetViper().Set("providers."+providerName, true)
	SaveConfig()
}

func SaveConfig() {
	err := viper.GetViper().WriteConfig()
	if err != nil {
		log.Warnf("Unable to save config: %s", err)
	}
}

// Returns a list of enabled providers from the currently loaded config
func GetEnabledProviders() []string {
	providers := viper.GetStringMap("providers")
	enabledProviders := make([]string, 0)
	for name, enabled := range providers {
		if enabled.(bool) {
			enabledProviders = append(enabledProviders, name)
		}
	}
	return enabledProviders
}

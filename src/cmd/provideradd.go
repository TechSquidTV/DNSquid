/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/techsquidtv/dnsquid/config"
	"github.com/techsquidtv/dnsquid/dnsproviders"
)

// provideraddCmd represents the provideradd command
var provideraddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new provider",
	Long:  `Configure a new provider to be used with dnsquid. This will prompt you for the necessary credentials to authenticate with the provider.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context().Value(dnsProviderContextKey).(*dnsproviders.DNSProviderContext)
		providerOptions := make([]huh.Option[string], 0, len(ctx.ListProviders()))
		providerName := ""
		for _, provider := range ctx.ListProviders() {
			providerOptions = append(providerOptions, huh.NewOption(provider, provider))
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select a provider").
					Options(
						providerOptions...,
					).Value(&providerName),
			),
		)

		err := form.Run()
		if err != nil {
			log.Fatalf("Unable to open provider selection form: %s", err)
		}

		ctx.RegisterProvider(providerName)

		provider, err := ctx.GetProvider(providerName)
		if err != nil {
			log.Fatalf("Unable to get provider: %s", err)
		}
		err = provider.PromptCredentials()
		if err != nil {
			log.Fatalf("Unable to collect credentials: %s", err)
		}
		err = provider.SaveCredentials()
		if err != nil {
			log.Fatalf("Unable to save credentials: %s", err)
		}
		config.EnableProvider(providerName)
		log.Infof("Enabled provider: %s", providerName)

	},
}

func init() {
	providerCmd.AddCommand(provideraddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// provideraddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// provideraddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

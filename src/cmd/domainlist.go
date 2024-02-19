package cmd

import (
	"github.com/spf13/cobra"
	"github.com/charmbracelet/log"
	"github.com/techsquidtv/dnsquid/dnsproviders"
)

// domainlistCmd represents the domainlist command
var domainlistCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context().Value(dnsProviderContextKey).(*dnsproviders.DNSProviderContext)
		fullDomainList := []string{}
		for _, provider := range ctx.RegisteredProviders {
			log.Debug("Listing domains for provider: %s", provider.Name())
			domains, err := provider.GetDomains()
			if err != nil {
				log.Warnf("Unable to get domains for provider: %s", err)
				return
			}
			fullDomainList = append(fullDomainList, domains...)
		}
		log.Infof("Domains: %s", fullDomainList)
	},
}

func init() {
	domainCmd.AddCommand(domainlistCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// domainlistCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// domainlistCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

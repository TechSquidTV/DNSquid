package cmd

import (
	"context"
	"os"
	"path"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/techsquidtv/dnsquid/config"
	"github.com/techsquidtv/dnsquid/dnsproviders"
	"gopkg.in/ini.v1"
)

// Define a new context key to be used for the DNSProviderContext
type contextKey struct{}

var dnsProviderContextKey = contextKey{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dnsquid",
	Short: "Manage all of your domains across multiple DNS providers in one common, local, and secure place.",
	Long:  `Manage all of your domains across multiple DNS providers in one common, local, and secure place.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ctx := dnsproviders.NewDNSProviderContext()
		ctx.Initialize()
		newCtx := context.WithValue(cmd.Context(), dnsProviderContextKey, ctx)
		cmd.SetContext(newCtx)

		// Register configured providers
		for _, providerName := range viper.GetStringSlice("providers") {
			ctx.RegisterProvider(providerName)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	var configPath string
	if configPath = os.Getenv("DNSQUID_CONFIG"); configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Unable to determine home directory: %s", err)
		}
		configPath = path.Join(homeDir, ".dnsquid.ini")
	}
	viper.SetConfigFile(configPath)

	//Initialize the config file if it does not exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Warnf("Config file does not exist: %s", configPath)
		_, err := os.Create(configPath)
		if err != nil {
			log.Fatalf("Unable to create config file: %s", err)
		}
		cfg := ini.Empty()
		cfg.ReflectFrom(&config.Config{})
		cfg.SaveTo(configPath)
		log.Infof("Created config file: %s", configPath)
	}

	// Load the config file
	var cfg config.Config
	err := viper.ReadInConfig()
	if err != nil {
		log.Warnf("Unable to load config file: %s", err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Warnf("Unable to parse config file: %s", err)
	}
	log.Debugf("Config: %+v", cfg)
}

func init() {
	log.SetLevel(log.DebugLevel)// TODO: Remove this line. Used for early local testing
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dnsquid.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cobra.OnInitialize(initConfig)
}

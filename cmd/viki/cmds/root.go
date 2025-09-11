package cmds

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "viki",
	Short: "Viki is a tool for turning Obsidian vaults into websites.",
}

var configFilePath string

func init() {
	cobra.OnInitialize(initConfig)

	rootFlags := RootCmd.PersistentFlags()

	rootFlags.StringSliceP("include-patterns", "i", []string{""}, "Gitignore-style glob patterns for files to include in the conversion. If included, only files matching these patterns are included.")

	viper.BindPFlags(rootFlags)

	// Special case for config file flag
	RootCmd.PersistentFlags().StringVar(&configFilePath, "config", "", "Path to config file (optional)")
}

func initConfig() {
	if configFilePath != "" {
		viper.SetConfigFile(configFilePath)
	}

	viper.SetEnvPrefix("VIKI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Ignore errors here because we don't necessarily need a config file
	_ = viper.ReadInConfig()

	err := viper.Unmarshal(&config)

	if err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
}

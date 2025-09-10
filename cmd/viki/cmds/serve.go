package cmds

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve [directory]",
	Short: "Serve an Obsidian vault as a website locally.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]
		fmt.Println("Serving directory:", dir)
		fmt.Println("Host:", config.Serve.Host)
		fmt.Println("Port:", config.Serve.Port)
	},
}

func init() {
	serveFlags := serveCmd.Flags()

	serveFlags.IntP("serve.port", "p", 8080, "Port to serve the website on")
	serveFlags.String("serve.host", "localhost", "Host to serve the website on")

	viper.BindPFlags(serveFlags)

	RootCmd.AddCommand(serveCmd)
}

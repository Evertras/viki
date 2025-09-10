package cmds

import (
	"fmt"
	"log"
	"net/http"

	"github.com/evertras/viki/lib/viki"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve an Obsidian vault in the current directory as a website locally.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		baseAddress := fmt.Sprintf("%s:%d", config.Serve.Host, config.Serve.Port)

		log.Printf("Listening on http://%s", baseAddress)

		inputFs := afero.NewOsFs()
		outputFs := afero.NewMemMapFs()

		converter := viki.NewConverter(viki.ConverterOptions{})

		err := converter.Convert(inputFs, ".", outputFs, "/")

		if err != nil {
			log.Println("Error during conversion:", err)
		}

		httpFs := afero.NewHttpFs(outputFs)
		http.Handle("/", http.FileServer(httpFs.Dir("/")))
		err = http.ListenAndServe(fmt.Sprintf("%s:%d", config.Serve.Host, config.Serve.Port), nil)

		if err != nil {
			log.Fatalf("Failed to serve website: %v", err)
		}
	},
}

func init() {
	serveFlags := serveCmd.Flags()

	serveFlags.IntP("serve.port", "p", 8080, "Port to serve the website on")
	serveFlags.String("serve.host", "localhost", "Host to serve the website on")

	viper.BindPFlags(serveFlags)

	RootCmd.AddCommand(serveCmd)
}

package cmds

import (
	"log"

	"github.com/evertras/viki/lib/viki"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate <src> <dst>",
	Short: "Generate a static website from an Obsidian vault.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		src := args[0]
		dst := args[1]

		log.Printf("Generating website from %s to %s", src, dst)

		inputFs := afero.NewBasePathFs(afero.NewOsFs(), src)
		outputFs := afero.NewBasePathFs(afero.NewOsFs(), dst)

		converter := viki.NewConverter(generateVikiConfig())

		err := converter.Convert(inputFs, outputFs)

		if err != nil {
			log.Fatalln("Error during conversion:", err)
		}

		log.Println("Generation complete")
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)
}

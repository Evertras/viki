package cmds

import "github.com/evertras/viki/lib/viki"

type Config struct {
	Serve struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"serve"`

	IncludePatterns []string `mapstructure:"include-patterns"`
}

var config Config

func generateVikiConfig() viki.ConverterOptions {
	return viki.ConverterOptions{
		IncludePatterns: config.IncludePatterns,
	}
}

package cmds

type Config struct {
	Serve struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"serve"`

	IncludePatterns []string `mapstructure:"include-patterns"`
}

var config Config

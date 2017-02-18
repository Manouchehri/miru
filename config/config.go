package config

// The name of the configuration file to look for.
const config_filename string = "config.json"

// Config contains global configuration information for the entire application.
// It is very likely that most request handler implementations will want to
// have a reference to a copy of this.
type Config struct {
  BindAddress string
  TemplateDir string
}

// MustLoad tries to load a configuration and panics if it cannot do so.
// A `CONFIG_DIR` environment variable can be set to specify the directory
// to read `config_filename` from.
func MustLoad() Config {
  return Config {
    BindAddress: "127.0.0.1:3000",
    TemplateDir: "templates",
  }
}

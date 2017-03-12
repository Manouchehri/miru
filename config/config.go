package config

import (
	"encoding/json"
	"os"
	"path"
)

// The name of the configuration file to look for.
const configFilename string = "config.json"

// Config contains global configuration information for the entire application.
// It is very likely that most request handler implementations will want to
// have a reference to a copy of this.
type Config struct {
	BindAddress string `json:"bindAddress"`      // The address and port to bind the server to.
	TemplateDir string `json:"templateDir"`      // The directory containing HTML page templates.
	Database    string `json:"database"`         // The connection string for the database.
	ScriptDir   string `json:"scriptDir"`        // The directory to save monitor scripts to.
	MGDomain    string `json:"mailgunDomain"`    // The domain registered to Mailgun to send emails through.
	MGAPIKey    string `json:"mailgunAPIKey"`    // Your Mailgun API key.
	MGPublicKey string `json:"mailgunPublicKey"` // Your Mailgun public key.
}

// MustLoad tries to load a configuration and panics if it cannot do so.
// A `CONFIG_DIR` environment variable can be set to specify the directory
// to read `configFilename` from.
func MustLoad() Config {
	c := Config{}
	configDir := os.Getenv("CONFIG_DIR")
	if configDir == "" {
		configDir = "config"
	}
	f, err := os.Open(path.Join(configDir, configFilename))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	decodeErr := decoder.Decode(&c)
	if decodeErr != nil {
		panic(decodeErr)
	}
	return c
}

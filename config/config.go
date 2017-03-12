package config

// The name of the configuration file to look for.
const configFilename string = "config.json"

// Config contains global configuration information for the entire application.
// It is very likely that most request handler implementations will want to
// have a reference to a copy of this.
type Config struct {
	BindAddress string // The address and port to bind the server to.
	TemplateDir string // The directory containing HTML page templates.
	Database    string // The connection string for the database.
	ScriptDir   string // The directory to save monitor scripts to.
	MGDomain    string // The domain registered to Mailgun to send emails through.
	MGAPIKey    string // Your Mailgun API key.
	MGPublicKey string // Your Mailgun public key.
}

// MustLoad tries to load a configuration and panics if it cannot do so.
// A `CONFIG_DIR` environment variable can be set to specify the directory
// to read `configFilename` from.
func MustLoad() Config {
	return Config{
		BindAddress: "127.0.0.1:3000",
		TemplateDir: "templates",
		Database:    "./miru.db",
		ScriptDir:   "monitorscripts",
		MGDomain:    "",
		MGAPIKey:    "",
		MGPublicKey: "",
	}
}

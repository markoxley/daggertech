package daggertech

// Config is the database configuration
type Config struct {
	host      string
	name      string
	user      string
	password  string
	deletable bool
}

// CreateConfig prepares the configuration for dagger
func CreateConfig(host string, name, user, password string, allowDelete bool) *Config {
	return &Config{
		host:      host,
		name:      name,
		user:      user,
		password:  password,
		deletable: allowDelete,
	}
}

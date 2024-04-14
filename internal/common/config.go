package common

type Database struct {
	Name      string `yaml:"name"`
	Url       string `yaml:"url"`
	AuthToken string `yaml:"auth_token"`
}

type CommonConfig struct {
	TeamName string
	Password string
	Database Database `yaml:"database"`
}

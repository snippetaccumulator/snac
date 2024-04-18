package common

type Database struct {
	Name      string `yaml:"name" json:"name"`
	Url       string `yaml:"url" json:"url"`
	AuthToken string `yaml:"auth_token" json:"auth_token"`
}

type CommonConfig struct {
	TeamName string   `yaml:"team_name" json:"team_name"`
	Password string   `yaml:"password" json:"password"`
	Database Database `yaml:"database" json:"database"`
}

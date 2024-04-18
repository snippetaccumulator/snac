package cli

import "github.com/snippetaccumulator/snac/internal/common"

type Config struct {
	common.CommonConfig
	LogLevel string `yaml:"log_level" json:"log_level"`
}

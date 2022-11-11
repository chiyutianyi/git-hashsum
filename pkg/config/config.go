package config

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

//Config is the configuration
type Config struct {
	//GitDir returns git dir
	GitDir string

	//LogLevel returns LogLevel
	LogLevel string
}

var Cfg Config

// GetLogLevel returns log level
func (c Config) GetLogLevel() (logLevel log.Level) {
	logLevel, err := log.ParseLevel(strings.ToLower(c.LogLevel))
	if err != nil {
		return log.InfoLevel
	}
	return logLevel
}

// GetGitDir return git dir
func (c Config) GetGitDir() string {
	if c.GitDir != "" && !filepath.IsAbs(c.GitDir) {
		c.GitDir = filepath.Join(os.Getenv("PWD"), c.GitDir)
	}
	if c.GitDir == "" {
		log.Fatalf("git-dir not set")
	}
	return c.GitDir
}

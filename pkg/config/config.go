package config

import (
	"fmt"
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
		c.GitDir = filepath.Join(GetPwd(), c.GitDir)
	}
	if c.GitDir == "" {
		gitdir := fmt.Sprintf("%s/.git", GetPwd())
		if st, err := os.Stat(gitdir); !os.IsNotExist(err) && st.IsDir() {
			return gitdir
		}
		if hasDir("objects") && hasDir("info") {
			return GetPwd()
		}
		log.Fatal("git-dir not found")
	}
	return c.GitDir
}

// GetPwd return PWD
func GetPwd() string {
	if pwd := os.Getenv("PWD"); pwd != "" {
		return pwd
	}
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Getwd", err)
	}
	return pwd
}

func hasDir(dir string) bool {
	st, err := os.Stat(fmt.Sprintf("%s/%s", GetPwd(), dir))
	return !os.IsNotExist(err) && st.IsDir()
}

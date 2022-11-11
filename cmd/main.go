package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/chiyutianyi/git-hashsum/pkg/config"
	"github.com/chiyutianyi/git-hashsum/pkg/version"
)

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:     "git-hashsum",
	Short:   "Shasum for git repositories.",
	Version: version.Version(),
}

func main() {
	if err := Cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "git-hashsum failed: %+v", err)
		os.Exit(1)
	}
}

func init() {
	flags := Cmd.PersistentFlags()
	flags.StringVarP(&config.Cfg.LogLevel, "log-level", "", "info", "log level")
	bindGitDir(flags, &config.Cfg.GitDir)
}

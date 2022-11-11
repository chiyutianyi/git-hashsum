package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/chiyutianyi/git-hashsum/pkg/config"
	"github.com/chiyutianyi/git-hashsum/pkg/lock"
)

type initCmd struct {
}

func (cmd *initCmd) Run(_ *cobra.Command, args []string) error {
	log.SetLevel(config.Cfg.GetLogLevel())

	l := lock.New(fmt.Sprintf("%s/info/checksum", config.Cfg.GetGitDir()))
	if err := l.Lock(); err != nil {
		return fmt.Errorf("failed to create lock: %v", err)
	}
	defer l.Unlock()

	rs, err := initShaSum()
	if err != nil {
		return err
	}
	log.Debugf("hashsum %x", rs)
	l.Write([]byte(fmt.Sprintf("%x", rs)))
	l.Flush()
	return nil
}

func init() {
	init := &initCmd{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init hashsum for current repository",
		Run: func(cmd *cobra.Command, args []string) {
			if err := init.Run(cmd, args); err != nil {
				log.Fatalf("init failed", err)
			}
		},
	}
	Cmd.AddCommand(cmd)
}

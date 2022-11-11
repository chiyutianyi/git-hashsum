package main

import (
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/chiyutianyi/git-hashsum/pkg/config"
	"github.com/chiyutianyi/git-hashsum/pkg/hashsum"
)

type initCmd struct {
}

func (cmd *initCmd) Run(_ *cobra.Command, args []string) {
	log.SetLevel(config.Cfg.GetLogLevel())

	c := exec.Command("git", "show-ref", "--head")

	out, err := c.CombinedOutput()
	if err != nil {
		log.Fatalf("failed to get ref head: %v", err)
		os.Exit(1)
	}

	if len(out) == 0 {
		return
	}

	var rs [32]byte

	for _, pair := range strings.Split(string(out), "\n") {
		if pair == "" {
			continue
		}
		a := strings.Split(pair, " ")
		if len(a) != 2 {
			log.Errorf("bad reference: %v", pair)
			continue
		}
		log.Debugf("reference: %v -> %v", a[1], a[0])
		// old_checksum XOR hash(refname,  newvalue)
		rs = hashsum.Sum(rs, a[0], a[1])
	}
	log.Infof("hashsum %x", rs)
}

func init() {
	init := &initCmd{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init hashsum for current repository",
		Run:   init.Run,
	}
	Cmd.AddCommand(cmd)
}

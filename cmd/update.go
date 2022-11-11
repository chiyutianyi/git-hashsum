package main

import (
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/chiyutianyi/git-hashsum/pkg/config"
	"github.com/chiyutianyi/git-hashsum/pkg/hashsum"
)

type updateCmd struct {
}

func (cmd *updateCmd) Run(c *cobra.Command, args []string) {
	log.SetLevel(config.Cfg.GetLogLevel())

	buffer, err := ioutil.ReadAll(c.InOrStdin())
	if err != nil {
		log.Fatal("read stdin", err)
	}

	var rs [32]byte

	for _, pair := range strings.Split(string(buffer), "\n") {
		if pair == "" {
			continue
		}
		a := strings.Split(pair, " ")
		if len(a) != 3 {
			log.Errorf("bad reference: %v", pair)
			continue
		}
		log.Debugf("reference: %v -> %v", a[1], a[0])
		// old_checksum XOR hash(refname, oldvalue)
		rs = hashsum.Sum(rs, a[0], a[2])
		// old_checksum XOR hash(refname,  newvalue)
		rs = hashsum.Sum(rs, a[1], a[2])
	}
	log.Infof("hashsum %x", rs)
}

func init() {
	update := &updateCmd{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update hashsum for current repository",
		Run:   update.Run,
	}
	Cmd.AddCommand(cmd)
}

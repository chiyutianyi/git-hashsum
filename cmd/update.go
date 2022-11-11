package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/chiyutianyi/git-hashsum/pkg/config"
	"github.com/chiyutianyi/git-hashsum/pkg/hashsum"
	"github.com/chiyutianyi/git-hashsum/pkg/lock"
)

type updateCmd struct {
}

func (cmd *updateCmd) Run(c *cobra.Command, args []string) error {
	log.SetLevel(config.Cfg.GetLogLevel())

	l := lock.New(fmt.Sprintf("%s/info/checksum", config.Cfg.GetGitDir()))
	old, err := l.LockAndRead()
	if err != nil {
		return fmt.Errorf("failed to create lock: %v", err)
	}
	defer l.Unlock()

	oldSum, err := hex.DecodeString(old)
	if err != nil {
		return fmt.Errorf("failed to decode old checksum %v: %v", old, err)
	}

	var rs [32]byte

	if len(oldSum) == 0 {
		rs, err = initShaSum()
		if err != nil {
			return err
		}
	} else {
		if len(oldSum) != 32 {
			return fmt.Errorf("bad old checksum %v", old)
		}
		for i := 0; i < 32; i++ {
			rs[i] = oldSum[i]
		}
	}

	buffer, err := ioutil.ReadAll(c.InOrStdin())
	if err != nil {
		return fmt.Errorf("failed to read stdin: %v", err)
	}

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
	log.Debugf("hashsum %x", rs)
	l.Write([]byte(fmt.Sprintf("%x", rs)))
	l.Flush()
	return nil
}

func init() {
	update := &updateCmd{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update hashsum for current repository",
		Run: func(cmd *cobra.Command, args []string) {
			if err := update.Run(cmd, args); err != nil {
				log.Fatalf("init failed", err)
			}
		},
	}
	Cmd.AddCommand(cmd)
}

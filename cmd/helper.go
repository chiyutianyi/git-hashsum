package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chiyutianyi/git-hashsum/pkg/hashsum"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func bindGitDir(flags *pflag.FlagSet, gitdir *string) {
	flags.StringVarP(gitdir, "git-dir", "C", "", "git dir")
}

func initShaSum() ([32]byte, error) {
	var rs [32]byte
	c := exec.Command("git", "show-ref", "--head")
	c.Env = os.Environ()

	out, err := c.CombinedOutput()
	if err != nil {
		return rs, fmt.Errorf("failed to get ref head: %v", err)
	}

	if len(out) == 0 {
		return rs, nil
	}

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
	return rs, nil
}

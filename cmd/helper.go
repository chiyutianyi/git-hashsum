package main

import "github.com/spf13/pflag"

func bindGitDir(flags *pflag.FlagSet, gitdir *string) {
	flags.StringVarP(gitdir, "git-dir", "C", "", "git dir")
}

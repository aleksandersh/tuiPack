package cli

import "github.com/alexflint/go-arg"

type Args struct {
	Config  string `arg:"-c,--config" default:"./commandPack.toml" help:"path to config file"`
	Script  string `arg:"-s,--script" help:"run script by the alias"`
	Aliases bool   `arg:"-a,--aliases" help:"print list of aliases for available scripts"`
}

func GetArgs() *Args {
	args := &Args{}
	arg.MustParse(args)
	return args
}

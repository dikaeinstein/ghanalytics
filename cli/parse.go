package cli

import (
	"bytes"
	"flag"
	"fmt"
)

type Config struct {
	help bool

	// args are the positional (non-flag) command-line arguments.
	args []string
}

const usage = `GhAnalytics is a CLI tool which analyzes Github event data for 1 hour.

Usage:
  ghanalytics [command]

Available Commands:
  topTenUsers			Top 10 active users sorted by amount of PRs created and commits.
  top10ReposByCommitsPushed	Top 10 repositories sorted by amount of commits pushed.
  top10ReposByWatchEvents	Top 10 repositories sorted by amount of watch events.

Flags:
  -h, -help	Show help`

// ParseFlags parses the command-line arguments provided to the program.
// Typically os.Args[0] is provided as 'progname' and os.args[1:] as 'args'.
// Returns the Config in case parsing succeeded, or an error. In any case, the
// output of the flag.Parse is returned in output.
// A special case is usage requests with -h or -help: then the error
// flag.ErrHelp is returned and output will contain the usage message.
func parseArgs(progname string, args []string) (config *Config, output string, err error) {
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)
	flags.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), usage)
	}

	var conf Config
	flags.BoolVar(&conf.help, "help", false, "Show help")
	flags.BoolVar(&conf.help, "h", false, "Show help")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}

	conf.args = flags.Args()
	return &conf, buf.String(), nil
}

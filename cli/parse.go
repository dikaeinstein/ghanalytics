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

const usage = `ghanalytics is a CLI tool which analyzes Github event data for 1 hour.

Usage of ghanalytics:

$ ghanalytics listTopTenUsers
`

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

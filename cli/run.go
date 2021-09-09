package cli

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/dikaeinstein/ghanalytics/analytics"
	"github.com/dikaeinstein/ghanalytics/data"
)

func Run() int {
	conf, output, err := parseArgs(os.Args[0], os.Args[1:])
	if err == flag.ErrHelp {
		fmt.Println(output)
		return 2
	}

	if err != nil {
		fmt.Println("got error:", err)
		fmt.Println("output:\n", output)
		return 1
	}
	if len(conf.args) < 1 {
		fmt.Println(usage)
		return 2
	}

	if conf.help {
		fmt.Println(usage)
		return 0
	}

	if err := run(conf); err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}

func run(conf *Config) error {
	actorsCSVFile, err := os.Open("data/actors.csv")
	if err != nil {
		return err
	}
	commitsCSVFile, err := os.Open("data/commits.csv")
	if err != nil {
		return err
	}
	eventsCSVFile, err := os.Open("data/events.csv")
	if err != nil {
		return err
	}
	reposCSVFile, err := os.Open("data/repos.csv")
	if err != nil {
		return err
	}

	store, err := data.NewStore(actorsCSVFile, commitsCSVFile,
		eventsCSVFile, reposCSVFile)
	if err != nil {
		return err
	}

	an := analytics.New(store)

	switch conf.args[0] {
	case "listTop10Users":
		return handleListUsers(an)
	default:
		return fmt.Errorf("unknown subcommand: %s", conf.args[0])
	}
}

func handleListUsers(an *analytics.Analytics) error {
	users, err := an.ListUsers(
		analytics.Sort([]analytics.SortCriteria{
			analytics.CommitsPushed, analytics.PrCreated,
		}),
		analytics.Limit(10),
	)
	if err != nil {
		return err
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
	fmt.Fprintln(tw, "ID\tUsername\t")
	fmt.Fprintln(tw, "-\t-\t")
	for _, u := range users {
		fmt.Fprintf(tw, "%v\t%v\t\n", u.ID, u.Username)
	}
	return tw.Flush()
}

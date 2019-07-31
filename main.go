package main

import (
	"fmt"
	"github.com/ribbondz/gsv/cmd"
	"github.com/ribbondz/gsv/cmd/utility"
	"github.com/ribbondz/gsv/cmd_desc"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "gsv"
	app.Version = "0.0.4"
	app.Usage = "Csv toolkit focused on performance and parallel processing"

	app.Commands = []cli.Command{
		{
			Name:        "head",
			Usage:       "Show head n records of file",
			Description: cmd_desc.Head,
			Action: func(c *cli.Context) error {
				path := c.Args().First()
				n := c.Int("l")
				cmd.Head(path, n)
				return nil
			},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "length, l",
					Usage: "Number of records to display",
					Value: 20,
				},
			},
		},
		{
			Name:        "header",
			Usage:       "Show headers of CSV file",
			Description: cmd_desc.Header,
			Action: func(c *cli.Context) error {
				path := c.Args().First()
				sep := utility.SepArg(c.String("s"))
				cmd.Header(path, sep)
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "sep, s",
					Usage: "File separator",
					Value: ",",
				},
			},
		},
		{
			Name:        "count",
			Usage:       "Count total lines of file",
			Description: cmd_desc.Count,
			Action: func(c *cli.Context) error {
				path := c.Args().First()
				header := !c.Bool("n")
				cmd.Count(path, header)
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-header, n",
					Usage: "When set, the first row will NOT be interpreted as column names",
				},
			},
		},
		{
			Name:        "cat",
			Usage:       "Concatenate files in a directory",
			Description: cmd_desc.Cat,
			Action: func(c *cli.Context) error {
				path := c.Args().First()
				header := !c.Bool("n")
				pattern := c.String("p")
				cmd.Cat(path, header, pattern)
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-header, n",
					Usage: "When set, the first row will NOT be interpreted as column names",
				},
				cli.StringFlag{
					Name:  "pattern, p",
					Usage: "Pattern of files to concat, default to all files",
					Value: "*",
				},
			},
		},
		{
			Name:        "partition",
			Usage:       "Partitions CSV file into chunks based on a column value",
			Description: cmd_desc.Partition,
			Action: func(c *cli.Context) error {
				path := c.Args().First()
				header := !c.Bool("n")
				column := c.Int("c")
				sep := utility.SepArg(c.String("s"))
				summary := c.Bool("summary")
				cmd.Partition(path, header, column, sep, summary)
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-header, n",
					Usage: "When set, the first row will NOT be interpreted as column names",
				},
				cli.IntFlag{
					Name:  "column, c",
					Usage: "Partition by which column",
				},
				cli.StringFlag{
					Name:  "sep, s",
					Usage: "File separator",
					Value: ",",
				},
				cli.BoolFlag{
					Name:  "summary",
					Usage: "Generate a summary file tabling line counts for each column value",
				},
			},
		},
		{
			Name:        "stats",
			Usage:       "Show statistics (e.g., min, max, average, unique count, null) on every column",
			Description: cmd_desc.Stats,
			Action: func(c *cli.Context) error {
				path := c.Args().First()
				header := !c.Bool("n")
				sep := utility.SepArg(c.String("s"))
				cmd.Stats(path, header, sep)
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-header, n",
					Usage: "When set, the first row will NOT be interpreted as column names",
				},
				cli.StringFlag{
					Name:  "sep, s",
					Usage: "File separator",
					Value: ",",
				},
			},
		},
		{
			Name:        "frequency",
			Usage:       "Show frequency tables",
			Description: cmd_desc.Frequency,
			Action: func(c *cli.Context) error {
				path := c.Args().First()
				header := !c.Bool("n")
				sep := utility.SepArg(c.String("s"))
				col, err := utility.ParseColArg(c.String("c"))
				if err != nil {
					fmt.Println("column selection syntax error.")
					return nil
				}
				out := c.Bool("o")
				ascending := c.Bool("a")
				limit := c.Int("l")
				cmd.Frequency(path, header, sep, col, out, ascending, limit)
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-header, n",
					Usage: "When set, the first row will NOT be interpreted as column names",
				},
				cli.StringFlag{
					Name:  "sep, s",
					Usage: "File separator",
					Value: ",",
				},
				cli.StringFlag{
					Name:  "col, c",
					Usage: "Select a subset of columns, default to first column",
					Value: "0",
				},
				cli.IntFlag{
					Name:  "limit, l",
					Usage: "Limit the frequency table to the N most common items. Set to '0' to disable a limit",
					Value: 50,
				},
				cli.BoolFlag{
					Name:  "output, o",
					Usage: "Print the frequency table to an output file, instead of stdout",
				},
				cli.BoolFlag{
					Name:  "ascending, a",
					Usage: "Frequency table in ascending order, default to descending",
				},
			},
		},
		{
			Name:        "select",
			Usage:       "Select rows and columns based on filters",
			Description: cmd_desc.Select,
			Action: func(c *cli.Context) error {
				path := c.Args().First()
				header := !c.Bool("n")
				sep := utility.SepArg(c.String("s"))
				filter := c.String("f")
				col, err := utility.ParseColArg(c.String("c"))
				if err != nil {
					fmt.Println("column selection syntax error.")
					return nil
				}
				out := c.Bool("o")
				cmd.Select(path, header, sep, filter, col, out)
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-header, n",
					Usage: "When set, the first row will NOT be interpreted as column names",
				},
				cli.StringFlag{
					Name:  "sep, s",
					Usage: "File separator",
					Value: ",",
				},
				cli.StringFlag{
					Name:  "filter, f",
					Usage: "Filter criterion, see filter syntax in description",
				},
				cli.StringFlag{
					Name:  "col, c",
					Usage: "Select a subset of columns, default to all column",
					Value: "-1",
				},
				cli.BoolFlag{
					Name:  "output, o",
					Usage: "Print the frequency table to an output file, instead of stdout",
				},
			},
		},
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Printf("No matching command '%s', available commands are ['head', 'header', 'count', 'cat', 'frequency', 'partition', 'select', 'stats']", command)
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

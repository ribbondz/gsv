package main

import (
	"fmt"
	"github.com/ribbondz/gsv/cmd"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "gsv"
	app.Version = "0.0.1"
	app.Usage = "Csv toolkit focused on performance and parallel processing"

	app.Commands = []cli.Command{
		{
			Name:  "head",
			Usage: "Show head n records of file",
			Description: `examples:
	gsv head a.txt         // head 20 rows (default)
	gsv head -l 50 a.txt   // head 50 rows
`,
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
			Name:  "count",
			Usage: "Count total rows of the file",
			Description: `
	examples:
	gsv count a.txt
`,
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
			Name:  "cat",
			Usage: "Cat files in a directory",
			Description: `examples:
	gsv cat data_dir                // has header, all files in data_dir (default)
	gsv cat -n data_dir             // no header, all files
	gsv cat -n -p *.txt data_dir    // no header, all txt files
	gsv cat -p *.csv data_dir       // all csv files
`,
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
			Name:  "partition",
			Usage: "Partitions the given CSV data into chunks based on the value of a column",
			Description: `examples:
	gsv partition a.txt                          // has header, partition by first column, no summary file (default)
	gsv partition -n a.txt                       // no header
	gsv partition -c 0 a.txt                     // partition by first column
	gsv partition -c 1 a.txt                     // partition by second column
	gsv partition -s , a.txt                     // sep ,
	gsv partition -s \t a.txt                    // sep \t
	gsv partition -summary a.txt                 // generate a summary file
	gsv partition -n -c 1 -s , -summary a.txt    // all options
`,
			Action: func(c *cli.Context) error {
				path := c.Args().First()
				header := !c.Bool("n")
				column := c.Int("c")
				sep := c.String("s")
				if sep == "t" || sep == "\\t" {
					sep = "\t"
				}
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
					Usage: "Partitions by which column",
				},
				cli.StringFlag{
					Name:  "sep, s",
					Usage: "File separation",
					Value: ",",
				},
				cli.BoolFlag{
					Name:  "summary",
					Usage: "Generate a summary file stating how many records for each column value",
				},
			},
		},
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Printf("No matching command '%s', available commands are ['head', 'count', 'cat', 'partition']", command)
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

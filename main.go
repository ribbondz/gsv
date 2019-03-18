package main

import (
	"fmt"
	"github.com/ribbondz/gsv/cmd"
	"github.com/ribbondz/gsv/cmd/utility"
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
			Name:  "header",
			Usage: "Show headers of CSV file",
			Description: `examples:
	 gsv header a.txt         // separator "," (default)
	 gsv header -s \t a.txt   // separator tab
`,
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
			Name:  "count",
			Usage: "Count total lines of file",
			Description: `examples:
	 gsv count a.txt
	 gsv count --help           // help info 
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
			Usage: "Concatenate files in a directory",
			Description: `examples:
	 gsv cat data_dir                // has header, all files in data_dir (default)
	 gsv cat -n data_dir             // no header, all files
	 gsv cat -n -p *.txt data_dir    // no header, all txt files
	 gsv cat -p *.csv data_dir       // all csv files
	 gsv cat --help                  // help info 
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
			Usage: "Partitions CSV file into chunks based on a column value",
			Description: `examples:
	 gsv partition a.txt                          // has header, partition by first column, no summary file (default)
	 gsv partition -n a.txt                       // no header
	 gsv partition -c 0 a.txt                     // partition by first column
	 gsv partition -c 1 a.txt                     // partition by second column
	 gsv partition -s , a.txt                     // sep ,
	 gsv partition -s \t a.txt                    // sep \t
	 gsv partition -summary a.txt                 // generate a summary file
	 gsv partition -n -c 1 -s , -summary a.txt    // all options
	 gsv partition --help                         // help info 
`,
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
			Name:  "stats",
			Usage: "Show statistics (e.g., min, max, average, unique count, null) on every column",
			Description: `examples:
	 gsv stats a.txt           // has header, separator "," (default)
	 gsv stats -n a.txt        // no header
	 gsv stats -s \t a.txt     // tab separator
	 gsv stats --help          // help info
`,
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
			Name:  "frequency",
			Usage: "Show frequency tables",
			Description: `output fields:
	 Col,  Value,  Count
	 col_1,    a,     10
	 col_1,    b,     20

	 examples:
	 gsv frequency a.txt           // first column, has header, separator "," (default)
	 gsv frequency -n a.txt        // no header
	 gsv frequency -s \t a.txt     // tab separator
	 gsv frequency -c 0 a.txt      // frequency table on first column (default)
	 gsv frequency -c 1 a.txt      // frequency table on second column
	 gsv frequency -c 0,1 a.txt    // frequency table on first and second columns
	 gsv frequency -l 10 a.txt     // keep top 10 records
	 gsv frequency -a a.txt        // frequency table in ascending order, default to descending
	 gsv frequency -o a.txt        // Print the frequency table to output file named "a-current-time.txt"
	 gsv frequency --help          // help info

	 column selection syntax:
	 '1,2':   cols [1,2]
	 '1-3,6': cols [1,2,3,6]
	 '!1':    cols [all except col 1]
	 '-1':    cols [all]
`,
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
			Name:  "select",
			Usage: "Select rows and columns based on filters",
			Description: `examples:
	 gsv select -f 0=abc a.txt                       // has header, separator ",", first column is 'abc',
	                                                 // set FILTER criterion using -f flag
	 gsv select -f "0=abc|0=de"" a.txt               // first column is 'abc' or 'de'
	 gsv select -f "0=abc&1=de"" a.txt               // first column is 'abc' and second column is 'de'
	 gsv select -f 0=abc -c 0,1,2 a.txt              // output keeps only columns 0, 1, and 2
	 gsv select -f 0=abc -o a.txt                    // save result to a-select-current-time.txt
	 gsv select -n -s \t -f 0=abc -c 0,1,2 -o a.txt  // all options
	 gsv select --help                               // help info on other options
	
	 column filter syntax:
	 -f '0=abc':       first column equal to string 'abc'
	 -f '1=5.0':       second column equal to number 5.0
	 -f '1=5':         same as pre command, second column equal to number 5.0
	 -f '0=abc&1=5.0': first column is 'abc' AND second column is 5.0
	 -f '0=abc|1=5.0': first column is 'abc' OR second column is 5.0
	
	 NOTE: 1. more complex syntax with brackets 
	          such as '(0=abc|1=5.0)&c=1' is not supported.
	       2. one filter can only have & or |, but never both. 
	          This feature maybe be added in the future.
`,
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

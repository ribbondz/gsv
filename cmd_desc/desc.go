package cmd_desc

const (
	Head = `examples:
	 gsv head a.txt         // head 20 rows (default)
	 gsv head -l 50 a.txt   // head 50 rows
`

	Header = `examples:
	 gsv header a.txt         // separator "," (default)
	 gsv header -s \t a.txt   // separator tab
`

	Count = `examples:
	 gsv count a.txt
	 gsv count --help           // help info 
`

	Cat = `examples:
	 gsv cat data_dir                // has header, all files in data_dir (default)
	 gsv cat -n data_dir             // no header, all files
	 gsv cat -n -p *.txt data_dir    // no header, all txt files
	 gsv cat -p *.csv data_dir       // all csv files
	 gsv cat --help                  // help info 
`

	Partition = `examples:
	 gsv partition a.txt                          // has header, partition by first column, no summary file (default)
	 gsv partition -n a.txt                       // no header
	 gsv partition -c 0 a.txt                     // partition by first column
	 gsv partition -c 1 a.txt                     // partition by second column
	 gsv partition -s , a.txt                     // sep ,
	 gsv partition -s \t a.txt                    // sep \t
	 gsv partition -summary a.txt                 // generate a summary file
	 gsv partition -n -c 1 -s , -summary a.txt    // all options
	 gsv partition --help                         // help info 
`

	Stats = `examples:
	 gsv stats a.txt           // has header, separator "," (default)
	 gsv stats -n a.txt        // no header
	 gsv stats -s \t a.txt     // tab separator
	 gsv stats --help          // help info
`

	Frequency = `output fields:
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
`

	Select = `examples:
	 gsv select -f 0=abc a.txt                       // has header, separator ",", first column is 'abc',
	                                                 // set FILTER criterion using -f flag
	 gsv select -f "0=abc|0=de"" a.txt               // first column is 'abc' or 'de'
	 gsv select -f "0=abc&1=de"" a.txt               // first column is 'abc' and second column is 'de'
	 gsv select -f 0=abc -c 0,1,2 a.txt              // output keeps only columns 0, 1, and 2
	 gsv select -f 0=abc -o a.txt                    // save result to a-select-current-time.txt
	 gsv select -n -s \t -f 0=abc -c 0,1,2 -o a.txt  // all options
	 gsv select -c 0,1 -o a.txt                      // NO filter, only to select columns
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
           3. The filter option can be omitted to select all rows.
`
)

# csv toolkit written in golang
gsv is a command line program to deal with CSV files. Gsv has following features:

- fast and parallel processing
- real-time progress bar
- simple usage

## 1. Usage
download gsv.exe from release tab; and choose the either one:
- put gsv.exe in system path
- put gsv.exe and the data in same folder

## 2. Available commands
- **head** - Show head n lines of CSV file.
- **count** - Count the lines in CSV file.
- **cat** - Concatenate CSV files by row ***(with progress bar)***.
- **frequency** - Show frequency table on columns.
- **partition** - Split CSV file based on a column value ***(with progress bar)***.
- **select** - Select rows and columns from CSV file.
- **stats** - Show statistics (e.g., min, max, average, unique count, null) on every column.

## 3. Examples

- gsv head
```shell
gsv head a.txt        // default to first 20 rows
gsv head -l 30 a.txt  // first 30 rows
gsv head --help       // help info on all flags
```

- gsv count
```shell
gsv count a.txt      // default to have a header
gsv count -n a.txt   // no header
gsv count --help     // help info on all flags
```

- gsv cat
```shell
gsv cat data_dir            // concatenate all files in data_dir directory, 
                            // assume a header for all files,
                            // output file is named to data_dir-current-time.txt
gsv cat -n data_dir         // no header
gsv cat -p * data_dir       // file pattern, default to all files
gsv cat -p *.csv data_dir   // all csv files in the directory
gsv cat --help              // help info on all flags
```

- gsv frequency
```shell
gsv frequency a.txt           // first column, has header, separator "," (default)
gsv frequency -n a.txt        // no header
gsv frequency -s \t a.txt     // tab separator
gsv frequency -c 0 a.txt      // frequency table on first column (default)
gsv frequency -c 1 a.txt      // frequency table on second column
gsv frequency -c 0,1 a.txt    // frequency table on first and second columns
gsv frequency -l 10 a.txt     // keep top 10 records
gsv frequency -a a.txt        // frequency table in ascending order, default to descending
gsv frequency -o a.txt        // Print the frequency table to output file named "a-current-time.txt"
gsv frequency --help          // help info on all flags

column selection syntax:
-c "1,2"   -->    cols [1,2]
-c "1-3,6" -->    cols [1,2,3,6]
-c "!1"    -->    cols [all except col 1]
-c "-1"    -->    cols [all]

frequency table:
+-------+-------+-------+
|  COL  | VALUE | COUNT |
+-------+-------+-------+
| col_1 |     a |     2 |
| col_1 |     b |     2 |
| col_2 |     3 |     2 |
| col_2 |     2 |     1 |
| col_2 |     4 |     1 |
+-------+-------+-------+
```

- gsv partition
```shell
gsv partition a.txt            // default to split by first column, separator ",", with file header
gsv partition -n a.txt         // no header
gsv partition -c 0 a.txt       // split by first column (default)
gsv partition -c 1 a.txt       // split by second column
gsv partition -s , a.txt       // row separator is "," (default) 
gsv partition -s \t a.txt      // row separator is tab
gsv partition -summary a.txt   // generate a summary file tabling the number of lines for unique column values
gsv partition --help           // help info on all flags
```

- gsv select
```shell
gsv select -f 0=abc a.txt                       // has header, separator ",", first column is "abc"
                                                // set FILTER criterion using -f flag
gsv select -f "0=abc|0=de" a.txt                // first column is "abc" or "de"
gsv select -f "0=abc&1=de" a.txt                // first column is "abc" and second column is "de"
gsv select -f 0=abc -c 0,1,2 a.txt              // output keeps only columns 0, 1, and 2
gsv select -f 0=abc -o a.txt                    // save result to a-filter-current-time.txt
gsv select -n -s \t -f 0=abc -c 0,1,2 -o a.txt  // all options
gsv select --help                               // help info on other options
	
column filter syntax:
-f "0=abc"        -->  first column equal to string "abc"
-f "1=5.0"        -->  second column equal to number 5.0
-f "1=5"          -->  same as pre command, second column equal to number 5.0
-f "0=abc&1=5.0"  -->  first column is "abc" AND second column is 5.0
-f "0=abc|1=5.0"  -->  first column is "abc" OR second column is 5.0

NOTE: 1. more complex syntax with brackets 
	 such as "(0=abc|1=5.0)&c=1" is not supported.
      2. one filter can only have & or |, but never both. 
	 This feature maybe be added in the future.
	     
column selection syntax:
-c "1,2"   -->    cols [1,2]
-c "1-3,6" -->    cols [1,2,3,6]
-c "!1"    -->    cols [all except col 1]
-c "-1"    -->    cols [all]
```

- gsv stats
```shell
gsv stats a.txt           // has header, separator "," (default)
gsv stats -n a.txt        // no header
gsv stats -s \t a.txt     // tab separator
gsv partition --help      // help info on all flags

statistics table.
+------+--------+------+--------+---------------------+---------------------+----------+------------+------------+
| COL  |  TYPE  | NULL | UNIQUE |         MIN         |          MAX        |   MEAN   | MIN LENGTH | MAX LENGTH |
+------+--------+------+--------+---------------------+---------------------+----------+------------+------------+
| col1 | string |    0 | 965304 | 00000208bb80146803f | ffffebf8245861dd564 |        - |         32 |         32 |
| col2 |  float |    0 |      - |             30.1054 |             31.3370 |  30.6524 |          2 |          9 |
| col3 |  float |    0 |      - |            103.0818 |            104.8750 | 104.0399 |          3 |         10 |
| col4 |  float |    0 |      - |             30.1041 |             31.3370 |  30.6522 |          2 |          9 |
| col5 |  float |    0 |      - |            103.0839 |            104.8742 | 104.0392 |          3 |         10 |
| col6 | string |    0 | 566252 | 2016-11-07 00:00:00 | 2016-11-14 00:00:00 |        - |         23 |         23 |
| col7 | string |    0 | 586711 | 1900-01-01 00:00:00 | 2021-09-24 13:52:23 |        - |         23 |         23 |
| col8 |  float |    0 |      - |              0.0000 |             84.9298 |   2.0013 |          1 |         22 |
+------+--------+------+--------+---------------------+--------------- -----+----------+------------+------------+
Total records: 9703035
Time consumed: 6s
```

# 4. Next
new features will be added in the future.

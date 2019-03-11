# csv toolkit written in golang
gsv is a command line program to deal with CSV files. Gsv has following features:

- fast, parallel processing
- real-time progress bar
- simple

## 1. Usage
choose the either one:
- put the gsv.exe file in system path
- put gsv.exe and the data in the same folder

## 2. Available commands
- **head** - Show head n lines of CSV file.
- **count** - Count the lines in the CSV file.
- **cat** - Concatenate CSV files by row.
- **partition** - Split CSV file based on a column value.

## 3. Examples

- gsv head
```shell
gsv head a.txt        // default to first 20 rows
gsv head -l 30 a.txt  // first 30 rows
gsv head --help       // help info
```

- gsv count
```shell
gsv count a.txt      // default to have a header
gsv count -n a.txt   // no header
gsv count --help     // help info
```

- gsv cat
```shell
gsv cat data_dir            // concatenate all files in data_dir directory, 
                            // assume a header for all files,
                            // output file is named to data_dir-current-time.txt
gsv cat -n data_dir         // no header
gsv cat -p * data_dir       // file pattern, default to all files
gsv cat -p *.csv data_dir   // all csv files in the directory
gsv cat --help              // help info
```

- gsv partition
```shell
gsv partition a.txt            // default to split by first column, separator ",", with file header
gsv partition -n a.txt         // no header
gsv partition -c 0 a.txt       // split by first column (default)
gsv partition -c 1 a.txt       // split by second column
gsv partition -s , a.txt       // row separator is "," (default) 
gsv partition -s \t a.txt      // row separator is tab
gsv partition -summary a.txt   // generate a summary file tabling the number of lines for each unique column value
gsv partition --help           // help info
```

# 4. Next
new features will be added in the future.

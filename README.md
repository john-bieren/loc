# loc

A command line utility for recursively counting lines of code in a directory and its subdirectories by language.

### Example output

<pre>
<code><b>Language: loc | bytes | files</b>
7 langs: 13,112 | 714,683 | 53
Python: 12,206 | 684,094 | 44
C++: 371 | 14,326 | 2
C: 355 | 11,360 | 2
Go: 77 | 1,882 | 1
Shell: 39 | 938 | 1
Ruby: 36 | 1,144 | 1
Powershell: 28 | 939 | 2
</code></pre>

## Install

1. Clone the repository:
    ```
    git clone https://github.com/john-bieren/loc.git
    ```
2. Navigate to the project directory and compile the program:
    ```
    go build
    ```
3. Add the project directory to your PATH environment variable to use `loc` system-wide

## Usage

This usage information can be found with `loc --help`:

```
Usage: loc [options] [path]
         Options must come before path
         Path defaults to current working directory if no argument is given

Options:
        -d        Print loc by directory
        -ed str   Directories to exclude (use name or full path, i.e. "src,lib,C:/Users/user/loc")
        -ef str   Files to exclude (use name or full path, i.e. "index.js,utils.go,C:/Users/user/lib/main.py")
        -el str   Languages to exclude (i.e. "HTML,Plain Text,JSON")
        -f        Print loc by file
             -mf int   Maximum number of files to print per directory
        -id       Include dot directories
        -md int   Maximum depth of subdirectories to search
        -ml int   Maximum number of language loc totals to print per directory
        -p        Print loc as a percentage of overall total
        -s str    Choose how to sort results ["loc", "bytes", "files"] (defult: "loc")
        -v        Print version and exit
```

## Disclaimer

This is my first project in Go. It's a personal project, and there are more advanced solutions, like [scc](https://github.com/boyter/scc), with more features and better methods for counting lines of code. This program counts any non-empty line as a line of code, including comments.

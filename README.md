# loc

A command line utility for recursively counting lines of code in directories and their subdirectories by language.

### Example output

<pre>
<code><b>Language: loc | size | files</b>
7 langs: 12,231 | 714.7 kb | 53
Python: 11,392 | 684.1 kb | 44
C++: 351 | 14.3 kb | 2
C: 323 | 11.4 kb | 2
Go: 72 | 1.9 kb | 1
Ruby: 33 | 1.1 kb | 1
Shell: 32 | 938 b | 1
Powershell: 28 | 939 b | 2
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
Usage: loc [options] [paths]
         Options must come before paths
         Paths are the directories you wish to search (cwd by default)

Options:
        -d        Print loc by directory
        -ed str   Directories to exclude (use name or full path, i.e. "src,lib,C:/Users/user/loc")
        -ef str   Files to exclude (use name or full path, i.e. "index.js,utils.go,C:/Users/user/lib/main.py")
        -el str   Languages to exclude (i.e. "HTML,Plain Text,JSON")
        -f        Print loc by file
             -mf int   Maximum number of files to print per directory (default: 100,000)
        -id       Include dot directories
        -md int   Maximum depth of subdirectories to search (default: 1,000)
        -ml int   Maximum number of language loc totals to print per directory (default: 1,000)
        -p        Print loc as a percentage of overall total
        -s str    Choose how to sort results ["loc", "size", "files"] (default: "loc")
        -v        Print version and exit
```

## Disclaimer

I'm new to Go and this is just a personal project. There are more advanced solutions like [scc](https://github.com/boyter/scc) with more features and better methods for counting lines of code. This program counts docstrings and multi-line comments as lines of code.

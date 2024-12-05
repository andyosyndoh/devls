# **my-ls**

## **Overview**

`my-ls` is a custom implementation of the Unix `ls` command, written in Go. It replicates the core functionality of the `ls` command, which lists the contents of directories and files. The project supports multiple flags for enhanced functionality and is designed to mimic the behavior of the original `ls` command as closely as possible.

---

## **Features**

The `my-ls` command includes the following features:
- Displays the contents of directories and files.
- Supports various flags to modify the output:
  - `-l`: Displays detailed information about files and directories (long listing format).
  - `-R`: Recursively lists subdirectories.
  - `-a`: Includes hidden files (those starting with `.`).
  - `-r`: Reverses the order of the listing.
  - `-t`: Sorts files by modification time.
- Handles symbolic links correctly.
- Outputs formatted results similar to the system `ls` command.

---

## **Objectives**

This project aims to:
1. Provide an understanding of Unix system commands.
2. Explore methods of receiving and processing user input.
3. Improve skills in data manipulation and output formatting.
4. Work with file system structures and flags in Go.

---

## **Usage**

### **Installation**
1. Clone the repository:
   ```bash
   git clone https://learn.zone01kisumu.ke/git/aosindo/my-ls-1.git
   ```
2. Navigate to the project directory:
   ```bash
   cd my-ls-1
   ```

### **Running the Command**
Use `my-ls` just like the standard `ls` command. Examples:
- List files and directories in the current directory:
  ```bash
  go run .
  ```
- Use flags for specific behaviors:
  ```bash
  go run . -l
  go run . -a
  go run . -R
  go run . -t -r
  ```
- Specify a directory:
  ```bash
  go run . /path/to/directory
  ```

---

## **Example Outputs**

### Without Flags
```bash
$ go run .
file1  file2  folder1
```

### With `-l` Flag
```bash
$ go run . -l
drwxr-xr-x  2 user group 4096 Dec  5 10:00 folder1
-rw-r--r--  1 user group 1024 Dec  5 09:00 file1
-rw-r--r--  1 user group 2048 Dec  4 15:00 file2
```

### With `-R` Flag
```bash
$ go run . -R
folder1:
file3
file4
```

---


## **Code Structure**

### **Key Functions**
1. **`ShortList(files []string, flags map[string]bool)`**: Handles short listing of files and directories.
2. **`LongList(files []string, flags map[string]bool)`**: Handles detailed (long format) listing with `-l` flag.
3. **`sortFiles(files []string)`**: Sorts files and directories for organized output.
4. **`shouldShowFile(name string, showHidden bool)`**: Determines whether a file should be displayed based on the `-a` flag.

### **Directory Recursion**
The `-R` flag implementation ensures all subdirectories are recursively listed. Hidden files (`.` and `..`) are handled appropriately.

### **Color Formatting**
- Directories, files, and symbolic links are color-coded for better readability in terminal output.

---

## **Testing**

It is recommended to use unit tests to verify the functionality of your code. Create test cases for:
- Basic file and directory listing.
- Behavior with various flag combinations.
- Recursive directory traversal with `-R`.

To run tests:
```bash
go test ./...
```

---

## **Authors**
 
GitHub: [andyosyndoh](https://github.com/andyosyndoh)

GitHub: [hezronokwach](https://github.com/hezronokwach)
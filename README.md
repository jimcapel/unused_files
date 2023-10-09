# Unused files

Unused files search's javascript/typescript projects for files that are not imported by any other file(unused). Written to help find old files in a large react-native project after rewriting a lot of the codebase (but forgetting to delete old files).

# Usage

To use, you can

```
go install github.com/jimcapel/unused_files
```

to install the application at your go path, then simply run with the command

```
unused_files search {path_to_root_directory}
```

# Performance mode

You can add the flag:

```
    -- performance || --p
```

to log the amount of time taken for the program to run.

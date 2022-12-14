package rm

import (
  "github.com/braswelljr/rmx/pkg/utils"
  "os"
)

// ShowHelp - show the help message
func (rm *RM) ShowHelp() error {
  // show the help message
  println("rmx - remove files or directories")
  println("Usage: rmx [OPTION]... [FILE]...")
  println("Remove (unlink) the FILE(s).")
  println("  -f, --force           ignore nonexistent files and arguments, never prompt")
  println("  -i                    prompt before every removal")
  println("  -I                    prompt once before removing more than three files, or when removing recursively; less intrusive than -i, while still giving protection against most mistakes")
  println("  -r, -R                remove directories and their contents recursively")
  println("  -d                    remove empty directories")
  println("  -v                    explain what is being done")
  println("      --help     display this help and exit")
  println("      --version  output version information and exit")
  println("Note:  rm does not remove directories.  Use rmdir to remove directories.")
  return nil
}

func (rm *RM) Run(args []string) error {
  // check for empty args
  if len(args) < 1 {
    return rm.ShowHelp()
  }

  // get the directory
  dir, err := os.Open(rm.Directory)
  if err != nil {
    panic(err)
  }

  // get the directory info
  dirInfo, err := dir.Stat()
  if err != nil {
    panic(err)
  }

  // check if the directory is a directory
  if dirInfo.IsDir() {
    // get the directory entries
    dirEntries, err := dir.Readdir(0)
    if err != nil {
      panic(err)
    }

    var Directories Directories

    // loop through the directory entries
    for _, dirEntry := range dirEntries {
      // check if the directory entry is a directory
      if dirEntry.IsDir() {
        // get the directory path
        dirPath := rm.Directory + string(os.PathSeparator) + dirEntry.Name()

        // create a new directory
        directory := Directory{
          Name: dirEntry.Name(),
          Path: dirPath,
          Info: dirEntry,
        }

        // add the directory to the directories
        Directories = append(Directories, directory)
      }
    }

    // loop through the directories
    for _, directory := range Directories {
      // check if the directory is empty
      if utils.IsEmpty(directory.Path) {
        // remove the directory
        if err := os.Remove(directory.Path); err != nil {
          panic(err)
        }

        // check if the verbose flag is set
        if rm.Flags.V {
          // print the directory
          println(directory.Path)
        }
      }
    }
  }

  // close the directory
  return nil
}

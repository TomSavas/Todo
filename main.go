package main 

import (
	"os"
)

func main() {
	GetDatabasePath()
	OpenDatabase()
	defer CloseDatabase()
	DetectCommand(os.Args)
}
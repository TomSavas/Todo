package main 

import (
	"os"
)

func main() {
	OpenDatabase()
	defer CloseDatabase()
	DetectCommand(os.Args)
}
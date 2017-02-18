package main

import (
	"fmt"
	"flag"
	"strings"
	"os"
	"regexp"
)

func Split(src string) []string {
	query, _ := regexp.Compile("\\s*?;+?\\s*")
	return query.Split(src, -1)
}

func AddCommand() {
	priority := flag.String("p", "top", ADD_P_FLAG_INFO)
	status := flag.String("s", "not_started", ADD_S_FLAG_INFO)
	types := flag.String("t", "general", ADD_T_FLAG_INFO)
	notes := flag.String("n", "", ADD_N_FLAG_INFO)	
	flag.Parse()
}

func LsCommand() {
	priority := flag.String("p", "", LS_P_FLAG_INFO)
	status := flag.String("s", "", LS_S_FLAG_INFO)
	types := flag.String("t", "", LS_T_FLAG_INFO)
	flag.Parse()
	Print(GetTodos(Split(priority), Split(status), Split(types)))
}

func LsdCommand() {
	priority := flag.String("p", "none", LSD_P_FLAG_INFO)
	task := flag.String("t", "general", LSD_T_FLAG_INFO)
	flag.Parse()
	_, _ = priority, task
}

func AppendCommand() {
	priority := flag.String("n", "task", APPEND_N_FLAG_INFO)
	flag.Parse()
	_ = priority 
}

func RmCommand() {

}

func ChpriCommand() {

}

func ChstatCommand() {

}

func ChtypeCommand() {

}

func ChnoteCommand() {

}

func ChtaskCommand() {

}

func DoneCommand() {

}


func DetectCommand(args []string) {
	if len(args) < 2 {
		fmt.Println("No arguments given.", HINT_FOR_HELP)
		return
	} else if len(args) == 3 && strings.ToLower(args[1]) == "-h" {			
		PrintSpecificInfo(args[2])
		return
	}

	os.Args = os.Args[1:]

	switch strings.ToLower(args[1]) {
	case "add":
		AddCommand()
	case "ls":
		LsCommand()
	case "lsd":
		LsdCommand()
	case "append":
		AppendCommand()
	case "rm":
		RmCommand()
	case "chpri":
		ChpriCommand()
	case "chstat":
		ChstatCommand()
	case "chtype":
		ChtypeCommand()
	case "chnote":
		ChnoteCommand()
	case "chtask":
		ChtaskCommand()
	case "done":
		DoneCommand()
	case "help":
		fallthrough
	case "-h":
		fmt.Println(UsageHelp)
	default:
		fmt.Println(args[1], "command was not recognized.", HINT_FOR_HELP)
	}
}
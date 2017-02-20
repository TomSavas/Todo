package main

import (
	"fmt"
	"flag"
	"strings"
	"os"
	"regexp"
)

func Split(src string) []string {
	if len(src) == 0 {
		return []string{}
	}

	query, _ := regexp.Compile("\\s*?;+?\\s*")
	return query.Split(src, -1)
}

func ValidateID(id string) bool {
	if found, _ := regexp.MatchString("[^\\d\\s]+", id); found {
		fmt.Println(ID_CANT_CONTAIN_LETTERS, HINT_FOR_HELP)
		return false
	} else if found, _ = regexp.MatchString("\\s+", id); found{
		fmt.Println(ID_CANT_CONTAIN_SPACES, HINT_FOR_HELP)
	}
	return true
}

func ValidateIDs (ids []string) bool {
	for _, id := range ids {
		if !ValidateID(id){
			return false
		}
	}
	return true
}

func AddCommand() {
	priority := flag.String("p", "top", ADD_P_FLAG_INFO)
	status := flag.String("s", "not_started", ADD_S_FLAG_INFO)
	types := flag.String("t", "general", ADD_T_FLAG_INFO)
	notes := flag.String("n", "", ADD_N_FLAG_INFO)	
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println(ZeroArgumentsGiven)
		return
	}

	if len(*notes) != 0 {
		*notes += "\n"
	}

	todo := Todo{0, strings.Join(flag.Args(), " "), *priority, *status, *types, *notes}
	AddTodo(todo)
}

func LsCommand() {
	priority := flag.String("p", "", LS_P_FLAG_INFO)
	status := flag.String("s", "", LS_S_FLAG_INFO)
	types := flag.String("t", "", LS_T_FLAG_INFO)
	flag.Parse()
	Print(GetTodos(Split(*priority), Split(*status), Split(*types)))
}

func LsdCommand() {
	priority := flag.String("p", "", LSD_P_FLAG_INFO)
	types := flag.String("t", "", LSD_T_FLAG_INFO)
	flag.Parse()
	Print(GetTodos(Split(*priority), []string{"done"}, Split(*types)))
}

func LswCommand() {
	priority := flag.String("p", "", LSD_P_FLAG_INFO)
	types := flag.String("t", "", LSD_T_FLAG_INFO)
	flag.Parse()
	Print(GetTodos(Split(*priority), []string{"wip"}, Split(*types)))
}

func AppendCommand() {
	priority := flag.String("n", "task", APPEND_N_FLAG_INFO)
	flag.Parse()
	_ = priority 
}

func RmCommand() {
	fmt.Println(os.Args)
	if ValidateIDs(strings.Split(os.Args[1], " ")) {
		for _, id := range strings.Split(os.Args[1], " ") {
			RemoveTodo(id)
		}
	}
}

func ChCommand(field, value string) {
	if ValidateIDs(strings.Split(os.Args[1], " ")){
		for _, id := range strings.Split(os.Args[1], " ") {
			ChangeField(id, field, value)
		}
	}	
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
	case "lsw":
		LswCommand()
	case "append":
		AppendCommand()
	case "del":
		fallthrough
	case "rm":
		RmCommand()
	case "chpri":
		fallthrough
	case "chpriority":
		ChCommand("Priority", os.Args[2])
	case "chst":
		fallthrough
	case "chstatus":
		ChCommand("Status", os.Args[2])
	case "chtype":
		ChCommand("Type", os.Args[2])
	case "chnote":
		ChCommand("Note", os.Args[2])
	case "chtask":
		ChCommand("Task", os.Args[2])
	case "done":
		ChCommand("Status", "done")
	case "backup":
		CloseDatabase()
		BackupDatabase()
	case "restore":
		CloseDatabase()
		RestoreDatabase()
	case "help":
		fallthrough
	case "-h":
		fmt.Println(UsageHelp)
	default:
		fmt.Println(args[1], "command was not recognized.", HINT_FOR_HELP)
	}
}
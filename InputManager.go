package main

import (
	"fmt"
	"flag"
	"strings"
	"os"
	"time"
	"strconv"
)

func PrintDefaultAmountOfTodos() {
	PrintNTodos(GetDefaultPrintLength())
}

func ChangeDefaultPrintLength() {
	flag.Parse()
	if len(flag.Args()) == 1 {
		if flag.Args()[0][0] > 47 && flag.Args()[1][0] < 58 {
			SetDefaultPrintLength(flag.Args()[0])
		} else {
			fmt.Println(CHANGE_DEFAULT_PRINT_LENGTH_INVALID_NUMBER)			
		}
	} else {
		fmt.Println(TOO_MANY_ARGUMENTS, HINT_FOR_HELP)
	}
}

func AddCommand() {
	priority := flag.String("p", "top", ADD_P_FLAG_INFO)
	status := flag.String("s", "not_started", ADD_S_FLAG_INFO)
	types := flag.String("t", "general", ADD_T_FLAG_INFO)
	timeCap := flag.String("d", "-1", ADD_D_FLAG_INFO)
	notes := flag.String("n", "", ADD_N_FLAG_INFO)	
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println(ZERO_ARGUMENTS_GIVEN, HINT_FOR_HELP)
		return
	}

	if len(*notes) != 0 {
		*notes += "\n"
	}

	timeCapp, _ := strconv.ParseFloat(*timeCap, 64)
	todo := NewTodo(0, time.Now().Unix(), strings.Join(flag.Args(), " "), *priority, *status, *types, FloatToString(float64(time.Now().Unix() + (int64(timeCapp) * 86400))), "0", *notes)
	AddTodo(*todo)
}

func LsCommand() {
	priority := flag.String("p", "", LS_P_FLAG_INFO)
	status := flag.String("s", "", LS_S_FLAG_INFO)
	types := flag.String("t", "", LS_T_FLAG_INFO)
	flag.Parse()
	PrintTodos(false, GetTodos(SplitBySemicolons(*priority), SplitBySemicolons(*status), SplitBySemicolons(*types)))
}

func LsdCommand() {
	priority := flag.String("p", "", LSD_P_FLAG_INFO)
	types := flag.String("t", "", LSD_T_FLAG_INFO)
	flag.Parse()
	PrintTodos(false, GetTodos(SplitBySemicolons(*priority), []string{"done"}, SplitBySemicolons(*types)))
}

func LswCommand() {
	priority := flag.String("p", "", LSD_P_FLAG_INFO)
	types := flag.String("t", "", LSD_T_FLAG_INFO)
	flag.Parse()
	PrintTodos(false, GetTodos(SplitBySemicolons(*priority), []string{"wip"}, SplitBySemicolons(*types)))
}

func RmCommand() {
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
	if len(args) == 1 {
		PrintDefaultAmountOfTodos()
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
	case "del":
		fallthrough
	case "rm":
		RmCommand()
	case "chpri":
		fallthrough
	case "chpriority":
		ChCommand("Priority", strings.Join(os.Args[2:], " "))
	case "chst":
		fallthrough
	case "chstatus":
		ChCommand("Status", strings.Join(os.Args[2:], " "))
	case "chtype":
		ChCommand("Type", strings.Join(os.Args[2:], " "))
	case "chnote":
		ChCommand("Notes", strings.Join(os.Args[2:], " "))
	case "chtask":
		ChCommand("Task", strings.Join(os.Args[2:], " "))
	case "chd":
		integer, _ := strconv.ParseFloat(os.Args[2], 64)
		ChCommand("TimeCap", fmt.Sprintf("%v", int64(float64(time.Now().Unix()) + integer * 86400)))
	case "chnum":
		ChangeDefaultPrintLength()
	case "done":
		ChCommand("Status", "done")
	case "backup":
		CloseDatabase()
		BackupDatabase()
	case "restore":
		CloseDatabase()
		RestoreDatabase()
	case "reset":
		ResetDatabase()
	case "help":
		fallthrough
	case "-h":
		fmt.Println(UsageHelp)
	case "+":
		ChCommand("+Progress", strings.Join(os.Args[2:], " "))
	case "-":
		ChCommand("-Progress", strings.Join(os.Args[2:], " "))
	default:
		fmt.Println(args[1], "command was not recognized.", HINT_FOR_HELP)
	}
}
package main

import (
	"fmt"
	"strings"
)

const (
	ADD_P_FLAG_INFO = `Sets the priority of the task, system values (although custom can be typed in): top, mid, low.
             Default value - MID. Values are case insensitive.`
	ADD_S_FLAG_INFO = `Sets the state of the task, system values (although custom can be typed in): wip (work in progress), 
             not-started, done, dead. Default value - NOT_STARTED. Values are case insensitive.`	
	ADD_T_FLAG_INFO = `Sets the type of the task. Custom ones must be used, not that if multiple words are used, 
             they must be typed between quotes. Default value - GENERAL.
               e.g.
               todo add -t "School tasks" "Prepare for final maths exam"
               todo add -t Piano "Learn Maple leaf rag"`
	ADD_N_FLAG_INFO = `Sets custom notes for the task.
               e.g.
               todo add -n "Tasks 1.100 through 1.125" "Solve physics problems"`	
	ADD_COMMAND_USAGE = `add [OPTIONS] [TASK TEXT]
      Description:
        Adds a new task.
      Options:
        -p   ` + ADD_P_FLAG_INFO + `
        -s   ` + ADD_S_FLAG_INFO + `
        -t   ` + ADD_T_FLAG_INFO + `
        -n   ` + ADD_N_FLAG_INFO + `
`

	LS_P_FLAG_INFO = `Specify the priority, case insensitive, if priority contains multiple words, must be typed between quotes.
             Multiple priorities must be typed between quotes and separated with semicolons. If not given any priorities prints
             existing priorities.`
	LS_S_FLAG_INFO = `Specify the state, case insensitive if state contains multiple words, must bet typed between quotes. 
             Multiple states must be typed between quotes and separated with semicolons. If not given any states prints existing states.`
	LS_T_FLAG_INFO = `Specify the type, case insensitive, if type contains multiple words, must be typed between quotes.
             Multiple types must be typed between quotes and separated with semicolons. If not given any types prints existing types.
               e.g.
               todo ls -p TOP -s wip -t piano
               todo ls -p "top, mid" -t "piano, school tasks"`
	LS_COMMAND_USAGE = `ls [OPTIONS]
      Description:
          Lists todos specified by options. If not given any lists all todos.
      Options:
        -p   ` + LS_P_FLAG_INFO + `
        -s   ` + LS_S_FLAG_INFO + `
        -t   ` + LS_T_FLAG_INFO + `
`

	LSD_P_FLAG_INFO = `Specify the priority, case insensitive, if priority contains multiple words, must be typed between quotes.
             Multiple priorities must be typed between quotes and separated with semicolons. If not given any priorities
             prints existing priorities.`
	LSD_T_FLAG_INFO = `Specify the type, case insensitive, if type contains multiple words, must be typed between quotes. 
             Multiple types must be typed between quotes and separated with semicolons. If not given any types prints existing types.
               e.g.
               todo lsd -p TOP -t piano
               todo lsd -p "top; mid" -t "piano, school tasks"`
	LSD_COMMAND_USAGE = `lsd [OPTIONS]
      Description:
        Lists done todos. An alternative to "ls -s done".
      Options:
        -p   ` + LSD_P_FLAG_INFO + `
        -t   ` + LSD_T_FLAG_INFO + `
`

	LSW_P_FLAG_INFO = `Specify the priority, case insensitive, if priority contains multiple words, must be typed between quotes.
             Multiple priorities must be typed between quotes and separated with semicolons. If not given any priorities
             prints existing priorities.`
	LSW_T_FLAG_INFO = `Specify the type, case insensitive, if type contains multiple words, must be typed between quotes. 
             Multiple types must be typed between quotes and separated with semicolons. If not given any types prints existing types.
               e.g.
               todo lsw -p TOP -t piano
               todo lsw -p "top; mid" -t "piano, school tasks"`
	LSW_COMMAND_USAGE = `lsw [OPTIONS]
      Description:
        Lists wip (work in progress) todos. An alternative to "ls -s wip".
      Options:
        -p   ` + LSW_P_FLAG_INFO + `
        -t   ` + LSW_T_FLAG_INFO + `
`

	APPEND_N_FLAG_INFO = `Specifies that text will be appended to the note of the task.`
	APPEND_COMMAND_USAGE = `append [OPTIONS] [ID] [TEXT_TO_APPEND]
      Description:
          Appends to task text or note of a todo with given ID. By default appends to the task of the todo.      
      Options:
          -n   ` + APPEND_N_FLAG_INFO + `
`

	RM_COMMAND_USAGE = `rm [ID]
      Description:
        Deletes todo with a given id.
	`

	CHPRI_COMMAND_USAGE = `chpri [ID] [PRIORITY]
      Description:
        Changes priority of a todo with given ID. If priority is not specified, sets it to default.
	`

	CHSTAT_COMMAND_USAGE = `chstat [ID] [STATE]
      Description:
        Changes state of a todo with given ID. If state is not specified, sets it to default.
	`

	CHTYPE_COMMAND_USAGE = `chtype [ID] [TYPE]
      Description:
        Changes type of a todo with given ID. If state is not specified, sets it to default.
	`

	CHNOTE_COMMAND_USAGE = `chnote [ID] [NOTE]
      Description:
        Overrides the note of a todo with given ID. If note is not specified, removes it.
	`

	CHTASK_COMMAND_USAGE = `chtask [ID] [TASK]
      Description:
        Overrides the task of a todo with given ID. Task must be given.
	`

	DONE_COMMAND_USAGE = `done [ID]
      Description:
        Sets todo's state to done.
	`

	HINT_FOR_HELP = "Type \"todo help\" or \"todo -h\" for help. Type \"todo -h [COMMAND]\" for specific information about the command."
)

var ZeroArgumentsGiven string = "0 arguments given. " + HINT_FOR_HELP

var UsageHelp string = `Usage: todo [COMMAND]
  Commands:
    ` + ADD_COMMAND_USAGE + `
    ` + LS_COMMAND_USAGE + `
    ` + LSD_COMMAND_USAGE + `
    ` + LSW_COMMAND_USAGE + `		
    ` + APPEND_COMMAND_USAGE + `		
    ` + RM_COMMAND_USAGE + `
    ` + CHPRI_COMMAND_USAGE + `		 
    ` + CHSTAT_COMMAND_USAGE + `
    ` + CHTYPE_COMMAND_USAGE + `
    ` + CHNOTE_COMMAND_USAGE + `
    ` + CHTASK_COMMAND_USAGE + `
    ` + DONE_COMMAND_USAGE

func PrintSpecificInfo(command string) {
	switch strings.ToLower(command) {
	case "add":
		fmt.Println(ADD_COMMAND_USAGE)
	case "ls":
		fmt.Println(LS_COMMAND_USAGE)
	case "lsd":
		fmt.Println(LSD_COMMAND_USAGE)
	case "lsw":
		fmt.Println(LSW_COMMAND_USAGE)
	case "append":
		fmt.Println(APPEND_COMMAND_USAGE)
	case "rm":
		fmt.Println(RM_COMMAND_USAGE)
	case "chpri":
		fmt.Println(CHPRI_COMMAND_USAGE)
	case "chstat":
		fmt.Println(CHSTAT_COMMAND_USAGE)
	case "chtype":
		fmt.Println(CHTYPE_COMMAND_USAGE)
	case "chnote":
		fmt.Println(CHNOTE_COMMAND_USAGE)
	case "chtask":
		fmt.Println(CHTASK_COMMAND_USAGE)
	case "done":
		fmt.Println(DONE_COMMAND_USAGE)
	default:
		fmt.Println(command, "command was not recognized.", HINT_FOR_HELP)
	}	
}

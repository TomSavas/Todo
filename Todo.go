package main

import (
	"fmt"
	"strings"
	"time"
	"strconv"
	"math"
)

type Todo struct {
	Id int
	Time int64 //Unix timestamp 
	PriorityValue float64
	TimeCap int64
	Progress int
	AgeValue string
	Task string
	Priority string
	Status string
	Type string
	Note string
}

func NewTodo(id int, tiime int64, task, priority, status, types, timeCap, progress, notes string) *Todo {
	timeCapp, _ := strconv.ParseFloat(timeCap, 64)
	progresss, _ := strconv.ParseFloat(progress, 64)
	todo := &Todo{id, tiime, 0, int64(timeCapp), int(progresss), "", task, strings.ToUpper(priority), strings.ToUpper(status), types, notes} 
	todo.PriorityValue = todo.CalculatePriority();
	todo.AgeValue = ConvertToReadableTime(time.Now().Unix() - todo.Time);
	return todo
}

func (todo1 Todo) Equals(todo2 Todo) bool {
	if todo1.Id == todo2.Id && todo1.Task == todo2.Task &&
	   todo1.Priority == todo2.Priority && todo1.Status == todo2.Status && 
	   todo1.Type == todo2.Type && todo1.Note == todo2.Note {
	   	return true
	   }
	return false
}

func ConvertToReadableTime(time int64) string {
	switch {
		case time >= 31536000:
			return fmt.Sprintf("%.1f yrs", float64(time)/float64(31536000))		
		case time >= 2592000:
			return fmt.Sprintf("%.1f mon", float64(time)/float64(2592000))
		case time >= 604800:
			return fmt.Sprintf("%.1f w", float64(time)/float64(604800))
		case time >= 86400:
			return fmt.Sprintf("%.1f d", float64(time)/float64(86400))
		case time >= 3600:
			return fmt.Sprintf("%.1f h", float64(time)/float64(3600))
		case time >= 60:
			return fmt.Sprintf("%.1f min", float64(time)/float64(60))
		default:
			return fmt.Sprintf("%.1f sec", float64(time))
	}
}

func (todo *Todo) GetTimeCap() string {
	if todo.TimeCap <= 0 {
		return SelectSymbol("time", "-1")
	}
	return SelectSymbol("time", ConvertToReadableTime(int64(todo.TimeCap) - time.Now().Unix()))
}

//TODO: should be a more elegant solution
func (todo *Todo) CalculatePriority() float64 {
	ageRatio := func(isWip bool) float64 {
		currentDeltaTime := math.Abs(float64(GetOldestTime(isWip)) - float64(time.Now().Unix()))
		deltaTime := math.Abs(float64(GetOldestTime(isWip)) - float64(todo.Time))
		ratio := 1 - deltaTime / currentDeltaTime
		if ratio == 0 {
			ratio = 1
		}
		return Round(ratio, 0.1)
	}

	switch todo.Status {
	case "WIP":
		ratio := ageRatio(true)
		switch todo.Priority { 
		case "TOP":
			return ratio + 9
		case "MID":	
			return ratio * 2 + 7
		case "LOW":
			return ratio + 6
		default:
			return float64(0)
		}
	case "NOT_STARTED":
		ratio := ageRatio(false)
		switch todo.Priority {
		case "TOP":
			return ratio * 2 + 4
		case "MID":	
			return ratio * 2 + 2
		case "LOW":
			return ratio + 1
		default:
			return float64(0)
		}		
	case "DONE":
		fallthrough
	default:
		return float64(0)
	}
}

func SelectSymbol(field, value string) string {
	switch field {
	case "priority":
		switch strings.ToUpper(value) {
		case "TOP":
			return TOP_PRIORITY
		case "MID":
			return MID_PRIORITY
		case "LOW":
			return LOW_PRIORITY
		default:
			return NO_PRIORITY
		}	
	case "status":
		switch strings.ToUpper(value) {
		case "WIP":
			return WIP_STATUS
		case "DONE":
			return DONE_STATUS
		case "NOT_STARTED":
			fallthrough
		default:
			return NOT_STARTED_STATUS
		}
	case "time":
		switch {
		case value[0] == '-':
			//Because len(INFINITY) for some reason returns 3, which screws up formating
			return INFINITY + TerribleIndentationHack(8)
		default:
			return value
		}
	}
	return value
}

func GetPrintParameters(isShort bool, todo *Todo) (types, tasks []string, newlineTypePrintForm, newlineTaskPrintForm string) {
	types = FitStrings(SplitBySemicolons(strings.ToUpper(todo.Type)), 20, true)
	// tasks = SplitTextByNWords(todo.Task, 5)
	tasks = FitStrings([]string{todo.Task}, 25, false)
	
	if isShort {
		newlineTypePrintForm = NEWLINE_TYPE_SHORT_PRINT_FORM
		newlineTaskPrintForm = NEWLINE_TASK_SHORT_PRINT_FORM
	} else {
		newlineTypePrintForm = NEWLINE_TYPE_LONG_PRINT_FORM
		newlineTaskPrintForm = NEWLINE_TASK_LONG_PRINT_FORM
	}

	return
}

func (todo *Todo) GetProgressBar() string {
	progress := "["
	for i := 0; i < 10; i++ {
		if i < todo.Progress {
			progress += "#"
		} else {
			progress += "-"
		}
	}
	progress += "]"
	return progress
}

func PrintHeader(isShortVersion bool) {
	if isShortVersion {
		fmt.Println(fmt.Sprintf(TerribleIndentationHack(4) + ToBold(ToUnderline("ID")) + 
								TerribleIndentationHack(2) + ToBold(ToUnderline("PRI")) + 
								TerribleIndentationHack(1) + ToBold(ToUnderline("PROGRESS")) + 
								TerribleIndentationHack(5) + ToBold(ToUnderline("TIME CAP")) +  
								TerribleIndentationHack(1) + ToBold(ToUnderline("TYPE")) + 
								TerribleIndentationHack(17) + ToBold(ToUnderline("TASK"))))
	} else {
		fmt.Println(fmt.Sprintf(TerribleIndentationHack(4) + ToBold(ToUnderline("ID")) + 
								TerribleIndentationHack(2) + ToBold(ToUnderline("PRIORITY")) + 
								TerribleIndentationHack(2) + ToBold(ToUnderline("AGE")) + 
								TerribleIndentationHack(6) + ToBold(ToUnderline("PROGRESS")) +
								TerribleIndentationHack(5) + ToBold(ToUnderline("TIME CAP")) +  
								TerribleIndentationHack(1) + ToBold(ToUnderline("TYPE")) + 
								TerribleIndentationHack(17) + ToBold(ToUnderline("TASK"))))
	}
}

func PrintNotFitingTaskAndTypeText(types, tasks []string, newlineTypePrintForm, newlineTaskPrintForm string) {
	for i := 1; i < int(math.Max(float64(len(tasks)), float64(len(types)))); i++ {
		if i < len(types) && i >= len(tasks) {
			fmt.Println(fmt.Sprintf(newlineTypePrintForm, types[i]))
		} else if i >= len(types) && i < len(tasks) {
			fmt.Println(fmt.Sprintf(newlineTaskPrintForm, tasks[i])) 
		} else {
			fmt.Print(fmt.Sprintf(newlineTypePrintForm, types[i])) 
			fmt.Println(fmt.Sprintf("%v%v", TerribleIndentationHack(21 - len(types[i])), tasks[i]))
		}
	}
}

func PrintNote(isColored, isShort bool, todo *Todo) {
	if len(todo.Note) != 0 && !isShort {
		fmt.Println(TerribleIndentationHack(4) + ToUnderline("Note"))
		
		/* 
		 * Requires a fainting repetition because "NOTE" is underlined and that efect must be removed,
		 * which also removes the fainting effect.
		 */
		CheckAndFaint(isColored)
		defer EndModifiers()
		
		notes := SplitTextByNWords(todo.Note, 10) 
		for _, value := range(notes) {
			fmt.Println(TerribleIndentationHack(4) + value)
		}
	}
}

func (todo *Todo) Print(isColored, isShort bool) {
	types, tasks, newlineTypePrintForm, newlineTaskPrintForm := GetPrintParameters(isShort, todo)

	CheckAndFaint(isColored)
	defer EndModifiers()

	var firstTodoLine string
	if isShort {
		firstTodoLine = fmt.Sprintf("%v%v%v%v%v%v%v%v", 
									SelectSymbol("priority", todo.Priority) + TerribleIndentationHack(1), 
									SelectSymbol("status", todo.Status) + TerribleIndentationHack(1), 
									strconv.Itoa(todo.Id) + TerribleIndentationHack(4 - len(strconv.Itoa(todo.Id))), 
									FloatToString(todo.PriorityValue) + TerribleIndentationHack(4 - len(FloatToString(todo.PriorityValue))),
									todo.GetProgressBar() + TerribleIndentationHack(13 - len(todo.GetProgressBar())), 
									todo.GetTimeCap() + TerribleIndentationHack(9 - len(todo.GetTimeCap())),
									types[0] + TerribleIndentationHack(21 - len(types[0])),
									tasks[0])
	} else {
		firstTodoLine = fmt.Sprintf("%v%v%v%v%v%v%v%v%v", 
									SelectSymbol("priority", todo.Priority) + TerribleIndentationHack(1), 
									SelectSymbol("status", todo.Status) + TerribleIndentationHack(1), 
									strconv.Itoa(todo.Id) + TerribleIndentationHack(4 - len(strconv.Itoa(todo.Id))), 
									FloatToString(todo.PriorityValue) + TerribleIndentationHack(4 - len(FloatToString(todo.PriorityValue))) + 
									"[" + todo.Priority + "]" + TerribleIndentationHack(1),
									todo.AgeValue + TerribleIndentationHack(9 - len(todo.AgeValue)),
									todo.GetProgressBar() + TerribleIndentationHack(13 - len(todo.GetProgressBar())),
									todo.GetTimeCap() + TerribleIndentationHack(9 - len(todo.GetTimeCap())),
									types[0] + TerribleIndentationHack(21 - len(types[0])), tasks[0])
	}

	fmt.Println(firstTodoLine)
	PrintNotFitingTaskAndTypeText(types, tasks, newlineTypePrintForm, newlineTaskPrintForm)
	PrintNote(isColored, isShort, todo)
}

func PrintTodos(isShortVersion bool, todos []Todo) {
	Sort(0, len(todos)-1, todos)

	if len(todos) == 0 {
		fmt.Println(NO_TODOS_FOUND, "\n" + HINT_FOR_HELP)
		return
	}

	PrintHeader(isShortVersion)

	for i, todo := range(todos){
		todo.Print(i % 2 != 0, isShortVersion)
	}
}

func PrintNTodos(count int) {
	todos := GetTodosBy("", []string{})
	
	if len(todos) == 0 {
		fmt.Println(NO_TODOS_FOUND, "\n" + HINT_FOR_HELP)
		return
	}

	Sort(0, len(todos)-1, todos)
	if count > len(todos) - 1{
		count = len(todos)
	}
	PrintTodos(true, todos[:count])
}

func Sort(lower, upper int, todos []Todo) {
	pivot := (lower + upper) / 2
	l, u := lower, upper
	for l <= u {
		for todos[l].PriorityValue > todos[pivot].PriorityValue {
			l++;
		} 
		for todos[u].PriorityValue < todos[pivot].PriorityValue {
			u--;
		} 
		if l <= u {
			todos[u], todos[l] = todos[l], todos[u]
			l++; u--;
		} 
	}

	if lower < u {
		Sort(lower, u, todos)
	}
	if upper > l {
		Sort(l, upper, todos)
	}
}
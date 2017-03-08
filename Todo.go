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
	AgeValue string
	Task string
	Priority string
	Status string
	Type string
	Note string
}

func NewTodo(id int, time int64, task, priority, status, types, notes string) *Todo {
	todo := &Todo{id, time, 0, "", task, strings.ToUpper(priority), strings.ToUpper(status), types, notes} 
	todo.PriorityValue = todo.CalculatePriority();
	todo.AgeValue = todo.CalculateAge();
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

func (todo *Todo) CalculateAge() string {
	age := time.Now().Unix() - todo.Time
	switch {
		case age > 31536000:
			return fmt.Sprintf("%.1f yrs", Round(float64(age)/float64(31536000), 0.1))		
		case age > 2592000:
			return fmt.Sprintf("%.1f mon", Round(float64(age)/float64(2592000), 0.1))
		case age > 604800:
			return fmt.Sprintf("%.1f w", Round(float64(age)/float64(604800), 0.1))
		case age > 86400:
			return fmt.Sprintf("%.1f d", Round(float64(age)/float64(86400), 0.1))
		case age > 3600:
			return fmt.Sprintf("%.1f h", Round(float64(age)/float64(3600), 0.1))
		case age > 60:
			return fmt.Sprintf("%.1f min", Round(float64(age)/float64(60), 0.1))
		default:
			return fmt.Sprintf("%.1f sec", Round(float64(age), 0.1))
	}
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

func (todo *Todo) SelectPrioritySymbol() string {
	switch strings.ToUpper(todo.Priority){
	case "TOP":
		return TOP_PRIORITY
	case "MID":
		return MID_PRIORITY
	case "LOW":
		return LOW_PRIORITY
	default:
		return NO_PRIORITY
	}
}

func (todo *Todo) SelectStatusSymbol() string {
	switch strings.ToUpper(todo.Status){
	case "WIP":
		return WIP_STATUS
	case "DONE":
		return DONE_STATUS
	case "NOT_STARTED":
		fallthrough
	default:
		return NOT_STARTED_STATUS
	}
}

func GetPrintParameters(isColored, isShort bool, todo *Todo) (types, tasks []string, fullPrintForm, newlineTypePrintForm, newlineTaskPrintForm, newlineTaskAfterTypePrintForm string) {
	types = FitTypes(SplitBySemicolons(strings.ToUpper(todo.Type)))
	tasks = SplitTextByNths(todo.Task, 5)
	
	if isColored {
		StartFaint()
	}
	
	if isShort {
		fullPrintForm = FULL_SHORT_PRINT_FORM
		newlineTypePrintForm = NEWLINE_TYPE_SHORT_PRINT_FORM
		newlineTaskPrintForm = NEWLINE_TASK_SHORT_PRINT_FORM
		newlineTaskAfterTypePrintForm = NEWLINE_TASK_AFTER_TYPE_SHORT_PRINT_FORM
	} else {
		fullPrintForm = FULL_LONG_PRINT_FORM
		newlineTypePrintForm = NEWLINE_TYPE_LONG_PRINT_FORM
		newlineTaskPrintForm = NEWLINE_TASK_LONG_PRINT_FORM
		newlineTaskAfterTypePrintForm = NEWLINE_TASK_AFTER_TYPE_LONG_PRINT_FORM
	}

	return
}

func PrintHeader(isShortVersion bool) {
	if isShortVersion {
		fmt.Println(fmt.Sprintf("    " + ToBold(ToUnderline("ID")) + TerribleIndentationHack(4) + ToBold(ToUnderline("PRI\tAGE\t")) + "  " + ToBold(ToUnderline("TYPE\t\t\tTASK"))))
	} else {
		fmt.Println(fmt.Sprintf("    " + ToBold(ToUnderline("ID")) + TerribleIndentationHack(4) + ToBold(ToUnderline("PRIORITY\tAGE\t")) + "  " + ToBold(ToUnderline("TYPE\t\t\tTASK"))))
	}
}

func PrintNotFitingTaskAndTypeText(types, tasks []string, fullPrintForm, newlineTypePrintForm, newlineTaskPrintForm, newlineTaskAfterTypePrintForm string) {
	for i := 1; i < int(math.Max(float64(len(tasks)), float64(len(types)))); i++ {
		if i < len(types) && i >= len(tasks) {
			fmt.Println(fmt.Sprintf(newlineTypePrintForm, types[i]))
		} else if i >= len(types) && i < len(tasks) {
			fmt.Println(fmt.Sprintf(newlineTaskPrintForm, tasks[i])) 
		} else {
			fmt.Print(fmt.Sprintf(newlineTypePrintForm, types[i])) 
			fmt.Println(fmt.Sprintf("%v%v", TerribleIndentationHack(22 - len(types[i])), tasks[i]))
		}
	}
}

func PrintNote(isColored, isShort bool, todo *Todo) {
	if len(todo.Note) != 0 && !isShort {
		fmt.Println(TerribleIndentationHack(4) + ToUnderline("Note"))
		if isColored {
			StartFaint()
		}
		notes := SplitTextByNths(todo.Note, 10) 
		for _, value := range(notes) {
			fmt.Println(TerribleIndentationHack(4) + value)
		}
	}
}

func (todo *Todo) Print(isColored, isShort bool) {
	types, tasks, fullPrintForm, newlineTypePrintForm, newlineTaskPrintForm, newlineTaskAfterTypePrintForm := GetPrintParameters(isColored, isShort, todo)

	var firstTodoLine string
	if isShort {
		firstTodoLine = fmt.Sprintf(fullPrintForm, todo.SelectPrioritySymbol(), todo.SelectStatusSymbol(), 
									todo.Id, TerribleIndentationHack(6 - len(strconv.Itoa(todo.Id))), todo.PriorityValue,
									todo.AgeValue + TerribleIndentationHack(10 - len(todo.AgeValue)), types[0] + TerribleIndentationHack(22 - len(types[0])), tasks[0])
	} else {
		firstTodoLine = fmt.Sprintf(fullPrintForm, 
									todo.SelectPrioritySymbol(), todo.SelectStatusSymbol(), 
									todo.Id, TerribleIndentationHack(6 - len(strconv.Itoa(todo.Id))), 
									todo.PriorityValue, 
									TerribleIndentationHack(4-len(strconv.FormatFloat(todo.PriorityValue, 'f', -1, 32))) + "[" + todo.Priority + "]", 
									todo.AgeValue + TerribleIndentationHack(10 - len(todo.AgeValue)), types[0] + TerribleIndentationHack(22 - len(types[0])), tasks[0])
	}

	fmt.Println(firstTodoLine)
	PrintNotFitingTaskAndTypeText(types, tasks, fullPrintForm, newlineTypePrintForm, newlineTaskPrintForm, newlineTaskAfterTypePrintForm)
	PrintNote(isColored, isShort, todo)
	EndModifiers()
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
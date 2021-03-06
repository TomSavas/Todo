package main

import (
	"fmt"
	"github.com/mxk/go-sqlite/sqlite3"
	"os"
	"strings"
	"regexp"
	"os/exec"
	"os/user"	
)

var databasePath string

func GetDatabasePath() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}
	databasePath = usr.HomeDir + "/.todos.db"
}

func CheckIfFileExists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
  		return true
	}
	return false
}

func OpenDatabase(){
	if CheckIfFileExists(databasePath) {
		ResetDatabase()
		return
	}
	CreateDatabaseReference()
}

var db *sqlite3.Conn
func CreateDatabaseReference() {
	var err error
	db, err = sqlite3.Open(databasePath)
	_ = err
}

func CreateDatabase() {
	db.Exec("CREATE TABLE `todos` (`ID`	INTEGER, `Task`	TEXT, `Time`	INTEGER, `Priority`	TEXT, `Status`	TEXT, `Type`	TEXT, `TimeCap`	TEXT, `Progress`	TEXT, `Notes`	TEXT, PRIMARY KEY(ID));")
	db.Exec("CREATE TABLE `parameters` (`ID`	INTEGER, `DefaultLength`	INTEGER, PRIMARY KEY(ID));")
	db.Exec("INSERT INTO Parameters VALUES(1, 7)")
}

func DeleteDatabase() {
	CloseDatabase()
	os.Remove(databasePath)
}

func ResetDatabase() {
	DeleteDatabase()
	CreateDatabaseReference()
	CreateDatabase()
}

func CloseDatabase() {
	if db != nil {
		db.Close()
	}
}

func BackupDatabase(){
	command := exec.Command("cp", databasePath, databasePath+".bak")
	_ = command.Run()
}

func RestoreDatabase(){
	command := exec.Command("rm", databasePath)
	_ = command.Run()
	command = exec.Command("cp", databasePath+".bak", databasePath)
	_ = command.Run()
}

func GetTodosBy(parameter string, values []string) []Todo {
	var query string
	if len(values) == 0 || parameter == "" {
		query = "SELECT * FROM todos"
	} else {
		if NaiveSqlInjectionsCheck(values) {
			fmt.Println(POSSIBLE_SQL_INJECTION_ERROR)
			return []Todo{}
		}
		query = "SELECT * FROM todos WHERE " + parameter + " LIKE \"%" + strings.Join(values, "%\" OR Status LIKE \"%") + "%\""
	}
	return GetTodosFromQuery(query)
}

func GetTodosFromQuery(query string) []Todo {
	var todos []Todo 

	CastToString := func (element interface{}) string {
		if element == nil {
			return ""
		} else {
			return element.(string)	
		}
	}

	row := make(sqlite3.RowMap)
	for s, err := db.Query(query); err == nil; err = s.Next() {
		var rowid int64
		s.Scan(&rowid, row)
		// todos = append(todos, Todo{int(rowid), CastToString(row["Task"]), CastToString(row["Priority"]), CastToString(row["Status"]), CastToString(row["Type"]), CastToString(row["Note"])})
		todos = append(todos, *NewTodo(int(rowid), row["Time"].(int64), CastToString(row["Task"]), CastToString(row["Priority"]), CastToString(row["Status"]), CastToString(row["Type"]), CastToString(row["TimeCap"]), CastToString(row["Progress"]), CastToString(row["Notes"])))
	}

	return todos
}

func GetTodos(priorities, statuses, types []string) []Todo {
	Min := func (arrays [][]Todo, nums []int) int {
		min := 0
		minValue := arrays[0][nums[0]]
		for i := 1; i < len(nums); i++ {
			if(arrays[i][nums[i]].Id < minValue.Id){
				min = i
				minValue = arrays[i][nums[i]]
			}
		}
		return min
	}	

	Equality := func (arrays [][]Todo, nums []int) bool{
		for i := 1; i < len(nums); i++ {
			if !arrays[i-1][nums[i-1]].Equals(arrays[i][nums[i]]){
				return false
			}
		}
		return true
	}

	CheckIfDone := func (arrays [][]Todo, indexes []int) bool {
		for index, value := range indexes {
			if int(value) != len(arrays[index]) - 1 {
				return false
			} 
		}
		return true
	}

	IntersectArrays := func (arrays ...[]Todo) []Todo {
		var todos []Todo
		var indexes []int

		for _, todos := range arrays {
			if cap(todos) == 0 {
				return []Todo{}
			}
		}
			
		for _, _ = range arrays{
			indexes = append(indexes, 0)
		} 

		for i := 0; ; i++ {
			if Equality(arrays, indexes){
				todos = append(todos, arrays[0][indexes[0]])
			}

			if CheckIfDone(arrays, indexes) {
				break
			}

			min := Min(arrays, indexes)
			if(indexes[min] < len(arrays[min]) - 1){
				indexes[min]++	
			} else {
				break
			}
		}
		return todos
	}

	return IntersectArrays(GetTodosBy("Priority", priorities), GetTodosBy("Status", statuses), GetTodosBy("Type", types))
}

func AddTodo(todo Todo) {
	if NaiveSqlInjectionsCheck([]string{todo.Task, todo.Priority, todo.Status, todo.Type, todo.Note}) {
		fmt.Println(POSSIBLE_SQL_INJECTION_ERROR)
		return
	}
	db.Exec(fmt.Sprintf("INSERT INTO todos VALUES(%v, \"%v\", %v, \"%v\", \"%v\", \"%v\", \"%v\", \"%v\", \"%v\")", GetLastID() + 1, todo.Task, todo.Time, todo.Priority, todo.Status, todo.Type, todo.TimeCap, todo.Progress, todo.Note))
}

func RemoveTodo(id string) {
	db.Exec("DELETE FROM todos WHERE ID=" + id + ";")
}

func GetLastID() int {
	response, err := db.Query("SELECT * FROM todos WHERE ID = (SELECT MAX(ID) FROM todos);")
	if err != nil {
		return 0
	}

	var rowid int64
	response.Scan(&rowid, make(sqlite3.RowMap))
	return int(rowid)
}

func GetOldestTime(isWip bool) int64 {
	response, err := db.Query("SELECT * FROM todos WHERE Time = (SELECT MIN(Time) FROM (SELECT * FROM todos WHERE Status LIKE \"wip\"))")

	if !isWip {
		response, err = db.Query("SELECT * FROM todos WHERE Time = (SELECT MIN(Time) FROM (SELECT * FROM todos WHERE Status LIKE \"not_started\"))")
	}

	if err != nil {
		return 0
	}

	var rowid int64
	row := make(sqlite3.RowMap)
	response.Scan(&rowid, row)
	return row["Time"].(int64)
}

func GetDefaultPrintLength() int {
	response, err := db.Query("SELECT * FROM Parameters;")
	if err != nil {
		return 0
	}

	var rowid int64
	row := make(sqlite3.RowMap)
	response.Scan(&rowid, row)
	return int(row["DefaultLength"].(int64))
}

func SetDefaultPrintLength(length string) {
	if NaiveSqlInjectionCheck(length) {
		fmt.Println(POSSIBLE_SQL_INJECTION_ERROR)
		return
	}
	db.Exec("UPDATE parameters SET DefaultLength" + "=\"" + length + "\" WHERE ID = 1;")
}

func ChangeField(id, field, value string) {
	if NaiveSqlInjectionCheck(value) {
		fmt.Println(POSSIBLE_SQL_INJECTION_ERROR)
		return
	}
	if field == "+Progress" || field == "-Progress" {
		progress := GetTodosFromQuery("SELECT * FROM todos WHERE ID = " + id)[0].Progress
		if field[0] == '+' {
			if progress < 10 {
				progress++
			}
		} else {
			if progress >= 1 {
				progress--
			}
		}
		db.Exec(fmt.Sprintf("UPDATE todos SET Progress = \"%v\" WHERE ID = %v", progress, id))
	} else {
		db.Exec("UPDATE todos SET " + field + "=\"" + value + "\" WHERE ID = " + id + ";")	
	}
}

func NaiveSqlInjectionCheck(s string) bool {
	s = strings.ToUpper(s)
	if found, _ := regexp.MatchString("DROP\\sTABLE", s); found {
		return true
	} else if found, _ := regexp.MatchString("\\\"", s); found {
		return true
	} else if found, _ := regexp.MatchString("DELETE\\sFROM", s); found {
		return true
	} else if found, _ := regexp.MatchString("\\s--", s); found {
		return true
	}
	return false
}

func NaiveSqlInjectionsCheck(s []string) bool {
	for _, value := range s {
		if (NaiveSqlInjectionCheck(value)) {
			return true
		}
	}
	return false
}
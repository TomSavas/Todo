package main

import (
	"fmt"
	"github.com/mxk/go-sqlite/sqlite3"
	"os"
	"strings"
)

const databaseName = "todos.db"

func CheckIfFileExists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
  		return true
	}
	return false
}

func OpenDatabase(){
	if CheckIfFileExists(databaseName) {
		ResetDatabase()
		return
	}
	CreateDatabaseReference()
}

var db *sqlite3.Conn
func CreateDatabaseReference() {
	var err error
	db, err = sqlite3.Open(databaseName)
	_ = err
}

func CreateDatabase() {
	db.Exec("CREATE TABLE `todos` (`ID` INTEGER, `Task` TEXT, `Priority` TEXT, `Status` TEXT, `Type` TEXT, `Notes` TEXT, PRIMARY KEY(ID));")
}

func DeleteDatabase() {
	CloseDatabase()
	os.Remove(databaseName)
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

func GetTodoById(id int) Todo {
	query := "SELECT * FROM todos WHERE ID LIKE " + string(id)
	_ = query
	return Todo{}
}

func GetTodoByPriorities(priorities []string) []Todo {
	var query string
	if len(priorities) == 0 {
		query = "SELECT * FROM todos"
	} else {
		query = "SELECT * FROM todos WHERE Priority LIKE \"%" + strings.Join(priorities, "%\" OR Priority LIKE \"%") + "%\""
	}
	return GetTodosFromQuery(query)
}

func GetTodoByStatuses(statuses []string) []Todo {
	var query string
	if len(statuses) == 0 {
		query = "SELECT * FROM todos"
	} else {
		query = "SELECT * FROM todos WHERE Status LIKE \"%" + strings.Join(statuses, "%\" OR Status LIKE \"%") + "%\""
	}
	return GetTodosFromQuery(query)
}

func GetTodoByTypes(types []string) []Todo {
	var query string
	if len(types) == 0 {
		query = "SELECT * FROM todos"
	} else {
		query = "SELECT * FROM todos WHERE Type LIKE \"%" + strings.Join(types, "%\" OR Type LIKE \"%") + "%\""
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
		todos = append(todos, Todo{int(rowid), CastToString(row["Task"]), CastToString(row["Priority"]), CastToString(row["Status"]), CastToString(row["Type"]), CastToString(row["Notes"])})
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

	return IntersectArrays(GetTodoByPriorities(priorities), GetTodoByStatuses(statuses), GetTodoByTypes(types))
}

func AddTodo(todo Todo) {
	db.Exec(fmt.Sprintf("INSERT INTO todos VALUES(%v, \"%v\", \"%v\", \"%v\", \"%v\", \"%v\")", GetLastID() + 1, todo.Task, todo.Priority, todo.Status, todo.Type, todo.Notes))
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

func ChangeField(id, field, value string) {
	db.Exec("UPDATE todos SET " + field + "=\"" + value + "\" WHERE ID = " + id + ";")
}

func ChangePriority(id, priority string) {
	ChangeField(id, "Priority", priority)
}

func ChangeStatus(id, status string) {
	ChangeField(id, "Status", status)
}

func ChangeType(id, types string) {
	ChangeField(id, "Type", types)
}

func ChangeNote(id, note string) {
	ChangeField(id, "Note", note)
}

func ChangeTask(id, task string) {
	ChangeField(id, "Task", task)
}

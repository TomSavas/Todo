package main

import (
	"fmt"
)

type Todo struct {
	Id int64 
	Task string
	Priority string
	Status string
	Type string
	Notes string
}

func (todo1 Todo) Equals(todo2 Todo) bool {
	if todo1.Id == todo2.Id && todo1.Task == todo2.Task &&
	   todo1.Priority == todo2.Priority && todo1.Status == todo2.Status && 
	   todo1.Type == todo2.Type && todo1.Notes == todo2.Notes {
	   	return true
	   }
	return false
}

func (todo *Todo) Print() {
	fmt.Println(fmt.Sprintf("%v. %v \nPriority: %v	Type: %v	Status: %v \n%v", todo.Id, todo.Task, todo.Priority, todo.Type, todo.Status, todo.Notes))
}

func Print(todos []Todo) {
	for _, todo := range todos {
		todo.Print()
	}
}
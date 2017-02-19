package main

import (
	"fmt"
	"strings"
	"regexp"
)

type Todo struct {
	Id int 
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

func FormatParameters(src string) string {
	query, _ := regexp.Compile("\\s*?;+?\\s*")
	return "{" + strings.Join(query.Split(strings.ToUpper(src), -1), "} {") + "}"
}

func (todo *Todo) Print() {
	fmt.Println(fmt.Sprintf("%v. %v \nPriority: %v	Type: %v	Status: %v \n%v", todo.Id, todo.Task, FormatParameters(todo.Priority), FormatParameters(todo.Type), FormatParameters(todo.Status), todo.Notes))
}

func Print(todos []Todo) {
	if cap(todos) == 0 {
		fmt.Println("No todos found with such parameters.")
	}
	for _, todo := range todos {
		todo.Print()
	}
}
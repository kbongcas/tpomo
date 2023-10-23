package main

import (
	"encoding/json"
	"os"
)

type todo struct {
	Name string `json:"name"`
	Done bool   `json:"done"`
	Poms int    `json:"poms"`
}

func getTodos(fileName string) (todos []todo) {

	bytes, err := os.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, &todos)

	if err != nil {
		panic(err)
	}

	return todos
}

func saveTodos(filename string, todos []todo) {

	bytes, err := json.Marshal(todos)

	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filename, bytes, 0644)

	if err != nil {
		panic(err)
	}
}

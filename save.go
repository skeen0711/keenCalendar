package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func saveTasks() {
	file, err := os.Create("tasks.json")
	if err != nil {
		fmt.Println("Error saving tasks: ", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(tasks)
	if err != nil {
		fmt.Println("Error encoding tasks to file:", err)
	}
}

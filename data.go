package main

import (
	"awesomeProject/tasks"
	"encoding/json"
	"fmt"
	"os"
)

const taskFile = "tasks.json"

func saveTasks() {
	file, err := os.Create(taskFile)
	if err != nil {
		fmt.Println("Error saving tasks: ", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(tasks.tasks)
	if err != nil {
		fmt.Println("Error encoding tasks to file:", err)
	} else {
		fmt.Println("Tasks successfully saved to", taskFile)
	}
}

func loadTasks() {
	file, err := os.Open(taskFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No existing tasks found. Starting fresh!")
			return
		}
		fmt.Println("Error opening tasks file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks.tasks)
	if err != nil {
		fmt.Println("Error loading tasks:", err)
	} else {
		fmt.Printf("%d tasks loaded successfully from %s\n", len(tasks.tasks), taskFile)
	}
}

package main

import (
	"awesomeProject/tasks"
	"fmt"
)

func displayOptions() {
	fmt.Println("\n Todo CLI App")
	fmt.Println("1. Add Task")
	fmt.Println("2. List Tasks")
	fmt.Println("3. Complete Task")
	fmt.Println("4. Update Task")
	fmt.Println("5. Delete Task")
	fmt.Println("6. Save and Exit")
	fmt.Println("Select an Option")
}
func main() {
	// Load tasks from file upon startup
	loadTasks()

	for {
		displayOptions()

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			tasks.AddTask()
		case 2:
			tasks.ListTasks()
		case 3:
			tasks.CompleteTask()
		case 4:
			tasks.UpdateTask()
		case 5:
			tasks.DeleteTask()
		case 6:
			saveTasks()
			fmt.Println("Tasks saved! Exiting program...")
			return
		default:
			fmt.Println("Choose a valid option")
		}
	}
}

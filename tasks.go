package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Task struct {
	description    string
	completed      bool
	estimatedHours float64
	dueDate        time.Time
}

var tasks []Task

const dateFormat = "2006-01-02" // YYYY-MM-DD

func addTask() {
	fmt.Print("Enter task description: \n")
	reader := bufio.NewReader(os.Stdin)
	descr, _ := reader.ReadString('\n')

	descr = strings.TrimSpace(descr)
	fmt.Print("Enter task due date in format 'YYYY-MM-DD': ")
	dueDateStr, _ := reader.ReadString('\n')
	dueDateStr = strings.TrimSpace(dueDateStr)

	dueDate, err := time.Parse(dateFormat, dueDateStr)
	if err != nil {
		fmt.Println("Invalid date format, try again.")
		return
	}
	fmt.Print("Enter Estimated task hours: ")
	var estimatedHours float64
	fmt.Scan(&estimatedHours)
	tasks = append(tasks,
		Task{
			description:    descr,
			completed:      false,
			dueDate:        dueDate,
			estimatedHours: estimatedHours,
		})
	fmt.Println("Task added.")
}

func listTasks() {
	if len(tasks) == 0 {
		fmt.Println("You have no tasks!")
		return
	}
	fmt.Println("\n-----------------------\nTasks:")
	for i, task := range tasks {
		Complete := "No"
		if task.completed {
			Complete = "Yes"
		}

		fmt.Printf("%d. %s [Completed: %s || Due Date: %s || Est. Hours:  %.2f]\n",
			i+1,
			task.description,
			Complete,
			task.dueDate.Format(dateFormat),
			task.estimatedHours,
		)
	}
}

func completeTask() {
	listTasks()
	fmt.Println("Enter the task number to mark as completed: ")
	var taskNo int
	fmt.Scanln(&taskNo)

	if taskNo > 0 && taskNo <= len(tasks) {
		tasks[taskNo-1].completed = true
		updateDueDates(tasks[taskNo-1].dueDate)
		fmt.Println("Task marked as completed.")
	} else {
		fmt.Println("Invalid task number.")
	}
}

func updateDueDates(completedTaskDate time.Time) {
	for i, task := range tasks {
		if !task.completed {
			tasks[i].dueDate = task.dueDate.Add(24 * time.Hour)
		}
	}
	fmt.Println("Due dates for tasks updated.")
}

func updateTask() {
	listTasks()
	fmt.Print("Enter task number to update: ")
	var taskNo int
	fmt.Scan(&taskNo)

	if taskNo > 0 && taskNo <= len(tasks) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter new task description (leave blank to keep current): ")
		descr, _ := reader.ReadString('\n')
		descr = strings.TrimSpace(descr)

		// Update due date
		fmt.Print("Enter new due date (YYYY-MM-DD) or leave blank to keep current: ")
		dueDateStr, _ := reader.ReadString('\n')
		dueDateStr = strings.TrimSpace(dueDateStr)

		// Update Estimated Hours
		fmt.Print("Enter new estimated hours to complete (leave blank to keep current): ")
		var estimatedHours float64
		n, _ := fmt.Scanf("%f\n", &estimatedHours)

		if descr != "" {
			tasks[taskNo-1].description = descr
		}

		if dueDateStr != "" {
			dueDate, err := time.Parse(dateFormat, dueDateStr)
			if err != nil {
				fmt.Println("Invalid Date Format, task not updated")
				return
			}
			tasks[taskNo-1].dueDate = dueDate
		}

		if n == 1 {
			tasks[taskNo-1].estimatedHours = estimatedHours
		}

		fmt.Println("Task Details Updated.")
	} else {
		fmt.Println("Invalid Task number.")
	}
}

func deleteTask() {
	listTasks()
	fmt.Println("Enter the task number to delete: ")
	var taskNo int
	fmt.Scanln(&taskNo)

	if taskNo > 0 && taskNo <= len(tasks) {
		tasks = append(tasks[:taskNo-1], tasks[taskNo:]...)
		fmt.Println("Task deleted.")
	} else {
		fmt.Println("Invalid task number.")
	}
}

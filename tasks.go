package main

import (
	"bufio"
	"fmt"
	"os"
	//"strconv" -- pointing fmt.Scan at a float64 var appears
	//				to automatically convert it
	"strings"
	"time"
)

type Task struct {
	Description             string     `json:"description"`
	Completed               bool       `json:"completed"`
	TotalHoursRemaining     float64    `json:"total_hours_remaining"`
	TotalHoursCompleted     float64    `json:"total_hours_completed"`
	EstimatedCompletionDate time.Time  `json:"estimated_completion_date"`
	WorkSlots               []WorkSlot `json:"work_slots"`
	SubTasks                []Task     `json:"sub_tasks"`
	ParentTask              *Task      `json:"-"`
	// subTasks might be best as a map. Need to think about the best
	// data structure for this. It needs to hold
	// a. list of all tasks in the "tree task"
	// b. an order that they should be completed in, when applicable
	// 		^^^ Is this really necessary??? Need to think this through
}

type WorkLog struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type WorkSlot struct {
	Day             time.Weekday `json:"scheduled_day"`
	TimeStart       time.Time    `json:"time_start"`
	TimeEnd         time.Time    `json:"time_end"`
	PlannedDuration float64      `json:"planned_duration"`
	ActualDuration  float64      `json:"actual_duration"`
	WorkLogs        []WorkLog    `json:"work_logs"`
}

var tasks []Task

const dateFormat = "2006-01-02" // YYYY-MM-DD

func addTask() {
	fmt.Print("Enter task description: \n")
	reader := bufio.NewReader(os.Stdin)
	descr, _ := reader.ReadString('\n')

	descr = strings.TrimSpace(descr)
	//fmt.Print("Enter estimated hours to complete: ")
	//hoursToCompleteStr, _ := reader.ReadString('\n')
	//hoursToCompleteStr = strings.TrimSpace(hoursToCompleteStr)
	//
	//TotalHoursRemaining, err := strconv.Atoi(hoursToCompleteStr)
	//if err != nil {
	//	fmt.Println("Invalid number, enter in arabic form")
	//	return
	//}
	fmt.Print("Enter Estimated hours to complete: ")
	var TotalHoursRemaining float64
	fmt.Scan(&TotalHoursRemaining)
	tasks = append(tasks,
		Task{
			Description:         descr,
			Completed:           false,
			TotalHoursRemaining: TotalHoursRemaining,
			TotalHoursCompleted: 0,
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
		if task.Completed {
			Complete = "Yes"
		}

		fmt.Printf("%d. %s [Completed: %s || Hours Remaining: %s || Hours Completed:  %.2f || completion date: %d]\n",
			i+1,
			task.Description,
			Complete,
			task.TotalHoursRemaining,
			task.TotalHoursCompleted,
			calculateTaskCompletionDate(&task),
		)
	}
}

func completeTask() {
	listTasks()
	fmt.Println("Enter the task number to mark as completed: ")
	var taskNo int
	_, err := fmt.Scanln(&taskNo)
	if err != nil {
		fmt.Println("Invalid task number.")
		return
	}

	if taskNo > 0 && taskNo <= len(tasks) {
		tasks[taskNo-1].Completed = true
		//updateDueDates(tasks[taskNo-1].DueDate)
		fmt.Println("Task marked as completed.")
	} else {
		fmt.Println("Invalid task number.")
	}
}

//func updateDueDates(completedTaskDate time.Time) {
//  ???????? func needed???????????
//	for i, task := range tasks {
//		if !task.Completed {
//			tasks[i].DueDate = task.DueDate.Add(24 * time.Hour)
//		}
//	}
//	fmt.Println("Due dates for tasks updated.")
//}

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

		//// Update due date -- Due Dates should not be hard coded bu caluclated
		// based on progress
		//fmt.Print("Enter new due date (YYYY-MM-DD) or leave blank to keep current: ")
		//dueDateStr, _ := reader.ReadString('\n')
		//dueDateStr = strings.TrimSpace(dueDateStr)

		// Update Estimated Hours
		fmt.Print("Enter new estimated hours to complete (leave blank to keep current): ")
		var TotalHoursRemaining float64
		n, _ := fmt.Scanf("%f\n", &TotalHoursRemaining)

		if descr != "" {
			tasks[taskNo-1].Description = descr
		}
		// Again, moving to due Dates being calculated rather than
		// hard coded
		//if dueDateStr != "" {
		//	dueDate, err := time.Parse(dateFormat, dueDateStr)
		//	if err != nil {
		//		fmt.Println("Invalid Date Format, task not updated")
		//		return
		//	}
		//	tasks[taskNo-1].DueDate = dueDate
		//}

		if n == 1 {
			tasks[taskNo-1].TotalHoursRemaining = TotalHoursRemaining
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

func calculateTaskCompletionDate(task *Task) time.Time {
	remainingHours := task.TotalHoursRemaining
	now := time.Now()
	var lastSlotEnd time.Time

	for _, slot := range task.WorkSlots {
		if slot.TimeEnd.After(now) {
			remainingHours -= slot.PlannedDuration
			lastSlotEnd = slot.TimeEnd
			if remainingHours <= 0 {
				break
			}
		}
	}

	if remainingHours > 0 {
		lastSlotEnd = lastSlotEnd.Add(time.Duration(
			remainingHours * float64(time.Hour)))
	}
	return lastSlotEnd
}

func calculateParentCompletionDate(task *Task) time.Time {
	var lastSlotEnd time.Time
	for _, subTask := range task.SubTasks {
		subTaskCompletionDate := calculateParentCompletionDate(&subTask)
		if subTaskCompletionDate.After(lastSlotEnd) {
			lastSlotEnd = subTaskCompletionDate
		}
	}
	return lastSlotEnd
}

func updateParentProgress(task *Task) {
	if task.ParentTask != nil {
		parent := task.ParentTask
		parent.TotalHoursRemaining = 0
		parent.TotalHoursCompleted = 0

		// Aggregate hours from subtasks
		for _, subTask := range parent.SubTasks {
			parent.TotalHoursRemaining += subTask.TotalHoursRemaining
			parent.TotalHoursCompleted += subTask.TotalHoursCompleted
		}
	}
}

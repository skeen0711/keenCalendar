package tasks

import (
	"bufio"
	"fmt"
	"os"
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
}

type WorkLog struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type WorkSlot struct { // how do I identify the next open slot when calculating dates
	Day             time.Weekday `json:"scheduled_day"`
	StartDateTime   time.Time    `json:"start_date_time"` // The date the working slot becomes active
	EndDateTime     time.Time    `json:"end_date_time"`   // The date the working slot is no longer in use
	TimeStart       string       `json:"time_start"`      // expect 4 char military time, ex 0900 = 9am, 2200 = 10pm
	TimeEnd         string       `json:"time_end"`        // expect 4 char military time, ex 0900 = 9am, 2200 = 10pm
	PlannedDuration float64      `json:"planned_duration"`
	ActualDuration  float64      `json:"actual_duration"`
	WorkLogs        []WorkLog    `json:"work_logs"`
}

var tasks []Task

//const dateFormat = "2006-01-02" // YYYY-MM-DD

func AddTask() {
	fmt.Print("Enter task description: \n")
	reader := bufio.NewReader(os.Stdin)
	descr, _ := reader.ReadString('\n')

	descr = strings.TrimSpace(descr)

	fmt.Print("Enter Estimated hours to complete: ")
	var TotalHoursRemaining float64
	fmt.Scan(&TotalHoursRemaining)
	estimatedCompletionDate := CalculateTaskCompletionDate(&Task{})

	tasks = append(tasks,
		Task{
			Description:             descr,
			Completed:               false,
			TotalHoursRemaining:     TotalHoursRemaining,
			TotalHoursCompleted:     0,
			EstimatedCompletionDate: estimatedCompletionDate,
			WorkSlots:               []WorkSlot{},
			SubTasks:                []Task{},
			ParentTask:              nil,
		})
	fmt.Println("Task added.")
}

func AddSubTask(parentTask *Task) {
	fmt.Print("Enter subtask description: \n")
	reader := bufio.NewReader(os.Stdin)
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)
	parentTask.SubTasks = append(parentTask.SubTasks,
		Task{
			Description:         description,
			Completed:           false,
			TotalHoursRemaining: 0,
			TotalHoursCompleted: 0,
			ParentTask:          parentTask,
		})
	fmt.Println("Subtask added.")
}

func ListTasks() {
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
			CalculateTaskCompletionDate(&task),
		)
	}
}

func CompleteTask() {
	ListTasks()
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

func UpdateTask() {
	ListTasks()
	fmt.Print("Enter task number to update: ")
	var taskNo int
	fmt.Scan(&taskNo)

	if taskNo > 0 && taskNo <= len(tasks) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter new task description (leave blank to keep current): ")
		descr, _ := reader.ReadString('\n')
		descr = strings.TrimSpace(descr)

		// Update Estimated Hours
		fmt.Print("Enter new estimated hours to complete (leave blank to keep current): ")
		var TotalHoursRemaining float64
		n, _ := fmt.Scanf("%f\n", &TotalHoursRemaining)

		if descr != "" {
			tasks[taskNo-1].Description = descr
		}

		if n == 1 {
			tasks[taskNo-1].TotalHoursRemaining = TotalHoursRemaining
		}

		fmt.Println("Task Details Updated.")
	} else {
		fmt.Println("Invalid Task number.")
	}
}

func DeleteTask() {
	ListTasks()
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

func CalculateTaskCompletionDate(task *Task) time.Time {
	remainingHours := task.TotalHoursRemaining
	now := time.Now()
	var lastSlotEnd time.Time

	// Iterate the tasks assigned workslots and subtract that workslot from the total
	// hours remaining until TotalHoursRemaining is <= 0 or all workslots are used
	// If TotalHoursRemaining < 0, store remainder in new variable leftOverHours.
	// If there are sub tasks associated with the parent, begin recursively call
	// calculateTaskcompletionDate on next subtask first using the remaining hours,
	// then consuming any applicable remaining workslots

	for _, slot := range task.WorkSlots {
		if slot.EndDateTime.After(now) {
			remainingHours -= slot.PlannedDuration
			lastSlotEnd = slot.EndDateTime
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

func CalculateParentCompletionDate(task *Task) time.Time {
	var lastSlotEnd time.Time
	for _, subTask := range task.SubTasks {
		subTaskCompletionDate := CalculateParentCompletionDate(&subTask)
		if subTaskCompletionDate.After(lastSlotEnd) {
			lastSlotEnd = subTaskCompletionDate
		}
	}
	return lastSlotEnd
}

func UpdateParentProgress(task *Task) {
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

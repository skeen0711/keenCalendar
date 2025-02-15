package tasks

import "time"

func GenerateWorkSlots(day time.Weekday, startTime string,
	endTime string, plannedDuration int, durationType string) []WorkSlot {
	var workSlots []WorkSlot
	var endDate time.Time

	// set today as start date and endDate as plannedDuration from today
	now := time.Now()
	if durationType == "days" {
		endDate = now.AddDate(0, 0, plannedDuration)
	} else if durationType == "months" {
		endDate = now.AddDate(0, plannedDuration, 0)
	} else {
		endDate = now.AddDate(plannedDuration, 0, 0)
	}

	// Convert received military times to go time
	startTimeHour := int((startTime[0]-'0')*10 + (startTime[1] - '0'))
	startTimeMinute := int((startTime[2]-'0')*10 + (startTime[3] - '0'))

	endTimeHour := int((endTime[0]-'0')*10 + (endTime[1] - '0'))
	endTimeMinute := int((endTime[2]-'0')*10 + (endTime[3] - '0'))

	current := now
	for current.Weekday() != day {
		current = current.AddDate(0, 0, 1)
	}
	for current.Before(endDate) {
		startDateTime := time.Date(current.Year(), current.Month(), current.Day(), startTimeHour, startTimeMinute, 0, 0, current.Location())
		endDateTime := time.Date(current.Year(), current.Month(), current.Day(), endTimeHour, endTimeMinute, 0, 0, current.Location())

		// Create the WorkSlot
		workSlots = append(workSlots, WorkSlot{
			Day:           day,
			TimeStart:     startTime,
			TimeEnd:       endTime,
			StartDateTime: startDateTime,
			EndDateTime:   endDateTime,
			PlannedDuration: float64(endTimeHour-startTimeHour) + // calculate duration in hours
				float64(endTimeMinute-startTimeMinute)/60,
		})

		// Move to the next occurrence of the same weekday
		current = current.AddDate(0, 0, 7) // Add 7 days
	}

	return workSlots
}

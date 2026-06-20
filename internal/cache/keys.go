package cache

import "fmt"

func TaskByIDKey(id int64) string {
	return fmt.Sprintf("tasks:by-id:%d", id)
}

func TasksListKey() string {
	return "tasks:list"
}

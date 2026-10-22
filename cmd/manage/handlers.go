package main

import (
	"fmt"
	"strconv"

	"github.com/pt-main/manage/manager"
	tap "github.com/pt-main/tap/go"
	"github.com/pt-main/tap/go/color"
)

func Do(p *tap.Parser, s []string) error {
	m := manager.NewManager()
	if err := Load(m); err != nil {
		return err
	}
	theme, ok := m.Tasks[s[0]]
	if !ok {
		return fmt.Errorf("Invalid theme name: %v", s[0])
	}

	_, delete := p.Flags["delete"]

	hasError := false
	tasks := []int{}

	do := func(id int) {
		done := false
		newTasks := []manager.TaskData{}
		for _, task := range theme.Tasks {
			if task.GetId() == id {
				done = true
				task.Do()
				if !delete {
					newTasks = append(newTasks, task)
				}
			} else {
				newTasks = append(newTasks, task)
			}
		}

		theme.Tasks = newTasks

		if !done {
			hasError = true
			tasks = append(tasks, id)
			color.PrintlnColored(
				"[?RD]🔴 Task not found[?RT]: [?RD]Invalid task id[?RT]: %v",
				id)
		}
	}

	if len(s) == 1 {
		tasks := theme.Tasks
		if len(tasks) == 0 {
			color.PrintColored("[?RD]🔴 Tasks not found")
			hasError = true
		} else {
			last := tasks[len(tasks)-1]
			do(last.GetId())
		}
	}

	for _, taskId := range s[1:] {
		id, err := strconv.Atoi(taskId)
		if err != nil {
			return err
		}
		do(id)
	}

	if hasError {
		color.PrintColored("[?RD]🔴 Failed to do tasks[?RT]")
		if len(tasks) != 0 {
			color.PrintlnColored("%v", tasks)
		} else {
			print("\n")
		}
	} else {
		color.PrintlnColored("[?GN]✅ Tasks done![?RT]")
	}
	return Save(m)
}

func Set(p *tap.Parser, s []string) error {
	m := manager.NewManager()
	if err := Load(m); err != nil {
		return err
	}
	switch s[0] {
	case "priority":
		if len(s) != 3 {
			return fmt.Errorf("Invalid args length: need 3, got %v", len(s))
		}
		theme, ok := m.Tasks[s[1]]
		if !ok {
			return fmt.Errorf("Theme not found: Invalid theme name: '%v'", theme)
		}
		priority := s[2]
		theme.Priority = priority
	default:
		return fmt.Errorf("Invalid set command: %v", s[0])
	}
	return Save(m)
}

type listTheme struct {
	text     string
	priority string
}

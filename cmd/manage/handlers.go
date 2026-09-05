package main

import (
	"fmt"
	"strconv"

	"github.com/pt-main/manage/manager"
	"github.com/pt-main/tap"
	"github.com/pt-main/tap/color"
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

	for _, taskId := range s[1:] {
		var id int
		var err error
		id, err = strconv.Atoi(taskId)
		if err != nil {
			return err
		}

		done := false
		newTasks := []*manager.OneTask{}
		for _, task := range theme.Tasks {
			if task.Id == id {
				done = true
				task.Done = true
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

	if hasError {
		color.PrintlnColored("[?RD]🔴 Failed to do tasks[?RT]: %v", tasks)
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

package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/pt-main/manage/manager"
	tap "github.com/pt-main/tap/go"
	"github.com/pt-main/tap/go/color"
)

func (mc *ManageCli) Do(p *tap.Parser, s []string) error {
	m := mc.C
	theme, ok := m.Tasks[s[0]]
	if !ok {
		return fmt.Errorf("Invalid theme name: %v", s[0])
	}

	_, delete := p.Flags["delete"]

	hasError := false
	tasks := []int{}

	do := func(id int) error {
		done := false
		newTasks := []manager.TaskData{}
		for _, task := range theme.Tasks {
			if task.GetId() == id {
				done = true
				if err := task.Do(); err != nil {
					return fmt.Errorf("Can't do task %v: %v", task, err)
				}
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
				"[?RD]❌ Task not found[?RT]: [?RD]Invalid task id[?RT]: %v",
				id)
		}
		return nil
	}

	if len(s) == 1 {
		tasks := theme.Tasks
		if len(tasks) == 0 {
			color.PrintColored("[?RD]❌ Tasks not found")
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
		color.PrintColored("[?RD]❌ Failed to do tasks[?RT]")
		msg := ""
		last := len(tasks) - 1
		for idx, task := range tasks {
			msg += strconv.Itoa(task)
			if idx != last {
				msg += ", "
			}
		}
		return errors.New(msg)
	} else {
		color.PrintlnColored("[?GN]✅ Tasks done![?RT]")
	}
	return nil
}

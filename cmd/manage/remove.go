package main

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/pt-main/manage/manager"
	tap "github.com/pt-main/tap/go"
	"github.com/pt-main/tap/go/color"
)

// RemoveTheme: manage remove theme <name>
func (mc *ManageCli) RemoveTheme(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) != 1 {
		return fmt.Errorf("❌ 'remove theme' needs exactly 1 argument (name), got %v", len(s))
	}
	name := s[0]

	if _, ok := m.Tasks[name]; !ok {
		return fmt.Errorf("❌ Invalid theme name: %v", name)
	}
	delete(m.Tasks, name)

	for tableName, themes := range m.Tables {
		newThemes := []string{}
		removed := false
		for _, t := range themes {
			if t == name {
				removed = true
				continue
			}
			newThemes = append(newThemes, t)
		}
		if removed {
			m.Tables[tableName] = newThemes
		}
	}

	return nil
}

// RemoveTask: manage remove task <theme> <id...>
func (mc *ManageCli) RemoveTask(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) < 2 {
		return fmt.Errorf("❌ 'remove task' needs at least 2 arguments (theme + ids), got %v", len(s))
	}
	name := s[0]

	theme, ok := m.Tasks[name]
	if !ok {
		return fmt.Errorf("❌ Invalid theme name: %v", name)
	}

	ids := make(map[int]bool)
	for _, idStr := range s[1:] {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("❌ Invalid task id '%v': %w", idStr, err)
		}
		ids[id] = true
	}

	newTasks := []manager.TaskData{}
	removed := 0
	for _, task := range theme.Tasks {
		if ids[task.GetId()] {
			removed++
			continue
		}
		newTasks = append(newTasks, task)
	}
	theme.Tasks = newTasks

	if removed == 0 {
		return fmt.Errorf("❌ No tasks removed (none of the given ids exist in theme '%v')", name)
	}
	if removed < len(ids) {
		color.PrintlnColored("[?YW]❗️ Some ids were not found in theme: %v", name)
	}

	return nil
}

// RemoveTag: manage remove tag <theme> <tags...>
func (mc *ManageCli) RemoveTag(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) < 2 {
		return fmt.Errorf("❌ 'remove tag' needs at least 2 arguments (theme + tags), got %v", len(s))
	}
	name := s[0]
	toRemove := s[1:]

	theme, ok := m.Tasks[name]
	if !ok {
		return fmt.Errorf("❌ Invalid theme name: %v", name)
	}

	newTags := []string{}
	removedAny := false
	for _, t := range theme.Tags {
		if slices.Contains(toRemove, t) {
			removedAny = true
			continue
		}
		newTags = append(newTags, t)
	}
	if !removedAny {
		return fmt.Errorf("❌ No tags removed (none of the given tags exist on theme '%v')", name)
	}
	theme.Tags = newTags

	return nil
}

// RemoveTable: manage remove table <name>
func (mc *ManageCli) RemoveTable(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) != 1 {
		return fmt.Errorf("❌ 'remove table' needs exactly 1 argument (name), got %v", len(s))
	}
	name := s[0]

	if _, ok := m.Tables[name]; !ok {
		return fmt.Errorf("❌ Invalid table name: %v", name)
	}
	delete(m.Tables, name)

	return nil
}

// RemoveFromTable: manage remove from-table <name> <themes...>
func (mc *ManageCli) RemoveFromTable(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) < 2 {
		return fmt.Errorf("❌ 'remove from-table' needs at least 2 arguments (name + themes), got %v", len(s))
	}
	name := s[0]
	toRemove := s[1:]

	themes, ok := m.Tables[name]
	if !ok {
		return fmt.Errorf("❌ Invalid table name: %v", name)
	}

	newThemes := []string{}
	removedAny := false
	for _, t := range themes {
		if slices.Contains(toRemove, t) {
			removedAny = true
			continue
		}
		newThemes = append(newThemes, t)
	}
	if !removedAny {
		return fmt.Errorf("❌ No themes removed (none of the given names are in table '%v')", name)
	}
	m.Tables[name] = newThemes

	return nil
}

// RemovePriority: manage remove priority <p1> [p2...]
func (mc *ManageCli) RemovePriority(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) < 1 {
		return fmt.Errorf("❌ 'remove priority' needs at least 1 argument, got %v", len(s))
	}

	newOrder := []string{}
	for _, pr := range m.Settings.PriorityOrder {
		if !slices.Contains(s, pr) {
			newOrder = append(newOrder, pr)
		}
	}
	if len(newOrder) == 0 {
		return fmt.Errorf("❌ Can't remove all priorities (at least one must remain)")
	}
	if len(newOrder) == len(m.Settings.PriorityOrder) {
		return fmt.Errorf("❌ No priorities removed (none of the given values exist)")
	}
	m.Settings.PriorityOrder = newOrder

	return nil
}

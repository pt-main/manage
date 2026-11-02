package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/pt-main/manage/manager"
	tap "github.com/pt-main/tap/go"
)

// AddTheme: manage add theme <name> [--priority=A] [--tags="go api"] [--force]
func (mc *ManageCli) AddTheme(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) != 1 {
		return fmt.Errorf("❌ 'add theme' needs exactly 1 argument (name), got %v", len(s))
	}
	name := s[0]

	_, force := p.Flags["force"]

	_tags, ok := p.Flags["tags"]
	tags := []string{}
	if ok && _tags != "" {
		tags = strings.Split(_tags, " ")
	}

	priority, ok := p.Flags["priority"]
	if !ok || priority == "" {
		priority = m.Settings.DefaultPriority
	}

	t := manager.NewTask(tags, priority, name)
	if err := m.AddTheme(name, t, force); err != nil {
		return err
	}
	return nil
}

// AddTask: manage add task <theme> <text> [--id=N] [--force]
func (mc *ManageCli) AddTask(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) != 2 {
		return fmt.Errorf("❌ 'add task' needs exactly 2 arguments (theme, text), got %v", len(s))
	}
	name, text := s[0], s[1]
	_, force := p.Flags["force"]

	theme, ok := m.Tasks[name]
	if !ok {
		return fmt.Errorf("❌ Invalid theme name: %v", name)
	}

	appendMode := true
	var to = theme.Last
	if _to, ok := p.Flags["id"]; ok {
		v, err := strconv.Atoi(_to)
		if err != nil {
			return fmt.Errorf("❌ Invalid --id value '%v': %w", _to, err)
		}
		to = v
		appendMode = false
	}

	switch {
	case appendMode && !force:
		if err := m.AddTask(name, text); err != nil {
			return err
		}
	case !appendMode && force:
		found := false
		for _, task := range theme.Tasks {
			if task.GetId() == to {
				found = true
				if !task.SetAbout(text) {
					return fmt.Errorf("❌ Can't edit '%v' task about", task.Type())
				}
			}
		}
		if !found {
			return fmt.Errorf("❌ Invalid id: %v", to)
		}
	case !appendMode && !force:
		found := false
		for _, task := range theme.Tasks {
			if task.GetId() == to {
				found = true
				if !task.Add(text) {
					return fmt.Errorf("❌ Can't edit '%v' task about", task.Type())
				}
			}
		}
		if !found {
			return fmt.Errorf("❌ Invalid id: %v", to)
		}
	case appendMode && force:
		return fmt.Errorf("❌ --force requires --id")
	}

	return nil
}

// AddTag: manage add tag <theme> <tags...> ]
func (mc *ManageCli) AddTag(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) < 2 {
		return fmt.Errorf("❌ 'add tag' needs at least 2 arguments (theme + tags), got %v", len(s))
	}
	name := s[0]
	tags := s[1:]

	theme, ok := m.Tasks[name]
	if !ok {
		return fmt.Errorf("❌ Invalid theme name: %v", name)
	}

	newTags := slices.Clone(theme.Tags)
	for _, t := range tags {
		if !slices.Contains(newTags, t) {
			newTags = append(newTags, t)
		}
	}

	theme.Tags = newTags
	return nil
}

// AddTable: manage add table <name> <themes...> [--force]
func (mc *ManageCli) AddTable(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) < 2 {
		return fmt.Errorf("❌ 'add table' needs at least 2 arguments (name + themes), got %v", len(s))
	}
	name := s[0]
	themes := s[1:]
	_, force := p.Flags["force"]

	t, ok := m.Tables[name]
	if !ok {
		force = true
	}

	if force {
		m.Tables[name] = themes
	} else {
		t = append(t, themes...)
		m.Tables[name] = t
	}

	return nil
}

// AddPriority: manage add priority <p1> [p2 ...]
func (mc *ManageCli) AddPriority(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) < 1 {
		return fmt.Errorf("❌ 'add priority' needs at least 1 argument, got %v", len(s))
	}
	for _, pr := range s {
		if !slices.Contains(m.Settings.PriorityOrder, pr) {
			m.Settings.PriorityOrder = append(m.Settings.PriorityOrder, pr)
		}
	}
	return nil
}

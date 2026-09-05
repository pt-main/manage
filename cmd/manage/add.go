package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/pt-main/manage/manager"
	"github.com/pt-main/tap"
)

func Add(p *tap.Parser, s []string) error {
	m := manager.NewManager()
	if err := Load(m); err != nil {
		return err
	}
	name := s[1]
	_, force := p.Flags["force"]

	lenArgs := func(n int) error {
		if len(s) != n {
			return fmt.Errorf(
				"❌ Invalid argument length for add: need %v, got %v",
				n, len(s))
		}
		return nil
	}

	switch s[0] {
	case "theme":
		if err := lenArgs(2); err != nil {
			return err
		}
		_tags, ok := p.Flags["tags"]
		tags := strings.Split(_tags, " ")
		if !ok {
			tags = make([]string, 0)
		}
		priority, ok := p.Flags["priority"]
		if priority == "" {
			priority = m.Settings.DefaultPriority
		}
		t := manager.NewTask(tags, priority, name)
		if err := m.AddTheme(name, t, force); err != nil {
			return err
		}
		return Save(m)

	case "task":
		if err := lenArgs(3); err != nil {
			return err
		}

		text := s[2]

		task, ok := m.Tasks[name]
		if !ok {
			return fmt.Errorf("❌ Invalid theme name: %v", name)
		}

		append := true
		var to = task.Last
		_to, ok := p.Flags["id"]
		if ok {
			var err error
			to, err = strconv.Atoi(_to)
			if err != nil {
				return err
			}
			append = false
		}

		if append && !force {
			if err := m.AddTask(name, text); err != nil {
				return err
			}
		} else if !append && force {
			added := false
			for _, task := range task.Tasks {
				if task.Id == to {
					added = true
					task.About = text
				}
			}
			if !added {
				return fmt.Errorf("❌ Invalid id: %v", to)
			}
		} else if !append && !force {
			added := false
			for _, task := range task.Tasks {
				if task.Id == to {
					added = true
					task.About += text
				}
			}
			if !added {
				return fmt.Errorf("❌ Invalid id: %v", to)
			}
		}

	case "tag":
		if len(s) < 3 {
			return fmt.Errorf(
				"❌ Invalid argument format for 'add tag': need 3 or more arguments, got %v",
				len(s))
		}
		tags := s[2:]

		task, ok := m.Tasks[name]
		if !ok {
			return fmt.Errorf("❌ Invalid theme name: %v", name)
		}

		add := []string{}
		rm := []string{}
		for _, tag := range tags {
			if strings.HasPrefix(tag, "!") {
				rm = append(rm, tag[1:])
			} else {
				add = append(add, tag)
			}
		}

		newTags := []string{}
		for _, tag := range task.Tags {
			if !slices.Contains(rm, tag) {
				newTags = append(newTags, tag)
			}
		}

		for _, addtag := range add {
			if !slices.Contains(newTags, addtag) {
				newTags = append(newTags, addtag)
			}
		}

		task.Tags = newTags

	case "table":
		if len(s) < 3 {
			return fmt.Errorf(
				"❌ Invalid argument format for 'add table': need 3 or more arguments, got %v",
				len(s))
		}

		name := s[1]
		themes := s[2:]
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

	default:
		return fmt.Errorf("❌ Invalid mode: %v", s[0])
	}
	return Save(m)
}

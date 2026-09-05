package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/pt-main/manage/manager"
	"github.com/pt-main/tap"
	"github.com/pt-main/tap/color"
	"github.com/pt-main/tap/utils"
)

func List(p *tap.Parser, s []string) error {
	m := manager.NewManager()
	if err := Load(m); err != nil {
		return err
	}

	_table, has_table := p.Flags["table"]
	var ok bool
	table := []string{}
	if has_table {
		table, ok = m.Tables[_table]
		if !ok {
			return fmt.Errorf("Table is not found: Invalid table name: '%v'", _table)
		}
	}
	_tags, has_tags := p.Flags["tags"]
	tags := []string{}
	if has_tags {
		tags = strings.Split(_tags, " ")
	}
	_priority, has_priority := p.Flags["priority"]
	priority := []string{}
	if has_priority {
		priority = strings.Split(_priority, " ")
	}
	_themes, has_themes := p.Flags["themes"]
	if has_themes {
		table = append(table, strings.Split(_themes, " ")...)
	}

	filter, _ := p.Flags["filter"]
	state, _ := p.Flags["state"]

	if state == "" {
		_, todo := p.Flags["todo"]
		_, done := p.Flags["done"]
		if todo && !done {
			state = manager.StateTodo
		} else if done && !todo {
			state = manager.StateDone
		} else if done && todo {
			new := manager.StateAll
			color.PrintlnColored("[?RD]Can't use todo and done flags in one filter. State = %v", new)
			state = new
		}
	}

	filtered := m.Filter(tags, priority, table, state, filter)

	textRaw := map[string]*listTheme{}

	for _, theme := range filtered {
		tasks := []string{}
		for _, task := range theme.Tasks {
			state := "todo"
			if task.Done {
				state = "done"
			}
			tasks = append(tasks, fmt.Sprintf(
				"[?RT]- [?GN]%v[?BBK] ([?GN]%v[?RT][?BBK]): [[?YW]%v[?BBK]]:[?RT]\n  %v",
				task.Id, task.CreatedAt,
				state, task.About,
			))
		}
		res := color.Set(strings.Join(tasks, "\n"))
		withoutColors := color.ReplaceColors(strings.Join(tasks, "\n"))
		if strings.TrimSpace(withoutColors) != "" {
			textRaw[theme.Name] = &listTheme{
				text:     res,
				priority: theme.Priority,
			}
		} else {
			color.PrintlnColored("[?YW]🟡 Empty theme: %v", theme.Name)
		}
	}

	outed := 0

	for _, priority := range m.Settings.PriorityOrder {
		text := map[string]string{}
		for name, theme := range textRaw {
			if theme.priority == priority {
				text[fmt.Sprintf("[?YW]%v[?RT]", name)] = theme.text
			}
		}
		if len(slices.Collect(maps.Keys(text))) != 0 {
			outed += 1
			color.PrintlnColored(utils.FramedTextMap(
				"GN", fmt.Sprintf("[?BGN]Priority %v[?RT]", priority), "", text,
			))
		}
	}

	if outed == 0 {
		color.PrintlnColored("[?RD]❗️ Has no themes[?RT]")
	}

	return nil
}

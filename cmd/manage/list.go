package main

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/pt-main/manage/manager"
	tap "github.com/pt-main/tap/go"
	"github.com/pt-main/tap/go/color"
	"github.com/pt-main/tap/go/utils"
)

type listTheme struct {
	text     string
	priority string
}

func (mc *ManageCli) List(p *tap.Parser, s []string) error {
	m := mc.C
	mc.NoSave = true

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
	for _, pr := range priority {
		if !slices.Contains(m.Settings.PriorityOrder, pr) {
			return fmt.Errorf("Priority is not found: %v", pr)
		}
	}

	_themes, has_themes := p.Flags["themes"]
	if has_themes {
		table = append(slices.Clone(table), strings.Split(_themes, " ")...)
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

	space := func(s string) string {
		return "|  " + strings.ReplaceAll(s, "\n", "\n|  ")
	}

	if _, has := p.Flags["code"]; has {
		for _, task := range manager.FindTodosInDir(".") {
			textRaw["[CODE]:"+strconv.Itoa(task.Id)] = &listTheme{
				priority: m.Settings.PriorityOrder[0],
				text:     space(task.Render(color.ColorEnabled)),
			}
		}
	}

	for _, theme := range filtered {
		tasks := map[string]string{}
		raw := []string{}
		for _, task := range theme.Tasks {
			// text := utils.FramedTextMap(
			// 	"CN", color.Set("[?BCN]"+string(task.Type())+" [?CN]TASK[?RT]"),
			// 	"", map[string]string{"": task.Render(color.ColorEnabled)},
			// )
			text := space(task.Render(color.ColorEnabled))
			tasks[strconv.Itoa(task.GetId())] = text
			raw = append(raw, text)
		}
		res := utils.FramedTextMap("CN", "TASKS", "", tasks)
		withoutColors := color.ReplaceColors(strings.Join(raw, "\n"))
		if strings.TrimSpace(withoutColors) != "" {
			textRaw[theme.Name] = &listTheme{
				text:     res,
				priority: theme.Priority,
			}
		} else {
			color.PrintlnColored("[?YW]❗️ Empty theme: %v", theme.Name)
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

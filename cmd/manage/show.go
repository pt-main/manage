package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/pt-main/manage/manager"
	tap "github.com/pt-main/tap/go"
	"github.com/pt-main/tap/go/color"
)

// ---------- show themes ----------

func (mc *ManageCli) ShowThemes(p *tap.Parser, s []string) error {
	m := mc.C
	mc.NoSave = true

	if len(m.Tasks) == 0 {
		color.PrintlnColored("[?RD]❗️ Has no themes[?RT]")
		return nil
	}

	themes := make([]*manager.Task, 0, len(m.Tasks))
	for _, t := range m.Tasks {
		themes = append(themes, t)
	}
	slices.SortFunc(themes, func(a, b *manager.Task) int {
		ai, bi := priorityIndex(m, a.Priority), priorityIndex(m, b.Priority)
		if ai != bi {
			return ai - bi
		}
		return strings.Compare(a.Name, b.Name)
	})

	color.PrintlnColored("[?BGN]Themes:[?RT]")
	fmt.Println()
	for _, theme := range themes {
		done := 0
		for _, t := range theme.Tasks {
			if t.IsDone() {
				done++
			}
		}
		color.PrintlnColored(
			"  [?YW]%v[?RT] [?BBK][%v][?RT] - [?BGN]%v[?RT] tasks ([?GN]%v[?RT] done, [?RD]%v[?RT] todo) %v",
			theme.Name, theme.Priority,
			len(theme.Tasks), done, len(theme.Tasks)-done,
			formatTags(theme.Tags),
		)
	}
	return nil
}

// ---------- show priorities ----------

func (mc *ManageCli) ShowPriorities(p *tap.Parser, s []string) error {
	m := mc.C
	mc.NoSave = true

	byPrio := map[string][]string{}
	for name, theme := range m.Tasks {
		byPrio[theme.Priority] = append(byPrio[theme.Priority], name)
	}
	for _, v := range byPrio {
		slices.Sort(v)
	}

	color.PrintlnColored("[?BGN]Priorities (in priority_order):[?RT]")
	fmt.Println()
	for _, prio := range m.Settings.PriorityOrder {
		themes, ok := byPrio[prio]
		if !ok || len(themes) == 0 {
			color.PrintlnColored("  [?YW]%v[?RT] [?BBK](no themes)[?RT]", prio)
			continue
		}
		color.PrintlnColored(
			"  [?YW]%v[?RT] - [?BGN]%v[?RT] theme(s): %v",
			prio, len(themes), strings.Join(themes, ", "),
		)
	}

	unknown := []string{}
	for prio := range byPrio {
		if !slices.Contains(m.Settings.PriorityOrder, prio) {
			unknown = append(unknown, prio)
		}
	}
	if len(unknown) > 0 {
		slices.Sort(unknown)
		fmt.Println()
		color.PrintlnColored("[?BGN]Unknown priorities (not in priority_order):[?RT]")
		fmt.Println()
		for _, prio := range unknown {
			color.PrintlnColored(
				"  [?RD]%v[?RT] - [?BGN]%v[?RT] theme(s): %v",
				prio, len(byPrio[prio]), strings.Join(byPrio[prio], ", "),
			)
		}
	}

	return nil
}

// ---------- show tags ----------

func (mc *ManageCli) ShowTags(p *tap.Parser, s []string) error {
	m := mc.C
	mc.NoSave = true

	tagThemes := map[string][]string{}
	for name, theme := range m.Tasks {
		for _, tag := range theme.Tags {
			tagThemes[tag] = append(tagThemes[tag], name)
		}
	}

	if len(tagThemes) == 0 {
		color.PrintlnColored("[?RD]❗️ Has no tags[?RT]")
		return nil
	}

	tags := slices.Collect(maps.Keys(tagThemes))
	slices.Sort(tags)

	color.PrintlnColored("[?BGN]Tags:[?RT]")
	fmt.Println()
	for _, tag := range tags {
		themes := tagThemes[tag]
		slices.Sort(themes)
		color.PrintlnColored(
			"  [?CN]#%v[?RT] - [?BGN]%v[?RT] theme(s): %v",
			tag, len(themes), strings.Join(themes, ", "),
		)
	}
	return nil
}

// ---------- show theme <name> ----------

func (mc *ManageCli) ShowTheme(p *tap.Parser, s []string) error {
	m := mc.C
	mc.NoSave = true

	if len(s) != 1 {
		return fmt.Errorf("❌ 'show theme' needs exactly 1 argument (theme name), got %v", len(s))
	}
	theme, ok := m.Tasks[s[0]]
	if !ok {
		return fmt.Errorf("❌ Invalid theme name: %v", s[0])
	}

	total := len(theme.Tasks)
	doneCount := 0
	for _, t := range theme.Tasks {
		if t.IsDone() {
			doneCount++
		}
	}

	color.PrintlnColored("[?BGN]--- Theme: %v ---[?RT]", theme.Name)
	fmt.Println()
	color.PrintlnColored("[?YW]Priority:[?RT] [?BGN]%v[?RT]", theme.Priority)
	color.PrintlnColored("[?YW]Tags:[?RT]     %v", formatTags(theme.Tags))
	color.PrintlnColored(
		"[?YW]Tasks:[?RT]    [?BGN]%v[?RT] total, [?GN]%v[?RT] done, [?RD]%v[?RT] todo",
		total, doneCount, total-doneCount,
	)
	color.PrintlnColored("[?YW]Last ID:[?RT]  [?BGN]%v[?RT]", theme.Last)

	if inTables := tablesForTheme(m, theme.Name); len(inTables) > 0 {
		color.PrintlnColored("[?YW]Tables:[?RT]   %v", strings.Join(inTables, ", "))
	}
	fmt.Println()

	return nil
}

// ---------- show info ----------

func (mc *ManageCli) ShowInfo(p *tap.Parser, s []string) error {
	m := mc.C
	mc.NoSave = true

	if len(m.Tasks) == 0 {
		color.PrintlnColored("[?RD]❗️ Has no themes[?RT]")
		return nil
	}

	byPriority := map[string][]*manager.Task{}
	for _, theme := range m.Tasks {
		byPriority[theme.Priority] = append(byPriority[theme.Priority], theme)
	}
	for _, list := range byPriority {
		slices.SortFunc(list, func(a, b *manager.Task) int {
			return strings.Compare(a.Name, b.Name)
		})
	}

	totalThemes, totalTasks, totalDone := 0, 0, 0

	printGroup := func(priority string, themes []*manager.Task) {
		color.PrintlnColored("[?BGN]--- Priority %v ---[?RT]", priority)
		fmt.Println()
		for _, theme := range themes {
			totalThemes++
			done := 0
			for _, t := range theme.Tasks {
				if t.IsDone() {
					done++
				}
			}
			totalTasks += len(theme.Tasks)
			totalDone += done
			color.PrintlnColored(
				"  [?YW]%v[?RT] - [?BGN]%v[?RT] tasks ([?GN]%v[?RT] done, [?RD]%v[?RT] todo) %v",
				theme.Name, len(theme.Tasks), done, len(theme.Tasks)-done,
				formatTags(theme.Tags),
			)
		}
		fmt.Println()
	}

	shown := map[string]bool{}
	for _, priority := range m.Settings.PriorityOrder {
		list, ok := byPriority[priority]
		if !ok || len(list) == 0 {
			continue
		}
		shown[priority] = true
		printGroup(priority, list)
	}

	remaining := []string{}
	for priority := range byPriority {
		if !shown[priority] && len(byPriority[priority]) > 0 {
			remaining = append(remaining, priority)
		}
	}
	slices.Sort(remaining)
	for _, priority := range remaining {
		printGroup(priority+" (not in priority_order)", byPriority[priority])
	}

	color.PrintlnColored(
		"[?BGN]Total:[?RT] [?BGN]%v[?RT] themes, [?BGN]%v[?RT] tasks ([?GN]%v[?RT] done, [?RD]%v[?RT] todo)",
		totalThemes, totalTasks, totalDone, totalTasks-totalDone,
	)
	return nil
}

// ---------- helpers ----------

func priorityIndex(m *manager.Manager, priority string) int {
	if i := slices.Index(m.Settings.PriorityOrder, priority); i != -1 {
		return i
	}
	return len(m.Settings.PriorityOrder)
}

func tablesForTheme(m *manager.Manager, name string) []string {
	result := []string{}
	for tableName, themes := range m.Tables {
		if slices.Contains(themes, name) {
			result = append(result, tableName)
		}
	}
	slices.Sort(result)
	return result
}

func formatTags(tags []string) string {
	if len(tags) == 0 {
		return "[?BBK]<no tags>[?RT]"
	}
	parts := make([]string, len(tags))
	for i, t := range tags {
		parts[i] = fmt.Sprintf("[?CN]#%v[?RT]", t)
	}
	return strings.Join(parts, " ")
}

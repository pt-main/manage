package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/pt-main/manage/manager"
	tap "github.com/pt-main/tap/go"
	"github.com/pt-main/tap/go/color"
	"github.com/pt-main/tap/go/utils"
)

// ---------- show themes ----------

func (mc *ManageCli) ShowThemes(p *tap.Parser, s []string) error {
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

	printGroup := func(priority string, themes []*manager.Task, known bool) {
		text := map[string]string{}
		for _, theme := range themes {
			text[fmt.Sprintf("[?YW]%v[?RT]", theme.Name)] = themeSummary(theme)
		}
		title := fmt.Sprintf("[?BGN]Priority %v[?RT]", priority)
		colorCode := "GN"
		if !known {
			title = fmt.Sprintf("[?RD]Priority %v (not in priority_order)[?RT]", priority)
			colorCode = "RD"
		}
		color.PrintlnColored(utils.FramedTextMap(colorCode, title, "", text))
	}

	shown := map[string]bool{}
	for _, priority := range m.Settings.PriorityOrder {
		list, ok := byPriority[priority]
		if !ok || len(list) == 0 {
			continue
		}
		shown[priority] = true
		printGroup(priority, list, true)
	}

	remaining := []string{}
	for priority := range byPriority {
		if !shown[priority] && len(byPriority[priority]) > 0 {
			remaining = append(remaining, priority)
		}
	}
	slices.Sort(remaining)
	for _, priority := range remaining {
		printGroup(priority, byPriority[priority], false)
	}

	return nil
}

// ---------- show priorities ----------

func (mc *ManageCli) ShowPriorities(p *tap.Parser, s []string) error {
	m := mc.C
	mc.NoSave = true

	byPriority := map[string][]string{}
	for name, theme := range m.Tasks {
		byPriority[theme.Priority] = append(byPriority[theme.Priority], name)
	}
	for _, v := range byPriority {
		slices.Sort(v)
	}

	order := m.Settings.PriorityOrder

	type row struct {
		label  string
		themes []string
		empty  bool
	}
	rows := []row{}
	for i := 0; i < len(order); {
		themes := byPriority[order[i]]
		if len(themes) > 0 {
			rows = append(rows, row{label: order[i], themes: themes})
			i++
			continue
		}
		j := i
		for j < len(order) && len(byPriority[order[j]]) == 0 {
			j++
		}
		rows = append(rows, row{label: collapseRange(order[i:j]), empty: true})
		i = j
	}

	if len(rows) > 0 {
		width := 0
		for _, r := range rows {
			if w := utf8.RuneCountInString(r.label); w > width {
				width = w
			}
		}
		lines := make([]string, 0, len(rows))
		for _, r := range rows {
			label := padLabel(r.label, width)
			if r.empty {
				lines = append(lines, fmt.Sprintf(
					"[?YW]%s[?RT]  [?BBK](no themes)[?RT]", label,
				))
			} else {
				lines = append(lines, fmt.Sprintf(
					"[?YW]%s[?RT]  [?BGN]%d[?RT] theme(s): %s",
					label, len(r.themes), strings.Join(r.themes, ", "),
				))
			}
		}
		color.PrintlnColored(utils.FramedTextMap(
			"GN",
			"[?BGN]Priority order[?RT]",
			"",
			map[string]string{
				fmt.Sprintf("[?BBK]%d priorities[?RT]", len(order)): strings.Join(lines, "\n"),
			},
		))
	}

	unknown := []string{}
	for priority := range byPriority {
		if !slices.Contains(order, priority) {
			unknown = append(unknown, priority)
		}
	}
	if len(unknown) > 0 {
		slices.Sort(unknown)
		width := 0
		for _, p := range unknown {
			if w := utf8.RuneCountInString(p); w > width {
				width = w
			}
		}
		lines := make([]string, 0, len(unknown))
		for _, priority := range unknown {
			themes := byPriority[priority]
			label := padLabel(priority, width)
			lines = append(lines, fmt.Sprintf(
				"[?RD]%s[?RT]  [?BGN]%d[?RT] theme(s): %s",
				label, len(themes), strings.Join(themes, ", "),
			))
		}
		fmt.Println()
		color.PrintlnColored(utils.FramedTextMap(
			"RD",
			"[?RD]Unknown priorities[?RT]",
			"[?BBK]not in priority_order[?RT]",
			map[string]string{
				fmt.Sprintf("[?RD]%d[?RT]", len(unknown)): strings.Join(lines, "\n"),
			},
		))
	}

	return nil
}

// collapseRange formats a run of consecutive empty priorities.
//   - single element        -> "high"
//   - single-letter bounds  -> "A-Z"
//   - otherwise             -> "high..low"
func collapseRange(ps []string) string {
	if len(ps) == 1 {
		return ps[0]
	}
	first, last := ps[0], ps[len(ps)-1]
	if isSingleLetter(first) && isSingleLetter(last) {
		return first + "-" + last
	}
	return first + ".." + last
}

func isSingleLetter(s string) bool {
	if len(s) != 1 {
		return false
	}
	c := s[0]
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

// padLabel pads by runes, not bytes, so Cyrillic/CJK priorities stay aligned.
func padLabel(s string, width int) string {
	n := utf8.RuneCountInString(s)
	if n >= width {
		return s
	}
	return s + strings.Repeat(" ", width-n)
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

	text := map[string]string{}
	for _, tag := range tags {
		themes := tagThemes[tag]
		slices.Sort(themes)
		text[fmt.Sprintf("[?CN]#%v[?RT]", tag)] = fmt.Sprintf(
			"[?BGN]%v[?RT] theme(s): %v",
			len(themes), strings.Join(themes, ", "),
		)
	}

	color.PrintlnColored(utils.FramedTextMap("CN", "[?BGN]Tags[?RT]", "", text))
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

	done := 0
	for _, t := range theme.Tasks {
		if t.IsDone() {
			done++
		}
	}

	text := map[string]string{
		"[?YW]Priority[?RT]": fmt.Sprintf("[?BGN]%v[?RT]", theme.Priority),
		"[?YW]Tags[?RT]":     formatTags(theme.Tags),
		"[?YW]Tasks[?RT]": fmt.Sprintf(
			"[?BGN]%v[?RT] total, [?GN]%v[?RT] done, [?RD]%v[?RT] todo",
			len(theme.Tasks), done, len(theme.Tasks)-done,
		),
		"[?YW]Last ID[?RT]": fmt.Sprintf("[?BGN]%v[?RT]", theme.Last),
	}

	if inTables := tablesForTheme(m, theme.Name); len(inTables) > 0 {
		text["[?YW]Tables[?RT]"] = strings.Join(inTables, ", ")
	}

	color.PrintlnColored(utils.FramedTextMap(
		"CN", fmt.Sprintf("[?BGN]Theme: %v[?RT]", theme.Name), "", text,
	))
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

	printGroup := func(priority string, themes []*manager.Task, known bool) {
		text := map[string]string{}
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
			text[fmt.Sprintf("[?YW]%v[?RT]", theme.Name)] = fmt.Sprintf(
				"[?BGN]%v[?RT] tasks ([?GN]%v[?RT] done, [?RD]%v[?RT] todo) %v",
				len(theme.Tasks), done, len(theme.Tasks)-done,
				formatTags(theme.Tags),
			)
		}
		title := fmt.Sprintf("[?BGN]Priority %v[?RT]", priority)
		colorCode := "GN"
		if !known {
			title = fmt.Sprintf("[?RD]Priority %v (not in priority_order)[?RT]", priority)
			colorCode = "RD"
		}
		color.PrintlnColored(utils.FramedTextMap(colorCode, title, "", text))
	}

	shown := map[string]bool{}
	for _, priority := range m.Settings.PriorityOrder {
		list, ok := byPriority[priority]
		if !ok || len(list) == 0 {
			continue
		}
		shown[priority] = true
		printGroup(priority, list, true)
	}

	remaining := []string{}
	for priority := range byPriority {
		if !shown[priority] && len(byPriority[priority]) > 0 {
			remaining = append(remaining, priority)
		}
	}
	slices.Sort(remaining)
	for _, priority := range remaining {
		printGroup(priority, byPriority[priority], false)
	}

	totals := map[string]string{
		"[?YW]Themes[?RT]": fmt.Sprintf("[?BGN]%v[?RT]", totalThemes),
		"[?YW]Tasks[?RT]": fmt.Sprintf(
			"[?BGN]%v[?RT] ([?GN]%v[?RT] done, [?RD]%v[?RT] todo)",
			totalTasks, totalDone, totalTasks-totalDone,
		),
	}
	color.PrintlnColored(utils.FramedTextMap("GN", "[?BGN]Total[?RT]", "", totals))
	return nil
}

// ---------- helpers ----------

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

func themeSummary(theme *manager.Task) string {
	done := 0
	for _, t := range theme.Tasks {
		if t.IsDone() {
			done++
		}
	}
	return fmt.Sprintf(
		"[?BGN]%v[?RT] tasks ([?GN]%v[?RT] done, [?RD]%v[?RT] todo) %v",
		len(theme.Tasks), done, len(theme.Tasks)-done,
		formatTags(theme.Tags),
	)
}

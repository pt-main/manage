package main

import (
	"fmt"
	"slices"
	"strings"

	tap "github.com/pt-main/tap/go"
	"github.com/pt-main/tap/go/color"
)

// SetPriority: manage set priority <theme> <value>
func (mc *ManageCli) SetPriority(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) != 2 {
		return fmt.Errorf("❌ 'set priority' needs exactly 2 arguments (theme, value), got %v", len(s))
	}
	name, priority := s[0], s[1]

	if priority == "" {
		return fmt.Errorf("❌ priority can't be empty")
	}
	if strings.ContainsRune(priority, ' ') {
		return fmt.Errorf("❌ priority can't contain spaces: %q", priority)
	}

	theme, ok := m.Tasks[name]
	if !ok {
		return fmt.Errorf("❌ Invalid theme name: %v", name)
	}

	theme.Priority = priority

	if !slices.Contains(m.Settings.PriorityOrder, priority) {
		color.PrintlnColored(
			"[?YW]⚠ priority '%v' is not in priority_order - theme won't appear in 'list'[?RT]",
			priority,
		)
		color.PrintlnColored(
			"[?BBK]add it with: [?YW]manage add priority %v[?RT]",
			priority,
		)
	}
	return nil
}

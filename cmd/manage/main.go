package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/pt-main/manage/config"
	"github.com/pt-main/manage/manager"
	tap "github.com/pt-main/tap/go"
	"github.com/pt-main/tap/go/utils"
)

var Version = "1.4.0"

const docInit = `[?GN]Initialize[?RT] config file.

[?BBE]Usage:[?RT]
  [?YW]manage init[?RT]

[?BBE]Description:[?RT]
  Creates [?BBK]~/manager_config.tess[?RT] with default settings.
  This is the first command you need to run before using any other commands.
  Warning: overwrites an existing config file.`

const docAddTheme = `[?GN]Add a new theme.[?RT]

[?BBE]Usage:[?RT]
  [?YW]manage add theme[?RT] [?BGN]<name>[?RT] [?BBK]--priority=[?BGN]<letter>[?BBK] [?BBK]--tags=[?BGN]"tag1 tag2 ..."[?RT] [?BBK]--force[?RT]

[?BBE]Flags:[?RT]
  --priority  priority (default = settings.default_priority)
  --tags      space-separated tags
  --force     replace existing theme with the same name`

const docAddTask = `[?GN]Add or edit a task in a theme.[?RT]

[?BBE]Usage:[?RT]
  [?YW]manage add task[?RT] [?BGN]<theme>[?RT] [?BGN]<text>[?RT] [?BBK]--id=[?BGN]<num>[?BBK] [?BBK]--force[?RT]

[?BBE]Description:[?RT]
  * without --id : create new task with autoincrement id
  * with --id    : append text to existing task
  * with --id --force : replace text of existing task`

const docAddTag = `[?GN]Add or remove tags on a theme.[?RT]

[?BBE]Usage:[?RT]
  [?YW]manage add tag[?RT] [?BGN]<theme>[?RT] [?BGN]<tags...>[?RT]

[?BBE]Examples:[?RT]
  [?BE]manage add tag backend api[?RT]`

const docAddTable = `[?GN]Add or update a named table of themes.[?RT]

[?BBE]Usage:[?RT]
  [?YW]manage add table[?RT] [?BGN]<name>[?RT] [?BGN]<themes...>[?RT] [?BBK]--force[?RT]

[?BBE]Description:[?RT]
  * without --force : append themes to the table (creates if missing)
  * with    --force : replace the whole table content`

const docAddPriority = `[?GN]Append priorities to the priority order.[?RT]

[?BBE]Usage:[?RT]
  [?YW]manage add priority[?RT] [?BGN]<p1>[?RT] [?BGN][p2...][?RT]

[?BBE]Examples:[?RT]
  [?BE]manage add priority S[?RT]
  [?BE]manage add priority S T U[?RT]`

const docDo = `[?GN]Mark[?RT] tasks as done or delete them.

[?BBE]Usage:[?RT]
  [?YW]manage do[?RT] [?BGN]<theme>[?RT] [?BGN][task-id...][?RT] [?BBK]--delete[?RT]

[?BBE]Description:[?RT]
  Marks one or more tasks as completed. With [?BBK]--delete[?RT], completed tasks are removed from the list;
  otherwise they stay visible with a "done" mark.

[?BBE]Examples:[?RT]
  [?BE]manage do math 1 2 3[?RT]
  [?BE]manage do math 1 --delete[?RT]
  [?BE]manage do math [?BBK]# do last task`

const docSetPriority = `[?GN]Change the priority of a theme.[?RT]

[?BBE]Usage:[?RT]
  [?YW]manage set priority[?RT] [?BGN]<theme>[?RT] [?BGN]<value>[?RT]

[?BBE]Description:[?RT]
  Changes the priority used to sort the theme in [?YW]manage list[?RT].
  The value should be present in [?BBK]priority_order[?RT] of the config
  (otherwise the theme won't be shown in the list).

[?BBE]Examples:[?RT]
  [?BE]manage set priority backend A[?RT]
  [?BE]manage set priority bio C[?RT]`

const docList = `[?GN]List[?RT] themes and tasks with optional filters.

[?BBE]Usage:[?RT]
  [?YW]manage list[?RT] [?BBK]--priority="A B C"[?BBK] [?BBK]--tags="tag1 tag2"[?BBK]
              [?BBK]--themes="math bio"[?BBK] [?BBK]--table=[?BGN]<day>[?BBK]
              [?BBK]--state=[?BGN](todo|done|all)[?BBK] [?BBK]--todo[?BBK] [?BBK]--done[?BBK]
              [?BBK]--filter=[?BGN](and|or)[?BBK] [?BBK]--code[?RT]

[?BBE]Description:[?RT]
  Displays themes and their tasks, sorted by priority. Filters:
    * [?BBK]--priority[?RT]   - filter by priority levels (space-separated)
    * [?BBK]--tags[?RT]       - filter by tags (space-separated)
    * [?BBK]--themes[?RT]     - filter by theme names (space-separated)
    * [?BBK]--table[?RT]      - filter by a saved schedule table
    * [?BBK]--state[?RT]      - show tasks by state: todo, done, all (default all)
    * [?BBK]--todo[?RT]       - shortcut for --state=todo
    * [?BBK]--done[?RT]       - shortcut for --state=done
    * [?BBK]--filter[?RT]     - logical mode: and (all filters) or or (any filter) (default and)
    * [?BBK]--code[?RT]       - include TODO comments found in the current directory

  If no filters are given, all themes are shown.

[?BBE]Examples:[?RT]
  [?BE]manage list[?RT]
  [?BE]manage list --priority="A B"[?RT]
  [?BE]manage list --tags="school exam"[?RT]
  [?BE]manage list --themes="math"[?RT]
  [?BE]manage list --table=monday[?RT]
  [?BE]manage list --todo --priority=A[?RT]
  [?BE]manage list --done --tags=exam[?RT]
  [?BE]manage list --filter=or --tags=exam --priority=A[?RT]`

const docRemoveTheme = `[?GN]Remove a theme and all its tasks.[?RT]

[?BBE]Usage:[?RT]
  [?YW]manage remove theme[?RT] [?BGN]<name>[?RT]

[?BBE]Description:[?RT]
  Deletes the theme together with its tasks.
  The theme is also removed from every table it belongs to.`

const docRemoveTask = `[?GN]Remove one or more tasks from a theme.[?RT]

[?BBE]Usage:[?RT]
  [?YW]manage remove task[?RT] [?BGN]<theme>[?RT] [?BGN]<id...>[?RT]

[?BBE]Description:[?RT]
  Deletes the given tasks. Task IDs are not renumbered,
  so the remaining tasks keep their identifiers.`

const docRemoveTag = `[?GN]Remove tags from a theme.[?RT]

[?BBE]Usage:[?RT]
  [?YW]manage remove tag[?RT] [?BGN]<theme>[?RT] [?BGN]<tags...>[?RT]

[?BBE]Examples:[?RT]
  [?BE]manage remove tag backend go[?RT]
  [?BE]manage remove tag backend go api[?RT]`

const docRemoveTable = `[?GN]Remove a named table.[?RT]

[?BBE]Usage:[?RT]
  [?YW]manage remove table[?RT] [?BGN]<name>[?RT]

[?BBE]Description:[?RT]
  Deletes the entire table. Themes inside it are not affected.`

const docRemoveFromTable = `[?GN]Remove themes from a table.[?RT]

[?BBE]Usage:[?RT]
  [?YW]manage remove from-table[?RT] [?BGN]<name>[?RT] [?BGN]<themes...>[?RT]

[?BBE]Examples:[?RT]
  [?BE]manage remove from-table monday backend[?RT]`

const docRemovePriority = `[?GN]Remove priorities from the priority order.[?RT]

[?BBE]Usage:[?RT]
  [?YW]manage remove priority[?RT] [?BGN]<p1>[?RT] [?BGN][p2...][?RT]

[?BBE]Examples:[?RT]
  [?BE]manage remove priority Z[?RT]
  [?BE]manage remove priority Y Z[?RT]`

const docShowThemes = `[?GN]List[?RT] all themes with priority, task counters, and tags.

[?BBE]Usage:[?RT]
  [?YW]manage show themes[?RT]

[?BBE]Examples:[?RT]
  [?BE]manage show themes[?RT]`

const docShowPriorities = `[?GN]List[?RT] priorities in priority_order and the themes in each.

[?BBE]Usage:[?RT]
  [?YW]manage show priorities[?RT]

[?BBE]Examples:[?RT]
  [?BE]manage show priorities[?RT]`

const docShowTags = `[?GN]List[?RT] all tags and the themes they belong to.

[?BBE]Usage:[?RT]
  [?YW]manage show tags[?RT]

[?BBE]Examples:[?RT]
  [?BE]manage show tags[?RT]`

const docShowTheme = `[?GN]Show[?RT] a summary card for a single theme.

[?BBE]Usage:[?RT]
  [?YW]manage show theme[?RT] [?BGN]<name>[?RT]

[?BBE]Description:[?RT]
  Prints the theme's priority, tags, task counters (total / done / todo),
  the last assigned task ID, and the list of tables that contain the theme.

[?BBE]Examples:[?RT]
  [?BE]manage show theme backend[?RT]`

const docShowInfo = `[?GN]Show[?RT] an overall summary, grouped by priority.

[?BBE]Usage:[?RT]
  [?YW]manage show info[?RT]

[?BBE]Description:[?RT]
  Prints every theme grouped by priority with task counters
  and tags, then the total counters.

[?BBE]Examples:[?RT]
  [?BE]manage show info[?RT]`

func NewCli(mc *ManageCli) *tap.Parser {
	p := tap.NewParser("manage", utils.FramedTextMap("CN", "[?YW]Manage[?RT]", "", map[string]string{
		"": `[?GN]Powerful task manager.[?RT]
[?BBK]Supports themes, tasks and tables.[?RT]
[?BBK]Very good for education, but is also 
[?BBK]useful for other tasks.[?RT]`,
	})+`

[?GN]Commands:[?RT]
  init    Initialize config file
  add     Add themes, tasks, tags, or tables
  remove
  do      Mark tasks as done
  set     Set any param
  list    List tasks with filters
  show    Show any info`, []string{"-h", "help"}, tap.DefaultParserConfig())

	p.AddCommand("init", func(p *tap.Parser, s []string) error {
		if err := config.Save(manager.NewManager()); err != nil {
			return err
		}
		mc.NoSave = true
		return nil
	}, docInit, nil, nil, false)

	add := tap.NewParser("add", "[?GN]Add[?RT] themes, tasks, tags, tables, or priorities.",
		[]string{"-h", "help"}, tap.DefaultParserConfig())

	add.AddCommand("theme", mc.AddTheme, docAddTheme, []string{"name"}, nil, false)
	add.AddCommand("task", mc.AddTask, docAddTask, []string{"theme", "text"}, nil, false)
	add.AddCommand("tag", mc.AddTag, docAddTag, []string{"theme", "tags"}, nil, true)
	add.AddCommand("table", mc.AddTable, docAddTable, []string{"name", "themes"}, nil, true)
	add.AddCommand("priority", mc.AddPriority, docAddPriority, []string{"priorities"}, nil, true)

	p.AddSubcommand("add", add)

	p.AddCommand("do", mc.Do, docDo, []string{"theme"}, nil, true)

	remove := tap.NewParser("remove", "[?GN]Remove[?RT] themes, tasks, tags, tables, or priorities.",
		[]string{"-h", "help"}, tap.DefaultParserConfig())

	remove.AddCommand("theme", mc.RemoveTheme, docRemoveTheme, []string{"name"}, nil, false)
	remove.AddCommand("task", mc.RemoveTask, docRemoveTask, []string{"theme", "ids"}, nil, true)
	remove.AddCommand("tag", mc.RemoveTag, docRemoveTag, []string{"theme", "tags"}, nil, true)
	remove.AddCommand("table", mc.RemoveTable, docRemoveTable, []string{"name"}, nil, false)
	remove.AddCommand("from-table", mc.RemoveFromTable, docRemoveFromTable, []string{"name", "themes"}, nil, true)
	remove.AddCommand("priority", mc.RemovePriority, docRemovePriority, []string{"priorities"}, nil, true)

	p.AddSubcommand("remove", remove)

	p.AddCommand("list", mc.List, docList, nil, nil, false)

	show := tap.NewParser(
		"show",
		"[?GN]Show[?RT] themes, priorities, tags, or a detailed theme card.",
		[]string{"-h", "help"},
		tap.DefaultParserConfig(),
	)
	show.AddCommand("themes", mc.ShowThemes, docShowThemes, nil, nil, false)
	show.AddCommand("priorities", mc.ShowPriorities, docShowPriorities, nil, nil, false)
	show.AddCommand("tags", mc.ShowTags, docShowTags, nil, nil, false)
	show.AddCommand("theme", mc.ShowTheme, docShowTheme, []string{"name"}, nil, false)
	show.AddCommand("info", mc.ShowInfo, docShowInfo, nil, nil, false)

	p.AddSubcommand("show", show)

	set := tap.NewParser(
		"set",
		"[?GN]Set[?RT] various parameters.",
		[]string{"-h", "help"},
		tap.DefaultParserConfig(),
	)
	set.AddCommand("priority", mc.SetPriority, docSetPriority, []string{"theme", "value"}, nil, false)

	p.AddSubcommand("set", set)

	sys := tap.NewParser("sys", "System info and debug", []string{"-h", "help"}, tap.DefaultParserConfig())

	sys.AddCommand("version", func(p *tap.Parser, s []string) error {
		mc.NoSave = true
		fmt.Printf("Tap v%v, Manage v%v\n", tap.Version, Version)
		return nil
	}, "Show version", nil, nil, false)

	sys.AddCommand("conf", func(p *tap.Parser, s []string) error {
		mc.NoSave = true
		file, err := config.GetConf()
		if err != nil {
			return err
		}
		fmt.Println(file)
		return nil
	}, "Show config", nil, nil, false)

	sys.AddCommand("conf-path", func(p *tap.Parser, s []string) error {
		mc.NoSave = true
		fmt.Println("~/" + config.ConfigFile)
		return nil
	}, "Show config", nil, nil, false)

	p.AddSubcommand("sys", sys)

	return p
}

func main() {
	m := manager.NewManager()
	has_code_flag := slices.Contains(os.Args, "--code")
	err := config.Load(m, has_code_flag)
	if err != nil && len(os.Args) > 1 && os.Args[1] != "init" {
		fmt.Println(err)
		os.Exit(1)
	}
	mc := &ManageCli{
		C:      m,
		NoSave: false,
	}
	err = NewCli(mc).Main()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !mc.NoSave {
		err = config.Save(mc.C)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

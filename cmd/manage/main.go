package main

import (
	"fmt"

	"github.com/pt-main/manage/manager"
	"github.com/pt-main/tap"
	"github.com/pt-main/tap/utils"
)

const docInit = `[?GN]Initialize[?RT] config file.

[?BBE]Usage:[?RT]
  [?YW]manage init[?RT]

[?BBE]Description:[?RT]
  Creates the config directory [?BBK]~/.manage/[?RT] and a default configuration file.
  This is the first command you need to run before using any other commands.`

const docAdd = `[?GN]Add[?RT] themes, tasks, tags, or tables.

[?BBE]Usage:[?RT]
  [?YW]manage add theme[?RT] [?BGN]<name>[?RT] [?BBK]--priority=[?BGN]<letter>[?BBK] [?BBK]--tags=[?BGN]"tag1 tag2 ..."[?RT] [?BBK]--force[?RT]

  [?YW]manage add task[?RT] [?BGN]<theme>[?RT] [?BGN]<text>[?RT] [?BBK]--id=[?BGN]<num>[?BBK] [?BBK]--force[?RT]

  [?YW]manage add tag[?RT] [?BGN]<theme>[?RT] [?BGN]<tags...>[?RT]

  [?YW]manage add table[?RT] [?BGN]<name>[?RT] [?BGN]<themes...>[?RT] [?BBK]--force[?RT]

[?BBE]Description:[?RT]
  * Theme: adds a new theme with optional priority and tags.
    [?BBK]--force[?RT] replaces an existing theme.
  * Task: adds a task to a theme. If [?BBK]--id[?RT] is given:
      - with [?BBK]--force[?RT]: replaces the task text if [?BBK]--id[?RT] passed.
      - without [?BBK]--force[?RT]: appends to the existing task.
    If [?BBK]--id[?RT] is not given, creates a new task.
  * Tag: add or remove tags. Use "[?RD]![?RT]tag" to remove a tag.
  * Table: manages named lists of themes.
      - without [?BBK]--force[?RT]: appends themes to the table (creates if missing).
      - with [?BBK]--force[?RT]: replaces the whole table content.`

const docDo = `[?GN]Mark[?RT] tasks as done or delete them.

[?BBE]Usage:[?RT]
  [?YW]manage do[?RT] [?BGN]<theme>[?RT] [?BGN]<task-id>[?RT] [?BGN][task-id...][?RT] [?BBK]--delete[?RT]

[?BBE]Description:[?RT]
  Marks one or more tasks as completed. With [?BBK]--delete[?RT], completed tasks are removed from the list;
  otherwise they stay visible with a "done" mark.

[?BBE]Examples:[?RT]
  [?BE]manage do math 1 2 3[?RT]
  [?BE]manage do math 1 --delete[?RT]`

const docSet = `[?GN]Set[?RT] various parameters.

[?BBE]Usage:[?RT]
  [?YW]manage set priority[?RT] [?BGN]<theme>[?RT] [?BGN]<A-Z>[?RT]

[?BBE]Description:[?RT]
  Changes the priority of a theme. Priority determines the order in [?YW]list[?RT].

[?BBE]Examples:[?RT]
  [?BE]manage set priority math A[?RT]
  [?BE]manage set priority bio C[?RT]`

const docList = `[?GN]List[?RT] themes and tasks with optional filters.

[?BBE]Usage:[?RT]
  [?YW]manage list[?RT] [?BBK]--priority="A B C"[?BBK] [?BBK]--tags="tag1 tag2"[?BBK]
              [?BBK]--themes="math bio"[?BBK] [?BBK]--table=[?BGN]<day>[?BBK]
              [?BBK]--state=[?BGN](todo|done|all)[?BBK] [?BBK]--todo[?BBK] [?BBK]--done[?BBK]
              [?BBK]--filter=[?BGN](and|or)[?BBK]

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

func NewCli() *tap.Parser {
	p := tap.NewParser("manage", utils.FramedTextMap("CN", "[?YW]Manage[?RT]", "", map[string]string{
		"": `[?GN]Powerful task manager.[?RT]
[?BBK]Supports themes, tasks and tables.[?RT]
[?BBK]Very good for education, but is also 
[?BBK]useful for other tasks.[?RT]`,
	})+`

[?GN]Commands:[?RT]
  init    Initialize config file
  add     Add themes, tasks, tags, or tables
  do      Mark tasks as done
  set     Set priorities
  list    List tasks with filters`, []string{"-h", "help"}, tap.DefaultParserConfig())

	p.AddCommand("init", func(p *tap.Parser, s []string) error {
		if err := Save(manager.NewManager()); err != nil {
			return err
		}
		return nil
	}, docInit, nil, nil, false)

	p.AddCommand("add", Add, docAdd, []string{"mode", "arg"}, nil, true)

	p.AddCommand("do", Do, docDo, []string{"theme", "task"}, nil, true)

	p.AddCommand("set", Set, docSet, nil, nil, true)

	p.AddCommand("list", List, docList, nil, nil, false)

	return p
}

func main() {
	err := NewCli().Main()
	if err != nil {
		fmt.Println(err)
	}
}

# manage - command-line task manager

```bash
go install github.com/pt-main/manage/cmd/manage@latest
```

**Manage** is a CLI tool for working with tasks organized in a hierarchical structure: themes (projects, areas) contain tasks. Each theme has a priority, tags, and a list of tasks; each task has a completion status, description, and timestamp. Additionally, named tables are supported for grouping themes. The tool is not tied to any specific domain and is suitable for any scenario: from personal affairs to work projects.

---

## Quick start

Before first use, run initialization:

```bash
manage init
```

This will create a file `~/manager_config.tess` with default settings.

Let's add a theme and a task:

```bash
manage add theme backend --priority=A --tags="go api"
manage add task backend "Write tests"
```

Let's view the list:

```bash
manage list
```

Mark a task as done:

```bash
manage do backend 1
```

Or do the last one:

```bash
manage do backend
```

---

## Basic concepts

- **Theme** - a group of tasks (project, module, client).
  A theme has a name, priority, tags, and a list of tasks.
- **Task** - a unit of work.
  A task has an ID (auto-increment within the theme), description, status (done/not done), and creation time.
- **Priority** - any string without spaces (e.g., `A`, `high`, `critical`).
  The sort order is defined in the settings.
- **Tags** - labels for grouping themes.
- **Table** - a named list of themes for quick filtering.

---

## Commands

### `init`

```bash
manage init
```

Creates the configuration file `~/manager_config.tess`.
Without this, commands do not work. It stores tasks, settings, and tables in Tessera format.

**Warning:** running `init` again when a config already exists overwrites it.

---

### `add`

Adds a theme, task, tags, table, or priorities.

#### Adding a theme

```bash
manage add theme <name> [--priority=<value>] [--tags="tag1 tag2 ..."] [--force]
```

- `name` - theme name (required).
- `--priority` - any string without spaces. If not specified or empty, the value of `default_priority` from settings is used (default `F`).
- `--tags` - list of tags separated by spaces. If not specified, no tags are added.
- `--force` - replace an existing theme with the same name. Without this flag, the command will fail if the theme already exists (data will not be affected).

Examples:

```bash
manage add theme backend
manage add theme backend --priority=A --tags="go api"
manage add theme backend --priority=B --tags="rust" --force
```

#### Adding a task

```bash
manage add task <theme> <text> [--id=<num>] [--force]
```

- Without `--id` - a new task is created with a new ID. `--force` in this case is not applied and causes an error.
- With `--id` and without `--force` - the text is **appended** to the existing task.
- With `--id` and `--force` - the text **replaces** the existing task.

Examples:

```bash
manage add task backend "Write tests"
manage add task backend " + integration" --id=1
manage add task backend "Replace text" --id=1 --force
```

#### Adding tags

```bash
manage add tag <theme> <tags...>
```

Adds tags to the theme. Existing tags are not duplicated.

#### Adding a table

```bash
manage add table <name> <themes...> [--force]
```

- Without `--force` - adds themes to an existing table (if the table does not exist, creates it).
- With `--force` - replaces the entire content of the table with the new list.

Examples:

```bash
manage add table monday backend frontend
manage add table monday database          # adds database
manage add table monday backend --force   # replaces with a single theme backend
```

#### Adding priorities

```bash
manage add priority <p1> [p2 ...]
```

Adds priorities to the end of `priority_order` if they are not already there.

Examples:

```bash
manage add priority S
manage add priority S T U
```

---

### `remove`

Removes themes, tasks, tags, tables, and priorities.

#### Removing a theme

```bash
manage remove theme <name>
```

Removes the theme along with all tasks. The theme is also removed from all tables.

#### Removing a task

```bash
manage remove task <theme> <id...>
```

Removes the specified tasks. IDs of remaining tasks are not renumbered.

If none of the specified IDs are found, the command fails.

#### Removing tags

```bash
manage remove tag <theme> <tags...>
```

Removes tags from the theme. If none of the specified tags are found, the command fails.

#### Removing a table

```bash
manage remove table <name>
```

Removes the entire table. Themes inside it are not affected.

#### Removing themes from a table

```bash
manage remove from-table <name> <themes...>
```

Removes the specified themes from the table without deleting the themes themselves.

#### Removing priorities

```bash
manage remove priority <p1> [p2 ...]
```

Removes priorities from `priority_order`. You cannot remove all priorities at once.

---

### `do`

Marks tasks as done or removes them.

```bash
manage do <theme> [id...] [--delete]
```

- Without `--delete` - the task is marked as done but remains in the list.
- With `--delete` - the task is removed.
- Without `id...` - the last task of the theme is done.

If there are no tasks in the theme or some of the specified IDs are not found, an error is displayed. The remaining tasks are still processed.

Example:

```bash
manage do backend 1 2 3
manage do backend 1 --delete
manage do backend # do the last task
```

If you run the command with the global flag `--code` (see the "Code mode" section), you can do tasks of the `.code` theme. For code tasks, `do` does not simply mark the task as done but **removes the `TODO:` marker** from the corresponding line in the source file:

```bash
manage --code do .code        # remove TODO from the last found line
manage --code do .code 2 5    # remove TODO from lines with id 2 and 5
```

---

### `set`

```bash
manage set priority <theme> <value>
```

Changes the priority of a theme. If the new priority is not in `priority_order`, a warning is displayed - the theme will not be shown in `list` until the priority is added.

Example:

```bash
manage set priority backend A
```

---

### `list`

Displays a list of themes and tasks with filtering. Empty themes are shown as a warning before the list.

```bash
manage list [--priority="A B"] [--tags="tag1 tag2"] [--themes="theme1 theme2"] [--table=<name>] [--state=(todo|done|all)] [--todo] [--done] [--filter=(and|or)]
```

> The global flag `--code` does not work as a `list` filter but as a separate mode that adds the `.code` theme to all commands. Description in the "Code mode" section.

#### Filters

- `--priority` - space-separated list of priorities. Values are checked against `priority_order`; if absent, an error is displayed.
- `--tags` - space-separated list of tags.
- `--themes` - space-separated list of theme names.
- `--table` - table name (themes from it will be shown). The table must exist, otherwise an error is displayed.
- `--state` - task state:
  - `todo` - only incomplete.
  - `done` - only completed.
  - `all` - all (default).
- `--todo` - shorthand for `--state=todo`.
- `--done` - shorthand for `--state=done`.
- `--filter` - logic for combining conditions:
  - `and` - all conditions must match (default).
  - `or` - any condition is sufficient.

`--table` and `--themes` are combined into a single filter by theme names.

If both `--todo` and `--done` are specified, `all` is used with a warning.

#### Sorting

Themes are displayed in the order of priorities defined in settings (by default from `A` to `Z`).
Tasks are displayed within their themes.

#### Examples

```bash
manage list
manage list --priority="A B" --tags="go api" --filter=or
manage list --table=monday --todo
manage list --done --tags=exam
```

---

### `show`

Displays reference information.

```bash
manage show themes
manage show priorities
manage show tags
manage show theme <name>
manage show info
```

- `themes` - all themes with priority, task counts, and tags, grouped by priority.
- `priorities` - priorities from `priority_order` and themes in each of them. Several consecutive empty priorities are collapsed into a range (`A-Z`, `high..low`). Priorities not in `priority_order` are also shown in a separate block.
- `tags` - all tags and the themes they belong to.
- `theme <name>` - theme card: priority, tags, task counts (total / done / remaining), last issued ID, and list of tables the theme belongs to. **The task list is not displayed in the card** - use `manage list --themes=<name>` or `manage show info` for that.
- `info` - summary by priorities: all themes grouped by priority with task counts and tags, then total sums by themes and tasks.

If the command is run with the `--code` flag, the output of `show themes`, `show priorities`, `show tags`, and `show info` will include the transient `.code` theme (see the "Code mode" section).

---

### `sys`

System information and debug commands.

```bash
manage sys version
manage sys conf
manage sys conf-path
```

- `version` - displays Tap and Manage versions.
- `conf` - displays the contents of the current configuration file.
- `conf-path` - displays the path to the configuration file (`~/manager_config.tess`).

---

## Code mode

`--code` is a **global flag** recognized at the very beginning of the program, and its position in the command line does not matter: `manage list --code` and `manage --code list` are equivalent.

What `--code` does:

1. **Scans the current working directory** (recursively) and collects all lines after `TODO:` (not necessarily after a comment or from the beginning of the line).
2. Creates a **transient theme** `.code` with priority `A` in memory.
3. Each found `TODO:` line becomes a task of the `.code` theme:
   - `About` - text after `TODO:`;
   - `File` - relative path to the file;
   - `CodeLine` - line number.
4. The `.code` theme exists only during the execution of the current command and is **never saved** to `~/manager_config.tess` (and saving a theme named .code is impossible even if the user added it).

So `--code` is a **mode that adds an artificial `.code` theme to all commands**, with which you can work like a regular theme: filter, sort, show, mark tasks as done.

### Scanned files

- Text files by extension: `.go`, `.py`, `.js`, `.ts`, `.rs`, `.c/.h/.cc/.cpp`, `.java`, `.rb`, `.php`, `.cs`, `.swift`, `.kt`, `.sh`, `.md`, `.txt`, `.yaml`, `.json`, `.html`, `.css`, `.sql`, `.lua`, etc.
- Specific basenames: `Makefile`, `Dockerfile`, `go.mod`, `go.sum`, `.gitignore`, `.env`, etc.
- **Skipped** directories: `.git`, `.svn`, `.hg`, `node_modules`, `vendor`, `target`, `build`, `dist`, `.idea`, `.vscode`, `__pycache__`, `.venv`, `bin`, `obj`, etc.
- Line size limited to 2 MiB.

### Task identifiers

IDs for code tasks are assigned after sorting (file -> line -> text), starting from `0`. The ID **does not match** the line number - it is simply the sequential number in the scanned list.

### Doing code tasks

For code tasks, `manage do` works differently than for regular ones:

- Instead of setting the `done` flag in the config, it **physically removes the `TODO:` marker** from the corresponding line of the file.
- The `TODO:` marker is removed case-insensitively (`todo:` also matches).
- If there are multiple `TODO:` in a line, the first one whose text matches the task's `About` is removed.
- The line is rewritten in place, `\r\n` / `\n` are preserved.
- If the file or line has changed since scanning, the command returns an error and the file is not modified.
- `--delete` has no meaning for code tasks: the result is the same - the marker is removed from the file.

### Examples

```bash
# Show all TODOs in the project (theme .code with priority A)
manage list --code

# Only code tasks
manage list --code --themes=.code

# Among other things, filter by tags/priorities
manage list --code --priority="A" --todo

# Do the last found TODO (remove marker from file)
manage --code do .code

# Do specific code tasks by their IDs
manage --code do .code 2 5 7

# View the card of the .code theme
manage show theme .code --code

# Summary including code tasks
manage show info --code
```

### Important nuances

- The theme is strictly named `.code`; the user cannot add a theme with this name.
- The priority of `.code` is always `A`. If you removed `A` from `priority_order`, the `.code` theme will be shown in the "Unknown priorities" block.
- IDs of code tasks are unstable between runs: when files change, numbering may change. Rely on content (file + line + text), not on ID.
- `--code` does not work with the `init` command and does not save scan results.

---

## Configuration

The file `~/manager_config.tess` is created on `init`.
It is stored in Tessera format. You can edit it manually for bulk changes, but usually all operations are performed via the CLI.

Settings:

- `priority_order` - priority order (default from A to Z).
- `default_priority` - value used when adding a theme without specifying `--priority`.

---

## License

MIT - details in [LICENCE](https://github.com/pt-main/manage/blob/main/LICENCE).

---

**manage** is written using [tap](https://github.com/pt-main/tap), [tessera](https://github.com/pt-main/tessera) and [tycl](https://github.com/pt-main/tycl).
Author: Pt.
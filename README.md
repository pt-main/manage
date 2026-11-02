# manage - command-line task manager

```bash
go install github.com/pt-main/manage/cmd/manage@latest
```

**Manage** is a CLI tool for working with tasks organized in a hierarchical structure: themes (projects, areas) contain tasks. Each theme has a priority, tags, and a list of tasks; each task has a completion status, description, and timestamp. Named tables are additionally supported for grouping themes. The tool is not tied to any specific domain and is suitable for any scenario: from personal affairs to work projects.

---

## Quick start

Before first use, run initialization:

```bash
manage init
```

This creates the file `~/manager_config.tess` with default settings.

Let's add a theme and a task:

```bash
manage add theme backend --priority=A --tags="go api"
manage add task backend "Write tests"
```

Let's view the list:

```bash
manage list
```

Mark the task as done:

```bash
manage do backend 1
```

Or do the last one:

```bash
manage do backend
```

---

## Core concepts

- **Theme** - a group of tasks (project, module, client).
  A theme has a name, priority, tags, and a list of tasks.
- **Task** - a unit of work.
  A task has an ID (auto-increment within the theme), description, status (done/not done), and creation time.
- **Priority** - any string without spaces (for example, `A`, `high`, `critical`).
  Sort order is defined in the settings.
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
- `--priority` - any string without spaces. If omitted or empty, the `default_priority` value from settings is used (default `F`).
- `--tags` - list of tags separated by spaces. If omitted, no tags are added.
- `--force` - replace an existing theme with the same name. Without this flag, the command fails if the theme already exists (data is not affected).

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
- With `--force` - replaces the entire contents of the table with the new list.

Examples:

```bash
manage add table monday backend frontend
manage add table monday database          # will add database
manage add table monday backend --force   # will replace with a single theme backend
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

Removes the specified tasks. The IDs of remaining tasks are not renumbered.

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

Marks tasks as done or deletes them.

```bash
manage do <theme> [id...] [--delete]
```

- Without `--delete` - the task is marked as done but remains in the list.
- With `--delete` - the task is deleted.
- Without `id...` - the last task of the theme is processed.

If there are no tasks in the theme or some of the specified IDs are not found, an error is displayed. The remaining tasks are still processed.

Example:

```bash
manage do backend 1 2 3
manage do backend 1 --delete
manage do backend # do the last task
```

---

### `set`

```bash
manage set priority <theme> <value>
```

Changes the theme's priority. If the new priority is not in `priority_order`, a warning is displayed - the theme will not be shown in `list` until the priority is added.

Example:

```bash
manage set priority backend A
```

---

### `list`

Displays a list of themes and tasks with filtering. Empty themes are output as a warning before the list.

```bash
manage list [--priority="A B"] [--tags="tag1 tag2"] [--themes="theme1 theme2"] [--table=<name>] [--state=(todo|done|all)] [--todo] [--done] [--filter=(and|or)] [--code]
```

#### Filters

- `--priority` - list of priorities separated by spaces. Values are checked against `priority_order`; if absent, an error is displayed.
- `--tags` - list of tags separated by spaces.
- `--themes` - list of theme names separated by spaces.
- `--table` - table name (themes from it will be shown). The table must exist, otherwise an error is displayed.
- `--state` - task state:
  - `todo` - only incomplete.
  - `done` - only completed.
  - `all` - all (default).
- `--todo` - shorthand for `--state=todo`.
- `--done` - shorthand for `--state=done`.
- `--filter` - logic for combining conditions:
  - `and` - all conditions must match (default).
  - `or` - any condition matching is sufficient.
- `--code` - additionally include `TODO:` comments from the current directory.

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
manage list --code
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

- `themes` - all themes with priority, task counters, and tags.
- `priorities` - priorities from `priority_order` and the themes in each of them. Also shows priorities that are not in `priority_order`.
- `tags` - all tags and the themes they relate to.
- `theme <name>` - theme card: priority, tags, task counters, last issued ID, tables containing this theme.
- `info` - summary by priority: all themes grouped by priority with task counters, then total sums.

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

## Configuration

The file `~/manager_config.tess` is created by `init`.
It is stored in Tessera format. You can edit it manually for bulk changes, but usually all operations are performed via the CLI.

Settings:

- `priority_order` - priority order (by default from A to Z).
- `default_priority` - value used when adding a theme without specifying `--priority`.

---

## License

MIT - details in [LICENCE](https://github.com/pt-main/manage/blob/main/LICENCE).

---

**manage** is written using [tap](https://github.com/pt-main/tap), [tessera](https://github.com/pt-main/tessera), and [tycl](https://github.com/pt-main/tycl).
Author: Pt.
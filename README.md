# manage — command-line task manager

```bash
go install github.com/pt-main/manage/cmd/manage@latest
```

**Manage** is a CLI tool for handling tasks organised in a hierarchical structure: topics (projects, areas) contain tasks. Each topic has a priority, tags, and a task list; each task has a completion status, a description, and a timestamp. Named tables are also supported for grouping topics. The tool is domain‑agnostic and fits any scenario, from personal errands to work projects.

---

## Quick start

Before the first use, run initialisation:

```bash
manage init
```

This creates a `~/manager_config.tycl` file with default settings.

Let's add a topic and a task:

```bash
manage add theme backend --priority=A --tags="go api"
manage add task backend "Write tests"
```

View the list:

```bash
manage list
```

Mark a task as done:

```bash
manage do backend 1
```

---

## Core concepts

- **Topic** – a group of tasks (project, module, client).  
  A topic has a name, priority, tags, and a task list.
- **Task** – a unit of work.  
  A task has an ID (auto‑incremented within the topic), a description, a status (done/not done), and a creation timestamp.
- **Priority** – any string without spaces (e.g., `A`, `high`, `critical`).  
  The sorting order is defined in the configuration.
- **Tags** – labels for grouping topics.
- **Table** – a named list of topics for quick filtering.

---

## Commands

### `init`

```bash
manage init
```

Creates the configuration file `~/manager_config.tycl`.  
Without this, commands won't work. The file stores tasks, settings, and tables in the contract‑fixed format of the configuration language used.

---

### `add`

Adds a topic, task, tags, or table.

#### Adding a topic

```bash
manage add theme <name> [--priority=<value>] [--tags="tag1 tag2 ..."] [--force]
```

- `name` – topic name (required).
- `--priority` – arbitrary string without spaces. If omitted, the `default_priority` from settings is used (default `F`).
- `--tags` – space‑separated list of tags. If omitted, no tags are added.
- `--force` – replace an existing topic with the same name. Without this flag, the command fails if the topic already exists (data is untouched).

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

- Without `--id` – creates a new task with a new ID.
- With `--id` and without `--force` – appends the text to the existing task.
- With `--id` and `--force` – replaces the existing task with the new text.

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

Adds tags to a topic. To remove a tag, use `!` before the name:

```bash
manage add tag backend api !go
```

#### Adding a table

```bash
manage add table <name> <themes...> [--force]
```

- Without `--force` – adds topics to an existing table (creates the table if it doesn't exist).
- With `--force` – replaces the whole table content with the new list.

Examples:

```bash
manage add table monday backend frontend
manage add table monday database          # adds database
manage add table monday backend --force   # replaces with just backend
```

---

### `do`

Marks tasks as done or deletes them.

```bash
manage do <theme> <id...> [--delete]
```

- Without `--delete` – marks the task as done, but keeps it in the list.
- With `--delete` – deletes the task.

Example:

```bash
manage do backend 1 2 3
manage do backend 1 --delete
```

---

### `set`

```bash
manage set priority <theme> <value>
```

Changes the priority of a topic.

Example:

```bash
manage set priority backend A
```

---

### `list`

Displays the list of topics and tasks with filtering. Empty topics are shown as a warning before the list.

```bash
manage list [--priority="A B"] [--tags="tag1 tag2"] [--themes="theme1 theme2"] [--table=<name>] [--state=(todo|done|all)] [--todo] [--done] [--filter=(and|or)]
```

#### Filters

- `--priority` – space‑separated list of priorities.
- `--tags` – space‑separated list of tags.
- `--themes` – space‑separated list of topic names.
- `--table` – table name (shows topics from that table).
- `--state` – task status:
  - `todo` – only unfinished tasks.
  - `done` – only completed tasks.
  - `all` – all (default).
- `--todo` – shorthand for `--state=todo`.
- `--done` – shorthand for `--state=done`.
- `--filter` – logic for combining conditions:
  - `and` – all conditions must match (default).
  - `or` – any condition is sufficient.

#### Sorting

Topics are displayed in the priority order defined in the settings (default from `A` to `Z`).  
Tasks within a topic are sorted by ID.

#### Examples

```bash
manage list
manage list --priority="A B" --tags="go api" --filter=or
manage list --table=monday --todo
manage list --done --tags=exam
```

---

## Configuration

The file `~/manager_config.tycl` is created by `init`.  
You can edit it manually for bulk changes, but usually all operations are done via the CLI.

Settings:

- `priority_order` – the priority order (default A to Z).
- `default_priority` – the value used when adding a topic without `--priority`.

---

## License

Apache 2.0 — details in [LICENSE](https://github.com/pt-main/manage/blob/main/LICENSE).

---

**manage** is written using [tap](https://github.com/pt-main/tap), [tycl](https://github.com/pt-main/tycl), and [lc](https://github.com/pt-main/lc).  
Author: Pt.
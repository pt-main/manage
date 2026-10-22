# Merged files

Всего файлов: 14

---

## `LICENCE`

```
MIT License

Copyright (c) 2026 Pt

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## `README-ru.md`

````markdown
# manage - менеджер задач для командной строки

```bash
go install github.com/pt-main/manage/cmd/manage@latest
```

**Manage** - это CLI-инструмент для работы с задачами, организованными в иерархическую структуру: темы (проекты, области) содержат задачи. Каждая тема имеет приоритет, теги и список задач, каждая задача - статус выполнения, описание и временную метку. Дополнительно поддерживаются именованные таблицы для группировки тем. Инструмент не привязан к какой-либо предметной области и подходит для любых сценариев: от личных дел до рабочих проектов.

---

## Быстрый старт

Перед первым использованием выполните инициализацию:

```bash
manage init
```

Это создаст файл `~/manager_config.tess` с настройками по умолчанию.

Добавим тему и задачу:

```bash
manage add theme backend --priority=A --tags="go api"
manage add task backend "Написать тесты"
```

Посмотрим список:

```bash
manage list
```

Отметим задачу выполненной:

```bash
manage do backend 1
```

Или выполним последнюю: 
```bash
manage do backend
```

---

## Основные понятия

- **Тема** - группа задач (проект, модуль, клиент).  
  У темы есть имя, приоритет, теги и список задач.
- **Задача** - единица работы.  
  У задачи есть ID (автоинкремент в рамках темы), описание, статус (выполнена/нет), время создания.
- **Приоритет** - любая строка без пробелов (например, `A`, `high`, `critical`).  
  Порядок сортировки задаётся в настройках.
- **Теги** - метки для группировки тем.
- **Таблица** - именованный список тем для быстрой фильтрации.

---

## Команды

### `init`

```bash
manage init
```

Создаёт конфигурационный файл `~/manager_config.tess`.  
Без этого команды не работают. В нём в формате Tessera хранятся задачи, настройки и таблицы.

**Внимание:** повторный запуск `init` при существующем конфиге перезаписывает его.

---

### `add`

Добавляет тему, задачу, теги или таблицу.

#### Добавление темы

```bash
manage add theme <name> [--priority=<value>] [--tags="tag1 tag2 ..."] [--force]
```

- `name` - имя темы (обязательно).
- `--priority` - произвольная строка без пробелов. Если не указан или пуст, используется значение `default_priority` из настроек (по умолчанию `F`).
- `--tags` - список тегов, разделённых пробелами. Если не указан, теги не добавляются.
- `--force` - заменить существующую тему с таким же именем. Без этого флага команда завершится ошибкой, если тема уже существует (данные не будут затронуты).

Примеры:

```bash
manage add theme backend
manage add theme backend --priority=A --tags="go api"
manage add theme backend --priority=B --tags="rust" --force
```

#### Добавление задачи

```bash
manage add task <theme> <text> [--id=<num>] [--force]
```

- Без `--id` - создаётся новая задача с новым ID. `--force` в этом случае не действует.
- С `--id` и без `--force` - текст **дописывается** к существующей задаче.
- С `--id` и `--force` - текст **заменяет** существующую задачу.

Примеры:

```bash
manage add task backend "Написать тесты"
manage add task backend " + интеграционные" --id=1
manage add task backend "Заменить текст" --id=1 --force
```

#### Добавление тегов

```bash
manage add tag <theme> <tags...>
```

Добавляет теги к теме. Чтобы удалить тег, используйте `!` перед именем:

```bash
manage add tag backend api !go
```

#### Добавление таблицы

```bash
manage add table <name> <themes...> [--force]
```

- Без `--force` - добавляет темы к существующей таблице (если таблицы нет - создаёт).
- С `--force` - заменяет всё содержимое таблицы новым списком.

Примеры:

```bash
manage add table monday backend frontend
manage add table monday database          # добавит database
manage add table monday backend --force   # заменит на одну тему backend
```

---

### `do`

Отмечает задачи выполненными или удаляет их.

```bash
manage do <theme> [id...] [--delete]
```

- Без `--delete` - задача помечается как выполненная, но остаётся в списке.
- С `--delete` - задача удаляется.
- Без `id...` - выполняется последняя задача темы.

Если в теме нет задач или некоторые из указанных ID не найдены, выводится ошибка. Остальные задачи при этом обрабатываются.

Пример:

```bash
manage do backend 1 2 3
manage do backend 1 --delete
manage do backend # сделать последнюю задачу
```

---

### `set`

```bash
manage set priority <theme> <value>
```

Меняет приоритет темы.

Пример:

```bash
manage set priority backend A
```

---

### `list`

Выводит список тем и задач с фильтрацией. Пустые темы выводятся предупреждением до списка.

```bash
manage list [--priority="A B"] [--tags="tag1 tag2"] [--themes="theme1 theme2"] [--table=<name>] [--state=(todo|done|all)] [--todo] [--done] [--filter=(and|or)] [--code]
```

#### Фильтры

- `--priority` - список приоритетов через пробел.
- `--tags` - список тегов через пробел.
- `--themes` - список имён тем через пробел.
- `--table` - имя таблицы (будут показаны темы из неё).
- `--state` - состояние задач:
  - `todo` - только невыполненные.
  - `done` - только выполненные.
  - `all` - все (по умолчанию).
- `--todo` - сокращение для `--state=todo`.
- `--done` - сокращение для `--state=done`.
- `--filter` - логика объединения условий:
  - `and` - все условия должны совпадать (по умолчанию).
  - `or` - достаточно выполнения любого условия.
- `--code` - дополнительно включить комментарии `TODO:` из текущего каталога.

`--table` и `--themes` объединяются в единый фильтр по именам тем.

Если указаны и `--todo`, и `--done`, используется `all` с предупреждением.

#### Сортировка

Темы выводятся в порядке приоритетов, заданных в настройках (по умолчанию от `A` до `Z`).  
Задачи выводятся в составе своих тем.

#### Примеры

```bash
manage list
manage list --priority="A B" --tags="go api" --filter=or
manage list --table=monday --todo
manage list --done --tags=exam
manage list --code
```

---

### `sys`

Системная информация и отладочные команды.

```bash
manage sys version
manage sys conf
manage sys conf-path
```

- `version` - выводит версии Tap и Manage.
- `conf` - выводит содержимое текущего конфигурационного файла.
- `conf-path` - выводит путь к конфигурационному файлу (`~/manager_config.tess`).

---

## Конфигурация

Файл `~/manager_config.tess` создаётся при `init`.  
Он хранится в формате Tessera. Вы можете редактировать его вручную для массовых изменений, но обычно все операции выполняются через CLI.

Настройки:

- `priority_order` - порядок приоритетов (по умолчанию от A до Z).
- `default_priority` - значение, используемое при добавлении темы без указания `--priority`.

---

## Лицензия

MIT - подробности в [LICENCE](https://github.com/pt-main/manage/blob/main/LICENCE).

---

**manage** написан с использованием [tap](https://github.com/pt-main/tap), [tessera](https://github.com/pt-main/tessera) и [tycl](https://github.com/pt-main/tycl).  
Автор: Pt.
````

## `README.md`

````markdown
# manage - command-line task manager

```bash
go install github.com/pt-main/manage/cmd/manage@latest
```

**Manage** is a CLI tool for handling tasks organised in a hierarchical structure: topics (projects, areas) contain tasks. Each topic has a priority, tags, and a task list; each task has a completion status, a description, and a timestamp. Named tables are also supported for grouping topics. The tool is domain-agnostic and fits any scenario, from personal errands to work projects.

---

## Quick start

Before the first use, run initialisation:

```bash
manage init
```

This creates a `~/manager_config.tess` file with default settings.

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

Or do latest:

```bash
manage do backend
```

---

## Core concepts

- **Topic** - a group of tasks (project, module, client).  
  A topic has a name, priority, tags, and a task list.
- **Task** - a unit of work.  
  A task has an ID (auto-incremented within the topic), a description, a status (done/not done), and a creation timestamp.
- **Priority** - any string without spaces (e.g., `A`, `high`, `critical`).  
  The sorting order is defined in the configuration.
- **Tags** - labels for grouping topics.
- **Table** - a named list of topics for quick filtering.

---

## Commands

### `init`

```bash
manage init
```

Creates the configuration file `~/manager_config.tess`.  
Without this, commands won't work. The file stores tasks, settings, and tables in Tessera format.

**Warning:** running `init` when a config already exists overwrites it.

---

### `add`

Adds a topic, task, tags, or table.

#### Adding a topic

```bash
manage add theme <name> [--priority=<value>] [--tags="tag1 tag2 ..."] [--force]
```

- `name` - topic name (required).
- `--priority` - arbitrary string without spaces. If omitted or empty, the `default_priority` from settings is used (default `F`).
- `--tags` - space-separated list of tags. If omitted, no tags are added.
- `--force` - replace an existing topic with the same name. Without this flag, the command fails if the topic already exists (data is untouched).

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

- Without `--id` - creates a new task with a new ID. `--force` has no effect in this case.
- With `--id` and without `--force` - appends the text to the existing task.
- With `--id` and `--force` - replaces the existing task with the new text.

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

- Without `--force` - adds topics to an existing table (creates the table if it doesn't exist).
- With `--force` - replaces the whole table content with the new list.

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
manage do <theme> [id...] [--delete]
```

- Without `--delete` - marks the task as done, but keeps it in the list.
- With `--delete` - deletes the task.
- Without `id...` - marks the last task of the topic.

If the topic has no tasks, or some of the given IDs do not exist, an error is printed. Other tasks are still processed.

Example:

```bash
manage do backend 1 2 3
manage do backend 1 --delete
manage do backend # do latest task
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
manage list [--priority="A B"] [--tags="tag1 tag2"] [--themes="theme1 theme2"] [--table=<name>] [--state=(todo|done|all)] [--todo] [--done] [--filter=(and|or)] [--code]
```

#### Filters

- `--priority` - space-separated list of priorities.
- `--tags` - space-separated list of tags.
- `--themes` - space-separated list of topic names.
- `--table` - table name (shows topics from that table).
- `--state` - task status:
  - `todo` - only unfinished tasks.
  - `done` - only completed tasks.
  - `all` - all (default).
- `--todo` - shorthand for `--state=todo`.
- `--done` - shorthand for `--state=done`.
- `--filter` - logic for combining conditions:
  - `and` - all conditions must match (default).
  - `or` - any condition is sufficient.
- `--code` - additionally include `TODO:` comments found in the current directory.

`--table` and `--themes` are combined into a single topic-name filter.

If both `--todo` and `--done` are given, `all` is used with a warning.

#### Sorting

Topics are displayed in the priority order defined in the settings (default from `A` to `Z`).  
Tasks are shown grouped by their topic.

#### Examples

```bash
manage list
manage list --priority="A B" --tags="go api" --filter=or
manage list --table=monday --todo
manage list --done --tags=exam
manage list --code
```

---

### `sys`

System information and debug commands.

```bash
manage sys version
manage sys conf
manage sys conf-path
```

- `version` - prints Tap and Manage versions.
- `conf` - prints the raw contents of the current configuration file.
- `conf-path` - prints the configuration file path (`~/manager_config.tess`).

---

## Configuration

The file `~/manager_config.tess` is created by `init`.  
It uses Tessera format. You can edit it manually for bulk changes, but usually all operations are done via the CLI.

Settings:

- `priority_order` - the priority order (default A to Z).
- `default_priority` - the value used when adding a topic without `--priority`.

---

## License

MIT - details in [LICENCE](https://github.com/pt-main/manage/blob/main/LICENCE).

---

**manage** is written using [tap](https://github.com/pt-main/tap), [tessera](https://github.com/pt-main/tessera), and [tycl](https://github.com/pt-main/tycl).  
Author: Pt.
````

## `build.json`

```json
{
  "project_path": "./cmd/manage",
  "output_dir": "../build",
  "name_template": "manage-v{version}-{os}-{arch}",
  "version": "1.2.0",
  "clean_before_build": true,
  "cgo_enabled": 0,
  "go_build_args": ["-trimpath", "-buildvcs=false", "-gcflags=all=-l"],
  "ldflags": "-s -w -X main.version=${version} -X main.commit=${commit} -X main.buildTime=${date}T${time}",
  "tags": "",
  "platforms": "all",
  "exclude_platforms": ["android/*", "ios/*", "plan9/*"],
  "platform_config": {
    "linux/arm": {
      "goarm": "7",
      "ldflags": "-s -w -X main.arch=armv7"
    },
    "linux/amd64": {
      "goamd64": "v3",
      "ldflags": "-s -w -X main.arch=amd64v3"
    },
    "windows/amd64": {
      "cgo_enabled": 1,
      "tags": "windows",
      "env": {
        "CC": "x86_64-w64-mingw32-gcc"
      }
    },
    "darwin/amd64": {
      "goamd64": "v2",
      "env": {
        "CGO_ENABLED": "1",
        "SDKROOT": "/path/to/macOS.sdk"
      }
    }
  },
  "verbose": true
}
```

## `go.mod`

```
module github.com/pt-main/manage

go 1.24.13

require (
	github.com/pt-main/tap/go v1.5.9
	github.com/pt-main/tessera/impl/go v0.0.0-20261019162208-a9b257904ffd
	github.com/pt-main/tycl v1.3.8
)

require (
	github.com/eiannone/keyboard v0.0.0-20220611211555-0d226195f203 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
)
```

## `go.sum`

```
github.com/eiannone/keyboard v0.0.0-20220611211555-0d226195f203 h1:XBBHcIb256gUJtLmY22n99HaZTz+r2Z51xUPi01m3wg=
github.com/eiannone/keyboard v0.0.0-20220611211555-0d226195f203/go.mod h1:E1jcSv8FaEny+OP/5k9UxZVw9YFWGj7eI4KR/iOBqCg=
github.com/pt-main/tap/go v1.5.9 h1:u/CZDvNfOQUcGbi87XEVhuRgIUun64TQxoO2GyJTMdY=
github.com/pt-main/tap/go v1.5.9/go.mod h1:JhdaGAsrmcZVFSANOv8PlSuKcF0Uqx0Kad2XakiOPlM=
github.com/pt-main/tessera/impl/go v0.0.0-20261019162208-a9b257904ffd h1:+7uvVZY/7UGvGzltbaDR9nrwMdG+RpDaglNJZIIWefM=
github.com/pt-main/tessera/impl/go v0.0.0-20261019162208-a9b257904ffd/go.mod h1:T/G2WSiqoy6OQ3u6F4HDd+ZWHdK39MUnxv/sbslwMVA=
github.com/pt-main/tycl v1.3.8 h1:TU7/axhLYE7l3TGmbvth3AfGQOorwI5uDnSO8FxeD54=
github.com/pt-main/tycl v1.3.8/go.mod h1:zKmXw4/TpgHkz+gYo0D88HtofifVkcdcF3pOWSc4IIQ=
golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a h1:dGzPydgVsqGcTRVwiLJ1jVbufYwmzD3LfVPLKsKg+0k=
golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
```

## `cmd/manage/add.go`

```go
package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/pt-main/manage/manager"
	tap "github.com/pt-main/tap/go"
)

func Add(p *tap.Parser, s []string) error {
	m := manager.NewManager()
	if err := Load(m); err != nil {
		return err
	}
	name := s[1]
	_, force := p.Flags["force"]

	lenArgs := func(n int) error {
		if len(s) != n {
			return fmt.Errorf(
				"❌ Invalid argument length for add: need %v, got %v",
				n, len(s))
		}
		return nil
	}

	switch s[0] {
	case "theme":
		if err := lenArgs(2); err != nil {
			return err
		}
		_tags, ok := p.Flags["tags"]
		tags := strings.Split(_tags, " ")
		if !ok {
			tags = make([]string, 0)
		}
		priority, ok := p.Flags["priority"]
		if priority == "" {
			priority = m.Settings.DefaultPriority
		}
		t := manager.NewTask(tags, priority, name)
		if err := m.AddTheme(name, t, force); err != nil {
			return err
		}
		return Save(m)

	case "task":
		if err := lenArgs(3); err != nil {
			return err
		}

		text := s[2]

		task, ok := m.Tasks[name]
		if !ok {
			return fmt.Errorf("❌ Invalid theme name: %v", name)
		}

		append := true
		var to = task.Last
		_to, ok := p.Flags["id"]
		if ok {
			var err error
			to, err = strconv.Atoi(_to)
			if err != nil {
				return err
			}
			append = false
		}

		if append && !force {
			if err := m.AddTask(name, text); err != nil {
				return err
			}
		} else if !append && force {
			added := false
			for _, task := range task.Tasks {
				if task.GetId() == to {
					added = true
					if !task.SetAbout(text) {
						return fmt.Errorf("❌ Can't edit '%v' task about", task.Type())
					}
				}
			}
			if !added {
				return fmt.Errorf("❌ Invalid id: %v", to)
			}
		} else if !append && !force {
			added := false
			for _, task := range task.Tasks {
				if task.GetId() == to {
					added = true
					if !task.Add(text) {
						return fmt.Errorf("❌ Can't edit '%v' task about", task.Type())
					}
				}
			}
			if !added {
				return fmt.Errorf("❌ Invalid id: %v", to)
			}
		}

	case "tag":
		if len(s) < 3 {
			return fmt.Errorf(
				"❌ Invalid argument format for 'add tag': need 3 or more arguments, got %v",
				len(s))
		}
		tags := s[2:]

		task, ok := m.Tasks[name]
		if !ok {
			return fmt.Errorf("❌ Invalid theme name: %v", name)
		}

		add := []string{}
		rm := []string{}
		for _, tag := range tags {
			if strings.HasPrefix(tag, "!") {
				rm = append(rm, tag[1:])
			} else {
				add = append(add, tag)
			}
		}

		newTags := []string{}
		for _, tag := range task.Tags {
			if !slices.Contains(rm, tag) {
				newTags = append(newTags, tag)
			}
		}

		for _, addtag := range add {
			if !slices.Contains(newTags, addtag) {
				newTags = append(newTags, addtag)
			}
		}

		task.Tags = newTags

	case "table":
		if len(s) < 3 {
			return fmt.Errorf(
				"❌ Invalid argument format for 'add table': need 3 or more arguments, got %v",
				len(s))
		}

		name := s[1]
		themes := s[2:]
		_, force := p.Flags["force"]

		t, ok := m.Tables[name]
		if !ok {
			force = true
		}

		if force {
			m.Tables[name] = themes
		} else {
			t = append(t, themes...)
			m.Tables[name] = t
		}

	default:
		return fmt.Errorf("❌ Invalid mode: %v", s[0])
	}
	return Save(m)
}
```

## `cmd/manage/handlers.go`

```go
package main

import (
	"fmt"
	"strconv"

	"github.com/pt-main/manage/manager"
	tap "github.com/pt-main/tap/go"
	"github.com/pt-main/tap/go/color"
)

func Do(p *tap.Parser, s []string) error {
	m := manager.NewManager()
	if err := Load(m); err != nil {
		return err
	}
	theme, ok := m.Tasks[s[0]]
	if !ok {
		return fmt.Errorf("Invalid theme name: %v", s[0])
	}

	_, delete := p.Flags["delete"]

	hasError := false
	tasks := []int{}

	do := func(id int) {
		done := false
		newTasks := []manager.TaskData{}
		for _, task := range theme.Tasks {
			if task.GetId() == id {
				done = true
				task.Do()
				if !delete {
					newTasks = append(newTasks, task)
				}
			} else {
				newTasks = append(newTasks, task)
			}
		}

		theme.Tasks = newTasks

		if !done {
			hasError = true
			tasks = append(tasks, id)
			color.PrintlnColored(
				"[?RD]🔴 Task not found[?RT]: [?RD]Invalid task id[?RT]: %v",
				id)
		}
	}

	if len(s) == 1 {
		tasks := theme.Tasks
		if len(tasks) == 0 {
			color.PrintColored("[?RD]🔴 Tasks not found")
			hasError = true
		} else {
			last := tasks[len(tasks)-1]
			do(last.GetId())
		}
	}

	for _, taskId := range s[1:] {
		id, err := strconv.Atoi(taskId)
		if err != nil {
			return err
		}
		do(id)
	}

	if hasError {
		color.PrintColored("[?RD]🔴 Failed to do tasks[?RT]")
		if len(tasks) != 0 {
			color.PrintlnColored("%v", tasks)
		} else {
			print("\n")
		}
	} else {
		color.PrintlnColored("[?GN]✅ Tasks done![?RT]")
	}
	return Save(m)
}

func Set(p *tap.Parser, s []string) error {
	m := manager.NewManager()
	if err := Load(m); err != nil {
		return err
	}
	switch s[0] {
	case "priority":
		if len(s) != 3 {
			return fmt.Errorf("Invalid args length: need 3, got %v", len(s))
		}
		theme, ok := m.Tasks[s[1]]
		if !ok {
			return fmt.Errorf("Theme not found: Invalid theme name: '%v'", theme)
		}
		priority := s[2]
		theme.Priority = priority
	default:
		return fmt.Errorf("Invalid set command: %v", s[0])
	}
	return Save(m)
}

type listTheme struct {
	text     string
	priority string
}
```

## `cmd/manage/list.go`

```go
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

	if _, has := p.Flags["code"]; has {
		for _, task := range manager.FindTodosInDir(".") {
			textRaw[strconv.Itoa(task.Id)] = &listTheme{
				priority: m.Settings.PriorityOrder[0],
				text:     task.Render(color.ColorEnabled),
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
			text := "- " + strings.ReplaceAll(task.Render(color.ColorEnabled), "\n", "\n  ")
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
```

## `cmd/manage/main.go`

```go
package main

import (
	"fmt"

	"github.com/pt-main/manage/manager"
	tap "github.com/pt-main/tap/go"
	"github.com/pt-main/tap/go/utils"
)

var Version = "1.2.0"

const docInit = `[?GN]Initialize[?RT] config file.

[?BBE]Usage:[?RT]
  [?YW]manage init[?RT]

[?BBE]Description:[?RT]
  Creates [?BBK]~/manager_config.tess[?RT] with default settings.
  This is the first command you need to run before using any other commands.
  Warning: overwrites an existing config file.`

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
  [?YW]manage do[?RT] [?BGN]<theme>[?RT] [?BGN][task-id...][?RT] [?BBK]--delete[?RT]

[?BBE]Description:[?RT]
  Marks one or more tasks as completed. With [?BBK]--delete[?RT], completed tasks are removed from the list;
  otherwise they stay visible with a "done" mark.

[?BBE]Examples:[?RT]
  [?BE]manage do math 1 2 3[?RT]
  [?BE]manage do math 1 --delete[?RT]
  [?BE]manage do math [?BBK]# do last task`

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

	p.AddCommand("do", Do, docDo, []string{"theme"}, nil, true)

	p.AddCommand("set", Set, docSet, nil, nil, true)

	p.AddCommand("list", List, docList, nil, nil, false)

	sys := tap.NewParser("sys", "system info and debug", []string{"-h", "help"}, tap.DefaultParserConfig())

	sys.AddCommand("version", func(p *tap.Parser, s []string) error {
		fmt.Printf("Tap v%v, Manage v%v\n", tap.Version, Version)
		return nil
	}, "Show version", nil, nil, false)

	sys.AddCommand("conf", func(p *tap.Parser, s []string) error {
		file, err := GetConf()
		if err != nil {
			return err
		}
		println(file)
		return nil
	}, "Show config", nil, nil, false)

	sys.AddCommand("conf-path", func(p *tap.Parser, s []string) error {
		println("~/" + ConfigFile)
		return nil
	}, "Show config", nil, nil, false)

	p.AddSubcommand("sys", sys)

	return p
}

func main() {
	err := NewCli().Main()
	if err != nil {
		fmt.Println(err)
	}
}
```

## `cmd/manage/utils.go`

```go
package main

import (
	"fmt"
	"os"

	"github.com/pt-main/manage/config"
	"github.com/pt-main/manage/manager"
	"github.com/pt-main/tycl/utils"
)

func Hd() string {
	hd, err := os.UserHomeDir()
	if err != nil {
		panic("Can't find to user home dir: " + err.Error())
	}
	return hd
}

var backTo = ""

func ToHd() {
	var err error
	backTo, err = os.Getwd()
	if err != nil {
		panic("Can't get current dir: " + err.Error())
	}
	err = os.Chdir(Hd())
	if err != nil {
		panic("Can't go to user home dir: " + err.Error())
	}
}

func BackFromHd() {
	err := os.Chdir(backTo)
	if err != nil {
		panic("Can't back to previous dir: " + err.Error())
	}
	backTo = ""
}

var ConfigFile = "manager_config.tess"

func Save(m *manager.Manager) error {
	file, err := config.GenCfg(m)
	if err != nil {
		return fmt.Errorf("❗️ Save data: %v", err)
	}
	ToHd()
	err = utils.WriteF(ConfigFile, file)
	BackFromHd()
	return err
}

func GetConf() (string, error) {
	ToHd()
	file, err := utils.OpenF(ConfigFile)
	if err != nil {
		return "", fmt.Errorf("❗️ Load data: %v. Do you forget to use init?", err)
	}
	BackFromHd()
	return file, nil
}

func Load(m *manager.Manager) error {
	file, err := GetConf()
	if err != nil {
		return err
	}
	err = config.GenManager(file, m)
	if err != nil {
		return err
	}
	return nil
}
```

## `manager/find.go`

```go
package manager

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pt-main/tap/go/color"
)

const todoMarker = "TODO:"

var lowerTodoMarker = strings.ToLower(todoMarker)

type CodeTask struct {
	Id       int
	About    string
	Done     bool
	CodeLine int
	File     string
}

func (ct *CodeTask) Type() TaskDataType { return CodeTaskT }

func (ct *CodeTask) Render(colored bool) string {
	res := fmt.Sprintf(
		"[?BBK]CODE:[?BYW]TODO [?CN][f:[?GN]%v[?CN]] [?BBK](l:[?GN]%v[?BBK])[?RT]:",
		ct.File, ct.CodeLine,
	)
	if colored {
		res = color.Set(res)
	} else {
		res = color.ReplaceColors(res)
	}
	return res + "\n" + ct.About
}

func (ct *CodeTask) GetId() int       { return ct.Id }
func (ct *CodeTask) IsDone() bool     { return ct.Done }
func (ct *CodeTask) GetAbout() string { return ct.About }

func (ct *CodeTask) GetDate() string { return "" }

func (ct *CodeTask) Add(string) bool      { return false }
func (ct *CodeTask) SetAbout(string) bool { return false }

func (ct *CodeTask) Do() error {
	if ct.File == "" {
		return errors.New("code task: empty file path")
	}

	info, err := os.Stat(ct.File)
	if err != nil {
		return fmt.Errorf("code task: stat %s: %w", ct.File, err)
	}

	data, err := os.ReadFile(ct.File)
	if err != nil {
		return fmt.Errorf("code task: read %s: %w", ct.File, err)
	}

	content := string(data)
	crlf := strings.Contains(content, "\r\n")
	if crlf {
		content = strings.ReplaceAll(content, "\r\n", "\n")
	}

	lines := strings.Split(content, "\n")
	idx := ct.CodeLine - 1
	if idx < 0 || idx >= len(lines) {
		return fmt.Errorf("code task: line %d out of range in %s", ct.CodeLine, ct.File)
	}

	newLine, ok := stripTodoMarker(lines[idx], ct.About)
	if !ok {
		return fmt.Errorf(
			"code task: TODO not found (or content changed) at %s:%d",
			ct.File, ct.CodeLine,
		)
	}
	lines[idx] = newLine

	result := strings.Join(lines, "\n")
	if crlf {
		result = strings.ReplaceAll(result, "\n", "\r\n")
	}
	if err := os.WriteFile(ct.File, []byte(result), info.Mode()); err != nil {
		return fmt.Errorf("code task: write %s: %w", ct.File, err)
	}

	ct.Done = true
	return nil
}

func stripTodoMarker(line, expectedAbout string) (string, bool) {
	positions := findAllTodoPositions(line)
	if len(positions) == 0 {
		return line, false
	}

	pos := positions[0]
	if expectedAbout != "" {
		pos = -1
		for _, p := range positions {
			if strings.TrimSpace(line[p+len(todoMarker):]) == expectedAbout {
				pos = p
				break
			}
		}
		if pos == -1 {
			return line, false
		}
	}

	before := strings.TrimRight(line[:pos], " \t")
	after := strings.TrimLeft(line[pos+len(todoMarker):], " \t")

	switch {
	case before == "":
		return after, true
	case after == "":
		return before, true
	default:
		return before + " " + after, true
	}
}

func findAllTodoPositions(line string) []int {
	lower := strings.ToLower(line)
	var positions []int
	for from := 0; ; {
		i := strings.Index(lower[from:], lowerTodoMarker)
		if i == -1 {
			return positions
		}
		positions = append(positions, from+i)
		from += i + len(lowerTodoMarker)
	}
}

var textExtensions = map[string]struct{}{
	".go": {}, ".py": {}, ".js": {}, ".jsx": {}, ".ts": {}, ".tsx": {},
	".rs": {}, ".c": {}, ".h": {}, ".cc": {}, ".cpp": {}, ".cxx": {},
	".hpp": {}, ".hxx": {}, ".java": {}, ".rb": {}, ".php": {}, ".cs": {},
	".swift": {}, ".kt": {}, ".kts": {}, ".sh": {}, ".bash": {}, ".zsh": {},
	".fish": {}, ".md": {}, ".txt": {}, ".rst": {}, ".org": {},
	".yaml": {}, ".yml": {}, ".toml": {}, ".json": {}, ".jsonc": {}, ".xml": {},
	".html": {}, ".htm": {}, ".css": {}, ".scss": {}, ".sass": {}, ".less": {},
	".sql": {}, ".lua": {}, ".pl": {}, ".pm": {}, ".r": {}, ".scala": {}, ".sc": {},
	".ex": {}, ".exs": {}, ".erl": {}, ".hrl": {}, ".hs": {}, ".lhs": {},
	".ml": {}, ".mli": {}, ".clj": {}, ".cljs": {}, ".vim": {},
	".dart": {}, ".groovy": {}, ".gradle": {}, ".cmake": {},
}

var skipDirs = map[string]struct{}{
	".git": {}, ".svn": {}, ".hg": {},
	"node_modules": {}, "vendor": {}, "target": {}, "build": {}, "dist": {},
	".idea": {}, ".vscode": {}, ".vs": {},
	"__pycache__": {}, ".venv": {}, "venv": {}, "env": {},
	"bin": {}, "obj": {}, ".gradle": {}, ".mypy_cache": {}, ".pytest_cache": {},
}

func FindTodosInDir(dir string) []*CodeTask {
	result := make([]*CodeTask, 0)

	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if _, skip := skipDirs[d.Name()]; skip && path != dir {
				return filepath.SkipDir
			}
			return nil
		}
		if !isTextFile(path) {
			return nil
		}
		scanFileForTodos(path, &result)
		return nil
	})

	sort.Slice(result, func(i, j int) bool {
		if result[i].File != result[j].File {
			return result[i].File < result[j].File
		}
		return result[i].CodeLine < result[j].CodeLine
	})
	for i, t := range result {
		t.Id = i
	}

	return result
}

func isTextFile(path string) bool {
	if _, ok := textExtensions[strings.ToLower(filepath.Ext(path))]; ok {
		return true
	}
	switch strings.ToLower(filepath.Base(path)) {
	case "makefile", "dockerfile", "rakefile", "gemfile", "cargo.lock",
		"go.mod", "go.sum", "requirements.txt",
		".gitignore", ".dockerignore", ".editorconfig", ".env":
		return true
	}
	return false
}

func scanFileForTodos(path string, out *[]*CodeTask) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 2*1024*1024)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		positions := findAllTodoPositions(line)
		if len(positions) == 0 {
			continue
		}
		pos := positions[0]
		about := strings.TrimSpace(line[pos+len(todoMarker):])
		*out = append(*out, &CodeTask{
			About:    about,
			CodeLine: lineNum,
			File:     path,
		})
	}
}
```

## `manager/main.go`

```go
package manager

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/pt-main/tap/go/color"
)

type TaskDataType string

var (
	DefaultTaskT TaskDataType = "DEFAULT"
	CodeTaskT    TaskDataType = "CODE"
)

type TaskData interface {
	Type() TaskDataType
	Render(colored bool) string
	GetId() int
	Do() error
	IsDone() bool
	Add(s string) bool
	GetAbout() string
	SetAbout(string) bool
	GetDate() string
}

type OneTask struct {
	Id        int
	About     string
	Done      bool
	CreatedAt string
}

func (ot *OneTask) GetAbout() string {
	return ot.About
}

func (ot *OneTask) SetAbout(s string) bool {
	ot.About = s
	return true
}

func (ot *OneTask) GetDate() string {
	return ot.CreatedAt
}

func (ot *OneTask) Type() TaskDataType {
	return DefaultTaskT
}

func (ot *OneTask) Render(colored bool) string {
	status := "TOOO"
	if ot.Done {
		status = "DONE"
	}
	res := fmt.Sprintf(
		"[?BYW]%v [?CN][[?GN]%v[?CN]] [?BBK]([?GN]%v[?BBK])[?RT]:",
		status, ot.Id, ot.CreatedAt,
	)
	if colored {
		res = color.Set(res)
	} else {
		res = color.ReplaceColors(res)
	}
	return res + "\n" + ot.About
}

func (ot *OneTask) GetId() int {
	return ot.Id
}

func (ot *OneTask) Do() error {
	ot.Done = true
	return nil
}

func (ot *OneTask) IsDone() bool {
	return ot.Done
}

func (ot *OneTask) Add(s string) bool {
	ot.About += s
	return true
}

type Task struct {
	Tags     []string
	Priority string
	Tasks    []TaskData
	Name     string
	Last     int
}

func (task *Task) GetTaskById(id int) TaskData {
	for _, task := range task.Tasks {
		if task.GetId() == id {
			return task
		}
	}
	return nil
}

func (task *Task) RmTags(tags []string) *Task {
	newTags := []string{}
	for _, tag := range task.Tags {
		if !slices.Contains(tags, tag) {
			newTags = append(newTags, tag)
		}
	}
	task.Tags = newTags
	return task
}

func (task *Task) AddTags(tags []string) *Task {
	task.Tags = append(task.Tags, tags...)
	return task
}

func (task *Task) SetPriority(priority string) *Task {
	task.Priority = priority
	return task
}

func (task *Task) DoTask(n int) error {
	for _, t := range task.Tasks {
		if t.GetId() == n {
			return t.Do()
		}
	}
	return fmt.Errorf("can't find task %d", n)
}

func (task *Task) AppendToTask(n int, about string) error {
	appended := false
	for _, task := range task.Tasks {
		if task.GetId() == n {
			appended = true
			if ok := task.Add(about); !ok {
				return fmt.Errorf("Can't add about to task: %v", n)
			}
			break
		}
	}
	if !appended {
		return fmt.Errorf("Invalid id: %v", n)
	}
	return nil
}

func NewTask(tags []string, priority, name string) *Task {
	return &Task{
		Tags:     tags,
		Priority: priority,
		Name:     name,
		Last:     0,
		Tasks:    make([]TaskData, 0),
	}
}

type Settings struct {
	PriorityOrder   []string
	DefaultPriority string
}

type Manager struct {
	Tasks    map[string]*Task
	Tables   map[string][]string
	Settings *Settings
}

func (m *Manager) hasTheme(theme string) error {
	if !slices.Contains(slices.Collect(maps.Keys(m.Tasks)), theme) {
		return fmt.Errorf("Invalid theme name: %v", theme)
	}
	return nil
}

func (m *Manager) GetTheme(theme string) (*Task, error) {
	if err := m.hasTheme(theme); err != nil {
		return nil, err
	}
	return m.Tasks[theme], nil
}

func (m *Manager) AddTheme(name string, task *Task, force bool) error {
	if err := m.hasTheme(name); err == nil && !force {
		return fmt.Errorf(
			"Can't replace theme (theme %v is already created), use force flag for replace",
			force)
	}
	m.Tasks[name] = task
	return nil
}

func (m *Manager) addTask(theme string, task TaskData) error {
	if err := m.hasTheme(theme); err != nil {
		return err
	}
	m.Tasks[theme].Tasks = append(m.Tasks[theme].Tasks, task)
	return nil
}

func (m *Manager) AddTask(theme string, about string) error {
	if err := m.hasTheme(theme); err != nil {
		return err
	}
	m.Tasks[theme].Last += 1
	m.addTask(theme, &OneTask{
		Id:        m.Tasks[theme].Last,
		About:     about,
		CreatedAt: time.Now().Format("2006.01.02 15:04:05"),
	})
	return nil
}

const (
	StateDone = "done"
	StateTodo = "todo"
	StateAll  = "all"

	FilteringAnd = "and"
	FilteringOr  = "or"
)

func (m *Manager) Filter(tags, priority, table []string, state, filtering string) []*Task {
	if state == "" {
		state = StateAll
	}
	if !slices.Contains([]string{StateDone, StateTodo, StateAll}, state) {
		new := StateAll
		fmt.Printf("Invalid state: '%v', converting into '%v'\n", state, new)
		state = new
	}
	if filtering == "" {
		filtering = FilteringAnd
	}
	if !slices.Contains([]string{FilteringAnd, FilteringOr}, filtering) {
		new := FilteringAnd
		fmt.Printf("Invalid filtering mode: '%v', converting into '%v'\n", filtering, new)
		filtering = new
	}
	if tags == nil {
		tags = make([]string, 0)
	}
	if priority == nil {
		priority = make([]string, 0)
	}
	if table == nil {
		table = make([]string, 0)
	}
	defaultValue := true
	if filtering == FilteringOr {
		defaultValue = false
	}
	filtered := []*Task{}
	for name, task := range m.Tasks {
		add := false
		activated := false
		if len(tags) != 0 || len(priority) != 0 || len(table) != 0 {
			activated = true
		}
		c0 := false
		if len(tags) > 0 {
			for _, tag := range tags {
				if slices.Contains(task.Tags, tag) {
					c0 = true
				} else if filtering == FilteringAnd {
					c0 = false
					break
				}
			}
		} else {
			c0 = defaultValue
		}
		c1 := defaultValue
		if len(table) != 0 {
			c1 = slices.Contains(table, name)
		}
		c2 := defaultValue
		if len(priority) != 0 {
			c2 = slices.Contains(priority, task.Priority)
		}

		if filtering == FilteringAnd {
			if c0 && c1 && c2 {
				add = true
			}
		} else {
			if c0 || c1 || c2 {
				add = true
			}
		}

		if activated && add || !activated {
			tasks := []TaskData{}
			for _, oneT := range task.Tasks {
				switch state {
				case "todo":
					if !oneT.IsDone() {
						tasks = append(tasks, oneT)
					}
				case "done":
					if oneT.IsDone() {
						tasks = append(tasks, oneT)
					}
				default:
					tasks = append(tasks, oneT)
				}
			}
			task.Tasks = tasks
			filtered = append(filtered, task)
		}
	}
	return filtered
}

func NewManager() *Manager {
	return &Manager{
		Tasks:  make(map[string]*Task),
		Tables: make(map[string][]string),
		Settings: &Settings{
			PriorityOrder:   strings.Split("A B C D E F G H I J K L M N O P Q R S T U V W X Y Z", " "),
			DefaultPriority: "F",
		},
	}
}
```

## `config/main.go`

```go
package config

import (
	"fmt"
	"strconv"

	"github.com/pt-main/manage/manager"
	tessera "github.com/pt-main/tessera/impl/go"
)

func Load(cfg *tessera.Config, m *manager.Manager) {
	if m.Settings == nil {
		m.Settings = &manager.Settings{}
	}
	if m.Tables == nil {
		m.Tables = make(map[string][]string)
	}
	if m.Tasks == nil {
		m.Tasks = make(map[string]*manager.Task)
	}

	for i := 0; ; i++ {
		key := strconv.Itoa(i)
		taskCfg, ok := cfg.Inners[key]
		if !ok {
			break
		}

		var innerTasks []manager.TaskData
		for j := 0; ; j++ {
			innerKey := strconv.Itoa(j)
			innerCfg, ok := taskCfg.Inners[innerKey]
			if !ok {
				break
			}
			innerTasks = append(innerTasks, &manager.OneTask{
				Id:        int(innerCfg.Ints["id"]),
				About:     innerCfg.Strs["about"],
				Done:      innerCfg.Logicals["done"],
				CreatedAt: innerCfg.Strs["created_at"],
			})
		}

		var tags []string
		for j := 0; ; j++ {
			tagKey := strconv.Itoa(j)
			tagCfg, ok := taskCfg.Inners["tags_"+tagKey]
			if !ok {
				break
			}
			tags = append(tags, tagCfg.Strs["value"])
		}

		name := taskCfg.Strs["name"]
		m.Tasks[name] = &manager.Task{
			Name:     name,
			Tags:     tags,
			Priority: taskCfg.Strs["priority"],
			Last:     int(taskCfg.Ints["last"]),
			Tasks:    innerTasks,
		}
	}

	for i := 0; ; i++ {
		key := strconv.Itoa(i)
		tableCfg, ok := cfg.Inners["table_"+key]
		if !ok {
			break
		}
		var subjects []string
		for j := 0; ; j++ {
			subjKey := strconv.Itoa(j)
			subjCfg, ok := tableCfg.Inners["tasks_"+subjKey]
			if !ok {
				break
			}
			subjects = append(subjects, subjCfg.Strs["value"])
		}
		m.Tables[tableCfg.Strs["name"]] = subjects
	}

	if sett, ok := cfg.Inners["settings"]; ok {
		var priorityOrder []string
		for j := 0; ; j++ {
			pKey := strconv.Itoa(j)
			pCfg, ok := sett.Inners["priority_order_"+pKey]
			if !ok {
				break
			}
			priorityOrder = append(priorityOrder, pCfg.Strs["value"])
		}
		m.Settings = &manager.Settings{
			PriorityOrder:   priorityOrder,
			DefaultPriority: sett.Strs["default_priority"],
		}
	}
}

func Save(m *manager.Manager) *tessera.Config {
	cfg := &tessera.Config{
		Inners: make(map[string]*tessera.Config),
	}

	i := 0
	for _, task := range m.Tasks {
		taskCfg := &tessera.Config{
			Strs: map[string]string{
				"name":     task.Name,
				"priority": task.Priority,
			},
			Ints:   map[string]int64{"last": int64(task.Last)},
			Inners: make(map[string]*tessera.Config),
		}

		for j, tag := range task.Tags {
			tagCfg := &tessera.Config{Strs: map[string]string{"value": tag}}
			taskCfg.Inners["tags_"+strconv.Itoa(j)] = tagCfg
		}

		for j, one := range task.Tasks {
			innerCfg := &tessera.Config{
				Ints:     map[string]int64{"id": int64(one.GetId())},
				Strs:     map[string]string{"about": one.GetAbout(), "created_at": one.GetDate()},
				Logicals: map[string]bool{"done": one.IsDone()},
			}
			taskCfg.Inners[strconv.Itoa(j)] = innerCfg
		}

		cfg.Inners[strconv.Itoa(i)] = taskCfg
		i++
	}

	i = 0
	for name, subjects := range m.Tables {
		tableCfg := &tessera.Config{
			Strs:   map[string]string{"name": name},
			Inners: make(map[string]*tessera.Config),
		}
		for j, subj := range subjects {
			subjCfg := &tessera.Config{Strs: map[string]string{"value": subj}}
			tableCfg.Inners["tasks_"+strconv.Itoa(j)] = subjCfg
		}
		cfg.Inners["table_"+strconv.Itoa(i)] = tableCfg
		i++
	}

	if m.Settings != nil {
		settCfg := &tessera.Config{
			Strs:   map[string]string{"default_priority": m.Settings.DefaultPriority},
			Inners: make(map[string]*tessera.Config),
		}
		for j, p := range m.Settings.PriorityOrder {
			pCfg := &tessera.Config{Strs: map[string]string{"value": p}}
			settCfg.Inners["priority_order_"+strconv.Itoa(j)] = pCfg
		}
		cfg.Inners["settings"] = settCfg
	}

	return cfg
}

func GenCfg(m *manager.Manager) (string, error) {
	return tessera.TryToTessera(Save(m))
}

func GenManager(code string, m *manager.Manager) error {
	cfg, errs := tessera.Parse(code)
	if len(errs) > 0 {
		return fmt.Errorf("tessera parse errors: %v", errs)
	}
	Load(cfg, m)
	return nil
}
```

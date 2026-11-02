# Merged files

Всего файлов: 18

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

Добавляет тему, задачу, теги, таблицу или приоритеты.

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

- Без `--id` - создаётся новая задача с новым ID. `--force` в этом случае не применяется и вызывает ошибку.
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

Добавляет теги к теме. Уже существующие теги не дублируются.

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

#### Добавление приоритетов

```bash
manage add priority <p1> [p2 ...]
```

Добавляет приоритеты в конец `priority_order`, если их там ещё нет.

Примеры:

```bash
manage add priority S
manage add priority S T U
```

---

### `remove`

Удаляет темы, задачи, теги, таблицы и приоритеты.

#### Удаление темы

```bash
manage remove theme <name>
```

Удаляет тему вместе со всеми задачами. Тема также удаляется из всех таблиц.

#### Удаление задачи

```bash
manage remove task <theme> <id...>
```

Удаляет указанные задачи. ID оставшихся задач не перенумеровываются.

Если ни один из указанных ID не найден, команда завершается ошибкой.

#### Удаление тегов

```bash
manage remove tag <theme> <tags...>
```

Удаляет теги с темы. Если ни один из указанных тегов не найден, команда завершается ошибкой.

#### Удаление таблицы

```bash
manage remove table <name>
```

Удаляет всю таблицу. Темы внутри неё не затрагиваются.

#### Удаление тем из таблицы

```bash
manage remove from-table <name> <themes...>
```

Удаляет указанные темы из таблицы, не удаляя сами темы.

#### Удаление приоритетов

```bash
manage remove priority <p1> [p2 ...]
```

Удаляет приоритеты из `priority_order`. Нельзя удалить все приоритеты сразу.

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

Меняет приоритет темы. Если нового приоритета нет в `priority_order`, выводится предупреждение - тема не будет отображаться в `list`, пока приоритет не будет добавлен.

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

- `--priority` - список приоритетов через пробел. Значения проверяются по `priority_order`; при отсутствии выводится ошибка.
- `--tags` - список тегов через пробел.
- `--themes` - список имён тем через пробел.
- `--table` - имя таблицы (будут показаны темы из неё). Таблица должна существовать, иначе выводится ошибка.
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

### `show`

Выводит справочную информацию.

```bash
manage show themes
manage show priorities
manage show tags
manage show theme <name>
manage show info
```

- `themes` - все темы с приоритетом, счётчиками задач и тегами.
- `priorities` - приоритеты из `priority_order` и темы в каждом из них. Также показывает приоритеты, которых нет в `priority_order`.
- `tags` - все теги и темы, к которым они относятся.
- `theme <name>` - карточка темы: приоритет, теги, счётчики задач, последний выданный ID, таблицы с этой темой.
- `info` - сводка по приоритетам: все темы сгруппированы по приоритету со счётчиками задач, затем итоговые суммы.

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
````

## `build.json`

```json
{
  "project_path": "./cmd/manage",
  "output_dir": "../build",
  "name_template": "manage-v{version}-{os}-{arch}",
  "version": "1.3.0",
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
	github.com/pt-main/tap/go v1.6.0
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
github.com/pt-main/tap/go v1.6.0 h1:h2ieDRZK9NDKqDsBvKzHmYdyjFk3QFMIlAd44T54zkM=
github.com/pt-main/tap/go v1.6.0/go.mod h1:JhdaGAsrmcZVFSANOv8PlSuKcF0Uqx0Kad2XakiOPlM=
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

// AddTheme: manage add theme <name> [--priority=A] [--tags="go api"] [--force]
func (mc *ManageCli) AddTheme(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) != 1 {
		return fmt.Errorf("❌ 'add theme' needs exactly 1 argument (name), got %v", len(s))
	}
	name := s[0]

	_, force := p.Flags["force"]

	_tags, ok := p.Flags["tags"]
	tags := []string{}
	if ok && _tags != "" {
		tags = strings.Split(_tags, " ")
	}

	priority, ok := p.Flags["priority"]
	if !ok || priority == "" {
		priority = m.Settings.DefaultPriority
	}

	t := manager.NewTask(tags, priority, name)
	if err := m.AddTheme(name, t, force); err != nil {
		return err
	}
	return nil
}

// AddTask: manage add task <theme> <text> [--id=N] [--force]
func (mc *ManageCli) AddTask(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) != 2 {
		return fmt.Errorf("❌ 'add task' needs exactly 2 arguments (theme, text), got %v", len(s))
	}
	name, text := s[0], s[1]
	_, force := p.Flags["force"]

	theme, ok := m.Tasks[name]
	if !ok {
		return fmt.Errorf("❌ Invalid theme name: %v", name)
	}

	appendMode := true
	var to = theme.Last
	if _to, ok := p.Flags["id"]; ok {
		v, err := strconv.Atoi(_to)
		if err != nil {
			return fmt.Errorf("❌ Invalid --id value '%v': %w", _to, err)
		}
		to = v
		appendMode = false
	}

	switch {
	case appendMode && !force:
		if err := m.AddTask(name, text); err != nil {
			return err
		}
	case !appendMode && force:
		found := false
		for _, task := range theme.Tasks {
			if task.GetId() == to {
				found = true
				if !task.SetAbout(text) {
					return fmt.Errorf("❌ Can't edit '%v' task about", task.Type())
				}
			}
		}
		if !found {
			return fmt.Errorf("❌ Invalid id: %v", to)
		}
	case !appendMode && !force:
		found := false
		for _, task := range theme.Tasks {
			if task.GetId() == to {
				found = true
				if !task.Add(text) {
					return fmt.Errorf("❌ Can't edit '%v' task about", task.Type())
				}
			}
		}
		if !found {
			return fmt.Errorf("❌ Invalid id: %v", to)
		}
	case appendMode && force:
		return fmt.Errorf("❌ --force requires --id")
	}

	return nil
}

// AddTag: manage add tag <theme> <tags...> ]
func (mc *ManageCli) AddTag(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) < 2 {
		return fmt.Errorf("❌ 'add tag' needs at least 2 arguments (theme + tags), got %v", len(s))
	}
	name := s[0]
	tags := s[1:]

	theme, ok := m.Tasks[name]
	if !ok {
		return fmt.Errorf("❌ Invalid theme name: %v", name)
	}

	newTags := slices.Clone(theme.Tags)
	for _, t := range tags {
		if !slices.Contains(newTags, t) {
			newTags = append(newTags, t)
		}
	}

	theme.Tags = newTags
	return nil
}

// AddTable: manage add table <name> <themes...> [--force]
func (mc *ManageCli) AddTable(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) < 2 {
		return fmt.Errorf("❌ 'add table' needs at least 2 arguments (name + themes), got %v", len(s))
	}
	name := s[0]
	themes := s[1:]
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

	return nil
}

// AddPriority: manage add priority <p1> [p2 ...]
func (mc *ManageCli) AddPriority(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) < 1 {
		return fmt.Errorf("❌ 'add priority' needs at least 1 argument, got %v", len(s))
	}
	for _, pr := range s {
		if !slices.Contains(m.Settings.PriorityOrder, pr) {
			m.Settings.PriorityOrder = append(m.Settings.PriorityOrder, pr)
		}
	}
	return nil
}
```

## `cmd/manage/do.go`

```go
package main

import (
	"fmt"
	"strconv"

	"github.com/pt-main/manage/manager"
	tap "github.com/pt-main/tap/go"
	"github.com/pt-main/tap/go/color"
)

func (mc *ManageCli) Do(p *tap.Parser, s []string) error {
	m := mc.C
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
				"[?RD]❌ Task not found[?RT]: [?RD]Invalid task id[?RT]: %v",
				id)
		}
	}

	if len(s) == 1 {
		tasks := theme.Tasks
		if len(tasks) == 0 {
			color.PrintColored("[?RD]❌ Tasks not found")
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
		color.PrintColored("[?RD]❌ Failed to do tasks[?RT]")
		if len(tasks) != 0 {
			color.PrintlnColored("%v", tasks)
		} else {
			print("\n")
		}
	} else {
		color.PrintlnColored("[?GN]✅ Tasks done![?RT]")
	}
	return nil
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
```

## `cmd/manage/main.go`

```go
package main

import (
	"fmt"
	"os"

	"github.com/pt-main/manage/config"
	"github.com/pt-main/manage/manager"
	tap "github.com/pt-main/tap/go"
	"github.com/pt-main/tap/go/utils"
)

var Version = "1.3.0"

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

const docShowTheme = `[?GN]Show[?RT] detailed information about a single theme.

[?BBE]Usage:[?RT]
  [?YW]manage show theme[?RT] [?BGN]<name>[?RT]

[?BBE]Description:[?RT]
  Prints priority, tags, task counters, last assigned ID,
  tables that contain the theme, and the full task list
  with status, ID, date and description.

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

	p.AddCommand("list", mc.List, docList, nil, nil, false)

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
	err := config.Load(m)
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
```

## `cmd/manage/remove.go`

```go
package main

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/pt-main/manage/manager"
	tap "github.com/pt-main/tap/go"
	"github.com/pt-main/tap/go/color"
)

// RemoveTheme: manage remove theme <name>
func (mc *ManageCli) RemoveTheme(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) != 1 {
		return fmt.Errorf("❌ 'remove theme' needs exactly 1 argument (name), got %v", len(s))
	}
	name := s[0]

	if _, ok := m.Tasks[name]; !ok {
		return fmt.Errorf("❌ Invalid theme name: %v", name)
	}
	delete(m.Tasks, name)

	for tableName, themes := range m.Tables {
		newThemes := []string{}
		removed := false
		for _, t := range themes {
			if t == name {
				removed = true
				continue
			}
			newThemes = append(newThemes, t)
		}
		if removed {
			m.Tables[tableName] = newThemes
		}
	}

	return nil
}

// RemoveTask: manage remove task <theme> <id...>
func (mc *ManageCli) RemoveTask(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) < 2 {
		return fmt.Errorf("❌ 'remove task' needs at least 2 arguments (theme + ids), got %v", len(s))
	}
	name := s[0]

	theme, ok := m.Tasks[name]
	if !ok {
		return fmt.Errorf("❌ Invalid theme name: %v", name)
	}

	ids := make(map[int]bool)
	for _, idStr := range s[1:] {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("❌ Invalid task id '%v': %w", idStr, err)
		}
		ids[id] = true
	}

	newTasks := []manager.TaskData{}
	removed := 0
	for _, task := range theme.Tasks {
		if ids[task.GetId()] {
			removed++
			continue
		}
		newTasks = append(newTasks, task)
	}
	theme.Tasks = newTasks

	if removed == 0 {
		return fmt.Errorf("❌ No tasks removed (none of the given ids exist in theme '%v')", name)
	}
	if removed < len(ids) {
		color.PrintlnColored("[?YW]❗️ Some ids were not found in theme: %v", name)
	}

	return nil
}

// RemoveTag: manage remove tag <theme> <tags...>
func (mc *ManageCli) RemoveTag(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) < 2 {
		return fmt.Errorf("❌ 'remove tag' needs at least 2 arguments (theme + tags), got %v", len(s))
	}
	name := s[0]
	toRemove := s[1:]

	theme, ok := m.Tasks[name]
	if !ok {
		return fmt.Errorf("❌ Invalid theme name: %v", name)
	}

	newTags := []string{}
	removedAny := false
	for _, t := range theme.Tags {
		if slices.Contains(toRemove, t) {
			removedAny = true
			continue
		}
		newTags = append(newTags, t)
	}
	if !removedAny {
		return fmt.Errorf("❌ No tags removed (none of the given tags exist on theme '%v')", name)
	}
	theme.Tags = newTags

	return nil
}

// RemoveTable: manage remove table <name>
func (mc *ManageCli) RemoveTable(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) != 1 {
		return fmt.Errorf("❌ 'remove table' needs exactly 1 argument (name), got %v", len(s))
	}
	name := s[0]

	if _, ok := m.Tables[name]; !ok {
		return fmt.Errorf("❌ Invalid table name: %v", name)
	}
	delete(m.Tables, name)

	return nil
}

// RemoveFromTable: manage remove from-table <name> <themes...>
func (mc *ManageCli) RemoveFromTable(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) < 2 {
		return fmt.Errorf("❌ 'remove from-table' needs at least 2 arguments (name + themes), got %v", len(s))
	}
	name := s[0]
	toRemove := s[1:]

	themes, ok := m.Tables[name]
	if !ok {
		return fmt.Errorf("❌ Invalid table name: %v", name)
	}

	newThemes := []string{}
	removedAny := false
	for _, t := range themes {
		if slices.Contains(toRemove, t) {
			removedAny = true
			continue
		}
		newThemes = append(newThemes, t)
	}
	if !removedAny {
		return fmt.Errorf("❌ No themes removed (none of the given names are in table '%v')", name)
	}
	m.Tables[name] = newThemes

	return nil
}

// RemovePriority: manage remove priority <p1> [p2...]
func (mc *ManageCli) RemovePriority(p *tap.Parser, s []string) error {
	m := mc.C
	if len(s) < 1 {
		return fmt.Errorf("❌ 'remove priority' needs at least 1 argument, got %v", len(s))
	}

	newOrder := []string{}
	for _, pr := range m.Settings.PriorityOrder {
		if !slices.Contains(s, pr) {
			newOrder = append(newOrder, pr)
		}
	}
	if len(newOrder) == 0 {
		return fmt.Errorf("❌ Can't remove all priorities (at least one must remain)")
	}
	if len(newOrder) == len(m.Settings.PriorityOrder) {
		return fmt.Errorf("❌ No priorities removed (none of the given values exist)")
	}
	m.Settings.PriorityOrder = newOrder

	return nil
}
```

## `cmd/manage/set.go`

```go
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
```

## `cmd/manage/show.go`

```go
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
```

## `cmd/manage/struct.go`

```go
package main

import "github.com/pt-main/manage/manager"

type ManageCli struct {
	C      *manager.Manager
	NoSave bool
}
```

## `manager/code.go`

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
		"[?BBK]CODE:[?BYW]TODO [?CN]f:[[?GN]%v:%v[?CN]][?RT]:",
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
		if result[i].CodeLine != result[j].CodeLine {
			return result[i].CodeLine < result[j].CodeLine
		}
		return result[i].About < result[j].About
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
	status := "TODO"
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
			name)
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
	for name, _task := range m.Tasks {
		task := *_task
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
			filtered = append(filtered, &task)
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

func LoadC(cfg *tessera.Config, m *manager.Manager) {
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

func SaveC(m *manager.Manager) *tessera.Config {
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
	return tessera.TryToTessera(SaveC(m))
}

func GenManager(code string, m *manager.Manager) error {
	cfg, errs := tessera.Parse(code)
	if len(errs) > 0 {
		return fmt.Errorf("tessera parse errors: %v", errs)
	}
	LoadC(cfg, m)
	return nil
}
```

## `config/utils.go`

```go
package config

import (
	"fmt"
	"os"

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
	file, err := GenCfg(m)
	if err != nil {
		return fmt.Errorf("❗️ Save data: %v", err)
	}
	ToHd()
	defer BackFromHd()

	tmp := ConfigFile + ".tmp"
	defer os.Remove(tmp)
	if err := utils.WriteF(tmp, file); err != nil {
		return fmt.Errorf("❗️ Save data (1): %v", err)
	}
	if err := os.Rename(tmp, ConfigFile); err != nil {
		return fmt.Errorf("❗️ Save data (2): %v", err)
	}
	return nil
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
	err = GenManager(file, m)
	if err != nil {
		return err
	}
	return nil
}
```

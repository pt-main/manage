package manager

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"
)

type OneTask struct {
	Id        int
	About     string
	Done      bool
	CreatedAt string
}

type Task struct {
	Tags     []string
	Priority string
	Tasks    []*OneTask
	Name     string
	Last     int
}

func (task *Task) GetTaskById(id int) *OneTask {
	for _, task := range task.Tasks {
		if task.Id == id {
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
	done := false
	for _, task := range task.Tasks {
		if task.Id == n {
			done = true
			task.Done = true
		}
	}
	if !done {
		return fmt.Errorf("Can't find task n %v", n)
	}
	return nil
}

func (task *Task) AppendToTask(n int, about string) error {
	appended := false
	for _, task := range task.Tasks {
		if task.Id == n {
			appended = true
			task.About += about
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
		Tasks:    make([]*OneTask, 0),
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

func (m *Manager) addTask(theme string, task *OneTask) error {
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
			tasks := []*OneTask{}
			for _, oneT := range task.Tasks {
				switch state {
				case "todo":
					if !oneT.Done {
						tasks = append(tasks, oneT)
					}
				case "done":
					if oneT.Done {
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

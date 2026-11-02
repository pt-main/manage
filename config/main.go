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

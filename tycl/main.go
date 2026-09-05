package loadtycl

import (
	"github.com/pt-main/manage/manager"
	"github.com/pt-main/tycl"
	"github.com/pt-main/tycl/generation"
	"github.com/pt-main/tycl/shared"
)

var contract = `strict {
    tasks: objects = strict {
        name: string,
        tags: strings,
        priority: string,
        last: int,
        tasks: objects = strict {
            id: int,
            done: bool,
            about: string,
            created_at: string,
        },
    },
    tables: objects = strict {
        name: string,
        tasks: strings,
    },
    settings: object = strict {
        priority_order: strings,
        default_priority: string,
    },
}`

func Load(cfg *shared.Config, m *manager.Manager) {
	if m.Settings == nil {
		m.Settings = &manager.Settings{}
	}
	if m.Tables == nil {
		m.Tables = make(map[string][]string)
	}
	if m.Tasks == nil {
		m.Tasks = make(map[string]*manager.Task)
	}

	for _, task := range cfg.InnerArrV["tasks"] {
		var tasks []*manager.OneTask = []*manager.OneTask{}
		for _, innertask := range task.InnerArrV["tasks"] {
			tasks = append(tasks, &manager.OneTask{
				Id:        innertask.IntV["id"],
				About:     innertask.StringV["about"],
				Done:      innertask.BoolV["done"],
				CreatedAt: innertask.StringV["created_at"],
			})
		}
		m.Tasks[task.StringV["name"]] = &manager.Task{
			Name:     task.StringV["name"],
			Tags:     task.StringArrV["tags"],
			Priority: task.StringV["priority"],
			Tasks:    tasks,
			Last:     task.IntV["last"],
		}
	}

	for _, table := range cfg.MainConf.InnerArrV["tables"] {
		m.Tables[table.StringV["name"]] = table.StringArrV["tasks"]
	}

	sett := cfg.InnerV["settings"]
	m.Settings = &manager.Settings{
		PriorityOrder:   sett.StringArrV["priority_order"],
		DefaultPriority: sett.StringV["default_priority"],
	}
}

func Save(m *manager.Manager) *shared.Config {
	cfg := shared.NewNilConfig()

	tasksSlice := []*shared.Config{}
	for name, task := range m.Tasks {
		taskCfg := shared.NewNilConfig()
		taskCfg.StringV["name"] = name
		taskCfg.StringArrV["tags"] = task.Tags
		taskCfg.StringV["priority"] = task.Priority
		taskCfg.IntV["last"] = task.Last
		innerSlice := []*shared.Config{}
		for _, one := range task.Tasks {
			innerCfg := shared.NewNilConfig()
			innerCfg.IntV["id"] = one.Id
			innerCfg.StringV["about"] = one.About
			innerCfg.BoolV["done"] = one.Done
			innerCfg.StringV["created_at"] = one.CreatedAt
			innerSlice = append(innerSlice, innerCfg)
		}
		taskCfg.InnerArrV["tasks"] = innerSlice
		tasksSlice = append(tasksSlice, taskCfg)
	}
	cfg.InnerArrV["tasks"] = tasksSlice

	tablesSlice := []*shared.Config{}
	for name, subjects := range m.Tables {
		tableCfg := shared.NewNilConfig()
		tableCfg.StringV["name"] = name
		tableCfg.StringArrV["tasks"] = subjects
		tablesSlice = append(tablesSlice, tableCfg)
	}
	cfg.InnerArrV["tables"] = tablesSlice

	if m.Settings != nil {
		settingsCfg := shared.NewNilConfig()
		settingsCfg.StringArrV["priority_order"] = m.Settings.PriorityOrder
		settingsCfg.StringV["default_priority"] = m.Settings.DefaultPriority
		cfg.InnerV["settings"] = settingsCfg
	}

	cfg.MainConf = cfg

	return cfg
}

func GenCfg(m *manager.Manager) (string, error) {
	return generation.Tycl(Save(m))
}

func GenManager(code string, m *manager.Manager) error {
	cfg, err := tycl.Process(code, contract, true)
	if err != nil {
		return err
	}
	Load(cfg, m)
	return nil
}

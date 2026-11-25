package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pt-main/manage/manager"
	"github.com/pt-main/tycl/utils"
)

const CODE_THEME = ".code"
const CODE_PRIORITY = "A"

var ConfigFile = "manager_config.tess"

func ConfigPath() string {
	if envPath := os.Getenv("MANAGE_CONFIG"); envPath != "" {
		return envPath
	}
	hd, err := os.UserHomeDir()
	if err != nil {
		panic("Can't find user home dir: " + err.Error())
	}
	return filepath.Join(hd, ConfigFile)
}

func Save(m *manager.Manager) error {
	file, err := GenCfg(m)
	if err != nil {
		return fmt.Errorf("❗️ Save data: %v", err)
	}

	cfgPath := ConfigPath()
	tmp := cfgPath + ".tmp"
	defer os.Remove(tmp)

	if err := utils.WriteF(tmp, file); err != nil {
		return fmt.Errorf("❗️ Save data (1): %v", err)
	}
	if err := os.Rename(tmp, cfgPath); err != nil {
		return fmt.Errorf("❗️ Save data (2): %v", err)
	}
	return nil
}

func GetConf() (string, error) {
	file, err := utils.OpenF(ConfigPath())
	if err != nil {
		return "", fmt.Errorf("❗️ Load data: %v. Do you forget to use init?", err)
	}
	return file, nil
}

func Load(m *manager.Manager, add_code bool) error {
	file, err := GetConf()
	if err != nil {
		return err
	}
	err = GenManager(file, m)
	if err != nil {
		return err
	}

	if add_code {
		m.AddTheme(CODE_THEME, manager.NewTask([]string{}, CODE_PRIORITY, CODE_THEME), true)
		for _, task := range manager.FindTodosInDir(".") {
			theme := m.Tasks[CODE_THEME]
			theme.Tasks = append(theme.Tasks, task)
			theme.Last = max(task.Id, theme.Last)
		}
	}

	return nil
}

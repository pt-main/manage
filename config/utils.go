package config

import (
	"fmt"
	"os"

	"github.com/pt-main/manage/manager"
	"github.com/pt-main/tycl/utils"
)

const CODE_THEME = ".code"
const CODE_PRIORITY = "A"

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

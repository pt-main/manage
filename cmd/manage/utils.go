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

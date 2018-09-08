package util

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

// GenericFileAssignment will take a path and parse check it for specific
// matches for file paths, and then find the given paths.
// - Absolute paths are inputs starting with / and remain unmodified.
// - Relative paths are inputs starting with ./ and append the source or pwd.
// - Empty paths are unmodified.
// Results are checked against os.Stat and the results are returned with an error or nil.
func GenericFileAssignment(input, path string, check bool) (string, error) {

	// The playbook will be located on the host if the remote flag is enabled.
	if input != "" {
		if strings.HasPrefix(input, "/") {
			input = fmt.Sprintf("%v", input)
		}
		if strings.HasPrefix(input, "./") {
			input = strings.Replace(input, "./", "", 1)
		}
		pwd, _ := os.Getwd()
		if path == "" {
			input = fmt.Sprintf("%v/%v", pwd, input)
		} else {
			input = fmt.Sprintf("%v", input)
		}

		if check {
			_, err := os.Stat(input);
			os.IsNotExist(err)
			if err != nil {
				return input, errors.New("Specified file " + input + " does not exist")
			}
		}
	}

	return input, nil
}

// GenericPlaybookAssignment will return the path of the only available playbook file under
// the roles tests directory, to be used when the default playbook path is invalid.
func GenericPlaybookAssignment(input, path string) (string, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return input, errors.New("no playbook.yml file was found, please provide path to playbook.")
	}

	for _, file := range files {
		if strings.Contains(file, ".yml") {
			if strings.Contains(file, "playbook.yml") {
				pwd, _ := os.Getwd()
				return strings.Replace(file, pwd + string(os.PathSeparator), "", -1), nil
			}
		}
	}

	return input, nil
}

// MapPlaybook will adjust the playbook path for the appropriate
// path based on the configuration. ie remote or not, and
// guesswork based upon input. For example, paths starting with
// /, ./ or otherwise.
func MapPlaybook(config *AnsibleConfig) {

	playbook, err := GenericFileAssignment(config.PlaybookFile, config.HostPath, true)
	playbook, err = GenericPlaybookAssignment(config.PlaybookFile, config.HostPath)

	config.PlaybookFile = playbook

	if err != nil {
		log.Fatalf("Specified playbook file %v does not exist.", config.PlaybookFile)
	}

	if !config.Remote && config.PlaybookFile != "" {
		pwd, _ := os.Getwd()
		config.PlaybookFile = strings.Replace(config.PlaybookFile, pwd, config.RemotePath, -1)
		config.PlaybookFile = fmt.Sprintf("%v/%v", config.RemotePath, config.PlaybookFile)
	}

	if err == nil {
		if config.Remote && config.RemotePath == "" {
			pwd, _ := os.Getwd()
			config.RemotePath = pwd
		} else if !config.Remote && config.RemotePath == "" {
			config.RemotePath = "/etc/ansible/roles/role_under_test"
		}
	}

}

// MapInventory will adjust the inventory path for the appropriate
// path based on the configuration. ie remote or not, and
// guesswork based upon input. For example, paths starting with
// /, ./ or otherwise.
func MapInventory(CID string, config *AnsibleConfig) {

	inventory, err := GenericFileAssignment(config.Inventory, config.HostPath, false)
	config.Inventory = inventory

	if err != nil {
		log.Fatalf("Specified inventory file %v does not exist.", config.Inventory)
	}

	if !config.Remote && config.Inventory != "" {
		pwd, _ := os.Getwd()
		config.Inventory = strings.Replace(config.Inventory, pwd, config.RemotePath, -1)
		config.Inventory = fmt.Sprintf("%v/%v", config.RemotePath, config.Inventory)
	}

}

// MapRequirements will adjust the requirements path for the appropriate
// path based on the configuration. ie remote or not, and
// guesswork based upon input. For example, paths starting with
// /, ./ or otherwise.
func MapRequirements(config *AnsibleConfig) {

	requirements, err := GenericFileAssignment(config.RequirementsFile, config.HostPath, true)
	config.RequirementsFile = requirements

	if err != nil {
		log.Fatalf("Specified requirements file %v does not exist.", config.RequirementsFile)
	}

	if !config.Remote && config.RequirementsFile != "" {
		pwd, _ := os.Getwd()
		config.RequirementsFile = strings.Replace(config.RequirementsFile, pwd, config.RemotePath, -1)
		config.RequirementsFile = fmt.Sprintf("%v/%v", config.RemotePath, config.RequirementsFile)
	}

}

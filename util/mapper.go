package util

import (
	"strings"
	"os"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

// MapPlaybook will adjust the playbook path for the appropriate
// path based on the configuration. ie remote or not, and
// guesswork based upon input. For example, paths starting with
// /, ./ or otherwise.
func MapPlaybook(config *AnsibleConfig) {
	Source := config.RemotePath
	Playbook := config.PlaybookFile
	Quiet := config.Quiet
	// Adjust playbook path
	if config.Remote {
		// The playbook will be located on the host if the remote flag is enabled.
		if strings.HasPrefix(config.PlaybookFile, "./") {
			config.PlaybookFile = fmt.Sprintf("%v/%v", Source, Playbook)
		} else if strings.HasPrefix(config.PlaybookFile, "/") {
			config.PlaybookFile = fmt.Sprintf("%v", config.PlaybookFile)
		} else if !config.Remote {
			config.PlaybookFile = fmt.Sprintf("%v/tests/%v", Source, config.PlaybookFile)
		} else if config.Remote {
			config.PlaybookFile = fmt.Sprintf("%v/tests/%v", Source, config.PlaybookFile)
		}
		fp := fmt.Sprintf(config.PlaybookFile)
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			if !Quiet {
				log.Fatalf("Specified playbook file %v does not exist.", fp)
			}
		}
	} else {
		// The playbook will be located on the container (via mount) if the remote flag is not enabled.
		config.PlaybookFile = fmt.Sprintf("/etc/ansible/roles/role_under_test/%v", Playbook)
		file := fmt.Sprintf("%v/%v", Source, Playbook)
		fp := fmt.Sprintf(file)
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			if !Quiet {
				log.Fatalf("Specified playbook file %v does not exist.", fp)
			}
		}
	}
}

// MapInventory will adjust the inventory path for the appropriate
// path based on the configuration. ie remote or not, and
// guesswork based upon input. For example, paths starting with
// /, ./ or otherwise.
func MapInventory(CID string, config *AnsibleConfig) {
	Source := config.RemotePath
	Inventory := config.Inventory
	// Adjust playbook path
	if config.Remote {
		// The inventory will be overriden if the remote flag is enabled.
		config.Inventory = fmt.Sprintf("%v,", CID)
	} else {
		// The inventory will be located on the container (via mount) if the remote flag is not enabled.
		config.Inventory = fmt.Sprintf("/etc/ansible/roles/role_under_test/%v", Inventory)
		file := fmt.Sprintf("%v/%v", Source, Inventory)
		fp := fmt.Sprintf(file)
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			// The file does not exist, we should assume the user inputted a value, not a path.
			config.Inventory = Inventory
		}
	}
}

// MapRequirements will adjust the requirements path for the appropriate
// path based on the configuration. ie remote or not, and
// guesswork based upon input. For example, paths starting with
// /, ./ or otherwise.
func MapRequirements(config *AnsibleConfig) {
	Source := config.RemotePath
	Requirements := config.RequirementsFile
	Quiet := config.Quiet
	// Adjust requirements path
	if config.Remote {
		// The requirements file will be located on the host if the remote flag is enabled.
		if strings.HasPrefix(config.RequirementsFile, "./") {
			config.RequirementsFile = fmt.Sprintf("%v/%v", Source, Requirements)
		} else if strings.HasPrefix(config.RequirementsFile, "/") {
			config.RequirementsFile = fmt.Sprintf("%v", Requirements)
		} else if !config.Remote {
			config.RequirementsFile = fmt.Sprintf("%v/tests/%v", Source, Requirements)
		} else if config.Remote {
			config.RequirementsFile = fmt.Sprintf("%v/tests/%v", Source, Requirements)
		}
		fp := fmt.Sprintf(config.RequirementsFile)
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			if !Quiet {
				log.Warnf("Specified requirements file %v does not exist.", fp)
			}
		}
	} else {
		// The requirements will be located on the container (via mount) if the remote flag is not enabled.
		config.RequirementsFile = fmt.Sprintf("/etc/ansible/roles/role_under_test/%v", Requirements)
		file := fmt.Sprintf("%v/%v", Source, Requirements)
		fp := fmt.Sprintf(file)
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			if !Quiet {
				log.Warnf("Specified requirements file %v does not exist.", fp)
			}
		}
	}
}

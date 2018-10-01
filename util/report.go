package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// AnsibleReport will contain metadata about the run which will be, is and has executed.
type AnsibleReport struct {
	Meta struct {
		Timestamp    time.Time
		Repository   string
		CommitHash   string
		LocalChanges bool
		ReportFile   string
	}
	Ansible struct {
		Config       AnsibleConfig
		Distribution Distribution
		Hosts        []string
		Syntax       bool
		Requirements bool
		Run          struct {
			Result bool
			Time   time.Duration
		}
		Idempotence struct {
			Result bool
			Time   time.Duration
		}
	}
	Docker struct {
		Run     bool
		Kill    bool
		Volumes []string
	}
}

// gitCmd will run git commands in the specified directory.
func gitCmd(path string, args []string) (string, error) {
	// Find git.
	git, err := exec.LookPath("git")
	if err != nil {
		log.Errorln(err)
	}

	// Generate the command, based on input.
	cmd := exec.Cmd{}
	cmd.Path = git
	cmd.Dir = path
	cmd.Args = args

	// Create a buffer for the output.
	var out bytes.Buffer
	multi := io.MultiWriter(&out)

	// Assign the output to the writer.
	cmd.Stdout = multi

	// Check the errors, return as needed.
	var wg sync.WaitGroup
	wg.Add(1)
	if err := cmd.Run(); err != nil {
		log.Errorln(err)
		return out.String(), err
	}
	wg.Done()

	// Return out output as a string.
	return out.String(), nil
}

// isGit will identify if the path is a git repository.
func isGit(path string) (bool, error) {
	git, err := gitCmd(path, []string{"git", "status"})
	if strings.Contains(git, "Not a git repository") {
		return false, err
	}
	return true, err
}

// IsGit will abstract the private func from the report.
func (report *AnsibleReport) IsGit() bool {
	isGitRepo, _ := isGit(report.Ansible.Config.HostPath)
	return isGitRepo
}

// getGitChanges will identify if there are local changes in git.
func getGitChanges(path string) (bool, error) {
	// Pass arguments to common git command runner; return results.
	git, err := gitCmd(path, []string{"git", "status", "--porcelain"})
	if git != "" {
		return true, err
	}
	return false, err
}

// getGitRepo will return the URL path for the repository.
func getGitRepo(path string) (string, error) {

	// Pass arguments to common git command runner; return results.
	git, err := gitCmd(path, []string{"git", "ls-remote", "--get-url", "origin"})
	return git, err

}

// getGitRepo will return the commit hash for the repository.
func getGitCommit(path string) (string, error) {

	// Pass arguments to common git command runner; return results.
	git, err := gitCmd(path, []string{"git", "rev-parse", "HEAD"})
	return git, err

}

// NewReport will generate a new Report variable from the input configuration.
func NewReport(config *AnsibleConfig) AnsibleReport {

	// Create the variable.
	report := new(AnsibleReport)

	// Set appropriate values for git integration.
	if isGitRepo, _ := isGit(config.HostPath); isGitRepo == true {
		gitRepository, _ := getGitRepo(config.HostPath)
		gitRepository = strings.Replace(gitRepository, "\n", "", -1)
		gitCommit, _ := getGitCommit(config.HostPath)
		gitCommit = strings.Replace(gitCommit, "\n", "", -1)
		gitChanges, _ := getGitChanges(config.HostPath)

		report.Meta.CommitHash = gitCommit
		report.Meta.Repository = gitRepository
		report.Meta.LocalChanges = gitChanges
	}

	// Set appropriate defaults as needed.
	report.Meta.Timestamp = time.Now()
	report.Ansible.Config = *config
	report.Ansible.Syntax = false
	report.Ansible.Requirements = false
	report.Ansible.Run.Result = false
	report.Ansible.Run.Time = 0
	report.Ansible.Idempotence.Result = false
	report.Ansible.Idempotence.Time = 0
	report.Docker.Run = false
	report.Docker.Kill = false

	// Return the report
	return *report

}

// GetJSON will return an unmarhaled object as JSON.
func (report *AnsibleReport) GetJSON(data interface{}) ([]byte, error) {

	// Marshal as JSON
	result, err := json.Marshal(data)
	if err != nil {
		log.Errorln(err)
		return []byte{}, err
	}

	// Return the []bytes as a string
	return result, nil

}

// GetYAML will return an unmarhaled object as YAML.
func (report *AnsibleReport) GetYAML(data interface{}) ([]byte, error) {

	// Marshal as JSON
	result, err := yaml.Marshal(data)
	if err != nil {
		log.Errorln(err)
		return []byte{}, err
	}

	// Return the []bytes as a string
	return result, nil

}

// printFile will output the input data to the given filename.
// Intended for exclusive use by GetJSON and GetYAML.
func (report *AnsibleReport) printFile(data []byte) (err error) {

	filename := report.Meta.ReportFile

	// If the file already exists, we should delete it.
	if _, err := os.Stat(filename); err == nil {
		// The file exists.
		if err = os.Remove(filename); err != nil {
			// The file could not be deleted
			log.Errorf("failed to delete %v\n", filename)
			return err
		}
	}

	// Assume no file is available.
	if _, err := os.Stat(filename); err != nil {
		if file, err := os.Create(filename); err != nil {
			// File could not be created.
			log.Errorf("could not create file %v\n", filename)
			return err
		} else {
			// File was created, attempt to write to it
			if err = ioutil.WriteFile(filename, data, 0644); err != nil {
				// Could not write to file.
				log.Errorf("could not write data to %v\n", filename)
				return err
			} else {
				// Wrote to file successfully.
				log.Infof("Report data has been written to %v\n", filename)
			}
			// Close the file.
			file.Sync()
			defer file.Close()
		}
	}
	return

}

// Printf will print the report in a formatted way.
func (report *AnsibleReport) Printf() {

	fmt.Println()
	fmt.Println("----------------------------------------------------------")
	fmt.Println("Ansible Role Tester Report")
	fmt.Println("----------------------------------------------------------")
	fmt.Printf("Timestamp: \t\t\t%v\n", report.Meta.Timestamp)
	if report.IsGit() {
		fmt.Printf("Repository URL: \t\t%v\n", report.Meta.Repository)
		fmt.Printf("Repository commit: \t\t%v\n", report.Meta.CommitHash)
		fmt.Printf("Local changes: \t\t\t%v\n", report.Meta.LocalChanges)
	}
	fmt.Println("----------------------------------------------------------")
	fmt.Printf("Syntax check: \t\t\t%v\n", report.Ansible.Syntax)
	fmt.Printf("Requirements installed: \t%v\n", report.Ansible.Requirements)
	fmt.Printf("Run result: \t\t\t%v\n", report.Ansible.Run.Result)
	fmt.Printf("Run time: \t\t\t%v\n", report.Ansible.Run.Time)
	fmt.Printf("Idempotence result: \t\t%v\n", report.Ansible.Idempotence.Result)
	fmt.Printf("Idempotence time: \t\t%v\n", report.Ansible.Idempotence.Time)
	fmt.Println("----------------------------------------------------------")
	fmt.Printf("Docker run: \t\t\t%v\n", report.Docker.Run)
	fmt.Printf("Docker kill: \t\t\t%v\n", report.Docker.Kill)
	fmt.Println("----------------------------------------------------------")
	fmt.Println()

	if strings.HasSuffix(report.Meta.ReportFile, ".yaml") {
		yamlReport, _ := report.GetYAML(report)
		report.printFile(yamlReport)
	}
	if strings.HasSuffix(report.Meta.ReportFile, ".yml") {
		yamlReport, _ := report.GetYAML(report)
		report.printFile(yamlReport)
	}
	if strings.HasSuffix(report.Meta.ReportFile, ".json") {
		jsonReport, _ := report.GetJSON(report)
		report.printFile(jsonReport)
	}

}

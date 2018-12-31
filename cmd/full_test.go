package cmd

import (
	"os"
	"testing"

	"io/ioutil"

	"bytes"

	"fmt"

	log "github.com/sirupsen/logrus"

	"bou.ke/monkey"
	"github.com/fubarhouse/ansible-role-tester/util"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/afero"
)

func TestFullCmd(t *testing.T) {

	Convey("Setup", t, func() {
		// Send testing output to /dev/null
		log.SetOutput(ioutil.Discard)

		// Here we use monkey to patch the os.Exit call and run asserts on it.
		// NOTE: If there are multiple calls to os.Exit in a single test case,
		// we will only see the last one here. This is a case where real world
		// behavior and test behavior differs.
		var exitCode int // flag for exit code.
		fakeExit := func(i int) {
			exitCode = i
		}
		patch := monkey.Patch(os.Exit, fakeExit)
		defer patch.Unpatch()

		// Clone testing repo.
		fs := afero.NewOsFs()
		baseDir, err := afero.TempDir(fs, "", "ansible-role-tester")
		if err != nil {
			t.Fatal("error creating tempDir", err)
		}
		defer fs.RemoveAll(baseDir)
		_, err = util.GitCmd(baseDir, []string{"git", "clone", "https://github.com/issmirnov/ansible-role-art-tester.git"})
		So(err, ShouldBeNil)
		artRepo := fmt.Sprintf("%s/%s", baseDir, "ansible-role-art-tester")

		// fmt.Printf("Setup complete, initiating testing.")

		Convey("Verify all exit code scenarios", func() {
			Convey("Correct runs return 0", func() {
				fullCmd := InitFullCmdForTest(baseDir)
				buf := new(bytes.Buffer)
				fullCmd.SetOutput(buf)
				fullCmd.SetArgs([]string{
					fmt.Sprintf("--playbook=%s", "tests/playbook-simple.yml"),
					fmt.Sprintf("--source=%s", artRepo),
					fmt.Sprintf("--destination=%s", "/etc/ansible/roles/ansible-role-art-tester"),
					"--quiet",
				})
				err := fullCmd.Execute()
				So(err, ShouldBeNil)
				So(exitCode, ShouldEqual, util.OKCode)
			})

			Convey(fmt.Sprintf("Syntax errors in playbook return %d", util.AnsibleSyntaxCode), func() {
				fullCmd := InitFullCmdForTest(baseDir)
				buf := new(bytes.Buffer)
				fullCmd.SetOutput(buf)
				fullCmd.SetArgs([]string{
					fmt.Sprintf("--playbook=%s", "tests/playbook-syntax-fail.yml"),
					fmt.Sprintf("--source=%s", artRepo),
					fmt.Sprintf("--destination=%s", "/etc/ansible/roles/ansible-role-art-tester"),
					"--quiet",
				})
				err := fullCmd.Execute()
				So(err, ShouldBeNil)
				So(exitCode, ShouldEqual, util.AnsibleSyntaxCode)
			})

			Convey(fmt.Sprintf("Idempotency failures return %d", util.AnsibleIdempotenceCode), func() {
				fullCmd := InitFullCmdForTest(baseDir)
				buf := new(bytes.Buffer)
				fullCmd.SetOutput(buf)
				fullCmd.SetArgs([]string{
					fmt.Sprintf("--playbook=%s", "tests/playbook-idempotency-fail.yml"),
					fmt.Sprintf("--source=%s", artRepo),
					fmt.Sprintf("--destination=%s", "/etc/ansible/roles/ansible-role-art-tester"),
					"--quiet",
				})
				err := fullCmd.Execute()
				So(err, ShouldBeNil)
				So(exitCode, ShouldEqual, util.AnsibleIdempotenceCode)
			})

			// Currently disabled. The PostRun hooks ends up overwriting our exit code.
			// Convey(fmt.Sprintf("Running outside of a role directory returns %d", util.NotARoleCode), func() {
			// 	fullCmd := InitFullCmdForTest("/tmp")
			// 	fullCmd.SetArgs([]string{})
			// 	err := fullCmd.Execute()
			// 	So(err, ShouldBeNil)
			// 	So(exitCode, ShouldEqual, util.NotARoleCode)
			// })
		})
	})
}

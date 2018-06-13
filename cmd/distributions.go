package cmd

import (
	"errors"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// A Distribution declares the options to
// pass to Docker to dockerRun and test the container.
type Distribution struct {

	// Initialise is the initialisation command
	Initialise string

	// Name is the identifying name of the distribution
	Name string

	// Privileged is a boolean to indicate to use privileged
	Privileged bool

	// The volume parameter when running a container.
	Volume string

	// The fully qualified container name in the format:
	// name/image:version - ie fubarhouse/docker-ansible:bionic
	Container string

	User   string
	Distro string
}

// CentOS 6
var CentOS6 = Distribution{
	"/sbin/init",
	"centos6",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:centos-6",
	"fubarhouse",
	"centos6",
}

// CentOS 7
var CentOS7 = Distribution{
	"/usr/lib/systemd/systemd",
	"centos7",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:centos-7",
	"fubarhouse",
	"centos7",
}

// Wheezy
var DebianWheezy = Distribution{
	"/bin/systemd",
	"wheezy",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:wheezy",
	"fubarhouse",
	"debian7",
}

// Jessie
var DebianJessie = Distribution{
	"/lib/systemd",
	"jessie",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:jessie",
	"fubarhouse",
	"debian8",
}

// Stretch
var DebianStretch = Distribution{
	"/bin/systemd",
	"stretch",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:stretch",
	"fubarhouse",
	"debian9",
}

// Buster
var DebianBuster = Distribution{
	"/bin/systemd",
	"buster",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:buster",
	"fubarhouse",
	"debian10",
}

// Fedora 24
var Fedora24 = Distribution{
	"/usr/lib/systemd/systemd",
	"fedora24",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:fedora-24",
	"fubarhouse",
	"fedora24",
}

// Fedora 25
var Fedora25 = Distribution{
	"/usr/lib/systemd/systemd",
	"fedora25",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:fedora-25",
	"fubarhouse",
	"fedora25",
}

// Fedora 26
var Fedora26 = Distribution{
	"/usr/lib/systemd/systemd",
	"fedora26",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:fedora-26",
	"fubarhouse",
	"fedora26",
}

// Fedora 27
var Fedora27 = Distribution{
	"/usr/lib/systemd/systemd",
	"fedora27",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:fedora-27",
	"fubarhouse",
	"fedora27",
}

// Fedora 28
var Fedora28 = Distribution{
	"/usr/lib/systemd/systemd",
	"fedora28",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:fedora-28",
	"fubarhouse",
	"fedora28",
}

// Ubuntu 12.04
var Ubuntu1204 = Distribution{
	"/sbin/init",
	"ubuntu1204",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:precise",
	"fubarhouse",
	"ubuntu1204",
}

// Ubuntu 12.10
var Ubuntu1210 = Distribution{
	"/sbin/init",
	"ubuntu1210",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:quantal",
	"fubarhouse",
	"ubuntu1210",
}

// Ubuntu 13.04
var Ubuntu1304 = Distribution{
	"/sbin/init",
	"ubuntu1304",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:raring",
	"fubarhouse",
	"ubuntu1304",
}

// Ubuntu 13.10
var Ubuntu1310 = Distribution{
	"/sbin/init",
	"ubuntu1310",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:saucy",
	"fubarhouse",
	"ubuntu1310",
}

// Ubuntu 14.04
var Ubuntu1404 = Distribution{
	"/sbin/init",
	"ubuntu1404",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:trusty",
	"fubarhouse",
	"ubuntu1404",
}

// Ubuntu 14.10
var Ubuntu1410 = Distribution{
	"/sbin/init",
	"ubuntu1410",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:utopic",
	"fubarhouse",
	"ubuntu1410",
}

// Ubuntu 15.04
var Ubuntu1504 = Distribution{
	"/sbin/init",
	"ubuntu1504",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:vivid",
	"fubarhouse",
	"ubuntu1504",
}

// Ubuntu 15.10
var Ubuntu1510 = Distribution{
	"/sbin/init",
	"ubuntu1510",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:wily",
	"fubarhouse",
	"ubuntu1510",
}

// Ubuntu 16.04
var Ubuntu1604 = Distribution{
	"/sbin/init",
	"ubuntu1604",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:xenial",
	"fubarhouse",
	"ubuntu1604",
}

// Ubuntu 16.10
var Ubuntu1610 = Distribution{
	"/sbin/init",
	"ubuntu1610",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:yakkety",
	"fubarhouse",
	"ubuntu1610",
}

// Ubuntu 17.04
var Ubuntu1704 = Distribution{
	"/sbin/init",
	"ubuntu1704",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:zesty",
	"fubarhouse",
	"ubuntu1704",
}

// Ubuntu 17.10
var Ubuntu1710 = Distribution{
	"/sbin/init",
	"ubuntu1710",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:artful",
	"fubarhouse",
	"ubuntu1710",
}

// Ubuntu 18.04
var Ubuntu1804 = Distribution{
	"/sbin/init",
	"ubuntu1804",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:bionic",
	"fubarhouse",
	"ubuntu1804",
}

// CentOS 6
var JeffCentOS6 = Distribution{
	"/sbin/init",
	"centos6",
	true,
	"",
	"geerlingguy/docker-centos6-ansible:latest",
	"geerlingguy",
	"centos6",
}

// CentOS 7
var JeffCentOS7 = Distribution{
	"/usr/lib/systemd/systemd",
	"centos7",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"geerlingguy/docker-centos7-ansible:latest",
	"geerlingguy",
	"centos7",
}

// Ubuntu 14.04
var JeffUbuntu1404 = Distribution{
	"/sbin/init",
	"ubuntu1404",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"geerlingguy/docker-ubuntu1404-ansible:latest",
	"geerlingguy",
	"ubuntu1404",
}

// Ubuntu 16.04
var JeffUbuntu1604 = Distribution{
	"/lib/systemd/systemd",
	"ubuntu1604",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"geerlingguy/docker-ubuntu1604-ansible:latest",
	"geerlingguy",
	"ubuntu1604",
}

// Ubuntu 18.04
var JeffUbuntu1804 = Distribution{
	"/lib/systemd/systemd",
	"ubuntu1804",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"geerlingguy/docker-ubuntu1804-ansible:latest",
	"geerlingguy",
	"ubuntu1804",
}

// A slice of distributions.
var Distributions = []Distribution{
	CentOS6,
	CentOS7,
	DebianWheezy,
	DebianJessie,
	DebianStretch,
	DebianBuster,
	Fedora24,
	Fedora25,
	Fedora26,
	Fedora27,
	Fedora28,
	Ubuntu1204,
	Ubuntu1210,
	Ubuntu1304,
	Ubuntu1310,
	Ubuntu1404,
	Ubuntu1410,
	Ubuntu1504,
	Ubuntu1510,
	Ubuntu1604,
	Ubuntu1610,
	Ubuntu1704,
	Ubuntu1710,
	Ubuntu1804,
	JeffCentOS6,
	JeffCentOS7,
	JeffUbuntu1404,
	JeffUbuntu1604,
	JeffUbuntu1804,
}

// getDistribution will get the distribution object to allow dynamic
// loading of different distributions. A suitable struct will be compiled
// from the inputs and returned with an error if the specified container
// cannot be found.
func getDistribution(container, target, init, volume, user, distro string) (Distribution, error) {

	// We will search for the exact container.
	for _, dist := range Distributions {
		// Check for explicit matches using image.
		if dist.Container == container {
			return dist, nil
		}
		// Check for explicit matches for user and distro.
		if dist.User == user && dist.Distro == distro {
			return dist, nil
		}
	}

	c, _ := dockerExec([]string{
		"images",
		container,
	}, false)

	if !strings.Contains(c, container) {
		log.Errorf("no valid image was found for '%v'\n", container)
		os.Exit(1)
	}

	return Distribution{
		init,
		target,
		true,
		volume,
		container,
		user,
		distro,
	},
	errors.New("could not find matching distribution, returned a compatible data structure")
}

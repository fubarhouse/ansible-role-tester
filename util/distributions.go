package util

import (
	"errors"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// A Distribution declares the options to
// pass to Docker to dockerRun and test the container.
type Distribution struct {

	// CID is the name/id of the container.
	CID string

	// Name is the identifying name of the distribution
	Name string

	// Privileged is a boolean to indicate to use privileged
	Privileged bool

	// The fully qualified container name in the format:
	// name/image:version - ie fubarhouse/docker-ansible:bionic
	Container string

	// User is the user associated to the image file, used
	// when searching for a user from the command line tool.
	User string

	// Distro is the distro associated to the image file, used
	// when searching for a distro from the command line tool.
	Distro string

	// Family associated to this distribution.
	Family Family
}

// Family is a set of characteristics describing a family of linux distributions.
// For example, ubuntu, centos, debian or fedora.
type Family struct {
	Name string
	Initialise string
	Volume string
}

var CentOS = Family{
	"CentOS",
	"/sbin/init",
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
}

var Debian = Family{
	"Debian",
	"/bin/systemd",
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
}
var Fedora = Family{
	"Fedora",
	"/usr/lib/systemd/systemd",
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
}
var Ubuntu = Family{
	"Ubuntu",
	"/sbin/init",
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
}

// CentOS 6
var CentOS6 = Distribution{
	"",
	"centos6",
	true,
	"fubarhouse/docker-ansible:centos-6",
	"fubarhouse",
	"centos6",
	CentOS,
}

// CentOS 7
var CentOS7 = Distribution{
	"",
	"centos7",
	true,
	"fubarhouse/docker-ansible:centos-7",
	"fubarhouse",
	"centos7",
	CentOS,
}

// Wheezy
var DebianWheezy = Distribution{
	"",
	"wheezy",
	true,
	"fubarhouse/docker-ansible:wheezy",
	"fubarhouse",
	"debian7",
	Debian,
}

// Jessie
var DebianJessie = Distribution{
	"",
	"jessie",
	true,
	"fubarhouse/docker-ansible:jessie",
	"fubarhouse",
	"debian8",
	Debian,
}

// Stretch
var DebianStretch = Distribution{
	"",
	"stretch",
	true,
	"fubarhouse/docker-ansible:stretch",
	"fubarhouse",
	"debian9",
	Debian,
}

// Buster
var DebianBuster = Distribution{
	"",
	"buster",
	true,
	"fubarhouse/docker-ansible:buster",
	"fubarhouse",
	"debian10",
	Debian,
}

// Fedora 24
var Fedora24 = Distribution{
	"",
	"fedora24",
	true,
	"fubarhouse/docker-ansible:fedora-24",
	"fubarhouse",
	"fedora24",
	Fedora,
}

// Fedora 25
var Fedora25 = Distribution{
	"",
	"fedora25",
	true,
	"fubarhouse/docker-ansible:fedora-25",
	"fubarhouse",
	"fedora25",
	Fedora,
}

// Fedora 26
var Fedora26 = Distribution{
	"",
	"fedora26",
	true,
	"fubarhouse/docker-ansible:fedora-26",
	"fubarhouse",
	"fedora26",
	Fedora,
}

// Fedora 27
var Fedora27 = Distribution{
	"",
	"fedora27",
	true,
	"fubarhouse/docker-ansible:fedora-27",
	"fubarhouse",
	"fedora27",
	Fedora,
}

// Fedora 28
var Fedora28 = Distribution{
	"",
	"fedora28",
	true,
	"fubarhouse/docker-ansible:fedora-28",
	"fubarhouse",
	"fedora28",
	Fedora,
}

// Ubuntu 12.04
var Ubuntu1204 = Distribution{
	"",
	"ubuntu1204",
	true,
	"fubarhouse/docker-ansible:precise",
	"fubarhouse",
	"ubuntu1204",
	Ubuntu,
}

// Ubuntu 12.10
var Ubuntu1210 = Distribution{
	"",
	"ubuntu1210",
	true,
	"fubarhouse/docker-ansible:quantal",
	"fubarhouse",
	"ubuntu1210",
	Ubuntu,
}

// Ubuntu 13.04
var Ubuntu1304 = Distribution{
	"",
	"ubuntu1304",
	true,
	"fubarhouse/docker-ansible:raring",
	"fubarhouse",
	"ubuntu1304",
	Ubuntu,
}

// Ubuntu 13.10
var Ubuntu1310 = Distribution{
	"",
	"ubuntu1310",
	true,
	"fubarhouse/docker-ansible:saucy",
	"fubarhouse",
	"ubuntu1310",
	Ubuntu,
}

// Ubuntu 14.04
var Ubuntu1404 = Distribution{
	"",
	"ubuntu1404",
	true,
	"fubarhouse/docker-ansible:trusty",
	"fubarhouse",
	"ubuntu1404",
	Ubuntu,
}

// Ubuntu 14.10
var Ubuntu1410 = Distribution{
	"",
	"ubuntu1410",
	true,
	"fubarhouse/docker-ansible:utopic",
	"fubarhouse",
	"ubuntu1410",
	Ubuntu,
}

// Ubuntu 15.04
var Ubuntu1504 = Distribution{
	"",
	"ubuntu1504",
	true,
	"fubarhouse/docker-ansible:vivid",
	"fubarhouse",
	"ubuntu1504",
	Ubuntu,
}

// Ubuntu 15.10
var Ubuntu1510 = Distribution{
	"",
	"ubuntu1510",
	true,
	"fubarhouse/docker-ansible:wily",
	"fubarhouse",
	"ubuntu1510",
	Ubuntu,
}

// Ubuntu 16.04
var Ubuntu1604 = Distribution{
	"",
	"ubuntu1604",
	true,
	"fubarhouse/docker-ansible:xenial",
	"fubarhouse",
	"ubuntu1604",
	Ubuntu,
}

// Ubuntu 16.10
var Ubuntu1610 = Distribution{
	"",
	"ubuntu1610",
	true,
	"fubarhouse/docker-ansible:yakkety",
	"fubarhouse",
	"ubuntu1610",
	Ubuntu,
}

// Ubuntu 17.04
var Ubuntu1704 = Distribution{
	"",
	"ubuntu1704",
	true,
	"fubarhouse/docker-ansible:zesty",
	"fubarhouse",
	"ubuntu1704",
	Ubuntu,
}

// Ubuntu 17.10
var Ubuntu1710 = Distribution{
	"",
	"ubuntu1710",
	true,
	"fubarhouse/docker-ansible:artful",
	"fubarhouse",
	"ubuntu1710",
	Ubuntu,
}

// Ubuntu 18.04
var Ubuntu1804 = Distribution{
	"",
	"ubuntu1804",
	true,
	"fubarhouse/docker-ansible:bionic",
	"fubarhouse",
	"ubuntu1804",
	Ubuntu,
}

// CentOS 6
var JeffCentOS6 = Distribution{
	"",
	"centos6",
	true,
	"geerlingguy/docker-centos6-ansible:latest",
	"geerlingguy",
	"centos6",
	CentOS,
}

// CentOS 7
var JeffCentOS7 = Distribution{
	"",
	"centos7",
	true,
	"geerlingguy/docker-centos7-ansible:latest",
	"geerlingguy",
	"centos7",
	CentOS,
}

// Ubuntu 14.04
var JeffUbuntu1404 = Distribution{
	"",
	"ubuntu1404",
	true,
	"geerlingguy/docker-ubuntu1404-ansible:latest",
	"geerlingguy",
	"ubuntu1404",
	Ubuntu,
}

// Ubuntu 16.04
var JeffUbuntu1604 = Distribution{
	"",
	"ubuntu1604",
	true,
	"geerlingguy/docker-ubuntu1604-ansible:latest",
	"geerlingguy",
	"ubuntu1604",
	Ubuntu,
}

// Ubuntu 18.04
var JeffUbuntu1804 = Distribution{
	"",
	"ubuntu1804",
	true,
	"geerlingguy/docker-ubuntu1804-ansible:latest",
	"geerlingguy",
	"ubuntu1804",
	Ubuntu,
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
func GetDistribution(container, target, init, volume, user, distro string) (Distribution, error) {

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

	c, _ := DockerExec([]string{
		"images",
		container,
	}, false)

	if !strings.Contains(c, container) {
		log.Errorf("no valid image was found for '%v'\n", container)
		os.Exit(1)
	}

	return Distribution{},
	errors.New("could not find matching distribution")
}

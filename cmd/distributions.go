package cmd

// A Distribution declares the options to
// pass to Docker to run and test the container.
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
}

// CentOS 6
var CentOS6 = Distribution{
	"/sbin/init",
	"centos6",
	true,
	"",
	"geerlingguy/docker-centos6-ansible:latest",
}

// CentOS 7
var CentOS7 = Distribution{
	"/usr/lib/systemd/systemd",
	"centos7",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"geerlingguy/docker-centos7-ansible:latest",
}

// Ubuntu 14.04
var Ubuntu1404 = Distribution{
	"/sbin/init",
	"ubuntu1404",
	true,
	"",
	"geerlingguy/docker-ubuntu1404-ansible:latest",
}

// Ubuntu 16.04
var Ubuntu1604 = Distribution{
	"/lib/systemd/systemd",
	"ubuntu1604",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"geerlingguy/docker-ubuntu1604-ansible:latest",
}

// Ubuntu 18.04
var Ubuntu1804 = Distribution{
	"/lib/systemd/systemd",
	"ubuntu1804",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"geerlingguy/docker-ubuntu1804-ansible:latest",
}

// A slice of distributions.
var Distributions = []Distribution{
	CentOS6,
	CentOS7,
	Ubuntu1404,
	Ubuntu1604,
	Ubuntu1804,
}

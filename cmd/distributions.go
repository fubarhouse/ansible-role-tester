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
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:centos-6",
}

// CentOS 7
var CentOS7 = Distribution{
	"/usr/lib/systemd/systemd",
	"centos7",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:centos-7",
}

// Wheezy
var DebianWheezy = Distribution{
	"/bin/systemd",
	"wheezy",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:wheezy",
}

// Jessie
var DebianJessie = Distribution{
	"/lib/systemd",
	"jessie",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:jessie",
}

// Stretch
var DebianStretch = Distribution{
	"/bin/systemd",
	"stretch",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:stretch",
}

// Buster
var DebianBuster = Distribution{
	"/bin/systemd",
	"buster",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:buster",
}

// Fedora 24
var Fedora24 = Distribution{
	"/usr/lib/systemd/systemd",
	"fedora24",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:fedora-24",
}

// Fedora 25
var Fedora25 = Distribution{
	"/usr/lib/systemd/systemd",
	"fedora25",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:fedora-25",
}

// Fedora 26
var Fedora26 = Distribution{
	"/usr/lib/systemd/systemd",
	"fedora26",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:fedora-26",
}

// Fedora 27
var Fedora27 = Distribution{
	"/usr/lib/systemd/systemd",
	"fedora27",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:fedora-27",
}

// Fedora 28
var Fedora28 = Distribution{
	"/usr/lib/systemd/systemd",
	"fedora28",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:fedora-28",
}

// Ubuntu 12.04
var Ubuntu1204 = Distribution{
	"/sbin/init",
	"ubuntu1204",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:precise",
}

// Ubuntu 12.10
var Ubuntu1210 = Distribution{
	"/sbin/init",
	"ubuntu1210",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:quantal",
}

// Ubuntu 13.04
var Ubuntu1304 = Distribution{
	"/sbin/init",
	"ubuntu1304",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:raring",
}

// Ubuntu 13.10
var Ubuntu1310 = Distribution{
	"/sbin/init",
	"ubuntu1310",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:saucy",
}

// Ubuntu 14.04
var Ubuntu1404 = Distribution{
	"/sbin/init",
	"ubuntu1404",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:trusty",
}

// Ubuntu 14.10
var Ubuntu1410 = Distribution{
	"/sbin/init",
	"ubuntu1410",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:utopic",
}

// Ubuntu 15.04
var Ubuntu1504 = Distribution{
	"/sbin/init",
	"ubuntu1504",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:vivid",
}

// Ubuntu 15.10
var Ubuntu1510 = Distribution{
	"/sbin/init",
	"ubuntu1510",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:wily",
}

// Ubuntu 16.04
var Ubuntu1604 = Distribution{
	"/sbin/init",
	"ubuntu1604",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:xenial",
}

// Ubuntu 16.10
var Ubuntu1610 = Distribution{
	"/sbin/init",
	"ubuntu1610",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:yakkety",
}

// Ubuntu 17.04
var Ubuntu1704 = Distribution{
	"/sbin/init",
	"ubuntu1704",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:zesty",
}

// Ubuntu 17.10
var Ubuntu1710 = Distribution{
	"/sbin/init",
	"ubuntu1710",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:artful",
}

// Ubuntu 18.04
var Ubuntu1804 = Distribution{
	"/sbin/init",
	"ubuntu1804",
	true,
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
	"fubarhouse/docker-ansible:bionic",
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
}

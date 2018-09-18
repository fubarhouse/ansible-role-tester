package util

import (
	"errors"
	"os"
	"strings"

	"fmt"
	"reflect"

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
	Name       string
	Initialise string
	Volume     string
}

// CentOS Family Distribution Identifier
var CentOS = Family{
	"CentOS",
	"/sbin/init",
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
}

// Debian Family Distribution Identifier
var Debian = Family{
	"Debian",
	"/bin/systemd",
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
}

// Fedora Family Distribution Identifier
var Fedora = Family{
	"Fedora",
	"/usr/lib/systemd/systemd",
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
}

// Ubuntu Family Distribution Identifier
var Ubuntu = Family{
	"Ubuntu",
	"/sbin/init",
	"/sys/fs/cgroup:/sys/fs/cgroup:ro",
}

// CentOS6 Distribution declaration
var CentOS6 = Distribution{
	"",
	"centos6",
	true,
	"fubarhouse/docker-ansible:centos-6",
	"fubarhouse",
	"centos6",
	CentOS,
}

// CentOS7 Distribution declaration
var CentOS7 = Distribution{
	"",
	"centos7",
	true,
	"fubarhouse/docker-ansible:centos-7",
	"fubarhouse",
	"centos7",
	CentOS,
}

// DebianWheezy Distribution declaration
var DebianWheezy = Distribution{
	"",
	"wheezy",
	true,
	"fubarhouse/docker-ansible:wheezy",
	"fubarhouse",
	"debian7",
	Debian,
}

// DebianJessie Distribution declaration
var DebianJessie = Distribution{
	"",
	"jessie",
	true,
	"fubarhouse/docker-ansible:jessie",
	"fubarhouse",
	"debian8",
	Debian,
}

// DebianStretch Distribution declaration
var DebianStretch = Distribution{
	"",
	"stretch",
	true,
	"fubarhouse/docker-ansible:stretch",
	"fubarhouse",
	"debian9",
	Debian,
}

// DebianBuster Distribution declaration
var DebianBuster = Distribution{
	"",
	"buster",
	true,
	"fubarhouse/docker-ansible:buster",
	"fubarhouse",
	"debian10",
	Debian,
}

// Fedora24 Distribution declaration
var Fedora24 = Distribution{
	"",
	"fedora24",
	true,
	"fubarhouse/docker-ansible:fedora-24",
	"fubarhouse",
	"fedora24",
	Fedora,
}

// Fedora25 Distribution declaration
var Fedora25 = Distribution{
	"",
	"fedora25",
	true,
	"fubarhouse/docker-ansible:fedora-25",
	"fubarhouse",
	"fedora25",
	Fedora,
}

// Fedora26 Distribution declaration
var Fedora26 = Distribution{
	"",
	"fedora26",
	true,
	"fubarhouse/docker-ansible:fedora-26",
	"fubarhouse",
	"fedora26",
	Fedora,
}

// Fedora27 Distribution declaration
var Fedora27 = Distribution{
	"",
	"fedora27",
	true,
	"fubarhouse/docker-ansible:fedora-27",
	"fubarhouse",
	"fedora27",
	Fedora,
}

// Fedora28 Distribution declaration
var Fedora28 = Distribution{
	"",
	"fedora28",
	true,
	"fubarhouse/docker-ansible:fedora-28",
	"fubarhouse",
	"fedora28",
	Fedora,
}

// Ubuntu1204 Distribution declaration
var Ubuntu1204 = Distribution{
	"",
	"ubuntu1204",
	true,
	"fubarhouse/docker-ansible:precise",
	"fubarhouse",
	"ubuntu1204",
	Ubuntu,
}

// Ubuntu1210 Distribution declaration
var Ubuntu1210 = Distribution{
	"",
	"ubuntu1210",
	true,
	"fubarhouse/docker-ansible:quantal",
	"fubarhouse",
	"ubuntu1210",
	Ubuntu,
}

// Ubuntu1304 Distribution declaration
var Ubuntu1304 = Distribution{
	"",
	"ubuntu1304",
	true,
	"fubarhouse/docker-ansible:raring",
	"fubarhouse",
	"ubuntu1304",
	Ubuntu,
}

// Ubuntu1310 Distribution declaration
var Ubuntu1310 = Distribution{
	"",
	"ubuntu1310",
	true,
	"fubarhouse/docker-ansible:saucy",
	"fubarhouse",
	"ubuntu1310",
	Ubuntu,
}

// Ubuntu1404 Distribution declaration
var Ubuntu1404 = Distribution{
	"",
	"ubuntu1404",
	true,
	"fubarhouse/docker-ansible:trusty",
	"fubarhouse",
	"ubuntu1404",
	Ubuntu,
}

// Ubuntu1410 Distribution declaration
var Ubuntu1410 = Distribution{
	"",
	"ubuntu1410",
	true,
	"fubarhouse/docker-ansible:utopic",
	"fubarhouse",
	"ubuntu1410",
	Ubuntu,
}

// Ubuntu1504 Distribution declaration
var Ubuntu1504 = Distribution{
	"",
	"ubuntu1504",
	true,
	"fubarhouse/docker-ansible:vivid",
	"fubarhouse",
	"ubuntu1504",
	Ubuntu,
}

// Ubuntu1510 Distribution declaration
var Ubuntu1510 = Distribution{
	"",
	"ubuntu1510",
	true,
	"fubarhouse/docker-ansible:wily",
	"fubarhouse",
	"ubuntu1510",
	Ubuntu,
}

// Ubuntu1604 Distribution declaration
var Ubuntu1604 = Distribution{
	"",
	"ubuntu1604",
	true,
	"fubarhouse/docker-ansible:xenial",
	"fubarhouse",
	"ubuntu1604",
	Ubuntu,
}

// Ubuntu1610 Distribution declaration
var Ubuntu1610 = Distribution{
	"",
	"ubuntu1610",
	true,
	"fubarhouse/docker-ansible:yakkety",
	"fubarhouse",
	"ubuntu1610",
	Ubuntu,
}

// Ubuntu1704 Distribution declaration
var Ubuntu1704 = Distribution{
	"",
	"ubuntu1704",
	true,
	"fubarhouse/docker-ansible:zesty",
	"fubarhouse",
	"ubuntu1704",
	Ubuntu,
}

// Ubuntu1710 Distribution declaration
var Ubuntu1710 = Distribution{
	"",
	"ubuntu1710",
	true,
	"fubarhouse/docker-ansible:artful",
	"fubarhouse",
	"ubuntu1710",
	Ubuntu,
}

// Ubuntu1804 Distribution declaration
var Ubuntu1804 = Distribution{
	"",
	"ubuntu1804",
	true,
	"fubarhouse/docker-ansible:bionic",
	"fubarhouse",
	"ubuntu1804",
	Ubuntu,
}

// Ubuntu1810 Distribution declaration
var Ubuntu1810 = Distribution{
	"",
	"ubuntu1810",
	true,
	"fubarhouse/docker-ansible:cosmic",
	"fubarhouse",
	"ubuntu1810",
	Ubuntu,
}

// JeffCentOS6 Distribution declaration
var JeffCentOS6 = Distribution{
	"",
	"centos6",
	true,
	"geerlingguy/docker-centos6-ansible:latest",
	"geerlingguy",
	"centos6",
	CentOS,
}

// JeffCentOS7 Distribution declaration
var JeffCentOS7 = Distribution{
	"",
	"centos7",
	true,
	"geerlingguy/docker-centos7-ansible:latest",
	"geerlingguy",
	"centos7",
	CentOS,
}

// JeffUbuntu1204 Distribution declaration
var JeffUbuntu1204 = Distribution{
	"",
	"ubuntu1204",
	true,
	"geerlingguy/docker-ubuntu1204-ansible:latest",
	"geerlingguy",
	"ubuntu1204",
	Ubuntu,
}

// JeffUbuntu1404 Distribution declaration
var JeffUbuntu1404 = Distribution{
	"",
	"ubuntu1404",
	true,
	"geerlingguy/docker-ubuntu1404-ansible:latest",
	"geerlingguy",
	"ubuntu1404",
	Ubuntu,
}

// JeffUbuntu1604 Distribution declaration
var JeffUbuntu1604 = Distribution{
	"",
	"ubuntu1604",
	true,
	"geerlingguy/docker-ubuntu1604-ansible:latest",
	"geerlingguy",
	"ubuntu1604",
	Ubuntu,
}

// JeffUbuntu1804 Distribution declaration
var JeffUbuntu1804 = Distribution{
	"",
	"ubuntu1804",
	true,
	"geerlingguy/docker-ubuntu1804-ansible:latest",
	"geerlingguy",
	"ubuntu1804",
	Ubuntu,
}

// JeffDebian8 Distribution declaration
var JeffDebian8 = Distribution{
	"",
	"debian8",
	true,
	"geerlingguy/docker-debian8-ansible:latest",
	"geerlingguy",
	"debian8",
	Debian,
}

// JeffDebian9 Distribution declaration
var JeffDebian9 = Distribution{
	"",
	"debian9",
	true,
	"geerlingguy/docker-debian9-ansible:latest",
	"geerlingguy",
	"debian9",
	Debian,
}

// JeffFedora24 Distribution declaration
var JeffFedora24 = Distribution{
	"",
	"fedora24",
	true,
	"geerlingguy/docker-fedora24-ansible:latest",
	"geerlingguy",
	"fedora24",
	Fedora,
}

// JeffFedora27 Distribution declaration
var JeffFedora27 = Distribution{
	"",
	"fedora27",
	true,
	"geerlingguy/docker-fedora27-ansible:latest",
	"geerlingguy",
	"fedora27",
	Fedora,
}

// Distributions is a slice of all distributions listed above.
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
	Ubuntu1810,
	JeffCentOS6,
	JeffCentOS7,
	JeffUbuntu1204,
	JeffUbuntu1404,
	JeffUbuntu1604,
	JeffUbuntu1804,
	JeffDebian8,
	JeffDebian9,
	JeffFedora24,
	JeffFedora27,
}

// NewCustomDistribution will return an empty distribution.
func NewCustomDistribution() *Distribution {
	return new(Distribution)
}

// CustomDistributionValueSet will set a field to a given value from a Distribution.
func CustomDistributionValueSet(dist *Distribution, key, value string) error {
	v := reflect.ValueOf(dist).Elem().FieldByName(key)
	if v.IsValid() {
		v.SetString(value)
		return nil
	}
	return errors.New("invalid key/value pair was specified")

}

// CustomFamilyValueSet will set a field to a given value from a Family.
func CustomFamilyValueSet(family *Family, key, value string) error {
	v := reflect.ValueOf(family).Elem().FieldByName(key)
	if v.IsValid() {
		v.SetString(value)
		return nil
	}
	return errors.New("invalid key/value pair was specified")

}

// CustomDistributionValueGet will get a field value from a Distribution.
func CustomDistributionValueGet(dist *Distribution, key string) (string, error) {
	s := reflect.ValueOf(dist).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		if typeOfT.Field(i).Name == key {
			f := s.Field(i)
			return fmt.Sprintf("%s", f.Interface()), nil
		}
	}
	return "", errors.New("could not find the specified field")
}

// CustomFamilyValueGet will get a field value from a Family.
func CustomFamilyValueGet(family *Family, key string) (string, error) {
	s := reflect.ValueOf(family).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		if typeOfT.Field(i).Name == key {
			f := s.Field(i)
			return fmt.Sprintf("%s", f.Interface()), nil
		}
	}
	return "", errors.New("could not find the specified field")
}

// GetDistribution will get the distribution object to allow dynamic
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

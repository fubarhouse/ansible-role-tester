# Ansible Role Tester

[![Build Status](https://img.shields.io/travis/fubarhouse/ansible-role-tester/master.svg?style=for-the-badge)](https://travis-ci.org/fubarhouse/ansible-role-tester)
![stability-experimental](https://img.shields.io/badge/stability-experimental-orange.svg?style=for-the-badge)
[![Go Report Card](https://img.shields.io/badge/go%20report-A%2B-brightgreen.svg?style=for-the-badge)](https://goreportcard.com/report/github.com/fubarhouse/ansible-role-tester)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg?style=for-the-badge)](https://raw.githubusercontent.com/fubarhouse/ansible-role-tester/master/LICENSE)

This program will allow users to mount [Ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html) roles into isolated Docker images for testing typically achieved using CI tools such as Travis CI.

The goal of this project is to facilitate better localized testing, which will secure greater success rates for development and idempotence before pushing to a public state.

This self-contained command line tool is inspired from a [script](https://gist.github.com/samdoran/c3d392ee697881fa33a1d1a65814a07b) by [@samdoran](https://github.com/samdoran) originally inspired by [@geerlingguy](https://github.com/geerlingguy)'s ansible role tester [script](https://gist.github.com/geerlingguy/73ef1e5ee45d8694570f334be385e181).

## Requirements

  * [Docker](https://www.docker.com/)

## Dependencies

None.

## Installation

  * If you have [Go](https://golang.org/) installed, you can install it using `go get github.com/fubarhouse/ansible-role-tester`
  * Releases will accompany compiled binaries excluding windows executables.
  
## Usage

Ansible-Role-Tester is based upon the [Cobra command line framework](https://github.com/spf13/cobra).

You can use it by typing `ansible-role-tester`, and by adding `-h` or `--help` you can see additional help.

## Selecting containers for testing

By entering user names and distributions available declared in the `distributions.go` file, you can selectively target any of them.

**Example use**:

````
ansible-role-tester full -u $USER -t $DISTRO
````

### Custom containers

Currently you are able to declare all parameters for the containers other than the intialize command and volume parameter.

You can do that with the following example:

````
ansible-role-tester full -i $IMAGE
````

### Available distributions

| user        | distro     | image                                        |
| ----------- | ---------- | -------------------------------------------- |
| fubarhouse  | centos6    | fubarhouse/docker-ansible:centos-6           |
| fubarhouse  | centos7    | fubarhouse/docker-ansible:centos-7           |
| fubarhouse  | debian7    | fubarhouse/docker-ansible:wheezy             |
| fubarhouse  | debian8    | fubarhouse/docker-ansible:jessie             |
| fubarhouse  | debian9    | fubarhouse/docker-ansible:stretch            |
| fubarhouse  | debian10   | fubarhouse/docker-ansible:buster             |
| fubarhouse  | fedora24   | fubarhouse/docker-ansible:fedora-24          |
| fubarhouse  | fedora25   | fubarhouse/docker-ansible:fedora-25          |
| fubarhouse  | fedora26   | fubarhouse/docker-ansible:fedora-26          |
| fubarhouse  | fedora27   | fubarhouse/docker-ansible:fedora-27          |
| fubarhouse  | fedora28   | fubarhouse/docker-ansible:fedora-28          |
| fubarhouse  | ubuntu1204 | fubarhouse/docker-ansible:precise            |
| fubarhouse  | ubuntu1210 | fubarhouse/docker-ansible:quantal            |
| fubarhouse  | ubuntu1304 | fubarhouse/docker-ansible:raring             |
| fubarhouse  | ubuntu1310 | fubarhouse/docker-ansible:saucy              |
| fubarhouse  | ubuntu1404 | fubarhouse/docker-ansible:trusty             |
| fubarhouse  | ubuntu1410 | fubarhouse/docker-ansible:utopic             |
| fubarhouse  | ubuntu1504 | fubarhouse/docker-ansible:vivid              |
| fubarhouse  | ubuntu1510 | fubarhouse/docker-ansible:wily               |
| fubarhouse  | ubuntu1604 | fubarhouse/docker-ansible:xenial             |
| fubarhouse  | ubuntu1610 | fubarhouse/docker-ansible:yakkety            |
| fubarhouse  | ubuntu1704 | fubarhouse/docker-ansible:zesty              |
| fubarhouse  | ubuntu1710 | fubarhouse/docker-ansible:artful             |
| fubarhouse  | ubuntu1804 | fubarhouse/docker-ansible:bionic             |
| geerlingguy | centos6    | geerlingguy/docker-centos6-ansible:latest    |
| geerlingguy | centos7    | geerlingguy/docker-centos7-ansible:latest    |
| geerlingguy | debian8    | geerlingguy/docker-debian8-ansible:latest    |
| geerlingguy | debian9    | geerlingguy/docker-debian9-ansible:latest    |
| geerlingguy | fedora24   | geerlingguy/docker-fedora24-ansible:latest    |
| geerlingguy | fedora27   | geerlingguy/docker-fedora27-ansible:latest    |
| geerlingguy | ubuntu1204 | geerlingguy/docker-ubuntu1204-ansible:latest |
| geerlingguy | ubuntu1404 | geerlingguy/docker-ubuntu1404-ansible:latest |
| geerlingguy | ubuntu1604 | geerlingguy/docker-ubuntu1604-ansible:latest |
| geerlingguy | ubuntu1804 | geerlingguy/docker-ubuntu1804-ansible:latest |

## Interesting uses.

The following command will execute properly inside a [DrupalVM](https://github.com/geerlingguy/drupal-vm) clone, however it won't include the configuration variables, but it's an interesting case which proves how flexible this tool can be.

````
ansible-role-tester full -p ../provisioning/playbook.yml -u fubarhouse -t centos7 -r provisioning/requirements.yml 
````

## Screenshots

<img style="display:inline-block; width:100%; max-width:100%" src="https://github.com/fubarhouse/ansible-role-tester/raw/master/screenshot-1.png">
<p>&nbsp;</p>
<img style="display:inline-block; width:100%; max-width:100%" src="https://github.com/fubarhouse/ansible-role-tester/raw/master/screenshot-2.png">

## License

MIT

## Author Information

This tool was created in 2018 by [Karl Hepworth](https://twitter.com/fubarhouse).

**Special thanks** to [@samdoran](https://github.com/samdoran) and [@geerlingguy](https://github.com/geerlingguy) for the original bash scripts and ideas which were used during creation.
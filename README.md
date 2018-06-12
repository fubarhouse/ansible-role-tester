# Ansible Role Tester

[![Build Status](https://img.shields.io/travis/fubarhouse/ansible-role-tester/master.svg?style=for-the-badge)](https://travis-ci.org/fubarhouse/ansible-role-tester)
[![Go Report Card](https://goreportcard.com/badge/github.com/fubarhouse/ansible-role-tester?style=for-the-badge)](https://goreportcard.com/report/github.com/fubarhouse/ansible-role-tester)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg?style=for-the-badge)](https://raw.githubusercontent.com/fubarhouse/ansible-role-tester/master/LICENSE)

This program will allow users to mount [Ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html) roles into isolated Docker images for isolated testing typically replicating testing which is occuring via CI tools such as Travis CI. It should provide for better testing, and greater success rates before pushing to a public state.

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

It's generally encouraged to use the `full` command as we're working on refining the others and the `full` command performs end-to-end tests including container creation, requirements, idempotence, role running and container removal.

## Selecting containers for testing

By entering user names and distributions available declared in the `distributions.go` file, you can selectively target any of them.

**Example use**: `ansible-role-tester full -u fubarhouse -t ubuntu1804`

| user        | distro     |
| ----------- | ---------- |
| fubarhouse  | centos6    |
| fubarhouse  | centos7    |
| fubarhouse  | debian7    |
| fubarhouse  | debian8    |
| fubarhouse  | debian9    |
| fubarhouse  | debian10   |
| fubarhouse  | fedora24   |
| fubarhouse  | fedora25   |
| fubarhouse  | fedora26   |
| fubarhouse  | fedora27   |
| fubarhouse  | fedora28   |
| fubarhouse  | ubuntu1204 |
| fubarhouse  | ubuntu1210 |
| fubarhouse  | ubuntu1304 |
| fubarhouse  | ubuntu1310 |
| fubarhouse  | ubuntu1404 |
| fubarhouse  | ubuntu1410 |
| fubarhouse  | ubuntu1504 |
| fubarhouse  | ubuntu1510 |
| fubarhouse  | ubuntu1604 |
| fubarhouse  | ubuntu1610 |
| fubarhouse  | ubuntu1704 |
| fubarhouse  | ubuntu1710 |
| fubarhouse  | ubuntu1804 |
| geerlingguy | centos6    |
| geerlingguy | centos7    |
| geerlingguy | ubuntu1404 |
| geerlingguy | ubuntu1604 |
| geerlingguy | ubuntu1804 |

### Custom containers

Currently you are able to declare all parameters for the containers other than the intialize command and volume parameter.

You can do that with the following example:

`ansible-role-tester full -i fubarhouse/docker-ansible:bionic`

## License

MIT / BSD

## Author Information

This tool was created in 2018 by [Karl Hepworth](https://twitter.com/fubarhouse).

**Special thanks** to [@samdoran](https://github.com/samdoran) and [@geerlingguy](https://github.com/geerlingguy) for the original bash scripts and ideas which were used during creation.
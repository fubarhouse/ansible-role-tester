#!/bin/bash
# Test Ansile role inside container
#
# Inspired by the wonderful work of Jeff Geerling (@geerlingguy)
# https://gist.github.com/geerlingguy/73ef1e5ee45d8694570f334be385e181

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NEUTRAL='\033[0m'


function show_help {
    echo '-h --help             Print out help'
    echo '-D --destroy          Stop and destroy test container'
    echo '-r --role-path        Path to role. Defaults to current directory'
    echo '-i --image            Docker image OS to run. Defaults to samdoran/${DISTRO}-ansible:latest'
    echo '-d --distribution     Image distribution to run. Defaults to centos7'
    echo '-o --run-options      Run options passed to container'
    echo '-p --playbook         Name of playbook to run inside tests/.'
    echo '-s --shell            Log in to running container'
    echo '-f --friendly-name    Sets a friendly container ID.'
}


function show_error {
    echo "Unknown option: '$1'"
    exit 1
}

function show_ansible_version {
    docker exec --tty "$CONTAINER_ID" ansible --version
}

function run_container {
    # Sanity check to make sure this is a role
    if [[ ! -d "$ROLE_PATH/tasks" ]]; then
        echo "This does not appear to be a valid role."
        exit 1
    else
        running_containers=$(docker ps -f status=running --format '{{.Names}}')

        # Run container in detached state if it's not already running
        if [[ ! "$running_containers" =~ "$CONTAINER_ID" ]]; then
            echo "Running container $CONTAINER_ID"
            docker run --detach \
                --volume="$ROLE_PATH":"$ROLE_PATH_IN_CONTAINER":ro \
                --name $CONTAINER_ID \
                $RUN_OPTS \
                $IMAGE_NAME \
                $INIT
        else
            echo "Container $CONTAINER_ID already running. Continuing..."
        fi

        # Install dependencies
        if [[ -f "$ROLE_PATH/$REQUIREMENTS_FILE" ]]; then
            echo "Installing role dependencies"
            docker exec --tty "$CONTAINER_ID" ansible-galaxy install -r "$ROLE_PATH_IN_CONTAINER/$REQUIREMENTS_FILE"
        fi

    fi

}

function test_role {
    run_container
    show_ansible_version

    # Ansible syntax check.
    printf ${GREEN}"Checking role syntax..."${NEUTRAL}
    docker exec --tty "$CONTAINER_ID" ansible-playbook "${ROLE_PATH_IN_CONTAINER}/tests/${PLAYBOOK}" --syntax-check > /dev/null
    printf "${GREEN}PASS\n"${NEUTRAL}

    # Test role.
    printf "\n${GREEN}Running the role...${NEUTRAL}\n"
    docker exec --tty "$CONTAINER_ID" ansible-playbook "${ROLE_PATH_IN_CONTAINER}/tests/${PLAYBOOK}"

    # Test role idempotence.
    printf "\n${GREEN}Testing role idempotence...${NEUTRAL}\n\n"
    idempotence_log=$(mktemp)
    docker exec --tty "$CONTAINER_ID" ansible-playbook "${ROLE_PATH_IN_CONTAINER}/tests/${PLAYBOOK}" | tee -a $idempotence_log

    tail $idempotence_log | grep -q 'changed=0.*failed=0' \
        && (printf "\n${GREEN}Idempotence test: PASS${NEUTRAL}\n\n" && exit 0) \
        || (printf "\n${GREEN}Idempotence test: ${RED}FAIL${NEUTRAL}\n\n" && exit 1)

}


function destroy_container {
    echo "Stopping and removing $1"
    docker stop "$1" > /dev/null
    docker rm "$1" > /dev/null
}


function get_shell {
    run_container

    echo
    echo "Logging in to $1"
    echo "To run the test inside the container run"
    echo "ansible-playbook ${ROLE_PATH_IN_CONTAINER}/tests/${PLAYBOOK}"
    echo

    docker exec -it "$1" bash
}


# --- Main --- #

ROLE_PATH_IN_CONTAINER="/usr/share/ansible/roles/role_under_test"
PLAYBOOK=travis.yml
REQUIREMENTS_FILE="tests/requirements.yml"


# Get command line options
while [ "$#" -gt 0 ]; do
    case "$1" in
        -h|--help)
            show_help
            shift
            exit 0
            ;;

        -D|--destroy)
            ACTION='destroy'
            shift
            ;;

        -r|--role-path)
            ROLE_PATH="$2"
            shift
            shift
            ;;

        -o|--run-options)
            RUN_OPTS="$2"
            shift
            shift
            ;;

        -p|--playbook)
            PLAYBOOK="$2"
            shift
            shift
            ;;

        -i|--image)
            IMAGE_NAME="$2"
            shift
            shift
            ;;

        -d|--distribution)
            DISTRO="$2"
            shift
            shift
            ;;

        -s|--shell)
            ACTION='shell'
            shift
            ;;

        -t|--test-only)
            ACTION='test'
            shift
            ;;

        -f|--friendly-name)
            FRIENDLY_NAME=1
            shift
            ;;

        *)
            show_error "$1"
            ;;
    esac

done

ACTION=${ACTION:-'run'}
DISTRO=${DISTRO:-'centos7'}
ROLE_PATH="${ROLE_PATH:-$(pwd)}"
IMAGE_NAME=${IMAGE_NAME:-"samdoran/${DISTRO}-ansible:latest"}
CONTAINER_ID=${CONTAINER_ID:-$(date +%s)}

if [[ $FRIENDLY_NAME ]]; then
    CONTAINER_PREFIX=$(basename "$ROLE_PATH")
    CONTAINER_ID="$CONTAINER_PREFIX-$DISTRO-test"
fi

# Do the work
case "$ACTION" in
    shell)
        run_container
        get_shell "$CONTAINER_ID"
        ;;

    run)
        test_role
        destroy_container "$CONTAINER_ID"
        ;;

    destroy)
        destroy_container "$CONTAINER_ID"
        ;;

    test)
        test_role
        ;;
esac
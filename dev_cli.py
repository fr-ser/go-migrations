#!/usr/bin/env python3

import argparse
import os
import subprocess
import sys
import time

def run_all_tests():
    """Run all tests (including integration tests). Takes longer"""

    os.system("docker-compose -f docker-compose.test.yaml up --detach")
    _wait_for_database()

    # using os.systems as otherwise (path?) problems arise
    print("")
    exit_code = os.system("LOG_LEVEL=DEBUG gotest ./...")
    print("")

    os.system(
        "docker-compose -f docker-compose.test.yaml down --remove-orphans --timeout 1 --volumes"
    )
    if exit_code != 0:
        # the code can be too large, resulting in inconsistent behavior
        sys.exit(1)

def run_unit_tests():
    """Run only unit tests (excluding integration tests). Is faster"""

    exit_code = os.system("LOG_LEVEL=DEBUG gotest ./... -tags=unit")
    if exit_code != 0:
        # the code can be too large, resulting in inconsistent behavior
        sys.exit(1)


def _wait_for_database():
    end_time = time.time() + 5
    error = None

    while time.time() < end_time:
        try:
            subprocess.run(
                "docker-compose -f docker-compose.test.yaml exec database pg_isready",
                shell=True, check=True, capture_output=True,
            )
            return
        except subprocess.CalledProcessError as err:
            error = err
            time.sleep(0.3)
    else:
        print("Database did not start up correctly")
        raise error



def install_go_tools():
    os.system(r"""cat tools.go | grep _ | awk -F'"' '{print $2}' | xargs -tI % go install %""")


def get_arguments():
    parser = argparse.ArgumentParser(description='Helper CLI for the development environment.')
    parser.add_argument("command", help="commands to execute", choices=["test", "install_tools"])
    parser.add_argument("--unit-only", "-u", action="store_true",
                        help="Only run unit tests - skip long running integration test (setup)")
    return parser.parse_args()


if __name__ == "__main__":
    args = get_arguments()

    if args.command == "test":
        if args.unit_only:
            run_unit_tests()
        else:
            run_all_tests()
    elif args.command == "install_tools":
        install_go_tools()
    else:
        print(f"Unsupported command {args.command}")

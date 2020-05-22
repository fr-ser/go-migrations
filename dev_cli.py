#!/usr/bin/env python3

import argparse
import os
import sys

def run_tests():
    os.system("docker-compose -f docker-compose.test.yaml up --detach")
    exit_code = os.system("LOG_LEVEL=DEBUG gotest ./...")
    os.system("docker-compose -f docker-compose.test.yaml down --remove-orphans --timeout 1")
    if exit_code != 0:
        # the code can be too large, resulting in inconsistent behavior
        sys.exit(1)


def install_go_tools():
    os.system(r"""cat tools.go | grep _ | awk -F'"' '{print $2}' | xargs -tI % go install %""")


def get_arguments():
    parser = argparse.ArgumentParser(description='Helper CLI for the development environment.')
    parser.add_argument("command", help="commands to execute", choices=["test", "install_tools"],)
    return parser.parse_args()


if __name__ == "__main__":
    args = get_arguments()

    if args.command == "test":
        run_tests()
    elif args.command == "install_tools":
        install_go_tools()
    else:
        print(f"Unsupported command {args.command}")

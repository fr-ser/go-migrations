#!/usr/bin/env python3

import argparse
import os


def run_tests():
    os.system("gotest ./...")


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

import ipaddress
import argparse
import random
import sys
import os

from typing import Iterable, List, TextIO, TypeVar


def parse_port_range(port_range: str) -> Iterable[int]:
    """
    Parse a port range string into a list of ports.
    """

    if "," in port_range:
        for part in port_range.split(","):
            yield from parse_port_range(part)

    elif "-" in port_range:
        start, end = map(int, port_range.split("-"))

        if end < start:
            raise ValueError(f"Invalid port range: {port_range}")

        yield from range(start, end + 1)

    else:

        yield int(port_range)


def parse_cidr(cidr: str) -> Iterable[str]:
    """
    Parse a CIDR string into a list of IP addresses.
    """
    yield from map(
        lambda k: k.compressed,
        ipaddress.ip_network(cidr).hosts(),
    )


T = TypeVar("T")


def randomize(target: List[T]) -> List[T]:
    """
    Randomize the order of a list.
    """
    random.shuffle(target)
    return target


def main(cidr: str, port_range: str, output: TextIO = sys.stdout) -> None:
    hosts: List[str] = []

    for ip in parse_cidr(cidr):
        hosts.extend(f"{ip}:{port}" for port in parse_port_range(port_range))

    for host in randomize(hosts):
        output.write(f"{host}\n")


if __name__ == "__main__":

    parser = argparse.ArgumentParser(
        description="Generate a list of IP addresses and ports."
    )

    parser.add_argument(
        "-p",
        "--ports",
        type=str,
        required=True,
        help="Port range to generate. (nmap fmt)",
        metavar="<ports>",
        dest="ports",
    )

    parser.add_argument(
        "-c",
        "--cidr",
        type=str,
        required=True,
        help="CIDR to generate.",
        metavar="<cidr>",
        dest="cidr",
    )

    parser.add_argument(
        "-o",
        "--output",
        type=argparse.FileType("w"),
        default=sys.stdout,
        help="Output file.",
        metavar="<file>",
        dest="output",
    )

    arguments = parser.parse_args()

    if os.path.exists(arguments.cidr) and os.path.isfile(arguments.cidr):

        with open(arguments.cidr, "r") as f:
            map(
                lambda x: main(x, arguments.ports, arguments.output),
                f.readlines(),
            )

    else:
        main(arguments.cidr, arguments.ports, arguments.output)

    sys.exit(1)

nano-storage
============

This is a fun experiment in storing arbitrary files on the Nano block lattice.

It uses the `representative` field to store 32-byte chunks of a file per block.

You will need a powerful GPU to generate blocks quickly if you want to store larger files.

Usage:

    -address string
            read file from NANO address
    -file string
            write file to a NANO address
    -rpc string
            RPC URL to use (default "https://mynano.ninja/api/node")

Example:

    nano-storage -file <path to file to write>
    nano-storage -address <NANO address previously written by this tool>

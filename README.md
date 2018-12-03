These are Go Utilities for LXD, that are are complementary to the lxc command.

They can be run on an LXD host, or in LXD containers that have access to the lxd daemon either via a unix socket, or by https.

- lxdtool:
    - rolling snapshots for containers.
    - profile import/export
    - a snapshot server that can be used by containers to manage their own snapshots
    
- snapshot:
    - a snapshot client to the snapshot server.  It can be run from inside any container to create/delete/list its own snapshots


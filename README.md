These are Go Utilities for LXD, that are are complementary to the lxc command.

They can be run on an LXD host, or in LXD containers that have access to the lxd daemon either via a unix socket, or by https.

- lxdtool:
    - rolling snapshots for containers.
    - profile import/export
    - a snapshot server that can be used by containers to manage their own snapshots
    
- snapshot:
    - a snapshot client to the snapshot server.  It can be run from inside any container to create/delete/list its own snapshots

## Examples:
- crontab lines for hourly, daily, and weekly snapshots, two each:
```
50 * * * * lxdtool snap create --all --running --period 1h --count 2 --name auto_hour
05 02 * * * lxdtool snap create -ar --period 1d --count 2 --name auto_day
10 02 * * 0 lxdtool snap create -ar --period 7d --count 2 --name auto_week
```

- find out which container a process belongs to, using its host pid:

`lxdtool container find <pid>`

- export all profiles to directory x:

`lxdtool profile export -a -d x`

- import profiles from files "a", "b"

`lxdtool profile import a b`

- run a snapshot server that allows all containers with the "snapshot" profile to connect via the snapshot command and create/list/delete snapshots of themselves:

`lxdtool snapshot-server --profile snapshot --listen :8080`
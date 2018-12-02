config: {}
description: ""
devices:
  lxdsocket:
    bind: container
    connect: unix:/var/snap/lxd/common/lxd/unix.socket
    gid: "0"
    listen: unix:/usr/local/lib/lxd.socket
    mode: "0660"
    security.gid: "999"
    security.uid: "65534"
    type: proxy
    uid: "0"
name: lxdtool

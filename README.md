# Container from Scratch in Go

## Run the Container

```bash
sudo go run main.go run /bin/bash
```

## List Namespaces

List all namespaces on the system:
```bash
lsns
```

List specific namespace types:
```bash
sudo lsns -t pid
sudo lsns -t uts
sudo lsns -t mnt
```

Check namespaces for a specific process:
```bash
sudo ls -la /proc/[PID]/ns/
```

Find container processes:
```bash
ps aux | grep "/bin/bash"
pgrep -a bash
```


## Testing Container Resource Limits

Inside the container, run a process and check cgroup membership:
```bash
sleep 1000 &
pidof sleep
cat /sys/fs/cgroup/mycontainer/cgroup.procs
```

## Key Concepts

**Cgroup and Proc Filesystems**: Pseudo filesystems used by the kernel to communicate with user space. User space can configure kernel behavior through these interfaces (e.g., via cgroup).

**Namespaces**: Control what the container can see (isolated view of system resources like PIDs, network, mounts, etc.)

**Cgroups**: Control the resources the container can use (CPU, memory, I/O limits, etc.)

## Rootless Containers

Rootless containers allow non-root users to create and run containers by using user namespaces (`CLONE_NEWUSER`). This maps the unprivileged user to root (UID 0) inside the container for isolation without requiring actual root privileges.

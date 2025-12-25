//go:build linux
// +build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"path/filepath"
	"strconv"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		fmt.Println("Unknown command")
		return
	}
	fmt.Println("Hello, World!")
}

func run() {
	fmt.Println("Running the application...", os.Args[2:])

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		Unshareflags: syscall.CLONE_NEWNS,
		UidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getuid(), Size: 1},
		},
		GidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getgid(), Size: 1},
		},
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Container PID: %d\n", cmd.Process.Pid)

	if err := cmd.Wait(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func child() {
	fmt.Println("Inside the container:", os.Args[2:])

	// cg() // Disabled for rootless container

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := syscall.Sethostname([]byte("container")); err != nil {
		fmt.Println("Error setting hostname:", err)
		os.Exit(1)
	}

	syscall.Chroot("/home/vignesh/containerfs")
	if err := os.Chdir("/"); err != nil {
		fmt.Println("Error changing directory:", err)
		os.Exit(1)
	}

	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		fmt.Println("Error mounting proc:", err)
		os.Exit(1)
	}
	defer syscall.Unmount("/proc", 0)

	if err := syscall.Mount("tmpfs", "/dev", "tmpfs", 0, "mode=755"); err != nil {
		fmt.Println("Error mounting tmpfs on /dev:", err)
		os.Exit(1)
	}
	defer syscall.Unmount("/dev", 0)

	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func cg() {
	cgroup := "/sys/fs/cgroup/mycontainer"
	
	// Create cgroup directory
	if err := os.Mkdir(cgroup, 0755); err != nil && !os.IsExist(err) {
		fmt.Println("Error creating cgroup:", err)
		os.Exit(1)
	}

	// Set memory limit (100MB) - cgroup v2 syntax
	if err := os.WriteFile(filepath.Join(cgroup, "memory.max"), []byte("100000000"), 0644); err != nil {
		fmt.Println("Error setting memory limit:", err)
		os.Exit(1)
	}

	// Add current process to cgroup
	pid := os.Getpid()
	if err := os.WriteFile(filepath.Join(cgroup, "cgroup.procs"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		fmt.Println("Error adding process to cgroup:", err)
		os.Exit(1)
	}
}
package ctfd

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
)

// getUnusedPort finds an available port in the given range
func getUnusedPort(start, end int) (int, error) {
	for port := start; port <= end; port++ {
		addr := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", addr)
		if err == nil {
			listener.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no unused ports available in the range %d-%d", start, end)
}

// createCTFdContainer creates a new CTFd container on an unused port
func createCTFdContainer(containerPrefix string) error {
	// Find an unused port
	startport := 10000
	endport := 65500
	port, err := getUnusedPort(startport, endport)
	if err != nil {
		return fmt.Errorf("failed to find unused port: %v", err)
	}
	fmt.Printf("Found unused port: %d\n", port)

	// Create a unique container name
	containerName := containerPrefix

	// Prepare docker run command
	dockerCmd := exec.Command("docker", "run",
		"--detach",
		"--name", containerName,
		"-p", fmt.Sprintf("%d:8000", port),
		"-e", fmt.Sprintf("CTFd_DATABASE_URL=sqlite:///%s.db", containerName),
		"--restart", "always",
		"ravenspar/ctfbattle:latest",
	)

	// Run the docker command
	output, err := dockerCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create container: %v\nOutput: %s", err, string(output))
	}

	fmt.Printf("CTFd container '%s' started on port %d\n", containerName, port)
	fmt.Printf("Access the platform at: http://localhost:%d\n", port)

	return nil
}

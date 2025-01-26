package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/cors"
)

func WrapHandler(handler http.Handler) http.Handler {

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	return corsHandler.Handler(handler)
}

func ExecuteBashScript(script string) (int, string, error) {

	cmd := exec.Command("bash", "-c", script)

	err := cmd.Start()
	if err != nil {
		return 0, "", fmt.Errorf("failed to start script: %v", err)
	}

	pid := cmd.Process.Pid

	time.Sleep(1 * time.Second)

	psCmd := exec.Command("ps", "-ef")
	psOutput, err := psCmd.Output()
	if err != nil {
		return 0, "", fmt.Errorf("failed to capture ps output: %v", err)
	}

	psOutputStr := string(psOutput)
	var rawCommand string
	for _, line := range strings.Split(psOutputStr, "\n") {
		if strings.Contains(line, fmt.Sprintf("%d", pid)) {
			rawCommand = line
			break
		}
	}

	log.Printf("Started script with PID: %d\nRaw command: %s", pid, rawCommand)

	err = cmd.Wait()
	if err != nil {
		return 0, "", fmt.Errorf("failed to wait for script: %v", err)
	}

	log.Printf("Script with PID %d finished executing", pid)

	return pid, rawCommand, nil
}

func StartService(serviceName string) error {
	var startCommand string
	switch serviceName {
	case "personalSite":
		startCommand = "npm run dev"
	case "basicCalculator":
		startCommand = "python3 app.py"
	default:
		return fmt.Errorf("unknown service: %s", serviceName)
	}

	_, rawCommand, err := ExecuteBashScript(startCommand)
	if err != nil {
		return fmt.Errorf("failed to start service %s: %v", serviceName, err)
	}

	log.Printf("Service %s started with command: %s", serviceName, rawCommand)

	return nil
}

func StopService(serviceName string) error {
	var startCommand string
	switch serviceName {
	case "personalSite":
		startCommand = "npm run dev"
	case "basicCalculator":
		startCommand = "python3 app.py"
	default:
		return fmt.Errorf("unknown service: %s", serviceName)
	}
	var stopCommand string
	stopCommand = fmt.Sprintf("pkill -f '%s'", startCommand)

	_, _, err := ExecuteBashScript(stopCommand)
	if err != nil {
		return fmt.Errorf("failed to stop service %s: %v", serviceName, err)
	}

	log.Printf("Service %s stopped successfully", serviceName)
	return nil
}

func RestartService(serviceName string) error {

	err := StopService(serviceName)
	if err != nil {
		return fmt.Errorf("failed to stop service %s: %v", serviceName, err)
	}

	err = StartService(serviceName)
	if err != nil {
		return fmt.Errorf("failed to start service %s: %v", serviceName, err)
	}

	log.Printf("Service %s restarted successfully", serviceName)
	return nil
}

func StartHandler(w http.ResponseWriter, r *http.Request) {
	serviceName := strings.TrimPrefix(r.URL.Path, "/api/start/")
	err := StartService(serviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Service %s started successfully", serviceName)
}

func StopHandler(w http.ResponseWriter, r *http.Request) {
	serviceName := strings.TrimPrefix(r.URL.Path, "/api/stop/")
	err := StopService(serviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Service %s stopped successfully", serviceName)
}

func RestartHandler(w http.ResponseWriter, r *http.Request) {
	serviceName := strings.TrimPrefix(r.URL.Path, "/api/restart/")
	err := RestartService(serviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Service %s restarted successfully", serviceName)
}

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Service manager is running...")

	http.Handle("/api/start/", WrapHandler(http.HandlerFunc(StartHandler)))
	http.Handle("/api/stop/", WrapHandler(http.HandlerFunc(StopHandler)))
	http.Handle("/api/restart/", WrapHandler(http.HandlerFunc(RestartHandler)))

	http.ListenAndServe(":8080", nil)
}

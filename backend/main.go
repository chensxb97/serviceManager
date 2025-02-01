package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/rs/cors"
)

type ActionPayload struct {
	Service string `json:"service"`
	Action  string `json:"action"`
}

var actionsQueue = make(chan ActionPayload, 10)
var serviceStates = map[string]string{
	"app1": "stopped",
	"app2": "stopped",
}
var mu = sync.RWMutex{}

func WrapHandler(handler http.Handler) http.Handler {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	return corsHandler.Handler(handler)
}

func checkServiceStates() error {
	mu.Lock()
	defer mu.Unlock()

	for service := range serviceStates {

		cmd := exec.Command("pgrep", "-f", service)
		output, err := cmd.Output()

		status := "stopped"
		if err == nil && len(output) > 0 {
			status = "running"
		}

		serviceStates[service] = status
		fmt.Printf("Updated State for service: %s: %s\n", service, status)
	}

	return nil
}

func ExecuteBashScript(script string) (int, string, error) {
	cmd := exec.Command("bash", "-c", script)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

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

	log.Printf("Executed script with PID: %d\nRaw command: %s", pid, rawCommand)

	return pid, rawCommand, nil
}

func StartService(serviceName string) error {
	var startCommand string
	switch serviceName {
	case "app1":
		startCommand = "cd /Users/username/workspace/app1 && npm run dev"
	case "app2":
		startCommand = "cd /Users/username/workspace/app2 && python3 app.py"
	default:
		return fmt.Errorf("unknown service: %s", serviceName)
	}

	_, rawCommand, err := ExecuteBashScript(startCommand)
	if err != nil {
		return fmt.Errorf("failed to start service %s: %v", serviceName, err)
	}

	fmt.Printf("Service %s started with command: %s\n", serviceName, rawCommand)

	return nil
}

func StopService(serviceName string) error {
	var stopRegex string
	switch serviceName {
	case "app1":
		stopRegex = "app1"
	case "app2":
		stopRegex = "app2"
	default:
		return fmt.Errorf("unknown service: %s", serviceName)
	}
	var stopCommand string
	stopCommand = fmt.Sprintf("pkill -f '%s'", stopRegex)

	_, _, err := ExecuteBashScript(stopCommand)
	if err != nil {
		return fmt.Errorf("failed to stop service %s: %v", serviceName, err)
	}

	fmt.Printf("Service %s stopped successfully\n", serviceName)
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

	fmt.Printf("Service %s restarted successfully\n", serviceName)
	return nil
}

func processAction(payload ActionPayload) error {
	switch payload.Action {
	case "start":
		err := StartService(payload.Service)
		if err != nil {
			return fmt.Errorf("Error while processing start action for service %s: %v", payload.Service, err)
		}
	case "stop":
		err := StopService(payload.Service)
		if err != nil {
			return fmt.Errorf("Error while processing stop action for service %s: %v", payload.Service, err)
		}
	case "restart":
		err := RestartService(payload.Service)
		if err != nil {
			return fmt.Errorf("Error while processing restart action for service %s: %v", payload.Service, err)
		}
	default:
		return fmt.Errorf("Unknown action: %s", payload.Action)
	}

	mu.Lock()
	defer mu.Unlock()

	switch payload.Action {
	case "start", "restart":
		serviceStates[payload.Service] = "running"
	case "stop":
		serviceStates[payload.Service] = "stopped"
	}

	return nil
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Service Manager is running safe and sound..."))
}

func ServicesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid HTTP method, only GET is allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := checkServiceStates(); err != nil {
		http.Error(w, fmt.Sprintf("Error while checking service states: %v", err), http.StatusInternalServerError)
		return
	}

	mu.RLock()
	defer mu.RUnlock()
	if err := json.NewEncoder(w).Encode(serviceStates); err != nil {
		http.Error(w, "Error while encoding payload in services handler", http.StatusInternalServerError)
	}

}

func ActionHandler(w http.ResponseWriter, r *http.Request) {
	var payload ActionPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "Invalid payload, please define compulsory service and action fields", http.StatusBadRequest)
		return
	}
	validActions := []string{"start", "stop", "restart"}
	valid := false
	for _, validAction := range validActions {
		if payload.Action == validAction {
			valid = true
			break
		}
	}
	if !valid {
		http.Error(w, fmt.Sprintf("Invalid action. Valid actions are: %v", validActions), http.StatusBadRequest)
		return
	}

	timeout := time.After(10 * time.Second)
	for {
		select {
		case actionsQueue <- payload:

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("Validated and Submitted Action %s for service %s", payload.Action, payload.Service)))
			return
		case <-timeout:

			http.Error(w, "Failed to queue action due to timeout", http.StatusInternalServerError)
			return
		default:

			time.Sleep(2 * time.Second)
		}
	}
}

func main() {
	log.Println("Service manager is running safe and sound...")

	if err := checkServiceStates(); err != nil {
		log.Printf("Error while checking service states")
	}

	go func() {
		for action := range actionsQueue {
			err := processAction(action)
			if err != nil {
				log.Printf("Error processing action for service %s: %v", action.Service, err)
			}
		}
	}()

	http.Handle("/", WrapHandler(http.HandlerFunc(RootHandler)))
	http.Handle("/api/services", WrapHandler(http.HandlerFunc(ServicesHandler)))
	http.Handle("/api/action", WrapHandler(http.HandlerFunc(ActionHandler)))

	http.ListenAndServe(":8080", nil)
}

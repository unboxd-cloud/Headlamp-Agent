package setup

import (
	"net"
	"os"
	"time"
)

type VerifyReport struct {
	OK        bool          `json:"ok"`
	ProfileID ProfileID     `json:"profileId"`
	Checks    []CheckResult `json:"checks"`
}

func VerifyLocalKubernetesOperator() VerifyReport {
	checks := []CheckResult{}
	checks = append(checks, LocalOperatorChecks()...)
	checks = append(checks, CheckCommand("ssh", "SSH Client", "ssh"))
	checks = append(checks, CheckKubeconfig())
	checks = append(checks, CheckTCPPort("surrealdb-port", "SurrealDB Port", "127.0.0.1:8000"))
	checks = append(checks, CheckTCPPort("opa-port", "OPA Port", "127.0.0.1:8181"))
	checks = append(checks, CheckTCPPort("node-agent-port", "Node Agent Port", "127.0.0.1:9080"))

	ok := true
	for _, check := range checks {
		if !check.OK {
			ok = false
			break
		}
	}

	return VerifyReport{
		OK:        ok,
		ProfileID: ProfileLocalKubernetesOperator,
		Checks:    checks,
	}
}

func CheckKubeconfig() CheckResult {
	path := os.Getenv("KUBECONFIG")
	if path == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			path = home + "/.kube/config"
		}
	}

	if path == "" {
		return CheckResult{ID: "kubeconfig", Name: "Kubeconfig", OK: false, Message: "could not resolve kubeconfig path"}
	}

	if _, err := os.Stat(path); err != nil {
		return CheckResult{ID: "kubeconfig", Name: "Kubeconfig", OK: false, Message: path + " not found"}
	}

	return CheckResult{ID: "kubeconfig", Name: "Kubeconfig", OK: true, Message: path}
}

func CheckTCPPort(id string, name string, address string) CheckResult {
	conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
	if err != nil {
		return CheckResult{ID: id, Name: name, OK: false, Message: address + " not reachable"}
	}
	_ = conn.Close()
	return CheckResult{ID: id, Name: name, OK: true, Message: address + " reachable"}
}

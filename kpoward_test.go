package kpoward_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goccy/kpoward"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func getKubeConfig() string {
	if v := os.Getenv("KUBECONFIG"); v != "" {
		return v
	}
	home := homedir.HomeDir()
	config := filepath.Join(home, ".kube", "config")
	if _, err := os.Stat(config); err == nil {
		return config
	}
	return ""
}

func loadConfig() (*rest.Config, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", getKubeConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}
	return cfg, nil
}

var restCfg *rest.Config

func init() {
	var err error
	restCfg, err = loadConfig()
	if err != nil {
		panic(err)
	}
}

// This test assumes that we have a pod deployed by `make deploy`.
func TestKpoward(t *testing.T) {
	ctx := context.Background()
	clientset, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		t.Fatal(err)
	}
	podList, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{
		LabelSelector: "app=echo",
	})
	if err != nil {
		t.Fatal(err)
	}
	if podList == nil || len(podList.Items) == 0 {
		t.Fatalf("failed to get pod")
	}
	podName := podList.Items[0].Name
	log.Printf("pod name: %s", podName)
	kpow := kpoward.New(restCfg, podName, 8080)
	if err := kpow.Run(ctx, func(ctx context.Context, localPort uint16) error {
		log.Printf("localPort: %d", localPort)
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/echo", localPort))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		got := strings.TrimRight(string(body), "\n")
		if got != "hello" {
			return fmt.Errorf("unexpected response: expected hello but got [%s]", got)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("cmd")

func printVersion() {
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
}

func main() {
	flag.Parse()

	// The logger instantiated here can be changed to any logger
	// implementing the logr.Logger interface. This logger will
	// be propagated through the whole operator, generating
	// uniform and structured logs.
	logf.SetLogger(logf.ZapLogger(false))

	printVersion()

	nodeName, ok := os.LookupEnv("NODE_NAME")
	if !ok {
		log.Error(errors.New("NODE_NAME env must be set"), "")
		os.Exit(1)
	}

	iface, ok := os.LookupEnv("INTERFACE")
	if !ok {
		log.Error(errors.New("INTERFACE env must be set"), "")
		os.Exit(1)
	}

	annotation, ok := os.LookupEnv("ANNOTATION")
	if !ok {
		log.Error(errors.New("ANNOTATION env must be set"), "")
		os.Exit(1)
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	intAddr, err := getInterfaceAddress(iface)
	if err != nil {
		errorAndExit(err)
	}

	node, err := clientset.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
	if err != nil {
		log.Error(err, "error getting node from api")
		os.Exit(1)
	}

	node.Annotations[annotation] = intAddr
	clientset.CoreV1().Nodes().Update(node)

	log.Info("Updated node annotation %s with IP %s", annotation, intAddr)

}

func errorAndExit(err error) {
	log.Error(err, "")
	os.Exit(1)
}

func getInterfaceAddress(name string) (string, error) {
	int, err := net.InterfaceByName(name)
	if err != nil {
		return "", err
	}

	addresses, err := int.Addrs()
	if err != nil {
		return "", err
	}

	if ipnet := addresses[0].(*net.IPNet).IP.String(); ipnet != "" {
		return ipnet, nil
	}

	return "", fmt.Errorf("could not get interface address")
}

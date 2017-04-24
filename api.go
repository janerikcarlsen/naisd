package main

import (
	"k8s.io/client-go/rest"
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"net/http"
	_"io/ioutil"
	"goji.io"
	"goji.io/pat"
	"io/ioutil"
	"k8s.io/client-go/pkg/util/json"
)

var clientset = initKubeConfiguration()

func initKubeConfiguration() *kubernetes.Clientset {
	kubeconfig := flag.String("kubeconfig", "", "Path to a kubeconfig file")
	flag.Parse()

	config, err := getClientConfig(*kubeconfig)

	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err.Error())
	}

	return clientset
}

func main() {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/pods"), listpods)
	mux.HandleFunc(pat.Post("/echo"), echo)

	serveLocation := "localhost:6969"
	fmt.Printf("serving @ %s\n", serveLocation)
	http.ListenAndServe(serveLocation, mux)
}

func listpods(w http.ResponseWriter, r *http.Request) {
	pods, err := clientset.CoreV1().Pods("").List(v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}

	output,_ := json.Marshal(pods.Items)

	fmt.Fprint(w, string(output))
}

func echo(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "%s", body)
}

// returns config using kubeconfig if provided, else from cluster context
// useful for testing locally w/minikube
func getClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}
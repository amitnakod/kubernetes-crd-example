package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/martin-helmich/kubernetes-crd-example/api/types/v1alpha1"
	clientV1alpha1 "github.com/martin-helmich/kubernetes-crd-example/clientset/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "/Users/amitnakod/.kube/config", "path to Kubernetes config file")
	flag.Parse()
}

func main() {
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		log.Printf("using in-cluster configuration")
		config, err = rest.InClusterConfig()
		fmt.Println("hhost:", os.Getenv("KUBERNETES_SERVICE_HOST"), "    PORT:", os.Getenv("KUBERNETES_SERVICE_PORT"))
		fmt.Println("Host", config.Host)
		fmt.Println("AuthProvider", config.AuthProvider)
		fmt.Println("bearerToken", config.BearerToken)
	} else {
		log.Printf("using configuration from '%s'", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		panic(err)
	}

	v1alpha1.AddToScheme(scheme.Scheme)

	clientSet, err := clientV1alpha1.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	project, err := clientSet.Projects("default").Get("example-project", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	// fmt.Printf("project found: %+v\n", project)
	project.SetName(project.GetObjectMeta().GetName())
	project.Spec.Replicas = 20
	// fmt.Printf("created project found: %+v\n", proj)
	_, err = clientSet.Projects("default").Update(project)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("new project found: %+v\n", updatedproject)

	proj := &v1alpha1.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name: "new-example1",
		},
		Spec: v1alpha1.ProjectSpec{
			Replicas: 10,
		},
	}

	// fmt.Printf("created project found: %+v\n", proj)
	newproject, err := clientSet.Projects("default").Create(proj)
	if err != nil {
		panic(err)
	}

	fmt.Printf("new project found: %+v\n", newproject)

	projects, err := clientSet.Projects("default").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("projects found: %+v\n", projects)

	store := WatchResources(clientSet)

	for {
		projectsFromStore := store.List()
		fmt.Printf("project in store: %d\n", len(projectsFromStore))

		time.Sleep(2 * time.Second)
	}
}

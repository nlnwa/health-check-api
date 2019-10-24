package k8s

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Query interface {
	GetHarvesterPodNames() ([]string, error)
}

func (c Client) GetHarvesterPodNames() ([]string, error) {
	// get pods in all the namespaces by omitting namespace
	// Or specify namespace to get pods in particular namespace
	pods, err := c.CoreV1().Pods("").List(meta.ListOptions{
		LabelSelector: meta.LabelSelector{
			MatchLabels: map[string]string{"component": "local-veidemann-harvester"},
		}.String(),
	})
	if err != nil {
		panic(err.Error())
	}

	podNames := make([]string, len(pods.Items))
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
	}

	return podNames, nil
	// Examples for error handling:
	// - Use helper functions e.g. errors.IsNotFound()
	// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	//_, err = clientset.CoreV1().Pods("default").Get("example-xxxxx", metav1.GetOptions{})
	//if errors.IsNotFound(err) {
	//	fmt.Printf("Pod example-xxxxx not found in default namespace\n")
	//} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
	//	fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
	//} else if err != nil {
	//	panic(err.Error())
	//} else {
	//	fmt.Printf("Found example-xxxxx pod in default namespace\n")
	//}
}

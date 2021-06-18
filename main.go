package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	"os"
)

var KubeClient *kubernetes.Clientset

func main() {
	err := getKubeConfig()
	if err != nil {
		panic(err)
	}
	nameSpace := os.Getenv("NAMESPACE")
	labelName := os.Getenv("LABEL_NAME")
	labelValue := os.Getenv("LABEL_VALUE")
	todoContext := context.TODO()
	deployClient := KubeClient.AppsV1().Deployments(nameSpace)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := getDeployment(labelName, labelValue, nameSpace)
		if getErr != nil {
			return getErr
		}
		refreshUUID, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		envVars := result.Spec.Template.Spec.Containers[0].Env
		// pop the old refresh
		for k, v := range envVars {
			if v.Name == "refresh" {
				envVars[k] = envVars[len(envVars)-1]      // Copy last element to index i.
				envVars[len(envVars)-1] = corev1.EnvVar{} // Erase last element (write zero value).
				envVars = envVars[:len(envVars)-1]        // Truncate slice.
				break
			}
		}
		envVar := corev1.EnvVar{
			Name:      "refresh",
			Value:     refreshUUID.String(),
			ValueFrom: nil,
		}
		envVars = append(envVars, envVar)
		result.Spec.Template.Spec.Containers[0].Env = envVars
		_, updateErr := deployClient.Update(todoContext, &result, metav1.UpdateOptions{})
		if updateErr != nil {
			return updateErr
		}
		return updateErr
	})
	if retryErr != nil {
		panic(err)
	}
}

func getDeployment(labelName, labelValue, namespace string) (appsv1.Deployment, error) {
	listOpts := metav1.ListOptions{LabelSelector: labelName}
	deployments, err := KubeClient.AppsV1().Deployments(namespace).List(context.TODO(), listOpts)
	if err != nil {
		return appsv1.Deployment{}, err
	}
	for _, deploy := range deployments.Items {
		if val, ok := deploy.Labels[labelName]; ok {
			if val == labelValue {
				return deploy, nil
			}
		}
	}
	return appsv1.Deployment{}, fmt.Errorf("no deployment with label %s and value %s exists", labelName, labelValue)
}

func getKubeConfig() error {
	config, err := rest.InClusterConfig()
	KubeClient, err = kubernetes.NewForConfig(config)
	return err
}

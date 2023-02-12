package controllers

import (
	"fmt"
	telcolabs1alpha1 "tlhelper_operator/api/v1alpha1"
)

func generatePodArgs(instance *telcolabs1alpha1.Tlhelper) []string {
	var args []string

	if instance.Spec.Format != "" {
		args = append(args, "-fmt", fmt.Sprintf("%v", instance.Spec.Format))
	}

	if instance.Spec.LogLevel != "" {
		args = append(args, "-lvl", fmt.Sprintf("%v", instance.Spec.LogLevel))
	}

	if instance.Spec.Remote {
		args = append(args, "-remote")
	}
	return args
}

func generateImage(repository string, tag string) string {
	return fmt.Sprintf("%v:%v", repository, tag)
}

package main

import (
	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
)

var patchValue []corev1.EnvVar = []corev1.EnvVar{{
	Name:  "Password",
	Value: "",
}}

func mutate(pod *apiv1.Pod) {
	securityValue := "false"
	for i, container := range pod.Spec.Containers {
		if pod.Name == "pod1" {
			patchValue[0].Value = "ThisIsATestPassword"
			securityValue = "true"
			pod.Spec.Containers[i].Env = append(container.Env, patchValue[0])
		}
	}

	if pod.ObjectMeta.Labels == nil {
		pod.ObjectMeta.Labels = make(map[string]string)
	}

	pod.ObjectMeta.Labels["SecurityOk"] = securityValue

}

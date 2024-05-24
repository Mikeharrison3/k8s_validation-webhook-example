package main

import (
	apiv1 "k8s.io/api/core/v1"
)

const nametoCheckAgainst string = "testname"

type validationstruct struct {
	IsAllowed bool
	Result    string
}

func IsNameOkay(pod apiv1.Pod, respx *validationstruct) /**validationstruct*/ {

	respx.IsAllowed = pod.Name != nametoCheckAgainst

	if !respx.IsAllowed {
		respx.Result = "Failed isNameOK line 19 validate.go"
	}

	//return respx
}

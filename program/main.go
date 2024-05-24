package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"fmt"

	"k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"

	"errors"

	"k8s.io/api/admission/v1beta1"

	"encoding/json"

	"github.com/wI2L/jsondiff"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServerParameters struct {
	port     int
	certFile string
	keyFile  string
}

type patchOperation struct {
	Op    string      `JSON:"OP"`
	Path  string      `json:"path"`
	Value []v1.EnvVar `json:"value"`
}

var parameters ServerParameters

var (
	universalDeserializer = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
)

var config *rest.Config
var clientSet *kubernetes.Clientset

func main() {

	http.HandleFunc("/", HandleRoot)
	http.HandleFunc("/mutate", HandleMutate)
	log.Print("Staring server.")
	log.Fatal(http.ListenAndServeTLS(":8443", "/etc/webhook/certs/tls.crt", "/etc/webhook/certs/tls.key", nil))
	log.Print("Server Started")

}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!"))
}

func HandleMutate(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	err := os.WriteFile("/tmp/request", body, 0644)

	var admissionReviewReq v1beta1.AdmissionReview

	if _, _, err := universalDeserializer.Decode(body, nil, &admissionReviewReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Errorf("could not deserialize request: %v", err)
	} else if admissionReviewReq.Request == nil {
		w.WriteHeader(http.StatusBadRequest)
		errors.New("malformed admission review: request is nil")
	}

	fmt.Printf("Type: %v \t Event: %v \t Name: %v \n",
		admissionReviewReq.Request.Kind,
		admissionReviewReq.Request.Operation,
		admissionReviewReq.Request.Name,
	)

	if admissionReviewReq.Request.Operation == "DELETE" {

		admissionReviewResponse := v1beta1.AdmissionReview{
			Response: &v1beta1.AdmissionResponse{
				UID:     admissionReviewReq.Request.UID,
				Allowed: true,
			},
		}

		bytes, err := json.Marshal(&admissionReviewResponse)

		if err != nil {
			fmt.Errorf("marshaling response: %v", err)
		}

		w.Write(bytes)

		return
	}

	var pod apiv1.Pod
	pT := v1beta1.PatchTypeJSONPatch

	err = json.Unmarshal(admissionReviewReq.Request.Object.Raw, &pod)

	//check name

	var validate validationstruct
	validate.IsAllowed = false
	validate.Result = "ok"

	IsNameOkay(pod, &validate)

	var mpod = pod.DeepCopy()

	if err != nil {
		log.Fatal("could not unmarshal pod on admission request: %v", err)
	}

	mutate(&pod)
	admissionReviewResponse := v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			UID:       admissionReviewReq.Request.UID,
			Allowed:   validate.IsAllowed,
			PatchType: &pT,
			Result: &metav1.Status{
				Code:    200,
				Message: validate.Result,
			},
		},
	}

	patch, err := jsondiff.Compare(mpod, pod)

	if patch != nil {
		admissionReviewResponse.Response.Patch, err = json.Marshal(patch)
	}

	if err != nil {
		fmt.Errorf("marshaling response: %v", err)
	}

	admissionReviewResponse.Response.AuditAnnotations = map[string]string{
		"mutateme": "Yeah, I was mutated",
	}

	bytes, err := json.Marshal(&admissionReviewResponse)

	if err != nil {
		fmt.Errorf("marshaling response: %v", err)
	}

	w.Write(bytes)

}

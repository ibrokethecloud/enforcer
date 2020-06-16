package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"k8s.io/api/admission/v1beta1"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var (
	// initalize deserializers
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
	// GlobalAnnotationsKey is annotation needed Namespace or other Objects to skip scanning
	GlobalAnnotationsKey = "disablescan.enforcer.io"
)

// Config is the parent struct to hold
// information about webhook server
type Config struct {
	CertFile   string
	KeyFile    string
	Port       string
	Prune      bool
	Severity   string
	IgnoreFile string
}

// Serve is the wrapper for webhook for validation.
func (c *Config) Serve() (err error) {
	// create a new mux router

	r := mux.NewRouter()
	r.HandleFunc("/validate", c.ValidationHandler)

	addr := "0.0.0.0:" + c.Port
	err = http.ListenAndServeTLS(addr, c.CertFile, c.KeyFile, r)
	if err != nil {
		log.Error(err)
	}
	return err
}

// ValidationHandler manages the validate requests
func (c *Config) ValidationHandler(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		logrus.Errorf("contentType=%s, expect application/json", contentType)
	}

	logrus.Infof("About to process request %s", body)

	// Initialize new decoder
	ds := scheme.Codecs.UniversalDeserializer()
	req := v1beta1.AdmissionReview{}
	_, _, err := ds.Decode(body, nil, &req)

	if err != nil {
		logrus.Errorf("Cannot decode request body: %v", err)

	}

	logrus.Info("Coming here to validate")

	rsp := c.processValidationRequest(req)
	rspBytes, err := json.Marshal(rsp)

	if _, err := w.Write(rspBytes); err != nil {
		logrus.Errorln(err)
	}

}

func (c Config) processValidationRequest(req v1beta1.AdmissionReview) (rsp v1beta1.AdmissionReview) {
	rsp = req

	status, message := c.validateRequest(req)

	rsp.Response = &v1beta1.AdmissionResponse{
		Allowed: status,
		Result: &metav1.Status{
			Message: message,
		},
	}

	return rsp
}

func (c Config) validateRequest(req v1beta1.AdmissionReview) (status bool, message string) {

	switch req.Request.Kind.Kind {
	case "Deployment":
		logrus.Info("Processing Deployment")
		status, message = c.processDeployment(req.Request.Object.Raw)
	case "Pod":
		logrus.Info("Processing Deployment")
		return c.processPod(req.Request.Object.Raw)
	/*case "DaemonSet":
	logrus.Info("Processing Deployment")
	return processDaemonset(req.Request.Object.Raw)*/
	default:
		return false, "unmatched k8s spec"
	}
	return status, message
}

func (c Config) processDeployment(rawObject []byte) (status bool, message string) {
	d := apps.Deployment{}
	err := json.Unmarshal(rawObject, &d)
	if err != nil {
		logrus.Error(err)
		return false, err.Error()
	}

	if checkNameSpace(d.GetNamespace()) {
		logrus.Info("skipping validation due to namespace annotations")
		return true, "Validation skipped due to annotations"
	}

	if parseAnnotations(d.GetAnnotations()) {
		logrus.Info("skipping validation due to workload annotations")
		return true, "Validation skipped due to annotations"
	}
	return c.processPodSpec(d.Spec.Template.Spec)
}

func (c Config) processPod(rawObject []byte) (status bool, message string) {
	p := core.Pod{}
	err := json.Unmarshal(rawObject, &p)
	if err != nil {
		logrus.Error(err)
		return false, err.Error()
	}

	if checkNameSpace(p.GetNamespace()) {
		logrus.Info("skipping validation due to namespace annotations")
		return true, "Validation skipped due to annotations"
	}

	if parseAnnotations(p.GetAnnotations()) {
		logrus.Info("skipping validation due to annotations")
		return true, "Validation skipped due to annotations"
	}

	return c.processPodSpec(p.Spec)
}

func (c Config) processPodSpec(p core.PodSpec) (status bool, message string) {
	for _, p := range p.Containers {
		logrus.Infof("Received Image: %s", p.Image)
		status, message = scanImage(p.Image, c.Severity, c.IgnoreFile)
	}

	return status, message
}

func checkNameSpace(namespace string) bool {
	config, err := rest.InClusterConfig()
	if err != nil {
		logrus.Error(err)
		return false
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Error(err)
		return false
	}

	ns, err := clientset.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})

	if err != nil {
		logrus.Error(err)
		return false
	}

	annotations := ns.GetAnnotations()

	return parseAnnotations(annotations)
}

func parseAnnotations(annotation map[string]string) (status bool) {
	if value, ok := annotation[GlobalAnnotationsKey]; ok {
		status, _ = strconv.ParseBool(value)
		//return status
	}

	return status
}

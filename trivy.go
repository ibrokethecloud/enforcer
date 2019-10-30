package main

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

// scanImage will use Trivy to scan
// the container image and parse the results.
func scanImage(image string, level string) (result bool, message string) {
	if len(image) == 0 {

		return false, "No image specified"
	}
	logrus.Infof("About to pull image: %s\n", image)
	dockerArgs := []string{"pull", "-q", image}
	out, err := exec.Command("docker", dockerArgs...).Output()
	if err != nil {
		logrus.Info(string(out))
		return false, err.Error()
	}
	logrus.Infof("About to scan image: %s\n", image)
	trivyArgs := []string{"--exit-code", "1", "--severity", level, "--quiet", "--format", "json", image}
	out, err = exec.Command("trivy", trivyArgs...).Output()
	if err != nil {
		logrus.Info(string(out))
		return false, err.Error()
	}
	logrus.Infof("No vulnerabilities detected in image - %s", image)
	return true, "Valid spec"
}

// cleanImages will prune images from the pod
func cleanImages() {
	pruneArgs := []string{"system", "prune", "-a"}
	out, err := exec.Command("docker", pruneArgs...).Output()
	if err != nil {
		logrus.Error(out)
		logrus.Error(err)
	}
}

package main

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

// scanImage will use Trivy to scan
// the container image and parse the results.
func scanImage(image string, level string, ignoreFile string) (result bool, message string) {
	if len(image) == 0 {

		return false, "No image specified"
	}

	logrus.Infof("About to scan image: %s\n", image)
	trivyArgs := []string{"--exit-code", "1", "--severity", level, "--quiet", "--format", "json", "--ignorefile", ignoreFile, image}
	out, err := exec.Command("trivy", trivyArgs...).Output()
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

func dbUpdate() {
	logrus.Infof("About to perform scheduled db refresh")
	trivyArgs := []string{"--download-db-only"}
	_, err := exec.Command("trivy", trivyArgs...).Output()
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("Trivy update completed")
}

## Enforcer: A validating webhook used to enforce deployment standards.

A sample validation webhook to enforce deployment standards.

Currently the webhook parses deployments and pod specs, and scans the image's specified.

The webhook uses aquasec/trivy to perform the container image scanning.

Based on level of vulnerabilities to check for, if the container image has no vulnerabilities then the image is allowed to be deployed.

However if the image has vulnerabilities then the webhook denies the request.

Use the chart available in the repo to deploy the webhook with a self signed certificate.

#### To do:
* PVC claims for trivy
* Pre-caching of vulnerabilties at predefined internvals to speed up image scanning.
* Add additional checks for DaemonSets, Jobs...
* Custom annotation to skip scanning for workloads or entire namespaces.
## Enforcer: A validating webhook used to enforce deployment standards.

A sample validation webhook to enforce deployment standards.

Currently the webhook parses deployments and pod specs, and scans the image's specified.

The webhook uses aquasec/trivy to perform the container image scanning.

Based on level of vulnerabilities to check for, if the container image has no vulnerabilities then the image is allowed to be deployed.

However if the image has vulnerabilities then the webhook denies the request.

Use the chart available in the repo to deploy the webhook with a self signed certificate.

### Ignore file
Currently enforcer helm chart creates a config map in the namespace enforcer is deployed to.

This is mounted as a file in the pod, and subsequently enforcer can use the same to pass an ignorefile to trivy.

At present cluster admins can update this configmap to bypass checks for specific CVE's.

```cassandraql
apiVersion: v1
kind: ConfigMap
metadata:
  name: ignorelist
data:
  # Details to be copied into trivy ignore
  trivyignore: |
    #CVE-2018-14618
``` 

### Skip scans
Users can skip scan's at the namespace or workload level by applying the specific annotation to the objects.

`disablescan.enforcer.io: true`

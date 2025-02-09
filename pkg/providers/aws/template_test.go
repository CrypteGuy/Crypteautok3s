package aws

import (
	"testing"
)

const (
	testOutput = `
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cloud-controller-manager
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: cloud-controller-manager:apiserver-authentication-reader
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: extension-apiserver-authentication-reader
subjects:
- apiGroup: ""
  kind: ServiceAccount
  name: cloud-controller-manager
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: system:cloud-controller-manager
rules:
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
  - patch
  - update
- apiGroups:
  - ""
  resources:
  - nodes
  verbs:
  - '*'
- apiGroups:
  - ""
  resources:
  - nodes/status
  verbs:
  - patch
- apiGroups:
  - ""
  resources:
  - services
  verbs:
  - list
  - patch
  - update
  - watch
- apiGroups:
  - ""
  resources:
  - services/status
  verbs:
  - list
  - patch
  - update
  - watch
- apiGroups:
  - ""
  resources:
  - serviceaccounts
  verbs:
  - create
- apiGroups:
  - ""
  resources:
  - persistentvolumes
  verbs:
  - get
  - list
  - update
  - watch
- apiGroups:
  - ""
  resources:
  - endpoints
  verbs:
  - create
  - get
  - list
  - watch
  - update
- apiGroups:
  - coordination.k8s.io
  resources:
  - leases
  verbs:
  - create
  - get
  - list
  - watch
  - update
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: system:cloud-controller-manager
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:cloud-controller-manager
subjects:
- apiGroup: ""
  kind: ServiceAccount
  name: cloud-controller-manager
  namespace: kube-system
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: aws-cloud-controller-manager
  namespace: kube-system
  labels:
    k8s-app: aws-cloud-controller-manager
spec:
  selector:
    matchLabels:
      k8s-app: aws-cloud-controller-manager
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        k8s-app: aws-cloud-controller-manager
    spec:
      nodeSelector:
        node-role.kubernetes.io/master: "true"
      tolerations:
      - key: node.cloudprovider.kubernetes.io/uninitialized
        value: "true"
        effect: NoSchedule
      - key: node-role.kubernetes.io/master
        effect: NoSchedule
      serviceAccountName: cloud-controller-manager
      containers:
        - name: aws-cloud-controller-manager
          image: gcr.io/k8s-staging-provider-aws/cloud-controller-manager:v1.23.1
          args:
            - --v=2
            - --cloud-provider=aws
            - --use-service-account-credentials=true
          resources:
            requests:
              cpu: 200m
      hostNetwork: true
`
)

func TestParseTemplate(t *testing.T) {
	manifest := getAWSCCMManifest("v1.23.9+k3s1")
	if manifest != testOutput {
		t.Fatal("template doesn't match target output")
	}
}

func TestGetCCMVersion(t *testing.T) {
	type testCase struct {
		k3sversion string
		ccmVersion string
	}
	for _, c := range []testCase{
		{
			k3sversion: "v1.20.15+k3s1",
			ccmVersion: ccmVersionMap["~1.20"],
		},
		{
			k3sversion: "v1.24.3+k3s1",
			ccmVersion: ccmVersionMap[">= 1.24"],
		},
	} {
		ccm, err := getCCMVersion(c.k3sversion)
		if err != nil {
			t.Fatalf("failed to get CCM version for k3s version %s", c.k3sversion)
		}
		if ccm != c.ccmVersion {
			t.Fatalf("ccm version is not match, want %s, but got %s", c.ccmVersion, ccm)
		}
	}
}

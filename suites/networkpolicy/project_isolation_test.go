package networkpolicy_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	normantypes "github.com/rancher/norman/types"
	"github.com/rancher/test-network-policy/utils"
	rclusterv3 "github.com/rancher/types/client/cluster/v3"
	rmgmtv3 "github.com/rancher/types/client/management/v3"
	rprojectv3 "github.com/rancher/types/client/project/v3"
	//"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var DefaultTimeout = 60

var _ = Describe("ProjectIsolation", func() {
	var (
		err                                error
		projAlpha, projBravo               *rmgmtv3.Project
		projAlphaClient, projBravoClient   *rprojectv3.Client
		ns1InAlpha, ns2InAlpha             *rclusterv3.Namespace
		ns1InBravo, ns2InBravo             *rclusterv3.Namespace
		w1InNS1ProjAlpha, w2InNS2ProjAlpha *rprojectv3.Workload
		w3InNS1ProjBravo, w4InNS2ProjBravo *rprojectv3.Workload
	)

	BeforeEach(func() {
		By("creating two different projects", func() {
			projAlpha, err = RancherServer.ManagementClient.Project.Create(&rmgmtv3.Project{
				Name:      "proj-alpha",
				ClusterId: RancherServer.DefaultCluster.ID,
			})
			Expect(err).NotTo(HaveOccurred(), "while creating project alpha")

			projBravo, err = RancherServer.ManagementClient.Project.Create(&rmgmtv3.Project{
				Name:      "proj-bravo",
				ClusterId: RancherServer.DefaultCluster.ID,
			})
			Expect(err).NotTo(HaveOccurred(), "while creating project bravo")

			Eventually(func() string {
				if p, err := RancherServer.ManagementClient.Project.ByID(projAlpha.ID); err == nil {
					return p.State
				}
				return ""
			}, DefaultTimeout).Should(Equal("active"), "waiting for project alpha to be active")

			Eventually(func() string {
				if p, err := RancherServer.ManagementClient.Project.ByID(projBravo.ID); err == nil {
					return p.State
				}
				return ""
			}, DefaultTimeout).Should(Equal("active"), "waiting for project bravo to be active")

			projAlphaClient, err = RancherServer.GetProjectClientByID(projAlpha.ID)
			Expect(err).NotTo(HaveOccurred(), "creating client for project alpha")

			projBravoClient, err = RancherServer.GetProjectClientByID(projBravo.ID)
			Expect(err).NotTo(HaveOccurred(), "creating client for project bravo")
		})

		By("creating namespaces in proj-alpha", func() {
			ns1InAlpha, err = RancherServer.DefaultClusterClient.Namespace.Create(&rclusterv3.Namespace{
				Name:      "ns1-in-proj-alpha",
				ProjectID: projAlpha.ID,
			})
			Expect(err).NotTo(HaveOccurred(), "while creating namespace ns1")

			ns2InAlpha, err = RancherServer.DefaultClusterClient.Namespace.Create(&rclusterv3.Namespace{
				Name:      "ns2-in-proj-alpha",
				ProjectID: projAlpha.ID,
			})
			Expect(err).NotTo(HaveOccurred(), "while creating namespace ns2 in project alpha")

			Eventually(func() string {
				if n, err := RancherServer.DefaultClusterClient.Namespace.ByID(ns1InAlpha.ID); err == nil {
					return n.State
				}
				return ""
			}, DefaultTimeout).Should(Equal("active"), "waiting for namespace ns1-in-proj-alpha to become active")

			Eventually(func() string {
				if n, err := RancherServer.DefaultClusterClient.Namespace.ByID(ns2InAlpha.ID); err == nil {
					return n.State
				}
				return ""
			}, DefaultTimeout).Should(Equal("active"), "waiting for namespace ns2-in-proj-alpha to become active")
		})

		By("creating namespaces in proj-bravo", func() {
			ns1InBravo, err = RancherServer.DefaultClusterClient.Namespace.Create(&rclusterv3.Namespace{
				Name:      "ns1-in-proj-bravo",
				ProjectID: projBravo.ID,
			})
			Expect(err).NotTo(HaveOccurred(), "while creating namespace ns1")

			ns2InBravo, err = RancherServer.DefaultClusterClient.Namespace.Create(&rclusterv3.Namespace{
				Name:      "ns2-in-proj-bravo",
				ProjectID: projBravo.ID,
			})
			Expect(err).NotTo(HaveOccurred(), "while creating namespace ns2 in project bravo")

			Eventually(func() string {
				if n, err := RancherServer.DefaultClusterClient.Namespace.ByID(ns1InBravo.ID); err == nil {
					return n.State
				}
				return ""
			}, DefaultTimeout).Should(Equal("active"), "waiting for namespace ns1-in-proj-bravo to become active")

			Eventually(func() string {
				if n, err := RancherServer.DefaultClusterClient.Namespace.ByID(ns2InBravo.ID); err == nil {
					return n.State
				}
				return ""
			}, DefaultTimeout).Should(Equal("active"), "waiting for namespace ns2-in-proj-bravo to become active")
		})

		By("creating workload in ns1-in-proj-alpha", func() {
			w1InNS1ProjAlpha = &rprojectv3.Workload{
				Name:        "workload-in-ns1-in-proj-alpha",
				NamespaceId: ns1InAlpha.Name,
				DNSPolicy:   "ClusterFirst",
				Containers: []rprojectv3.Container{
					{
						Name:  "workload-in-ns1-in-proj-alpha",
						Image: "leodotcloud/swiss-army-knife",
						Stdin: true,
						TTY:   true,
					},
				},
				DeploymentConfig: &rprojectv3.DeploymentConfig{
					MaxSurge:             intstr.FromInt(1),
					Strategy:             "RollingUpdate",
					RevisionHistoryLimit: func(i int64) *int64 { return &i }(10),
				},
			}
			w1InNS1ProjAlpha, err = projAlphaClient.Workload.Create(w1InNS1ProjAlpha)
			Expect(err).NotTo(HaveOccurred(), "while creating workload w1InNS1ProjAlpha")
		})

		By("creating workload in ns2-in-proj-alpha", func() {
			w2InNS2ProjAlpha = &rprojectv3.Workload{
				Name:        "workload-in-ns2-in-proj-alpha",
				NamespaceId: ns2InAlpha.Name,
				DNSPolicy:   "ClusterFirst",
				Containers: []rprojectv3.Container{
					{
						Name:  "workload-in-ns2-in-proj-alpha",
						Image: "leodotcloud/swiss-army-knife",
						Stdin: true,
						TTY:   true,
					},
				},
				DeploymentConfig: &rprojectv3.DeploymentConfig{
					MaxSurge:             intstr.FromInt(1),
					Strategy:             "RollingUpdate",
					RevisionHistoryLimit: func(i int64) *int64 { return &i }(10),
				},
			}
			w2InNS2ProjAlpha, err = projAlphaClient.Workload.Create(w2InNS2ProjAlpha)
			Expect(err).NotTo(HaveOccurred(), "while creating workload w2InNS2ProjAlpha")
		})

		By("creating workload in ns1-in-proj-bravo", func() {
			w3InNS1ProjBravo = &rprojectv3.Workload{
				Name:        "workload-in-ns1-in-proj-bravo",
				NamespaceId: ns1InBravo.Name,
				DNSPolicy:   "ClusterFirst",
				Containers: []rprojectv3.Container{
					{
						Name:  "workload-in-ns1-in-proj-bravo",
						Image: "leodotcloud/swiss-army-knife",
						Stdin: true,
						TTY:   true,
					},
				},
				DeploymentConfig: &rprojectv3.DeploymentConfig{
					MaxSurge:             intstr.FromInt(1),
					Strategy:             "RollingUpdate",
					RevisionHistoryLimit: func(i int64) *int64 { return &i }(10),
				},
			}
			w3InNS1ProjBravo, err = projBravoClient.Workload.Create(w3InNS1ProjBravo)
			Expect(err).NotTo(HaveOccurred(), "while creating workload w3InNS1ProjBravo")
		})

		By("creating workload in ns2-in-proj-bravo", func() {
			w4InNS2ProjBravo = &rprojectv3.Workload{
				Name:        "workload-in-ns2-in-proj-bravo",
				NamespaceId: ns2InBravo.Name,
				DNSPolicy:   "ClusterFirst",
				Containers: []rprojectv3.Container{
					{
						Name:  "workload-in-ns2-in-proj-bravo",
						Image: "leodotcloud/swiss-army-knife",
						Stdin: true,
						TTY:   true,
					},
				},
				DeploymentConfig: &rprojectv3.DeploymentConfig{
					MaxSurge:             intstr.FromInt(1),
					Strategy:             "RollingUpdate",
					RevisionHistoryLimit: func(i int64) *int64 { return &i }(10),
				},
			}
			w4InNS2ProjBravo, err = projBravoClient.Workload.Create(w4InNS2ProjBravo)
			Expect(err).NotTo(HaveOccurred(), "while creating workload w4InNS2ProjBravo")
		})

		Eventually(func() string {
			if w, err := projAlphaClient.Workload.ByID(w1InNS1ProjAlpha.ID); err == nil {
				return w.State
			}
			return ""
		}, DefaultTimeout).Should(Equal("active"), "waiting for workload to become active")

		Eventually(func() string {
			if w, err := projAlphaClient.Workload.ByID(w2InNS2ProjAlpha.ID); err == nil {
				return w.State
			}
			return ""
		}, DefaultTimeout).Should(Equal("active"), "waiting for workload to become active")

		Eventually(func() string {
			if w, err := projBravoClient.Workload.ByID(w3InNS1ProjBravo.ID); err == nil {
				return w.State
			}
			return ""
		}, DefaultTimeout).Should(Equal("active"), "waiting for workload to become active")

		Eventually(func() string {
			if w, err := projBravoClient.Workload.ByID(w4InNS2ProjBravo.ID); err == nil {
				return w.State
			}
			return ""
		}, DefaultTimeout).Should(Equal("active"), "waiting for workload to become active")

	})

	AfterEach(func() {
		By("deleting projects", func() {
			err = RancherServer.ManagementClient.Project.Delete(projAlpha)
			Expect(err).NotTo(HaveOccurred(), "while deleting project alpha")
			err = RancherServer.ManagementClient.Project.Delete(projBravo)
			Expect(err).NotTo(HaveOccurred(), "while deleting project bravo")
		})
	})

	It("projects should be isolated", func() {
		// pods
		w1PodCollection, err := projAlphaClient.Pod.List(&normantypes.ListOpts{
			Filters: map[string]interface{}{
				"workloadId": w1InNS1ProjAlpha.ID,
			},
		})
		Expect(err).NotTo(HaveOccurred(), "while fetching pods of workload w1")
		Expect(len(w1PodCollection.Data)).To(Equal(1), "expect only one pod in workload w1")
		Expect(len(w1PodCollection.Data[0].Containers)).To(Equal(1), "expect only one container in pod of workload w1")

		w2PodCollection, err := projAlphaClient.Pod.List(&normantypes.ListOpts{
			Filters: map[string]interface{}{
				"workloadId": w2InNS2ProjAlpha.ID,
			},
		})
		Expect(err).NotTo(HaveOccurred(), "while fetching pods of workload w2")
		Expect(len(w2PodCollection.Data)).To(Equal(1), "expect only one pod in workload w2")
		Expect(len(w2PodCollection.Data[0].Containers)).To(Equal(1), "expect only one container in pod of workload w2")

		w3PodCollection, err := projBravoClient.Pod.List(&normantypes.ListOpts{
			Filters: map[string]interface{}{
				"workloadId": w3InNS1ProjBravo.ID,
			},
		})
		Expect(err).NotTo(HaveOccurred(), "while fetching pods of workload w3")
		Expect(len(w3PodCollection.Data)).To(Equal(1), "expect only one pod in workload w3")
		Expect(len(w3PodCollection.Data[0].Containers)).To(Equal(1), "expect only one container in pod of workload w3")

		w4PodCollection, err := projBravoClient.Pod.List(&normantypes.ListOpts{
			Filters: map[string]interface{}{
				"workloadId": w4InNS2ProjBravo.ID,
			},
		})
		Expect(err).NotTo(HaveOccurred(), "while fetching pods of workload w4")
		Expect(len(w4PodCollection.Data)).To(Equal(1), "expect only one pod in workload w4")
		Expect(len(w4PodCollection.Data[0].Containers)).To(Equal(1), "expect only one container in pod of workload w4")

		w1Pod := w1PodCollection.Data[0]
		w2Pod := w2PodCollection.Data[0]
		w3Pod := w3PodCollection.Data[0]
		w4Pod := w4PodCollection.Data[0]

		var output, curlCommand, wsURL string

		// w2 -> w1 should succeed
		curlCommand = "curl --max-time 5 -s http://" + w1InNS1ProjAlpha.Name + "." + w1InNS1ProjAlpha.NamespaceId
		wsURL = utils.GetWSURL(RancherServer.URL, RancherServer.DefaultCluster.ID, w2Pod.NamespaceId, w2Pod.Name, w2Pod.Containers[0].Name, curlCommand)

		output, err = utils.RunExecCommand(wsURL, RancherServer.AccessKey, RancherServer.SecretKey, RancherServer.TokenKey)
		Expect(err).NotTo(HaveOccurred(), "while running command")
		Expect(output).Should(ContainSubstring(w1Pod.Name))

		// w4 -> w3 should succeed
		curlCommand = "curl --max-time 5 -s http://" + w3InNS1ProjBravo.Name + "." + w3InNS1ProjBravo.NamespaceId
		wsURL = utils.GetWSURL(RancherServer.URL, RancherServer.DefaultCluster.ID, w4Pod.NamespaceId, w4Pod.Name, w4Pod.Containers[0].Name, curlCommand)

		output, err = utils.RunExecCommand(wsURL, RancherServer.AccessKey, RancherServer.SecretKey, RancherServer.TokenKey)
		Expect(err).NotTo(HaveOccurred(), "while running command")
		Expect(output).Should(ContainSubstring(w3Pod.Name))

		// w4 -> w2 should fail
		curlCommand = "curl --max-time 5 -s http://" + w2InNS2ProjAlpha.Name + "." + w2InNS2ProjAlpha.NamespaceId
		wsURL = utils.GetWSURL(RancherServer.URL, RancherServer.DefaultCluster.ID, w4Pod.NamespaceId, w4Pod.Name, w4Pod.Containers[0].Name, curlCommand)

		output, err = utils.RunExecCommand(wsURL, RancherServer.AccessKey, RancherServer.SecretKey, RancherServer.TokenKey)
		Expect(err).NotTo(HaveOccurred(), "while running command")
		Expect(output).Should(ContainSubstring("non-zero exit code"))
		Expect(output).ShouldNot(ContainSubstring(w2Pod.Name))

		// w1 -> w3 should fail
		curlCommand = "curl --max-time 5 -s http://" + w3InNS1ProjBravo.Name + "." + w3InNS1ProjBravo.NamespaceId
		wsURL = utils.GetWSURL(RancherServer.URL, RancherServer.DefaultCluster.ID, w1Pod.NamespaceId, w1Pod.Name, w1Pod.Containers[0].Name, curlCommand)

		output, err = utils.RunExecCommand(wsURL, RancherServer.AccessKey, RancherServer.SecretKey, RancherServer.TokenKey)
		Expect(err).NotTo(HaveOccurred(), "while running command")
		Expect(output).Should(ContainSubstring("non-zero exit code"))
		Expect(output).ShouldNot(ContainSubstring(w3Pod.Name))

	})
})

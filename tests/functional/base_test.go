/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package functional

import (
	"fmt"

	maps "golang.org/x/exp/maps"

	. "github.com/onsi/gomega" //revive:disable:dot-imports

	corev1 "k8s.io/api/core/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	barbicanv1 "github.com/openstack-k8s-operators/barbican-operator/api/v1beta1"
	condition "github.com/openstack-k8s-operators/lib-common/modules/common/condition"
)

func CreateBarbicanSecret(namespace string, name string) *corev1.Secret {
	return th.CreateSecret(
		types.NamespacedName{Namespace: namespace, Name: name},
		map[string][]byte{
			"AdminPassword":            []byte("12345678"),
			"BarbicanPassword":         []byte("12345678"),
			"KeystoneDatabasePassword": []byte("12345678"),
			"BarbicanSimpleCryptoKEK":  []byte("sEFmdFjDUqRM2VemYslV5yGNWjokioJXsg8Nrlc3drU="),
		},
	)
}

func CreateCustomConfigSecret(namespace string, name string, contents map[string][]byte) *corev1.Secret {
	return th.CreateSecret(
		types.NamespacedName{Namespace: namespace, Name: name},
		contents,
	)
}

func GetDefaultBarbicanSpec() map[string]interface{} {
	return map[string]interface{}{
		"databaseInstance":          "openstack",
		"secret":                    SecretName,
		"simpleCryptoBackendSecret": SecretName,
		"customServiceConfig":       barbicanTest.BaseCustomServiceConfig,
	}
}

func CreateBarbican(name types.NamespacedName, spec map[string]interface{}) client.Object {
	raw := map[string]interface{}{
		"apiVersion": "barbican.openstack.org/v1beta1",
		"kind":       "Barbican",
		"metadata": map[string]interface{}{
			"name":      name.Name,
			"namespace": name.Namespace,
		},
		"spec": spec,
	}
	return th.CreateUnstructured(raw)
}

func GetBarbican(name types.NamespacedName) *barbicanv1.Barbican {
	instance := &barbicanv1.Barbican{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(ctx, name, instance)).Should(Succeed())
	}, timeout, interval).Should(Succeed())
	return instance
}

func CreateBarbicanMessageBusSecret(namespace string, name string) *corev1.Secret {
	s := th.CreateSecret(
		types.NamespacedName{Namespace: namespace, Name: name},
		map[string][]byte{
			"transport_url": []byte(fmt.Sprintf("rabbit://%s/fake", name)),
		},
	)
	logger.Info("Secret created", "name", name)
	return s
}

func BarbicanConditionGetter(name types.NamespacedName) condition.Conditions {
	instance := GetBarbican(name)
	return instance.Status.Conditions
}

func BarbicanAPINotExists(name types.NamespacedName) {
	Consistently(func(g Gomega) {
		instance := &barbicanv1.BarbicanAPI{}
		err := k8sClient.Get(ctx, name, instance)
		g.Expect(k8s_errors.IsNotFound(err)).To(BeTrue())
	}, timeout, interval).Should(Succeed())
}

func BarbicanWorkerNotExists(name types.NamespacedName) {
	Consistently(func(g Gomega) {
		instance := &barbicanv1.BarbicanWorker{}
		err := k8sClient.Get(ctx, name, instance)
		g.Expect(k8s_errors.IsNotFound(err)).To(BeTrue())
	}, timeout, interval).Should(Succeed())
}

func BarbicanKeystoneListenerNotExists(name types.NamespacedName) {
	Consistently(func(g Gomega) {
		instance := &barbicanv1.BarbicanKeystoneListener{}
		err := k8sClient.Get(ctx, name, instance)
		g.Expect(k8s_errors.IsNotFound(err)).To(BeTrue())
	}, timeout, interval).Should(Succeed())
}

func BarbicanExists(name types.NamespacedName) {
	Consistently(func(g Gomega) {
		instance := &barbicanv1.Barbican{}
		err := k8sClient.Get(ctx, name, instance)
		g.Expect(k8s_errors.IsNotFound(err)).To(BeFalse())
	}, timeout, interval).Should(Succeed())
}

func BarbicanAPIConditionGetter(name types.NamespacedName) condition.Conditions {
	instance := GetBarbicanAPI(name)
	return instance.Status.Conditions
}

func BarbicanAPIExists(name types.NamespacedName) {
	Consistently(func(g Gomega) {
		instance := &barbicanv1.BarbicanAPI{}
		err := k8sClient.Get(ctx, name, instance)
		g.Expect(k8s_errors.IsNotFound(err)).To(BeFalse())
	}, timeout, interval).Should(Succeed())
}

// GetBarbicanAPI - Returns BarbicanAPI subCR
func GetBarbicanAPI(name types.NamespacedName) *barbicanv1.BarbicanAPI {
	instance := &barbicanv1.BarbicanAPI{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(ctx, name, instance)).Should(Succeed())
	}, timeout, interval).Should(Succeed())
	return instance
}

// GetBarbicanKeystoneListener - Returns BarbicanKeystoneListener subCR
func GetBarbicanKeystoneListener(name types.NamespacedName) *barbicanv1.BarbicanKeystoneListener {
	instance := &barbicanv1.BarbicanKeystoneListener{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(ctx, name, instance)).Should(Succeed())
	}, timeout, interval).Should(Succeed())
	return instance
}

func BarbicanKeystoneListenerConditionGetter(name types.NamespacedName) condition.Conditions {
	instance := GetBarbicanKeystoneListener(name)
	return instance.Status.Conditions
}

// GetBarbicanWorker - Returns BarbicanWorker subCR
func GetBarbicanWorker(name types.NamespacedName) *barbicanv1.BarbicanWorker {
	instance := &barbicanv1.BarbicanWorker{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(ctx, name, instance)).Should(Succeed())
	}, timeout, interval).Should(Succeed())
	return instance
}

func BarbicanWorkerConditionGetter(name types.NamespacedName) condition.Conditions {
	instance := GetBarbicanWorker(name)
	return instance.Status.Conditions
}

// ========== TLS Stuff ==============
func GetTLSBarbicanSpec() map[string]interface{} {
	return map[string]interface{}{
		"databaseInstance":          "openstack",
		"secret":                    SecretName,
		"simpleCryptoBackendSecret": SecretName,
		"barbicanAPI":               GetTLSBarbicanAPISpec(),
		"customServiceConfig":       barbicanTest.BaseCustomServiceConfig,
		"defaultConfigOverwrite":    barbicanTest.BaseDefaultConfigOverwrite,
	}
}

func GetTLSBarbicanAPISpec() map[string]interface{} {
	spec := GetDefaultBarbicanAPISpec()
	maps.Copy(spec, map[string]interface{}{
		"tls": map[string]interface{}{
			"api": map[string]interface{}{
				"internal": map[string]interface{}{
					"secretName": InternalCertSecretName,
				},
				"public": map[string]interface{}{
					"secretName": PublicCertSecretName,
				},
			},
			"caBundleSecretName": CABundleSecretName,
		},
		"customServiceConfigSecrets": barbicanTest.APICustomServiceConfigSecrets,
	})
	return spec
}

// ========== End of TLS Stuff ============

// ========== PKCS11 Stuff ============
const PKCS11CustomData = `
[p11_crypto_plugin]
plugin_name = PKCS11
library_path = some_library_path
token_labels = some_partition_label
mkek_label = my_mkek_label
hmac_label = my_hmac_label
encryption_mechanism = CKM_AES_GCM
aes_gcm_generate_iv = true
hmac_key_type = CKK_GENERIC_SECRET
hmac_keygen_mechanism = CKM_GENERIC_SECRET_KEY_GEN
hmac_mechanism = CKM_SHA256_HMAC
key_wrap_mechanism = CKM_AES_CBC_PAD
key_wrap_generate_iv = true
always_set_cka_sensitive = true
os_locking_ok = false`

func GetPKCS11BarbicanSpec(hsmModel string) map[string]interface{} {
	spec := GetDefaultBarbicanSpec()
	maps.Copy(spec, map[string]interface{}{
		"customServiceConfig":      PKCS11CustomData,
		"enabledSecretStores":      []string{"pkcs11"},
		"globalDefaultSecretStore": "pkcs11",
		"pkcs11": map[string]interface{}{
			"clientDataPath":   PKCS11ClientDataPath[hsmModel],
			"loginSecret":      PKCS11LoginSecret,
			"clientDataSecret": PKCS11ClientDataSecret,
		},
	})
	return spec
}

func GetPKCS11BarbicanAPISpec(hsmModel string) map[string]interface{} {
	spec := GetPKCS11BarbicanSpec(hsmModel)
	maps.Copy(spec, GetDefaultBarbicanAPISpec())
	return spec
}

func CreatePKCS11LoginSecret(namespace string, name string) *corev1.Secret {
	return th.CreateSecret(
		types.NamespacedName{Namespace: namespace, Name: name},
		map[string][]byte{
			"PKCS11Pin": []byte("12345678"),
		},
	)
}

func CreatePKCS11ClientDataSecret(namespace string, name string, hsmModel string) *corev1.Secret {
	secretContents := make(map[string][]byte)
	if hsmModel == "luna" {
		secretContents = map[string][]byte{
			"Client.cfg": []byte("dummy-data"),
			"CACert.pem": []byte("dummy-data"),
			"Server.pem": []byte("dummy-data"),
			"Client.pem": []byte("dummy-data"),
			"Client.key": []byte("dummy-data"),
		}
	} else if hsmModel == "proteccio" {
		secretContents = map[string][]byte{
			"proteccio.rc": []byte("dummy-data"),
			"Server.CRT":   []byte("dummy-data"),
			"Client.crt":   []byte("dummy-data"),
			"Client.key":   []byte("dummy-data"),
		}
	}
	return th.CreateSecret(
		types.NamespacedName{Namespace: namespace, Name: name},
		secretContents,
	)
}

// ========== End of PKCS11 Stuff ============

func GetDefaultBarbicanAPISpec() map[string]interface{} {
	return map[string]interface{}{
		"secret":                    SecretName,
		"simpleCryptoBackendSecret": SecretName,
		"replicas":                  1,
		"databaseHostname":          barbicanTest.DatabaseHostname,
		"databaseInstance":          barbicanTest.DatabaseInstance,
		"containerImage":            barbicanTest.ContainerImage,
		"serviceAccount":            barbicanTest.BarbicanSA.Name,
		"transportURLSecret":        barbicanTest.RabbitmqSecretName,
		"customServiceConfig":       barbicanTest.APICustomServiceConfig,
		"defaultConfigOverwrite":    barbicanTest.APIDefaultConfigOverwrite,
	}
}

func CreateBarbicanAPI(name types.NamespacedName, spec map[string]interface{}) client.Object {
	// we get the parent CR and set ownership to the barbicanAPI CR
	raw := map[string]interface{}{
		"apiVersion": "barbican.openstack.org/v1beta1",
		"kind":       "BarbicanAPI",
		"metadata": map[string]interface{}{
			"annotations": map[string]interface{}{
				"keystoneEndpoint": "true",
			},
			"name":      name.Name,
			"namespace": name.Namespace,
		},
		"spec": spec,
	}

	return th.CreateUnstructured(raw)
}

// GetBarbicanAPISpec -
func GetBarbicanAPISpec(name types.NamespacedName) barbicanv1.BarbicanAPITemplate {
	instance := &barbicanv1.BarbicanAPI{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(ctx, name, instance)).Should(Succeed())
	}, timeout, interval).Should(Succeed())
	return instance.Spec.BarbicanAPITemplate
}

// GetBarbicanKeystoneListenerSpec -
func GetBarbicanKeystoneListenerSpec(name types.NamespacedName) barbicanv1.BarbicanKeystoneListenerTemplate {
	instance := &barbicanv1.BarbicanKeystoneListener{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(ctx, name, instance)).Should(Succeed())
	}, timeout, interval).Should(Succeed())
	return instance.Spec.BarbicanKeystoneListenerTemplate
}

// GetBarbicanWorkerSpec -
func GetBarbicanWorkerSpec(name types.NamespacedName) barbicanv1.BarbicanWorkerTemplate {
	instance := &barbicanv1.BarbicanWorker{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(ctx, name, instance)).Should(Succeed())
	}, timeout, interval).Should(Succeed())
	return instance.Spec.BarbicanWorkerTemplate
}

// GetSampleTopologySpec - A sample (and opinionated) Topology Spec used to
// test Barbican
// Note this is just an example that should not be used in production for
// multiple reasons:
// 1. It uses ScheduleAnyway as strategy, which is something we might
// want to avoid by default
// 2. Usually a topologySpreadConstraints is used to take care about
// multi AZ, which is not applicable in this context
func GetSampleTopologySpec(
	label string,
) (map[string]interface{}, []corev1.TopologySpreadConstraint) {
	// Build the topology Spec
	topologySpec := map[string]interface{}{
		"topologySpreadConstraints": []map[string]interface{}{
			{
				"maxSkew":           1,
				"topologyKey":       corev1.LabelHostname,
				"whenUnsatisfiable": "ScheduleAnyway",
				"labelSelector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"component": label,
					},
				},
			},
		},
	}
	// Build the topologyObj representation
	topologySpecObj := []corev1.TopologySpreadConstraint{
		{
			MaxSkew:           1,
			TopologyKey:       corev1.LabelHostname,
			WhenUnsatisfiable: corev1.ScheduleAnyway,
			LabelSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"component": label,
				},
			},
		},
	}
	return topologySpec, topologySpecObj
}

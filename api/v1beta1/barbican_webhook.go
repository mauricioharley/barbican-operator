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

//
// Generated by:
//
// operator-sdk create webhook --group barbican --version v1beta1 --kind Barbican --programmatic-validation --defaulting
//

package v1beta1

import (
	"fmt"
	"slices"

	"github.com/openstack-k8s-operators/lib-common/modules/common/service"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// BarbicanDefaults -
type BarbicanDefaults struct {
	APIContainerImageURL              string
	WorkerContainerImageURL           string
	KeystoneListenerContainerImageURL string
	BarbicanAPITimeout                int
}

var barbicanDefaults BarbicanDefaults

// log is for logging in this package.
var barbicanlog = logf.Log.WithName("barbican-resource")

// SetupBarbicanDefaults - initialize Barbican spec defaults for use with either internal or external webhooks
func SetupBarbicanDefaults(defaults BarbicanDefaults) {
	barbicanDefaults = defaults
	barbicanlog.Info("Barbican defaults initialized", "defaults", defaults)
}

// SetupWebhookWithManager sets up the webhook with the Manager
func (r *Barbican) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-barbican-openstack-org-v1beta1-barbican,mutating=true,failurePolicy=fail,sideEffects=None,groups=barbican.openstack.org,resources=barbicans,verbs=create;update,versions=v1beta1,name=mbarbican.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Barbican{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Barbican) Default() {
	barbicanlog.Info("default", "name", r.Name)

	r.Spec.Default()
}

// Default - set defaults for this Barbican spec
func (spec *BarbicanSpec) Default() {
	if spec.BarbicanAPI.ContainerImage == "" {
		spec.BarbicanAPI.ContainerImage = barbicanDefaults.APIContainerImageURL
	}

	if spec.BarbicanWorker.ContainerImage == "" {
		spec.BarbicanWorker.ContainerImage = barbicanDefaults.WorkerContainerImageURL
	}

	if spec.BarbicanKeystoneListener.ContainerImage == "" {
		spec.BarbicanKeystoneListener.ContainerImage = barbicanDefaults.KeystoneListenerContainerImageURL
	}
	spec.BarbicanSpecBase.Default()
}

// Default - for shared base validations
func (spec *BarbicanSpecBase) Default() {
	// no validations
}

// Default - set defaults for this BarbicanSpecBase. NOTE: this version is used by the OpenStackControlplane webhook
func (spec *BarbicanSpecCore) Default() {
	if spec.APITimeout == 0 {
		spec.APITimeout = barbicanDefaults.BarbicanAPITimeout
	}
	spec.BarbicanSpecBase.Default()
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-barbican-openstack-org-v1beta1-barbican,mutating=false,failurePolicy=fail,sideEffects=None,groups=barbican.openstack.org,resources=barbicans,verbs=create;update,versions=v1beta1,name=vbarbican.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Barbican{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Barbican) ValidateCreate() (admission.Warnings, error) {
	barbicanlog.Info("validate create", "name", r.Name)

	var allErrs field.ErrorList
	basePath := field.NewPath("spec")
	if err := r.Spec.ValidateCreate(basePath); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) != 0 {
		return nil, apierrors.NewInvalid(
			schema.GroupKind{Group: "barbican.openstack.org", Kind: "Barbican"},
			r.Name, allErrs)
	}

	return nil, nil
}

// ValidateCreate - Exported function wrapping non-exported validate functions,
// this function can be called externally to validate an barbican spec.
func (r *BarbicanSpec) ValidateCreate(basePath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	// validate the service override key is valid
	allErrs = append(allErrs, service.ValidateRoutedOverrides(
		basePath.Child("barbicanAPI").Child("override").Child("service"),
		r.BarbicanAPI.Override.Service)...)

	// pkcs11 verifications
	r.ValidatePKCS11(basePath, &allErrs)

	return allErrs
}

func (r *BarbicanSpec) ValidatePKCS11(basePath *field.Path, allErrs *field.ErrorList) {
	if slices.Contains(r.EnabledSecretStores, SecretStorePKCS11) {
                if r.PKCS11 == nil {
                        *allErrs = append(*allErrs, field.Required(basePath.Child("PKCS11"),
                                "PKCS11 specification is missing, PKCS11 is required when pkcs11 is an enabled SecretStore"),
                        )
                }
        }
}

func (r *BarbicanSpecCore) ValidateCreate(basePath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	// validate the service override key is valid
	allErrs = append(allErrs, service.ValidateRoutedOverrides(
		basePath.Child("barbicanAPI").Child("override").Child("service"),
		r.BarbicanAPI.Override.Service)...)

	return allErrs
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Barbican) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	barbicanlog.Info("validate update", "name", r.Name)

	oldBarbican, ok := old.(*Barbican)
	if !ok || oldBarbican == nil {
		return nil, apierrors.NewInternalError(fmt.Errorf("unable to convert existing object"))
	}

	var allErrs field.ErrorList
	basePath := field.NewPath("spec")

	if err := r.Spec.ValidateUpdate(oldBarbican.Spec, basePath); err != nil {
		allErrs = append(allErrs, err...)
	}

	if len(allErrs) != 0 {
		return nil, apierrors.NewInvalid(
			schema.GroupKind{Group: "barbican.openstack.org", Kind: "Barbican"},
			r.Name, allErrs)
	}

	return nil, nil
}

// ValidateUpdate - Exported function wrapping non-exported validate functions,
// this function can be called externally to validate an barbican spec.
func (r *BarbicanSpec) ValidateUpdate(old BarbicanSpec, basePath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	// validate the service override key is valid
	allErrs = append(allErrs, service.ValidateRoutedOverrides(
		basePath.Child("barbicanAPI").Child("override").Child("service"),
		r.BarbicanAPI.Override.Service)...)

	// pkcs11 verifications
	r.ValidatePKCS11(basePath, &allErrs)

	return allErrs
}

func (r *BarbicanSpecCore) ValidateUpdate(old BarbicanSpecCore, basePath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	// validate the service override key is valid
	allErrs = append(allErrs, service.ValidateRoutedOverrides(
		basePath.Child("barbicanAPI").Child("override").Child("service"),
		r.BarbicanAPI.Override.Service)...)

	return allErrs
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Barbican) ValidateDelete() (admission.Warnings, error) {
	barbicanlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

func (spec *BarbicanSpecCore) GetDefaultRouteAnnotations() (annotations map[string]string) {
	return map[string]string{
		"haproxy.router.openshift.io/timeout": fmt.Sprintf("%ds", barbicanDefaults.BarbicanAPITimeout),
	}
}

// SetDefaultRouteAnnotations sets HAProxy timeout values for Barbican API routes
func (spec *BarbicanSpecCore) SetDefaultRouteAnnotations(annotations map[string]string) {
	const haProxyAnno = "haproxy.router.openshift.io/timeout"
	// Use a custom annotation to flag when the operator has set the default HAProxy timeout
	// With the annotation func determines when to overwrite existing HAProxy timeout with the APITimeout
	const barbicanAnno = "api.barbican.openstack.org/timeout"
	valBarbicanAPI, okBarbicanAPI := annotations[barbicanAnno]
	valHAProxy, okHAProxy := annotations[haProxyAnno]

	// Human operator set the HAProxy timeout manually
	if !okBarbicanAPI && okHAProxy {
		return
	}
	// Human operator modified the HAProxy timeout manually without removing the Barbican flag
	if okBarbicanAPI && okHAProxy && valBarbicanAPI != valHAProxy {
		delete(annotations, barbicanAnno)
		return
	}

	timeout := fmt.Sprintf("%ds", spec.APITimeout)
	annotations[barbicanAnno] = timeout
	annotations[haProxyAnno] = timeout
}

/*
Copyright 2021.

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

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cachev1alpha1 "github.com/adamg/contract-operator/api/v1alpha1"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ContractReconciler reconciles a Contract object
type ContractReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cache.adamg.com,resources=contracts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.adamg.com,resources=contracts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.adamg.com,resources=contracts/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Contract object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *ContractReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the Contract instance
	contract := &cachev1alpha1.Contract{}

	err := r.Get(ctx, req.NamespacedName, contract)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Contract resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Contract")
		return ctrl.Result{}, err
	}

	// Fetch Title from CR spec
	pagetype := contract.Spec.PageType
	pagetitle := contract.Spec.PageTitle
	pagespace := contract.Spec.PageSpace
	pagebody := contract.Spec.PageBody
	fmt.Println("type :", pagetype)
	fmt.Println("title :", pagetitle)
	fmt.Println("space :", pagespace)
	fmt.Println("body :", pagebody)

	// your logic here
	httpposturl := "https://cryptoadglad.atlassian.net/wiki/rest/api/content/"
	fmt.Println("HTTP JSON POST URL:", httpposturl)

	jsonString := fmt.Sprintf(
		`{ 
                          "type": "`+"%s"+`",
			  "title": "`+"%s"+`", 
                          "space": {
                                "key": "`+"%s"+`"`+
			`},
                        "body": { "storage": {
                                       "value": "<p>`+"%s"+`</p>", 
                                 "representation": "storage" 
                                } 
                        } 
                }`, pagetype, pagetitle, pagespace, pagebody)

	fmt.Println("jsonString:", jsonString)

	var jsonData = []byte(jsonString)
	//var jsonData = []byte(`
	//{
	//        "type": "page",
	//        "title": "title"
	//        "space": {
	//                "key": "CRYPTO"
	//        },
	//        "body": { "storage": {
	//                       "value": "<p>This is a new page</p>",
	//                 "representation": "storage"
	//                }
	//        }
	//}`)

	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	request.SetBasicAuth("cryptoadglad@gmail.com", "v8yTzH6PzBsHEBHbEndw9EF3")
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ContractReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.Contract{}).
		Complete(r)
}

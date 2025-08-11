package fclient_test

import (
	"context"
	"fmt"

	ET "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
	RIOE "github.com/IBM/fp-go/readerioeither"
	"github.com/appthrust/fcr/pkg/fclient"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func ExampleGet() {
	resultToStr := func(result ET.Either[error, *corev1.ConfigMap]) string {
		return ET.Fold(
			func(err error) string { return fmt.Sprintf("Left(%T)", err) },
			func(cm *corev1.ConfigMap) string { return fmt.Sprintf("Right(%T)", cm) },
		)(result)
	}

	// Setup client
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	cl := fake.NewClientBuilder().
		WithScheme(scheme).
		// Emulate that the API has a configmap named "exists"
		WithObjects(&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "exists", Namespace: "default"}}).
		Build()

	// Setup environment for reader monad
	env := fclient.Env{Ctx: context.TODO(), Client: cl}

	// Example 1: Found → Right(*corev1.ConfigMap)
	params1 := fclient.ToGetParams(client.ObjectKey{Namespace: "default", Name: "exists"})
	result1 := fclient.Get[corev1.ConfigMap](params1)(env)()
	fmt.Printf("Example 1: %s\n", resultToStr(result1))

	// Example 2: NotFound → Left(*errors.StatusError)
	params2 := fclient.ToGetParams(client.ObjectKey{Namespace: "default", Name: "missing"})
	result2 := fclient.Get[corev1.ConfigMap](params2)(env)()
	fmt.Printf("Example 2: %s\n", resultToStr(result2))

	// Output:
	// Example 1: Right(*v1.ConfigMap)
	// Example 2: Left(*errors.StatusError)
}

func ExampleGetOption() {
	resultToStr := func(result ET.Either[error, O.Option[*corev1.ConfigMap]]) string {
		return ET.Fold(
			func(err error) string { return fmt.Sprintf("Left(%T)", err) },
			func(opt O.Option[*corev1.ConfigMap]) string {
				return O.Fold(
					func() string { return "Right(None)" },
					func(cm *corev1.ConfigMap) string { return fmt.Sprintf("Right(Some(%T))", cm) },
				)(opt)
			},
		)(result)
	}

	// Setup client
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	cl := fake.NewClientBuilder().
		WithScheme(scheme).
		// Emulate that the API has a configmap named "exists"
		WithObjects(&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "exists", Namespace: "default"}}).
		Build()

	// Setup environment for reader monad
	env := fclient.Env{Ctx: context.TODO(), Client: cl}

	// Example1: Found → Some
	params1 := fclient.ToGetParams(client.ObjectKey{Namespace: "default", Name: "exists"})
	result1 := fclient.GetOption[corev1.ConfigMap](params1)(env)()
	fmt.Printf("Example 1: %s\n", resultToStr(result1))

	// Example 2: NotFound → None
	params2 := fclient.ToGetParams(client.ObjectKey{Namespace: "default", Name: "missing"})
	result2 := fclient.GetOption[corev1.ConfigMap](params2)(env)()
	fmt.Printf("Example 2: %s\n", resultToStr(result2))

	// Output:
	// Example 1: Right(Some(*v1.ConfigMap))
	// Example 2: Right(None)
}

func ExampleIgnoreNotFound() {
	resultToStr := func(result ET.Either[error, O.Option[*corev1.ConfigMap]]) string {
		return ET.Fold(
			func(err error) string { return fmt.Sprintf("Left(%T)", err) },
			func(opt O.Option[*corev1.ConfigMap]) string {
				return O.Fold(
					func() string { return "Right(None)" },
					func(cm *corev1.ConfigMap) string { return fmt.Sprintf("Right(Some(%T))", cm) },
				)(opt)
			},
		)(result)
	}

	// Set up environment for reader monad
	env := fclient.Env{Ctx: context.TODO()}

	// Example 1: Right(*corev1.ConfigMap) → Right(Some(*corev1.ConfigMap))
	result1 := F.Pipe1(
		RIOE.Right[fclient.Env, error](&corev1.ConfigMap{ /*...*/ }),
		fclient.IgnoreNotFound,
	)(env)()
	fmt.Printf("Example 1: %s\n", resultToStr(result1))

	// Example 2: Left(*errors.StatusError=NotFound) → Right(None)
	result2 := F.Pipe1(
		RIOE.Left[fclient.Env, *corev1.ConfigMap, error](apierrors.NewNotFound(corev1.Resource("configmaps"), "not-exists-config")),
		fclient.IgnoreNotFound,
	)(env)()
	fmt.Printf("Example 2: %s\n", resultToStr(result2))

	// Example 3: Left(*errors.StatusError) → Left(*errors.StatusError)
	result3 := F.Pipe1(
		RIOE.Left[fclient.Env, *corev1.ConfigMap, error](apierrors.NewBadRequest("bad request")),
		fclient.IgnoreNotFound,
	)(env)()
	fmt.Printf("Example 3: %s\n", resultToStr(result3))

	// Output:
	// Example 1: Right(Some(*v1.ConfigMap))
	// Example 2: Right(None)
	// Example 3: Left(*errors.StatusError)
}

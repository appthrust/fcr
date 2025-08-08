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

func ExampleIgnoreNotFound() {
	resultToStr := func(result ET.Either[error, O.Option[*corev1.ConfigMap]]) string {
		return ET.Fold(
			func(err error) string { return "Left" },
			func(opt O.Option[*corev1.ConfigMap]) string {
				return O.Fold(
					func() string { return "Right(None)" },
					func(cm *corev1.ConfigMap) string { return "Right(Some)" },
				)(opt)
			},
		)(result)
	}

	env := fclient.Env{Ctx: context.TODO()}

	result1 := F.Pipe1(
		RIOE.Right[fclient.Env, error](&corev1.ConfigMap{ /*...*/ }),
		fclient.IgnoreNotFound,
	)(env)()
	fmt.Printf("result1: %v\n", resultToStr(result1))

	result2 := F.Pipe1(
		RIOE.Left[fclient.Env, *corev1.ConfigMap, error](apierrors.NewNotFound(corev1.Resource("configmaps"), "not-exists-config")),
		fclient.IgnoreNotFound,
	)(env)()
	fmt.Printf("result2: %v\n", resultToStr(result2))

	result3 := F.Pipe1(
		RIOE.Left[fclient.Env, *corev1.ConfigMap, error](apierrors.NewBadRequest("bad request")),
		fclient.IgnoreNotFound,
	)(env)()
	fmt.Printf("result3: %v\n", resultToStr(result3))

	// Output:
	// result1: Right(Some)
	// result2: Right(None)
	// result3: Left
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

	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	cl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "exists", Namespace: "default"}}).
		Build()

	env := fclient.Env{Ctx: context.TODO(), Client: cl}

	// Found → Some
	p1 := fclient.ToGetParams(client.ObjectKey{Namespace: "default", Name: "exists"})
	r1 := fclient.GetOption[corev1.ConfigMap](p1)(env)()
	fmt.Printf("r1: %v\n", resultToStr(r1))

	// NotFound → None
	p2 := fclient.ToGetParams(client.ObjectKey{Namespace: "default", Name: "missing"})
	r2 := fclient.GetOption[corev1.ConfigMap](p2)(env)()
	fmt.Printf("r2: %v\n", resultToStr(r2))

	// Output:
	// r1: Right(Some(*v1.ConfigMap))
	// r2: Right(None)
}

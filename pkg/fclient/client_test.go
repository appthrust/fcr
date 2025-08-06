package fclient_test

import (
	"context"
	"os"
	"testing"

	ET "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	RIOE "github.com/IBM/fp-go/readerioeither"
	v1 "github.com/appthrust/fcr/internal/api/v1"
	"github.com/appthrust/fcr/pkg/fclient"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/yaml"
)

func TestMyPackage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MyPackage Suite")
}

var _ = Describe(
	"Client", func() {

		var cl client.Client
		var createdObjects []client.Object

		initClient := func() {
			scheme := runtime.NewScheme()
			err := corev1.AddToScheme(scheme)
			Expect(err).NotTo(HaveOccurred())
			err = apiextv1.AddToScheme(scheme)
			Expect(err).NotTo(HaveOccurred())
			err = v1.AddToScheme(scheme)
			Expect(err).NotTo(HaveOccurred())
			cl, err = client.New(
				config.GetConfigOrDie(), client.Options{Scheme: scheme},
			)
			Expect(err).NotTo(HaveOccurred())
		}

		createObject := func(obj client.Object) {
			err := cl.Patch(
				context.TODO(), obj, client.Apply, client.FieldOwner("test-spec"),
			)
			Expect(err).NotTo(HaveOccurred())
			createdObjects = append(createdObjects, obj)
		}

		createConfigMap := func(name string, namespace string, data map[string]string) {
			createObject(&corev1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					APIVersion: corev1.SchemeGroupVersion.String(),
					Kind:       "ConfigMap",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Data: data,
			})
		}

		applyCRD := func() {
			data, err := os.ReadFile("../../internal/manifests/test.appthrust.com_cats.yaml")
			Expect(err).NotTo(HaveOccurred())
			var crd apiextv1.CustomResourceDefinition
			err = yaml.Unmarshal(data, &crd)
			Expect(err).NotTo(HaveOccurred())
			err = cl.Create(context.TODO(), &crd)
			if err != nil {
				Expect(client.IgnoreAlreadyExists(err)).To(Succeed())
			}
			createdObjects = append(createdObjects, &crd)

			// Wait for CRD to be established
			Eventually(func() bool {
				var updatedCRD apiextv1.CustomResourceDefinition
				err := cl.Get(context.TODO(), client.ObjectKey{Name: crd.Name}, &updatedCRD)
				if err != nil {
					return false
				}
				for _, condition := range updatedCRD.Status.Conditions {
					if condition.Type == apiextv1.Established && condition.Status == apiextv1.ConditionTrue {
						return true
					}
				}
				return false
			}, "10s", "1s").Should(BeTrue(), "CRD should be established")
		}

		createCat := func(name string) {
			cat := &v1.Cat{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "test.appthrust.com/v1",
					Kind:       "Cat",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: "default",
				},
				Spec: v1.CatSpec{},
			}
			createObject(cat)
		}

		cleanupObjects := func() {
			for _, obj := range createdObjects {
				// ignore errors
				_ = cl.Delete(context.TODO(), obj)
			}
		}

		Describe(
			"Get", func() {
				BeforeEach(
					func() {
						initClient()
						createConfigMap("my-config", "default", map[string]string{"foo": "bar"})
						createConfigMap("my-config-referencing", "default", map[string]string{"ref": "my-config"})
					},
				)

				AfterEach(
					func() {
						cleanupObjects()
					},
				)

				It(
					"should return the object", func() {
						env := fclient.Env{Client: cl, Ctx: context.TODO()}
						params := fclient.ToGetParams(client.ObjectKey{Name: "my-config", Namespace: "default"})
						var makeGetMyConfig fclient.ReaderIOEither[*corev1.ConfigMap] = fclient.Get[corev1.ConfigMap](params)
						var getMyConfig IOE.IOEither[error, *corev1.ConfigMap] = makeGetMyConfig(env)
						var result ET.Either[error, *corev1.ConfigMap] = getMyConfig()
						myConfig, err := ET.UnwrapError(result)
						Expect(err).NotTo(HaveOccurred())
						Expect(myConfig.Name).To(Equal("my-config"))
					},
				)

				It(
					"should work with Flow", func() {
						// Define building blocks
						getConfigFromConfig := func(ref *corev1.ConfigMap) fclient.ReaderIOEither[*corev1.ConfigMap] {
							// TODO: validate ref value
							params := fclient.ToGetParams(client.ObjectKey{Name: ref.Data["ref"], Namespace: ref.Namespace})
							return fclient.Get[corev1.ConfigMap](params)
						}

						// Define a flow
						var makeGetConfigMapFromConfigMapRIOE func(fclient.ReaderIOEither[fclient.GetParams]) fclient.ReaderIOEither[*corev1.ConfigMap] = F.Flow2(
							RIOE.Chain(fclient.Get[corev1.ConfigMap]),
							RIOE.Chain(getConfigFromConfig),
						)

						// Instantiate the defined flow with binding dynamic parameters (like user input)
						params := fclient.ToGetParams(client.ObjectKey{Name: "my-config-referencing", Namespace: "default"})
						var getConfigMapFromConfigMapRIOE fclient.ReaderIOEither[*corev1.ConfigMap] = makeGetConfigMapFromConfigMapRIOE(RIOE.Of[fclient.Env, error](params))

						// Bind the runtime environment like client and context
						env := fclient.Env{Client: cl, Ctx: context.TODO()}
						var getConfigMapFromConfigMap IOE.IOEither[error, *corev1.ConfigMap] = getConfigMapFromConfigMapRIOE(env)

						// Evaluate to occur side-effects
						var result ET.Either[error, *corev1.ConfigMap] = getConfigMapFromConfigMap()

						// Inspect the result
						myConfig, err := ET.UnwrapError(result)
						Expect(err).NotTo(HaveOccurred())
						Expect(myConfig.Name).To(Equal("my-config"))
					},
				)
			},
		)

		Describe(
			"List", func() {
				BeforeEach(
					func() {
						initClient()
						createConfigMap("list-config-1", "default", map[string]string{"type": "test"})
						createConfigMap("list-config-2", "default", map[string]string{"type": "test"})
						createConfigMap("list-config-3", "default", map[string]string{"type": "test"})
					},
				)

				AfterEach(
					func() {
						cleanupObjects()
					},
				)

				It(
					"should return all objects in namespace", func() {
						env := fclient.Env{Client: cl, Ctx: context.TODO()}
						params := fclient.ToListParams(client.InNamespace("default"))
						var makeListConfigMaps fclient.ReaderIOEither[*corev1.ConfigMapList] = fclient.List[corev1.ConfigMapList](params)
						var listConfigMaps IOE.IOEither[error, *corev1.ConfigMapList] = makeListConfigMaps(env)
						var result ET.Either[error, *corev1.ConfigMapList] = listConfigMaps()
						configList, err := ET.UnwrapError(result)
						Expect(err).NotTo(HaveOccurred())

						// Check that we have our test ConfigMaps
						var foundNames []string
						for _, cm := range configList.Items {
							foundNames = append(foundNames, cm.Name)
						}
						Expect(foundNames).To(ContainElement("list-config-1"))
						Expect(foundNames).To(ContainElement("list-config-2"))
						Expect(foundNames).To(ContainElement("list-config-3"))
					},
				)

				It(
					"should work with Flow for counting", func() {
						var countByType = func(configList *corev1.ConfigMapList) int {
							filtered := lo.Filter(configList.Items, func(item corev1.ConfigMap, _ int) bool {
								return item.Data["type"] == "test"
							})
							return len(filtered)
						}

						// Define a flow
						var makeCountConfigMapsRIOE func(fclient.ReaderIOEither[fclient.ListParams]) fclient.ReaderIOEither[int] = F.Flow2(
							RIOE.Chain(fclient.List[corev1.ConfigMapList]),
							RIOE.Map[fclient.Env, error](countByType),
						)

						// Instantiate the defined flow with binding dynamic parameters
						params := fclient.ToListParams(client.InNamespace("default"))
						var countConfigMapsRIOE fclient.ReaderIOEither[int] = makeCountConfigMapsRIOE(RIOE.Of[fclient.Env, error](params))

						// Bind the runtime environment
						env := fclient.Env{Client: cl, Ctx: context.TODO()}
						var countConfigMaps IOE.IOEither[error, int] = countConfigMapsRIOE(env)

						// Evaluate to occur side-effects
						var result ET.Either[error, int] = countConfigMaps()

						// Inspect the result
						configs, err := ET.UnwrapError(result)
						Expect(err).NotTo(HaveOccurred())
						Expect(configs).To(Equal(3))
					},
				)
			},
		)

		Describe(
			"Create", func() {
				BeforeEach(
					func() {
						initClient()
					},
				)

				AfterEach(
					func() {
						cleanupObjects()
					},
				)

				It(
					"should create a new object", func() {
						env := fclient.Env{Client: cl, Ctx: context.TODO()}
						configMap := &corev1.ConfigMap{
							TypeMeta: metav1.TypeMeta{
								APIVersion: corev1.SchemeGroupVersion.String(),
								Kind:       "ConfigMap",
							},
							ObjectMeta: metav1.ObjectMeta{
								Name:      "create-test-config",
								Namespace: "default",
							},
							Data: map[string]string{"create": "test"},
						}
						params := fclient.ToCreateParams(configMap)
						var makeCreateConfigMap fclient.ReaderIOEither[fclient.Unit] = fclient.Create(params)
						var createConfigMap IOE.IOEither[error, fclient.Unit] = makeCreateConfigMap(env)
						var result ET.Either[error, fclient.Unit] = createConfigMap()
						_, err := ET.UnwrapError(result)
						Expect(err).NotTo(HaveOccurred())

						// Add to cleanup list
						createdObjects = append(createdObjects, configMap)

						// Verify the object was created in the API server
						getParams := fclient.ToGetParams(client.ObjectKey{Name: "create-test-config", Namespace: "default"})
						result2 := fclient.Get[corev1.ConfigMap](getParams)(env)()
						Expect(ET.IsRight(result2)).To(BeTrue())
					},
				)
			},
		)

		Describe(
			"Delete", func() {
				BeforeEach(
					func() {
						initClient()
						createConfigMap("delete-test-config", "default", map[string]string{"delete": "test"})
					},
				)

				AfterEach(
					func() {
						// Clear the list since Delete will handle deletion
						createdObjects = []client.Object{}
					},
				)

				It(
					"should delete an existing object", func() {
						env := fclient.Env{Client: cl, Ctx: context.TODO()}
						configMap := &corev1.ConfigMap{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "delete-test-config",
								Namespace: "default",
							},
						}
						params := fclient.ToDeleteParams(configMap)
						var makeDeleteConfigMap fclient.ReaderIOEither[fclient.Unit] = fclient.Delete(params)
						var deleteConfigMap IOE.IOEither[error, fclient.Unit] = makeDeleteConfigMap(env)
						var result ET.Either[error, fclient.Unit] = deleteConfigMap()
						_, err := ET.UnwrapError(result)
						Expect(err).NotTo(HaveOccurred())

						// Verify the object was deleted from the API server
						getParams := fclient.ToGetParams(client.ObjectKey{Name: "delete-test-config", Namespace: "default"})
						result2 := fclient.Get[corev1.ConfigMap](getParams)(env)()
						_, err = ET.UnwrapError(result2)
						Expect(client.IgnoreNotFound(err)).To(Succeed())
					},
				)
			},
		)

		Describe(
			"Update", func() {
				BeforeEach(
					func() {
						initClient()
						createConfigMap("update-test-config", "default", map[string]string{"version": "1"})
					},
				)

				AfterEach(
					func() {
						cleanupObjects()
					},
				)

				It(
					"should update an existing object", func() {
						env := fclient.Env{Client: cl, Ctx: context.TODO()}

						// First get the object to update
						getParams := fclient.ToGetParams(client.ObjectKey{Name: "update-test-config", Namespace: "default"})
						var makeGetConfigMap fclient.ReaderIOEither[*corev1.ConfigMap] = fclient.Get[corev1.ConfigMap](getParams)
						var getConfigMap IOE.IOEither[error, *corev1.ConfigMap] = makeGetConfigMap(env)
						var getResult ET.Either[error, *corev1.ConfigMap] = getConfigMap()
						configMap, err := ET.UnwrapError(getResult)
						Expect(err).NotTo(HaveOccurred())

						// Update the object
						configMap.Data["version"] = "2"
						updateParams := fclient.ToUpdateParams(configMap)
						var makeUpdateConfigMap fclient.ReaderIOEither[fclient.Unit] = fclient.Update(updateParams)
						var updateConfigMap IOE.IOEither[error, fclient.Unit] = makeUpdateConfigMap(env)
						var updateResult ET.Either[error, fclient.Unit] = updateConfigMap()
						_, err = ET.UnwrapError(updateResult)
						Expect(err).NotTo(HaveOccurred())

						// Verify the object was updated in the API server
						verifyGetParams := fclient.ToGetParams(client.ObjectKey{Name: "update-test-config", Namespace: "default"})
						result2 := fclient.Get[corev1.ConfigMap](verifyGetParams)(env)()
						updatedConfigMap, err := ET.UnwrapError(result2)
						Expect(err).NotTo(HaveOccurred())
						Expect(updatedConfigMap.Data["version"]).To(Equal("2"))
					},
				)
			},
		)

		Describe(
			"Patch", func() {
				BeforeEach(
					func() {
						initClient()
						createConfigMap("patch-test-config", "default", map[string]string{"original": "value"})
					},
				)

				AfterEach(
					func() {
						cleanupObjects()
					},
				)

				It(
					"should patch an existing object", func() {
						env := fclient.Env{Client: cl, Ctx: context.TODO()}
						configMap := &corev1.ConfigMap{
							TypeMeta: metav1.TypeMeta{
								APIVersion: corev1.SchemeGroupVersion.String(),
								Kind:       "ConfigMap",
							},
							ObjectMeta: metav1.ObjectMeta{
								Name:      "patch-test-config",
								Namespace: "default",
							},
							Data: map[string]string{"patched": "value"},
						}
						params := fclient.ToPatchParams(configMap, client.Apply, client.FieldOwner("test-patch"))
						var makePatchConfigMap fclient.ReaderIOEither[fclient.Unit] = fclient.Patch(params)
						var patchConfigMap IOE.IOEither[error, fclient.Unit] = makePatchConfigMap(env)
						var result ET.Either[error, fclient.Unit] = patchConfigMap()
						_, err := ET.UnwrapError(result)
						Expect(err).NotTo(HaveOccurred())

						// Verify the object was patched in the API server
						getParams := fclient.ToGetParams(client.ObjectKey{Name: "patch-test-config", Namespace: "default"})
						result2 := fclient.Get[corev1.ConfigMap](getParams)(env)()
						patchedConfigMap, err := ET.UnwrapError(result2)
						Expect(err).NotTo(HaveOccurred())
						Expect(patchedConfigMap.Data["patched"]).To(Equal("value"))
					},
				)
			},
		)

		Describe(
			"DeleteAllOf", func() {
				BeforeEach(
					func() {
						initClient()
						createConfigMap("deleteallof-test-1", "default", map[string]string{"batch": "delete"})
						createConfigMap("deleteallof-test-2", "default", map[string]string{"batch": "delete"})
						createConfigMap("deleteallof-test-3", "default", map[string]string{"batch": "delete"})
					},
				)

				AfterEach(
					func() {
						// Clear the list since DeleteAllOf will handle deletion
						createdObjects = []client.Object{}
					},
				)

				It(
					"should delete all matching objects", func() {
						env := fclient.Env{Client: cl, Ctx: context.TODO()}
						params := fclient.ToDeleteAllOfParams(
							client.InNamespace("default"),
							client.MatchingLabels{},
						)
						var makeDeleteAllOfConfigMaps fclient.ReaderIOEither[fclient.Unit] = fclient.DeleteAllOf[corev1.ConfigMap](params)
						var deleteAllOfConfigMaps IOE.IOEither[error, fclient.Unit] = makeDeleteAllOfConfigMaps(env)
						var result ET.Either[error, fclient.Unit] = deleteAllOfConfigMaps()
						_, err := ET.UnwrapError(result)
						Expect(err).NotTo(HaveOccurred())

						// Verify the objects were deleted from the API server
						listParams := fclient.ToListParams(client.InNamespace("default"))
						listResult := fclient.List[corev1.ConfigMapList](listParams)(env)()
						configList, err := ET.UnwrapError(listResult)
						Expect(err).NotTo(HaveOccurred())
						foundNames := lo.Map(configList.Items, func(item corev1.ConfigMap, _ int) string {
							return item.Name
						})
						Expect(foundNames).NotTo(ContainElement("deleteallof-test-1"))
						Expect(foundNames).NotTo(ContainElement("deleteallof-test-2"))
						Expect(foundNames).NotTo(ContainElement("deleteallof-test-3"))
					},
				)
			},
		)

		Describe(
			"StatusUpdate", func() {
				BeforeEach(
					func() {
						initClient()
						applyCRD()
						createCat("status-update-test-cat")
					},
				)

				AfterEach(
					func() {
						cleanupObjects()
					},
				)

				It(
					"should update the status of an existing object", func() {
						env := fclient.Env{Client: cl, Ctx: context.TODO()}

						// First get the cat to update its status
						getParams := fclient.ToGetParams(client.ObjectKey{Name: "status-update-test-cat", Namespace: "default"})
						var makeGetCat fclient.ReaderIOEither[*v1.Cat] = fclient.Get[v1.Cat](getParams)
						var getCat IOE.IOEither[error, *v1.Cat] = makeGetCat(env)
						var getResult ET.Either[error, *v1.Cat] = getCat()
						cat, err := ET.UnwrapError(getResult)
						Expect(err).NotTo(HaveOccurred())

						// Update the status
						cat.Status = &v1.CatStatus{Sleepy: true}
						statusUpdateParams := fclient.ToStatusUpdateParams(cat)
						var makeStatusUpdateCat fclient.ReaderIOEither[fclient.Unit] = fclient.StatusUpdate(statusUpdateParams)
						var statusUpdateCat IOE.IOEither[error, fclient.Unit] = makeStatusUpdateCat(env)
						var updateResult ET.Either[error, fclient.Unit] = statusUpdateCat()
						_, err = ET.UnwrapError(updateResult)
						Expect(err).NotTo(HaveOccurred())

						// Verify the status was updated in the API server
						verifyGetParams := fclient.ToGetParams(client.ObjectKey{Name: "status-update-test-cat", Namespace: "default"})
						result2 := fclient.Get[v1.Cat](verifyGetParams)(env)()
						updatedCat, err := ET.UnwrapError(result2)
						Expect(err).NotTo(HaveOccurred())
						Expect(updatedCat.Status).NotTo(BeNil())
						Expect(updatedCat.Status.Sleepy).To(BeTrue())
					},
				)
			},
		)

		Describe(
			"StatusPatch", func() {
				BeforeEach(
					func() {
						initClient()
						applyCRD()
						createCat("status-patch-test-cat")
					},
				)

				AfterEach(
					func() {
						cleanupObjects()
					},
				)

				It(
					"should patch the status of an existing object", FlakeAttempts(5), func() {
						env := fclient.Env{Client: cl, Ctx: context.TODO()}

						// First get the existing cat
						getParams := fclient.ToGetParams(client.ObjectKey{Name: "status-patch-test-cat", Namespace: "default"})
						getResult := fclient.Get[v1.Cat](getParams)(env)()
						cat, err := ET.UnwrapError(getResult)
						Expect(err).NotTo(HaveOccurred())
						Expect(cat.Status).To(BeNil())

						// Create a cat object with status to patch
						statusPatchCat := &v1.Cat{
							TypeMeta: metav1.TypeMeta{
								APIVersion: "test.appthrust.com/v1",
								Kind:       "Cat",
							},
							ObjectMeta: metav1.ObjectMeta{
								Name:      cat.Name,
								Namespace: cat.Namespace,
							},
							Status: &v1.CatStatus{Sleepy: false},
						}
						params := fclient.ToStatusPatchParams(statusPatchCat, client.Apply, client.FieldOwner("test-status-patch"))
						var makeStatusPatch fclient.ReaderIOEither[fclient.Unit] = fclient.StatusPatch(params)
						var statusPatch IOE.IOEither[error, fclient.Unit] = makeStatusPatch(env)
						var result ET.Either[error, fclient.Unit] = statusPatch()
						_, err = ET.UnwrapError(result)
						Expect(err).NotTo(HaveOccurred())

						// Verify the status was patched in the API server
						verifyGetParams := fclient.ToGetParams(client.ObjectKey{Name: "status-patch-test-cat", Namespace: "default"})
						getResult2 := fclient.Get[v1.Cat](verifyGetParams)(env)()
						patchedCat, err := ET.UnwrapError(getResult2)
						Expect(err).NotTo(HaveOccurred())
						Expect(patchedCat.Status).NotTo(BeNil())
						Expect(patchedCat.Status.Sleepy).To(BeFalse())
					},
				)
			},
		)
	},
)

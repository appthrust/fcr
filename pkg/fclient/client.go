// Package fclient provides functional programming wrappers for Kubernetes controller-runtime [client.Client] operations.
package fclient

import (
	"context"

	ET "github.com/IBM/fp-go/either"
	IOE "github.com/IBM/fp-go/ioeither"
	RIOE "github.com/IBM/fp-go/readerioeither"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ReaderIOEither is a type alias for ReaderIOEither with Env and error types.
type ReaderIOEither[T any] = RIOE.ReaderIOEither[Env, error, T]

// IOEither is a type alias for IOEither with error type.
type IOEither[T any] = IOE.IOEither[error, T]

// Either is a type alias for Either with error type.
type Either[T any] = ET.Either[error, T]

// Unit represents the unit type (empty struct).
type Unit = struct{}

// UnitValue is the zero value of Unit type.
var UnitValue = Unit{}

// Env represents the environment containing context and client for operations.
type Env struct {
	Ctx    context.Context
	Client client.Client
}

// GetParams contains parameters for Get operations.
type GetParams struct {
	key  client.ObjectKey
	opts []client.GetOption
}

// ToGetParams creates GetParams from key and options.
func ToGetParams(key client.ObjectKey, opts ...client.GetOption) GetParams {
	return GetParams{key, opts}
}

// Get retrieves a Kubernetes object using the provided parameters.
func Get[T any, OP ObjectPointer[T]](p GetParams) ReaderIOEither[OP] {
	return readerize(func(env Env) (OP, error) {
		var obj T       // Initialize the object with zero value
		ptr := OP(&obj) // Cast to the pointer type
		return ptr, env.Client.Get(env.Ctx, p.key, ptr, p.opts...)
	})
}

// ListParams contains parameters for List operations.
type ListParams struct {
	opts []client.ListOption
}

// ToListParams creates ListParams from options.
func ToListParams(opts ...client.ListOption) ListParams {
	return ListParams{opts}
}

// List retrieves a list of Kubernetes objects using the provided parameters.
func List[T any, OLP ObjectListPointer[T]](p ListParams) ReaderIOEither[OLP] {
	return readerize(func(env Env) (OLP, error) {
		var obj T        // Initialize the object with zero value
		ptr := OLP(&obj) // Cast to the pointer type
		return ptr, env.Client.List(env.Ctx, ptr, p.opts...)
	})
}

// CreateParams contains parameters for Create operations.
type CreateParams struct {
	obj  client.Object
	opts []client.CreateOption
}

// ToCreateParams creates CreateParams from object and options.
func ToCreateParams(obj client.Object, opts ...client.CreateOption) CreateParams {
	return CreateParams{obj, opts}
}

// Create creates a Kubernetes object using the provided parameters.
func Create(p CreateParams) ReaderIOEither[Unit] {
	return readerize(func(env Env) (Unit, error) {
		return UnitValue, env.Client.Create(env.Ctx, p.obj, p.opts...)
	})
}

// DeleteParams contains parameters for Delete operations.
type DeleteParams struct {
	obj  client.Object
	opts []client.DeleteOption
}

// ToDeleteParams creates DeleteParams from object and options.
func ToDeleteParams(obj client.Object, opts ...client.DeleteOption) DeleteParams {
	return DeleteParams{obj, opts}
}

// Delete deletes a Kubernetes object using the provided parameters.
func Delete(p DeleteParams) ReaderIOEither[Unit] {
	return readerize(func(env Env) (Unit, error) {
		return UnitValue, env.Client.Delete(env.Ctx, p.obj, p.opts...)
	})
}

// UpdateParams contains parameters for Update operations.
type UpdateParams struct {
	obj  client.Object
	opts []client.UpdateOption
}

// ToUpdateParams creates UpdateParams from object and options.
func ToUpdateParams(obj client.Object, opts ...client.UpdateOption) UpdateParams {
	return UpdateParams{obj, opts}
}

// Update updates a Kubernetes object using the provided parameters.
func Update(p UpdateParams) ReaderIOEither[Unit] {
	return readerize(func(env Env) (Unit, error) {
		return UnitValue, env.Client.Update(env.Ctx, p.obj, p.opts...)
	})
}

// PatchParams contains parameters for Patch operations.
type PatchParams struct {
	obj   client.Object
	patch client.Patch
	opts  []client.PatchOption
}

// ToPatchParams creates PatchParams from object, patch, and options.
func ToPatchParams(obj client.Object, patch client.Patch, opts ...client.PatchOption) PatchParams {
	return PatchParams{obj, patch, opts}
}

// Patch patches a Kubernetes object using the provided parameters.
func Patch(p PatchParams) ReaderIOEither[Unit] {
	return readerize(func(env Env) (Unit, error) {
		return UnitValue, env.Client.Patch(env.Ctx, p.obj, p.patch, p.opts...)
	})
}

// DeleteAllOfParams contains parameters for DeleteAllOf operations.
type DeleteAllOfParams struct {
	opts []client.DeleteAllOfOption
}

// ToDeleteAllOfParams creates DeleteAllOfParams from options.
func ToDeleteAllOfParams(opts ...client.DeleteAllOfOption) DeleteAllOfParams {
	return DeleteAllOfParams{opts}
}

// DeleteAllOf deletes all objects of a specific type using the provided parameters.
func DeleteAllOf[T any, OP ObjectPointer[T]](p DeleteAllOfParams) ReaderIOEither[Unit] {
	return readerize(func(env Env) (Unit, error) {
		var obj T       // Initialize the object with zero value
		ptr := OP(&obj) // Cast to the pointer type
		return UnitValue, env.Client.DeleteAllOf(env.Ctx, ptr, p.opts...)
	})
}

// StatusUpdateParams contains parameters for status update operations.
type StatusUpdateParams struct {
	obj  client.Object
	opts []client.SubResourceUpdateOption
}

// ToStatusUpdateParams creates StatusUpdateParams from object and options.
func ToStatusUpdateParams(obj client.Object, opts ...client.SubResourceUpdateOption) StatusUpdateParams {
	return StatusUpdateParams{obj, opts}
}

// StatusUpdate updates the status of a Kubernetes object using the provided parameters.
func StatusUpdate(p StatusUpdateParams) ReaderIOEither[Unit] {
	return readerize(func(env Env) (Unit, error) {
		return UnitValue, env.Client.Status().Update(env.Ctx, p.obj, p.opts...)
	})
}

// StatusPatchParams contains parameters for status patch operations.
type StatusPatchParams struct {
	obj   client.Object
	patch client.Patch
	opts  []client.SubResourcePatchOption
}

// ToStatusPatchParams creates StatusPatchParams from object, patch, and options.
func ToStatusPatchParams(obj client.Object, patch client.Patch, opts ...client.SubResourcePatchOption) StatusPatchParams {
	return StatusPatchParams{obj, patch, opts}
}

// StatusPatch patches the status of a Kubernetes object using the provided parameters.
func StatusPatch(p StatusPatchParams) ReaderIOEither[Unit] {
	return readerize(func(env Env) (Unit, error) {
		return UnitValue, env.Client.Status().Patch(env.Ctx, p.obj, p.patch, p.opts...)
	})
}

// ObjectPointer is a type that constraints T to be a pointer type and implements [client.Object].
type ObjectPointer[T any] interface {
	client.Object // Rule 1: T must implement client.Object
	*T            // Rule 2: T must be a pointer type
}

// ObjectListPointer is a type that constraints T to be a pointer type and implements [client.ObjectList].
type ObjectListPointer[T any] interface {
	client.ObjectList // Rule 1: T must implement client.ObjectList
	*T                // Rule 2: T must be a pointer type
}

func readerize[T any](f func(env Env) (T, error)) ReaderIOEither[T] {
	return func(env Env) IOEither[T] {
		return IOE.TryCatchError(func() (T, error) {
			return f(env)
		})
	}
}

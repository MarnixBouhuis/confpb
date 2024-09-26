package config

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Overlay merges two protobuf message structs into one.
// It applies "overlay" on top of "base" and returns the resulting message.
// The original messages are not modified.
func Overlay[T any, TMsg ProtoMessage[T]](base TMsg, overlay TMsg) TMsg {
	if base == nil {
		// No base, return clone of overlay
		return proto.Clone(overlay).(TMsg) //nolint:forcetypeassert // proto.Clone always returns the return type. If not it's fine to panic.
	}

	clone := proto.Clone(base).(TMsg) //nolint:forcetypeassert // proto.Clone always returns the return type. If not it's fine to panic.
	if overlay != nil {
		// We have an overlay, apply it to the clone
		mergeOverlay(clone.ProtoReflect(), overlay.ProtoReflect())
	}
	return clone
}

// mergeOverlay merges two protobuf message structs into one.
// Its just like proto.Merge, but instead of merging things like lists and maps using appends it overwrites them.
// "base" and "overlay" should be the same proto message.
func mergeOverlay(base protoreflect.Message, overlay protoreflect.Message) {
	if overlay == nil {
		// No fields to overlay
		return
	}

	overlay.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		switch {
		case fd.Message() != nil && !fd.IsMap() && !fd.IsList():
			mergeOverlay(base.Mutable(fd).Message(), v.Message())
		case fd.Kind() == protoreflect.BytesKind:
			src := v.Bytes()
			clone := make([]byte, len(src))
			copy(clone, src)
			base.Set(fd, protoreflect.ValueOfBytes(clone))
		default:
			base.Set(fd, v)
		}
		return true
	})
}

package types

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// Attribute is an entity attribute, that holds specific data such as the health of the entity. Each attribute
// holds a default value, maximum and minimum value, name and its current value.
type Attribute struct {
	protocol.AttributeValue
	// Default is the default value of the attribute. It's not clear why this field must be sent to the
	// client, but it is required regardless.
	Default float32
	// Modifiers is a slice of AttributeModifiers that are applied to the attribute.
	Modifiers []protocol.AttributeModifier
}

// Marshal encodes/decodes an Attribute.
func (x *Attribute) Marshal(r protocol.IO) {
	r.Float32(&x.Min)
	r.Float32(&x.Max)
	r.Float32(&x.Value)
	r.Float32(&x.Default)
	r.String(&x.Name)
	protocol.Slice(r, &x.Modifiers)
}

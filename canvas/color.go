package canvas

var MAX_CHANNEL_VALUE uint16 = 65535

type Color struct {
	value uint64
}

func MakeColor(value uint64) Color {
	return Color{value}
}

func MakeColorRgba(r, g, b, a uint16) Color {
	v := uint64(r)
	v <<= 16
	v += uint64(g)
	v <<= 16
	v += uint64(b)
	v <<= 16
	v += uint64(a)

	return MakeColor(v)
}

func (p Color) Value() uint64 {
	return p.value
}

func (p Color) R() uint16 {
	return uint16(p.Value() >> 48)
}

func (p Color) G() uint16 {
	return uint16((p.Value() >> 32) & uint64(MAX_CHANNEL_VALUE))
}

func (p Color) B() uint16 {
	return uint16((p.Value() >> 16) & uint64(MAX_CHANNEL_VALUE))
}

func (p Color) A() uint16 {
	return uint16(p.Value() & uint64(MAX_CHANNEL_VALUE))
}

func (p Color) Rgba() (r, g, b, a uint16) {
	return p.R(), p.G(), p.B(), p.A()
}

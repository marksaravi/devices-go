package gpio

type Level = bool

const (
	// Low represents 0v.
	Low Level = false
	// High represents Vin, generally 3.3v or 5v.
	High Level = true
)

type GPIOPinOut interface {
	Out(Level)
}

type GPIOPinIn interface {
	Read() Level
}

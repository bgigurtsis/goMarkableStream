package remarkable

// DeviceModel represents the type of reMarkable device being used
type DeviceModel int

const (
	// UnknownDevice represents an unidentified reMarkable device
	UnknownDevice DeviceModel = iota
	// Remarkable2 represents the reMarkable 2 device
	Remarkable2
	// RemarkablePaperPro represents the reMarkable Paper Pro device
	RemarkablePaperPro
	// RemarkablePaperPure represents the reMarkable Paper Pure device
	RemarkablePaperPure
)

func (d DeviceModel) String() string {
	switch d {
	case Remarkable2:
		return "Remarkable2"
	case RemarkablePaperPro:
		return "RemarkablePaperPro"
	case RemarkablePaperPure:
		return "RemarkablePaperPure"
	default:
		return "UnknownDevice"
	}
}

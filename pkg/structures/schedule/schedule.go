package schedule

type Schedule struct {
	GatingCycle          float32 `yaml:"gating-cycle"`
	Isochronous          int     `yaml:"isochronous"`
	Cyclic               int     `yaml:"cyclic"`
	ProfileX             int     `yaml:"profile-x"`
	UserSpecific         int     `yaml:"user-specific"`
	StrictPriorityMixed  int     `yaml:"strict-priority-mixed"`
	BandwidthReservation int     `yaml:"bandwidth-reservation"`
}

package options

type MergePolicy int

const (
	// default
	//Always MergePolicy = iota	// No restrictions on when merge operations can occur (default)
	Always MergePolicy = iota
	Never   					// never merge files
	Windows 					// Merge operations occur during specified hours
)

type MergeWindow struct {
	Start int
	End int
}

type Merge struct {
	//Policy MergePolicy
	Policy MergePolicy
	Window MergeWindow
	Interval int
}

type TestType int

const (
	T1 TestType = iota
	T2
	T3
)

func NewDefaultMerge() Merge {

	return Merge{
		Policy: Always,
		Interval: 3600,
	}
}
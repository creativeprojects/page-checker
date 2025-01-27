package lib

type Flags struct {
	Verbose      bool
	ProfilerPort int
	Output       string
	URL          string
	Width        int
	Height       int
	Cookie       []string
	Useragent    string
	Timeout      int
	Sleep        int
	Quality      int      // flag for screenshot only
	Noidlewait   bool     // flag for screenshot only
	Nowait       bool     // flag for screenshot only
	Savepartial  bool     // flag for screenshot only
	Scale        float64  // flag for screenshot only
	ClickXPath   []string // flag for screenshot only
	Stealth      bool
}

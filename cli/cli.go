package cli

type Flags struct {
	Url          string
	NoRecursion  bool
	NoSubdomains bool
	Threads      uint64
}

package comm

type Originator struct {
	Identity string
	Platform string
}

type IngressMessage struct {
	Content    string
	Originator Originator
}

// Some other junk

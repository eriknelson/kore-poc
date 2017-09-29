package comm

type Originator struct {
	Identity string
	Platform string
}

type IngressMessage struct {
	Content    string
	Originator Originator
}

// AdapterIngressMessage - Similar to an Ingress, but one expected to be
// used by the adapters themselves. The Engine will annotate them with the
// full Originator type as an IngressMessage before passing on to the Plugin.
type AdapterIngressMessage struct {
	Identity string
	Content  string
}

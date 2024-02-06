package main

type Delegate struct{}

func (delegate *Delegate) NodeMeta(limit int) []byte {
	mb := make([]byte, limit)
	// TODO: Return the meta for this current node
	return mb
}

func (delegate *Delegate) NotifyMsg([]byte) {
	// TODO: Process message from peers
}

func (delegate *Delegate) GetBroadcasts(overhead, limit int) [][]byte {
	// TODO: Get broadcasts from the broadcast queue
	return [][]byte{}
}

func (delegate *Delegate) LocalState(join bool) []byte {
	// TODO: Return local state
	return []byte{}
}

func (delegate *Delegate) MergeRemoteState(buf []byte, join bool) {
	// TODO: Merge remote state
}

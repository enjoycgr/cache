package cache

import "cache/cachepb"

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

type PeerGetter interface {
	Get(in *cachepb.GetRequest, out *cachepb.GetResponse) error
	Set(in *cachepb.SetRequest, out *cachepb.SetResponse) error
}

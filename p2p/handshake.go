package p2p

type HandShakeFunc func(any) error

func NOPHandshakeFunc(any) error { return nil }

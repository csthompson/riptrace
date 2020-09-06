package types

type Publisher interface {
	Publish()
	Shutdown()
}

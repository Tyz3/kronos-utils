package lib

type Engine interface {
	Start() bool
	Stop() bool
}

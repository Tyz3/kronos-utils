package utils

type Engine interface {
	Start() bool
	Stop() bool
}

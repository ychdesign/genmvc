package models

type Server struct {
	Name       string
	Port       int
	enableLogs bool
	BaseDomain string
}

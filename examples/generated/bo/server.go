package bo

type Server struct {
	Name       string `json:"name"`
	Port       int    `json:"port"`
	enableLogs bool
	BaseDomain string `json:"base_domain"`
}

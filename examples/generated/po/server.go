package po

type Server struct {
	Name       string `gorm:"column:name"`
	Port       int    `gorm:"column:port"`
	enableLogs bool
	BaseDomain string `gorm:"column:base_domain"`
}

const ServerTableName = ""

func (Server) TableName() string {
	return ServerTableName
}

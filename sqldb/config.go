package sqldb

// Config stores configuration to sqldatabase
type Config struct {
	DriverName     string `json:"driverName"`
	DataSourceName string `json:"dataSourceName"`
	//ConnMaxIdleTime
}

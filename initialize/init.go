package initialize

func init() {
	go initapi()
	initlog()
	initdb()
}

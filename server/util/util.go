package util

// 错误检查
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

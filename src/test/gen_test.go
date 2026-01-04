package test

import (
	"testing"
	// "github.com/go-framework-v2/go-backnormal-gen/gen"
)

// 遇见测试不通过的问题，可能是包版本问题，执行以下：
// go mod tidy
// go get gorm.io/gorm@v1.25.0
// go get gorm.io/plugin/dbresolver@v1.4.7
func TestGen(t *testing.T) {
	// dsn := "root:dev123456@tcp(192.168.1.99:13301)/biz_db?charset=utf8mb4&parseTime=True&loc=Local"
	// tableList := []string{"conf", "conf_ext", "conf_his"}
	// poDir := "/Users/huanlema/Documents/Code/my_code/github_go-framework-v2/go-backnormal/src/internal/conf/dao/po"
	// boDir := "/Users/huanlema/Documents/Code/my_code/github_go-framework-v2/go-backnormal/src/internal/conf/model/bo"
	// poPath := "github.com/go-framework-v2/go-backnormal-grpc/src/internal/conf/dao/po"
	// daoDir := "/Users/huanlema/Documents/Code/my_code/github_go-framework-v2/go-backnormal/src/internal/conf/dao"
	// boPath := "github.com/go-framework-v2/go-backnormal-grpc/src/internal/conf/model/bo"
	// err := gen.GenPoBoDao(dsn, tableList, poDir, boDir, poPath, daoDir, boPath)
	// if err != nil {
	// 	t.Error(err)
	// }
}

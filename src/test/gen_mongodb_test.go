package test

import (
	"testing"
	// "github.com/go-framework-v2/go-backnormal-gen/bo"
	// "github.com/go-framework-v2/go-backnormal-gen/dao"
	// "github.com/go-framework-v2/go-backnormal-gen/po"
)

// 遇见测试不通过的问题，可能是包版本问题，执行以下：
// go mod tidy
// go get gorm.io/gorm@v1.25.0
// go get gorm.io/plugin/dbresolver@v1.4.7
func TestGen_MongoDB(t *testing.T) {
	// host := "127.0.0.1"
	// port := 27017
	// database := "test_1029"
	// username := "appuser"
	// password := "app123"
	// tables := []string{"users"}
	// poDir := "/Users/huanlema/Documents/Code/my_code/github_go-framework-v2/go-backnormal/src/internal/mongodb/dao/po"
	// boDir := "/Users/huanlema/Documents/Code/my_code/github_go-framework-v2/go-backnormal/src/internal/mongodb/model/bo"
	// poPath := "github.com/go-framework-v2/go-backnormal-grpc/src/internal/conf/mongodb/po"
	// daoDir := "/Users/huanlema/Documents/Code/my_code/github_go-framework-v2/go-backnormal/src/internal/mongodb/dao"
	// boPath := "github.com/go-framework-v2/go-backnormal-grpc/src/internal/mongodb/model/bo"

	// err := po.GenPo_MongoDB_WithConfig(host, port, database, username, password, tables, poDir)
	// if err != nil {
	// 	t.Error(err)
	// }
	// err = bo.GenBo_MongoDB_WithConfig(host, port, database, username, password, tables, boDir, poPath)
	// if err != nil {
	// 	t.Error(err)
	// }
	// err = dao.GenDao_MongoDB_WithConfig(host, port, database, username, password, tables, daoDir, boPath)
	// if err != nil {
	// 	t.Error(err)
	// }
}

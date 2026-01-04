package dao

import (
	"github.com/go-framework-v2/go-backnormal-grpc/src/cons"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/conf/model/bo"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// 设计要点：
// 接口组合实现扩展
// 结构体嵌入实现代码复用
type (
	// 基础接口 + 扩展方法
	ConfTypeDao interface {
		// 基础接口，基于单表
		Truncate() error                        // 清空表并重置自增主键
		FindOne(id int) (*bo.ConfTypeBo, error) // 根据id查询
		FindOneByUk()                           // 根据唯一键查询
		Insert(obj bo.ConfTypeBo) (*bo.ConfTypeBo, error)
		InsertBatch(objs []bo.ConfTypeBo) (int64, error) // 500条以内
		Update(obj bo.ConfTypeBo) (*bo.ConfTypeBo, error)

		// 扩展接口
		GetAllConfTypeList() ([]bo.ConfTypeBo, error) // 查询所有配置类型列表
		GetConfType(param bo.ConfTypeBo) (*bo.ConfTypeBo, error)
	}

	// 自定义模型实现
	customConfTypeDao struct {
		// 默认实现
		db *gorm.DB // 不带事务的基础数据库连接
		tx *gorm.DB // 带事务的数据库连接

		// 扩展实现
		cache *redis.Client // 新增redis缓存
	}
)

// 对外暴露扩展接口
func NewConfTypeDao(db *gorm.DB, tx *gorm.DB, cache *redis.Client) ConfTypeDao {
	return &customConfTypeDao{
		db:    db,
		tx:    tx,
		cache: cache,
	}
}

// 基础方法
// Truncate 清空表并重置自增主键
func (d *customConfTypeDao) Truncate() error {
	// 临时禁用外键检查
	if err := d.tx.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return err
	}

	// 1. 清空表数据
	if err := d.tx.Exec("TRUNCATE TABLE " + "conf_type").Error; err != nil {
		return err
	}

	// 2. 重置自增主键为1
	if err := d.tx.Exec("ALTER TABLE " + "conf_type" + " AUTO_INCREMENT = 1").Error; err != nil {
		return err
	}

	return nil
}

// FindOne 根据id查询
func (d *customConfTypeDao) FindOne(id int) (*bo.ConfTypeBo, error) {
	var obj bo.ConfTypeBo

	err := d.db.First(&obj, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &obj, err
}

// FindOneByUk 根据唯一键查询
func (d *customConfTypeDao) FindOneByUk() {

}

// Insert 插入 (事务内)
func (d *customConfTypeDao) Insert(obj bo.ConfTypeBo) (*bo.ConfTypeBo, error) {
	err := d.tx.Create(&obj).Error
	if err != nil {
		return nil, err
	}

	return &obj, err
}

// InsertBatch 插入 (事务内)
func (d *customConfTypeDao) InsertBatch(objs []bo.ConfTypeBo) (int64, error) {
	err := d.tx.Create(&objs).Error
	if err != nil {
		return 0, err
	}

	return int64(len(objs)), err
}

// Update 更新(事务内)
func (d *customConfTypeDao) Update(obj bo.ConfTypeBo) (*bo.ConfTypeBo, error) {
	// 只做字段值有值而不是默认值的更新
	err := d.tx.Model(&obj).Updates(obj).Error
	if err != nil {
		return nil, err
	}

	return &obj, err
}

// 扩展方法
func (d *customConfTypeDao) GetAllConfTypeList() ([]bo.ConfTypeBo, error) {
	var objs []bo.ConfTypeBo

	err := d.db.Where("valid = ?", cons.VALID_YES).Find(&objs).Error
	if err != nil {
		return nil, err
	}
	return objs, err
}

func (d *customConfTypeDao) GetConfType(param bo.ConfTypeBo) (*bo.ConfTypeBo, error) {
	var obj bo.ConfTypeBo

	dbm := d.db.Model(&bo.ConfTypeBo{})

	if param.ID != 0 {
		dbm = dbm.Where("id = ?", param.ID)
	}
	if param.TypeCode != "" {
		dbm = dbm.Where("type_code = ?", param.TypeCode)
	}
	if param.TypeName != "" {
		dbm = dbm.Where("type_name = ?", param.TypeName)
	}
	if param.Category != "" {
		dbm = dbm.Where("category = ?", param.Category)
	}
	dbm = dbm.Where("valid = ?", cons.VALID_YES)

	err := dbm.First(&obj).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &obj, err
}

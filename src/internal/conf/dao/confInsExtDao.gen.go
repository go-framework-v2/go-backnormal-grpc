package dao

import (
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/conf/model/bo"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// 设计要点：
// 接口组合实现扩展
// 结构体嵌入实现代码复用
type (
	// 基础接口 + 扩展方法
	ConfInsExtDao interface {
		// 基础接口，基于单表
		Truncate() error                          // 清空表并重置自增主键
		FindOne(id int) (*bo.ConfInsExtBo, error) // 根据id查询
		FindOneByUk()                             // 根据唯一键查询
		Insert(obj bo.ConfInsExtBo) (*bo.ConfInsExtBo, error)
		InsertBatch(objs []bo.ConfInsExtBo) (int64, error) // 500条以内
		Update(obj bo.ConfInsExtBo) (*bo.ConfInsExtBo, error)

		// 扩展接口
	}

	// 自定义模型实现
	customConfInsExtDao struct {
		// 默认实现
		db *gorm.DB // 不带事务的基础数据库连接
		tx *gorm.DB // 带事务的数据库连接

		// 扩展实现
		cache *redis.Client // 新增redis缓存
	}
)

// 对外暴露扩展接口
func NewConfInsExtDao(db *gorm.DB, tx *gorm.DB, cache *redis.Client) ConfInsExtDao {
	return &customConfInsExtDao{
		db:    db,
		tx:    tx,
		cache: cache,
	}
}

// 基础方法
// Truncate 清空表并重置自增主键
func (d *customConfInsExtDao) Truncate() error {
	// 临时禁用外键检查
	if err := d.tx.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return err
	}

	// 1. 清空表数据
	if err := d.tx.Exec("TRUNCATE TABLE " + "conf_ins_ext").Error; err != nil {
		return err
	}

	// 2. 重置自增主键为1
	if err := d.tx.Exec("ALTER TABLE " + "conf_ins_ext" + " AUTO_INCREMENT = 1").Error; err != nil {
		return err
	}

	return nil
}

// FindOne 根据id查询
func (d *customConfInsExtDao) FindOne(id int) (*bo.ConfInsExtBo, error) {
	var obj bo.ConfInsExtBo

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
func (d *customConfInsExtDao) FindOneByUk() {

}

// Insert 插入 (事务内)
func (d *customConfInsExtDao) Insert(obj bo.ConfInsExtBo) (*bo.ConfInsExtBo, error) {
	err := d.tx.Create(&obj).Error
	if err != nil {
		return nil, err
	}

	return &obj, err
}

// InsertBatch 插入 (事务内)
func (d *customConfInsExtDao) InsertBatch(objs []bo.ConfInsExtBo) (int64, error) {
	err := d.tx.Create(&objs).Error
	if err != nil {
		return 0, err
	}

	return int64(len(objs)), err
}

// Update 更新(事务内)
func (d *customConfInsExtDao) Update(obj bo.ConfInsExtBo) (*bo.ConfInsExtBo, error) {
	// 只做字段值有值而不是默认值的更新
	err := d.tx.Model(&obj).Updates(obj).Error
	if err != nil {
		return nil, err
	}

	return &obj, err
}

// 扩展方法

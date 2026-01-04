package dao

import (
	"github.com/go-framework-v2/go-backnormal-grpc/src/cons"
	"github.com/go-framework-v2/go-backnormal-grpc/src/internal/person/model/bo"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// 设计要点：
// 接口组合实现扩展
// 结构体嵌入实现代码复用
type (
	// 基础接口 + 扩展方法
	TestPersonDao interface {
		// 基础接口，基于单表
		Truncate() error                                           // 清空表并重置自增主键
		FindOne(id int) (*bo.TestPersonBo, error)                  // 根据id查询
		FindOneByUk(name, idcard string) (*bo.TestPersonBo, error) // 根据唯一键查询
		Insert(obj bo.TestPersonBo) (*bo.TestPersonBo, error)
		InsertBatch(objs []bo.TestPersonBo) (int64, error) // 500条以内
		Update(obj bo.TestPersonBo) (*bo.TestPersonBo, error)

		// 扩展接口
		GetTestPersonList(obj bo.TestPersonBo) ([]bo.TestPersonBo, error) // 条件查询
	}

	// 自定义模型实现
	customTestPersonDao struct {
		// 默认实现
		db *gorm.DB // 不带事务的基础数据库连接
		tx *gorm.DB // 带事务的数据库连接

		// 扩展实现
		cache *redis.Client // 新增redis缓存
	}
)

// 对外暴露扩展接口
func NewTestPersonDao(db *gorm.DB, tx *gorm.DB, cache *redis.Client) TestPersonDao {
	return &customTestPersonDao{
		db:    db,
		tx:    tx,
		cache: cache,
	}
}

// 基础方法
// Truncate 清空表并重置自增主键
func (d *customTestPersonDao) Truncate() error {
	// 临时禁用外键检查
	if err := d.tx.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return err
	}

	// 1. 清空表数据
	if err := d.tx.Exec("TRUNCATE TABLE " + "test_person").Error; err != nil {
		return err
	}

	// 2. 重置自增主键为1
	if err := d.tx.Exec("ALTER TABLE " + "test_person" + " AUTO_INCREMENT = 1").Error; err != nil {
		return err
	}

	return nil
}

// FindOne 根据id查询
func (d *customTestPersonDao) FindOne(id int) (*bo.TestPersonBo, error) {
	var obj bo.TestPersonBo

	err := d.db.Debug().First(&obj, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &obj, err
}

// FindOneByUk 根据唯一键查询
func (d *customTestPersonDao) FindOneByUk(name, idcard string) (*bo.TestPersonBo, error) {
	var obj bo.TestPersonBo

	err := d.db.Debug().First(&obj, "name = ? AND idcard = ? AND valid = ?", name, idcard, cons.VALID_YES).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &obj, err
}

// Insert 插入 (事务内)
func (d *customTestPersonDao) Insert(obj bo.TestPersonBo) (*bo.TestPersonBo, error) {
	err := d.tx.Debug().Create(&obj).Error
	if err != nil {
		return nil, err
	}

	return &obj, err
}

// InsertBatch 插入 (事务内)
func (d *customTestPersonDao) InsertBatch(objs []bo.TestPersonBo) (int64, error) {
	err := d.tx.Create(&objs).Error
	if err != nil {
		return 0, err
	}

	return int64(len(objs)), err
}

// Update 更新(事务内)
func (d *customTestPersonDao) Update(obj bo.TestPersonBo) (*bo.TestPersonBo, error) {
	// 只做字段值有值而不是默认值的更新
	err := d.tx.Model(&obj).Updates(obj).Error
	if err != nil {
		return nil, err
	}

	return &obj, err
}

// 扩展方法
func (d *customTestPersonDao) GetTestPersonList(obj bo.TestPersonBo) ([]bo.TestPersonBo, error) {
	var objs []bo.TestPersonBo

	dbm := d.db.Model(&bo.TestPersonBo{})

	// 条件查询
	if obj.Page > 0 && obj.PageSize > 0 {
		offset := (obj.Page - 1) * obj.PageSize
		dbm = dbm.Offset(offset).Limit(obj.PageSize)
	}

	// 最终查询
	err := dbm.Debug().Find(&objs).Error
	if err != nil {
		return nil, err
	}

	return objs, nil
}

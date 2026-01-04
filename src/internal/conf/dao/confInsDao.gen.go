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
	ConfInsDao interface {
		// 基础接口，基于单表
		Truncate() error                       // 清空表并重置自增主键
		FindOne(id int) (*bo.ConfInsBo, error) // 根据id查询
		FindOneByUk()                          // 根据唯一键查询
		Insert(obj bo.ConfInsBo) (*bo.ConfInsBo, error)
		InsertBatch(objs []bo.ConfInsBo) (int64, error) // 500条以内
		Update(obj bo.ConfInsBo) (*bo.ConfInsBo, error)

		// 扩展接口
		GetConfIns(param bo.ConfInsBo) (*bo.ConfInsBo, error) // 条件查询
	}

	// 自定义模型实现
	customConfInsDao struct {
		// 默认实现
		db *gorm.DB // 不带事务的基础数据库连接
		tx *gorm.DB // 带事务的数据库连接

		// 扩展实现
		cache *redis.Client // 新增redis缓存
	}
)

// 对外暴露扩展接口
func NewConfInsDao(db *gorm.DB, tx *gorm.DB, cache *redis.Client) ConfInsDao {
	return &customConfInsDao{
		db:    db,
		tx:    tx,
		cache: cache,
	}
}

// 基础方法
// Truncate 清空表并重置自增主键
func (d *customConfInsDao) Truncate() error {
	// 临时禁用外键检查
	if err := d.tx.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return err
	}

	// 1. 清空表数据
	if err := d.tx.Exec("TRUNCATE TABLE " + "conf_ins").Error; err != nil {
		return err
	}

	// 2. 重置自增主键为1
	if err := d.tx.Exec("ALTER TABLE " + "conf_ins" + " AUTO_INCREMENT = 1").Error; err != nil {
		return err
	}

	return nil
}

// FindOne 根据id查询
func (d *customConfInsDao) FindOne(id int) (*bo.ConfInsBo, error) {
	var obj bo.ConfInsBo

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
func (d *customConfInsDao) FindOneByUk() {

}

// Insert 插入 (事务内)
func (d *customConfInsDao) Insert(obj bo.ConfInsBo) (*bo.ConfInsBo, error) {
	err := d.tx.Create(&obj).Error
	if err != nil {
		return nil, err
	}

	return &obj, err
}

// InsertBatch 插入 (事务内)
func (d *customConfInsDao) InsertBatch(objs []bo.ConfInsBo) (int64, error) {
	err := d.tx.Create(&objs).Error
	if err != nil {
		return 0, err
	}

	return int64(len(objs)), err
}

// Update 更新(事务内)
func (d *customConfInsDao) Update(obj bo.ConfInsBo) (*bo.ConfInsBo, error) {
	// 只做字段值有值而不是默认值的更新
	err := d.tx.Model(&obj).Updates(obj).Error
	if err != nil {
		return nil, err
	}

	return &obj, err
}

// 扩展方法
func (d *customConfInsDao) GetConfIns(param bo.ConfInsBo) (*bo.ConfInsBo, error) {
	var obj bo.ConfInsBo

	// 构造条件查询
	query := d.db.Model(&bo.ConfInsBo{})
	if param.ID != 0 {
		query = query.Where("id = ?", param.ID)
	}
	if param.InsCode != "" {
		query = query.Where("ins_code = ?", param.InsCode)
	}
	if param.InsName != "" {
		query = query.Where("ins_name = ?", param.InsName)
	}
	if param.TypeID != 0 {
		query = query.Where("type_id = ?", param.TypeID)
	}

	err := query.First(&obj).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &obj, err
}

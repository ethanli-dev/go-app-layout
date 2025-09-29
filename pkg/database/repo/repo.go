/*
Copyright © 2025 lixw
*/
package repo

import "gorm.io/gorm"

type Cond struct {
	Query string
	Args  []any
}

type TxRepo interface {
	DB() *gorm.DB
	WithDB(db *gorm.DB) TxRepo
	WithTx(tx *gorm.DB, fn func(repo TxRepo) error) error
}

type CrudRepo[T any] interface {
	Create(entity *T) error
	GetByID(id interface{}) (*T, error)
	Update(entity *T) error
	Delete(id interface{}) error
	BatchCreate(entities []*T) error
	BatchUpdate(entities []*T, fields ...string) error
	BatchDelete(ids []interface{}) error
}

type QueryRepo[T any] interface {
	Find(conds ...Cond) (*T, error)
	FindAll(order string, conds ...Cond) ([]*T, error)
	FindPage(page, size int, order string, conds ...Cond) ([]*T, int64, error)
	Count(conds ...Cond) (int64, error)
}

type BaseRepo[T any] interface {
	TxRepo
	CrudRepo[T]
	QueryRepo[T]
}

type baseRepo[T any] struct {
	db *gorm.DB
}

func New[T any](db *gorm.DB) BaseRepo[T] {
	return &baseRepo[T]{
		db: db,
	}
}

// --- TxRepo ---
func (r *baseRepo[T]) DB() *gorm.DB {
	return r.db
}

func (r *baseRepo[T]) WithDB(db *gorm.DB) TxRepo {
	return &baseRepo[T]{db: db}
}

func (r *baseRepo[T]) WithTx(tx *gorm.DB, fn func(repo TxRepo) error) error {
	if tx != nil {
		return fn(r.WithDB(tx))
	}
	return r.db.Transaction(func(innerTx *gorm.DB) error {
		return fn(r.WithDB(innerTx))
	})
}

// --- CrudRepo ---
func (r *baseRepo[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *baseRepo[T]) GetByID(id interface{}) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, err
}

func (r *baseRepo[T]) Update(entity *T) error {
	if entity == nil {
		return gorm.ErrInvalidData
	}
	return r.db.Save(entity).Error
}

func (r *baseRepo[T]) Delete(id interface{}) error {
	var entity T
	return r.db.Delete(&entity, id).Error
}

func (r *baseRepo[T]) BatchCreate(entities []*T) error {
	if len(entities) == 0 {
		return nil
	}
	return r.db.Create(&entities).Error
}

func (r *baseRepo[T]) BatchUpdate(entities []*T, fields ...string) error {
	if len(entities) == 0 {
		return nil
	}
	if len(fields) == 0 {
		// 默认 Save 更新所有字段
		return r.db.Save(&entities).Error
	}
	return r.db.Select(fields).Save(&entities).Error
}

func (r *baseRepo[T]) BatchDelete(ids []interface{}) error {
	if len(ids) == 0 {
		return nil
	}
	var entity T
	return r.db.Delete(&entity, ids).Error
}

// --- QueryRepo ---
func (r *baseRepo[T]) Find(conds ...Cond) (*T, error) {
	var entity T
	tx := r.applyConds(r.db.Model(&entity), conds)
	err := tx.First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, err
}

func (r *baseRepo[T]) FindAll(order string, conds ...Cond) ([]*T, error) {
	var entities []*T
	tx := r.applyConds(r.db.Model(&entities), conds)
	if order != "" {
		tx = tx.Order(order)
	}
	err := tx.Find(&entities).Error
	return entities, err
}

func (r *baseRepo[T]) FindPage(page, size int, order string, conds ...Cond) ([]*T, int64, error) {
	var (
		entities []*T
		count    int64
	)
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10 // 默认分页大小
	}
	tx := r.applyConds(r.db.Model(&entities), conds)
	tx.Count(&count)

	if order != "" {
		tx = tx.Order(order)
	}
	offset := (page - 1) * size
	err := tx.Offset(offset).Limit(size).Find(&entities).Error
	return entities, count, err
}

func (r *baseRepo[T]) Count(conds ...Cond) (int64, error) {
	var count int64
	var entity T
	err := r.applyConds(r.db.Model(&entity), conds).Count(&count).Error
	return count, err
}

// --- 内部方法 ---
func (r *baseRepo[T]) applyConds(tx *gorm.DB, conds []Cond) *gorm.DB {
	for _, c := range conds {
		tx = tx.Where(c.Query, c.Args...)
	}
	return tx
}

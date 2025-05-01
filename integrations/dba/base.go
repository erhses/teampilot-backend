package dba

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func GetModel(db *gorm.DB, model interface{}) (*QueryBuilder, string) {
	db = db.Model(model)
	s, _ := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	return &QueryBuilder{
		DB: db,
	}, s.Table
}

// InTransactionFind is a function to find data in transaction
// The function InTransactionFind takes two parameters:
// tx, which is a GORM DB instance representing a transaction,
// and dest, which is the destination object where the query results will be stored.
// The function uses a deferred function
// to handle any panics that may occur during the execution of the function.
// If a panic occurs, the transaction is rolled back using tx.Rollback().
// If no panic occurs, the transaction is committed using tx.Commit().
// The function executes the Find method on the transaction tx
// with the dest object as the destination to store the query results.
// If an error occurs during the query execution,
// the transaction is rolled back and the function returns the transaction tx.
// If the query is successful and no errors occur, the transaction is committed using tx.Commit().
// Finally, the function returns the transaction tx, whether it has been committed or rolled back.
func InTransactionFind(tx *gorm.DB, dest interface{}) *gorm.DB {
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Find(dest).Error; err != nil {
		tx.Rollback()
		return tx
	}
	tx.Commit()
	return tx
}

func (db *QueryBuilder) FilterByBase(tableName string, params BaseFilter) *QueryBuilder {
	db = db.BetweenDate(fmt.Sprintf(`"%s".created_at`, tableName), params.CreatedAt)
	db = db.BetweenDate(fmt.Sprintf(`"%s".deleted_at`, tableName), params.DeletedAt)
	return db
}

func (db *QueryBuilder) FilterByBaseModifier(tableName string, params BaseModifierFilter) *QueryBuilder {
	db = db.BetweenDate(fmt.Sprintf(`"%s".updated_at`, tableName), params.UpdatedAt)
	db = db.Equal(fmt.Sprintf(`"%s".modifier_id`, tableName), params.ModifierID)
	return db
}

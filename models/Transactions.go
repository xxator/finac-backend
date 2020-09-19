package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Transaction struct {
	ID                uint64    `gorm:"primary_key;auto_increment" json:"id"`
	DateOfTransaction time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date_of_transaction"`
	Amount            uint32    `gorm:"not null;" json:"amount"`
	Author            User      `json:"author"`
	AuthorID          uint32    `gorm:"not null" json:"author_id"`
	CreatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (t *Transaction) Prepare() {
	t.ID = 0
	t.DateOfTransaction = time.Now()
	t.Author = User{}
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	t.Amount = 0
}

func (t *Transaction) Validate() error {

	if t.Amount == 0 {
		return errors.New("Required Amount")
	}
	if t.AuthorID < 1 {
		return errors.New("Required Author")
	}
	return nil
}

func (t *Transaction) SavePost(db *gorm.DB) (*Transaction, error) {
	var err error
	err = db.Debug().Model(&Transaction{}).Create(&t).Error
	if err != nil {
		return &Transaction{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.AuthorID).Take(&t.Author).Error
		if err != nil {
			return &Transaction{}, err
		}
	}
	return t, nil
}

func (t *Transaction) FindAllTransactions(db *gorm.DB) (*[]Transaction, error) {
	var err error
	transactions := []Transaction{}
	err = db.Debug().Model(&Transaction{}).Limit(100).Find(&transactions).Error
	if err != nil {
		return &[]Transaction{}, err
	}
	if len(transactions) > 0 {
		for i, _ := range transactions {
			err := db.Debug().Model(&User{}).Where("id = ?", transactions[i].AuthorID).Take(&transactions[i].Author).Error
			if err != nil {
				return &[]Transaction{}, err
			}
		}
	}
	return &transactions, nil
}

func (t *Transaction) FindTransactionByID(db *gorm.DB, pid uint64) (*Transaction, error) {
	var err error
	err = db.Debug().Model(&Transaction{}).Where("id = ?", pid).Take(&t).Error
	if err != nil {
		return &Transaction{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.AuthorID).Take(&t.Author).Error
		if err != nil {
			return &Transaction{}, err
		}
	}
	return t, nil
}

func (t *Transaction) UpdateATransaction(db *gorm.DB) (*Transaction, error) {

	var err error

	err = db.Debug().Model(&Transaction{}).Where("id = ?", t.ID).Updates(Transaction{DateOfTransaction: t.DateOfTransaction, Amount: t.Amount, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Transaction{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.AuthorID).Take(&t.Author).Error
		if err != nil {
			return &Transaction{}, err
		}
	}
	return t, nil
}

func (t *Transaction) DeleteATransaction(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Transaction{}).Where("id = ? and author_id = ?", pid, uid).Take(&Transaction{}).Delete(&Transaction{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Transaction not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

package mysql

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// BaseModel provides common ULID identifier handling for all entities.
type BaseModel struct {
	ID string `json:"id" gorm:"type:char(26);primaryKey"`
}

// SetID assigns the ULID to the model.
func (b *BaseModel) SetID(id string) {
	b.ID = id
}

// GetID returns the ULID of the model.
func (b *BaseModel) GetID() string {
	return b.ID
}

// Entity describes models that expose ULID identifiers.
type Entity interface {
	GetID() string
	SetID(string)
}

type User struct {
	BaseModel
	Name         string `json:"name" gorm:"size:255;not null"`
	Email        string `json:"email" gorm:"size:255;uniqueIndex;not null"`
	Password     string `json:"password,omitempty" gorm:"-"`
	PasswordHash string `json:"-" gorm:"column:password;size:255;not null"`
}

type Company struct {
	BaseModel
	Name string `json:"name" gorm:"size:255;not null"`
}

type RetailPoint struct {
	BaseModel
	CompanyID string `json:"company_id" gorm:"type:char(26);not null"`
	Name      string `json:"name" gorm:"size:255;not null"`
	Address   string `json:"address" gorm:"size:512"`
}

type Brand struct {
	BaseModel
	Name string `json:"name" gorm:"size:255;not null"`
}

type Category struct {
	BaseModel
	Name     string  `json:"name" gorm:"size:255;not null"`
	ParentID *string `json:"parent_id,omitempty" gorm:"type:char(26);"`
}

type Product struct {
	BaseModel
	Name       string `json:"name" gorm:"size:255;not null"`
	SKU        string `json:"sku" gorm:"size:64;uniqueIndex;not null"`
	BrandID    string `json:"brand_id" gorm:"type:char(26);not null"`
	CategoryID string `json:"category_id" gorm:"type:char(26);not null"`
}

type Visit struct {
	BaseModel
	UserID        string    `json:"user_id" gorm:"type:char(26);not null"`
	RetailPointID string    `json:"retail_point_id" gorm:"type:char(26);not null"`
	VisitedAt     time.Time `json:"visited_at"`
	Notes         string    `json:"notes" gorm:"size:1024"`
}

type VisitItem struct {
	BaseModel
	VisitID         string  `json:"visit_id" gorm:"type:char(26);not null"`
	ProductID       string  `json:"product_id" gorm:"type:char(26);not null"`
	PresentQuantity int     `json:"present_quantity"`
	StoreQuantity   int     `json:"store_quantity"`
	Price           float64 `json:"price"`
}

type UserToken struct {
	BaseModel
	UserID    string     `json:"user_id" gorm:"type:char(26);not null"`
	Token     string     `json:"token" gorm:"size:64;uniqueIndex;not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// BeforeSave hashes the password if a plain-text password has been provided.
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.PasswordHash = string(hashed)
		u.Password = ""
	}
	return nil
}

// CheckPassword verifies the provided password against the stored hash.
func (u *User) CheckPassword(plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plain))
}

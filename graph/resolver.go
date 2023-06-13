//go:generate go run github.com/99designs/gqlgen generate
package graph

import (
	"github.com/jinzhu/gorm"
)

type Resolver struct {
	Database *gorm.DB
}

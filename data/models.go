package data

import (
	"github.com/jinzhu/gorm"
)

type Topology struct {
	gorm.Model
	Name     string   `gorm:"not null;unique"`
	Active   bool     `gorm:"default:true"`
	Entities []Entity `gorm:"foreignkey:TopologyID"`
}

type Entity struct {
	gorm.Model
	Name       string `gorm:"not null;unique"`
	Active     bool   `gorm:"default:true"`
	TopologyID uint   `sql:"type:int REFERENCES topologies(id)"`
}

type Component struct {
	gorm.Model
	Name   string `gorm:"not null;unique"`
	Active bool   `gorm:"default:true"`
}

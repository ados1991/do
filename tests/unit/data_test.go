package unit

import (
	"app/data"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

var Db *gorm.DB

func Setup() func() {
	var err error
	Db, err = gorm.Open("sqlite3", "test.db")
	Db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("setup database for test")
	fmt.Println("migration")
	Db.AutoMigrate(&data.Topology{}, &data.Entity{}, &data.Component{})
	return func() {
		fmt.Println("teardown database for test")
		err_db := Db.Close()
		if err_db != nil {
			panic("failed to close database")
		}
		err_file := os.Remove("test.db")
		if err_file != nil {
			fmt.Printf("err_file=%v", err_file)
			panic("failed to close file")
		}
		return
	}
}

func TestCreateUpdateDeleteTopology(t *testing.T) {
	topo := data.Topology{Name: "test", Active: true}
	Db.Create(&topo)
	if Db.First(&topo).RecordNotFound() {
		t.Error(fmt.Sprintf("topo %v doesn't exist", topo))
	}
	Db.First(&topo)
	topo.Name = "test2"
	Db.Save(&topo)
	record_not_found := Db.Where("Name = ?", "test2").First(&data.Topology{}).RecordNotFound()
	if record_not_found {
		t.Error(fmt.Sprintf("topo %v has not updated", topo))
	}
	Db.Delete(&topo)
	record_not_deleted := Db.Where("Name = ?", "test2").First(&data.Topology{}).RecordNotFound()
	if !record_not_deleted {
		t.Error(fmt.Sprintf("topo %v has not deleted", topo))
	}
}

func TestCreateTwoTopologiesWithSameName(t *testing.T) {
	topo := data.Topology{Name: "test", Active: true}
	topo2 := data.Topology{Name: "test", Active: true}
	Db.Save(&topo)
	db := Db.Save(&topo2)
	if db.Error == nil {
		t.Error("topo2 must not be created")
	}
}

func TestCreateEntityWithoutTopology(t *testing.T) {
	entity := data.Entity{Name: "test_entity", Active: true}
	db := Db.Save(&entity)
	if db.Error == nil {
		t.Error("Entity can not be created without associated topology")
	}
}

func TestCreateEntityWithTopology(t *testing.T) {
	topo := data.Topology{Name: "topo_with_entity", Active: true}
	Db.Save(&topo)
	entity := data.Entity{Name: "test_entity", Active: true, TopologyID: topo.ID}
	db := Db.Save(&entity)
	if db.Error != nil {
		t.Error(fmt.Sprintf("Entity cannot be created with topo's id==%d", topo.ID))
	}
}

func TestGetTopologyWithItsEntities(t *testing.T) {
	topo := data.Topology{Name: "topo_with_entities", Active: true}
	Db.Save(&topo)
	for i := 1; i <= 2; i++ {
		Db.Save(&data.Entity{Name: fmt.Sprintf("entity_%d", i), Active: true, TopologyID: topo.ID})
	}
	Db.Preload("Entities").First(&topo)
	if len(topo.Entities) != 2 {
		t.Error(fmt.Sprintf("topo %d must have 2 entities", topo.ID))
	}
}

func TestMain(m *testing.M) {
	teardown := Setup()
	code := m.Run()
	fmt.Println("teardown()")
	teardown()
	os.Exit(code)
}

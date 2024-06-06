package main

import (
	"database/sql"
	"errors"
	"log"
	"testing"
	"time"

	data "cloth-service/internal/data" // Update with the correct path

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("postgres", "postgres://nurtileu:root@localhost/a.nurtileuDB?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

func tearDownTestDB(db *sql.DB) {
	db.Exec("TRUNCATE TABLE clothes_info RESTART IDENTITY CASCADE")
	db.Close()
}

type ClotheInfo struct {
	ID        int64
	ClothName string
	ClothCost float64
	ClothSize string
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

var ErrRecordNotFound = errors.New("record not found")

type Filters struct {
	Page     int
	PageSize int
	Sort     string
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

func (f Filters) sortColumn() string {
	return f.Sort
}

func (f Filters) sortDirection() string {
	return "ASC"
}
func TestInsertClothe(t *testing.T) {
	db := setupTestDB()
	defer tearDownTestDB(db)
	model := data.NewModels(db)

	clothe := &data.ClotheInfo{
		ClothName: "Test Cloth",
		ClothCost: 12,
		ClothSize: "M",
	}

	err := model.ClotheInfo.Insert(clothe)
	assert.NoError(t, err)
	assert.NotZero(t, clothe.ID)
	assert.WithinDuration(t, time.Now(), clothe.CreatedAt, 2*time.Second)
	assert.WithinDuration(t, time.Now(), clothe.UpdatedAt, 2*time.Second)
}

func TestGetClothe(t *testing.T) {
	db := setupTestDB()
	defer tearDownTestDB(db)
	model := data.NewModels(db)

	clothe := &data.ClotheInfo{
		ClothName: "Test Cloth",
		ClothCost: 15,
		ClothSize: "L",
	}
	explicitInsertion(t, model, clothe)

	fetchedClothe, err := model.ClotheInfo.Get(clothe.ID)
	assert.NoError(t, err)
	assert.Equal(t, clothe.ID, fetchedClothe.ID)
	assert.Equal(t, clothe.ClothName, fetchedClothe.ClothName)
	assert.Equal(t, clothe.ClothCost, fetchedClothe.ClothCost)
	assert.Equal(t, clothe.ClothSize, fetchedClothe.ClothSize)
}

func TestUpdateClothe(t *testing.T) {
	db := setupTestDB()
	defer tearDownTestDB(db)
	model := data.NewModels(db)

	clothe := &data.ClotheInfo{
		ClothName: "Test Cloth",
		ClothCost: 15.99,
		ClothSize: "L",
	}
	explicitInsertion(t, model, clothe)

	clothe.ClothCost = 17.50
	err := model.ClotheInfo.Update(clothe)
	assert.NoError(t, err)

	updatedClothe, err := model.ClotheInfo.Get(clothe.ID)
	assert.NoError(t, err)
	assert.Equal(t, 17.50, updatedClothe.ClothCost)
}

func TestDeleteClothe(t *testing.T) {
	db := setupTestDB()
	defer tearDownTestDB(db)
	model := data.NewModels(db)

	clothe := &data.ClotheInfo{
		ClothName: "Test Cloth",
		ClothCost: 15.99,
		ClothSize: "L",
	}
	explicitInsertion(t, model, clothe)

	err := model.ClotheInfo
}
func TestGetAllClothes(t *testing.T) {
	db := setupTestDB()
	defer tearDownTestDB(db)
	model := data.NewModels(db)

	clothe1 := &data.ClotheInfo{
		ClothName: "Test Cloth 1",
		ClothCost: 12,
		ClothSize: "M",
	}
	clothe2 := &data.ClotheInfo{
		ClothName: "Test Cloth 2",
		ClothCost: 13,
		ClothSize: "L",
	}

	explicitInsertion(t, model, clothe1)
	explicitInsertion(t, model, clothe2)

	filters := data.Filters{
		Page:     1,
		PageSize: 10,
		Sort:     "id",
	}

	clothes, err := model.ClotheInfo.GetAll("", "", filters)
	assert.NoError(t, err)
	assert.Len(t, clothes, 2)
	assert.Equal(t, clothe1.ClothName, clothes[0].ClothName)
	assert.Equal(t, clothe2.ClothName, clothes[1].ClothName)
}

func explicitInsertion(t *testing.T, model data.Models, clothe *data.ClotheInfo) {
	err := model.ClotheInfo.Insert(clothe)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}
}

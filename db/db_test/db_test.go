package db_test

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"

	. "github.com/sunreaver/antman/v3/db"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type TempTable struct {
	gorm.Model
	Name     string    `gorm:"column:name;size:10;not null;index:idx_name"`
	State    uint8     `gorm:"column:state;not null;index:idx_state"`
	TestJSON *TestJSON `gorm:"column:t_json;type:json;not null"`
}

func (t *TempTable) TableName() string {
	return "tmp_table"
}

type TestJSON struct {
	V1        uint8      `json:"v1,omitempty"`
	V2        string     `json:"v2,omitempty"`
	TestArray []TestJSON `json:"vs,omitempty"` // 子节点
}

// Value 实现方法.
func (t *TestJSON) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// Scan 实现方法.
func (t *TestJSON) Scan(input interface{}) error {
	bytes, ok := input.([]byte)
	if !ok {
		return fmt.Errorf("input is not []byte: %v", input)
	}

	result := TestJSON{}
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return err
	}
	*t = result

	return nil
}

func (t *TestJSON) GormDataType() string {
	return "json"
}

func (t *TestJSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	// use field.Tag, field.TagSettings gets field's tags
	// checkout https://github.com/go-gorm/gorm/blob/master/schema/field.go for all options

	// returns different database type based on driver name
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "JSON"
	case "postgres":
		return "JSONB"
	}

	return ""
}

var dbTest *Databases

func TestMain(m *testing.M) {
	dbFile := path.Join(os.TempDir(), "tmp.db")
	os.Remove(dbFile)
	cfg := Config{
		Type:      "sqlite",
		MasterURI: dbFile,
		SlaveURIs: []string{
			dbFile,
			dbFile,
		},
		LogMode:      false,
		MaxIdleConns: 1,
		MaxOpenConns: 10,
	}
	tmp, e := MakeSqliteClient(cfg, nil)
	if e != nil {
		fmt.Println("make db err:", e)
		os.Exit(1)
	}
	dbTest = tmp
	m.Run()
	dbTest.Free()
	os.Remove(dbFile)
	os.Exit(0)
}

func TestRun(t *testing.T) {
	e := dbTest.Master().AutoMigrate(&TempTable{})
	if e != nil {
		t.Fatalf("auto migrate: %v", e)
	}

	tmp1 := TempTable{
		Name:  "tmp1",
		State: 1,
		TestJSON: &TestJSON{
			V1: 0,
			V2: "0",
			TestArray: []TestJSON{
				{
					V1:        1,
					V2:        "1",
					TestArray: nil,
				},
			},
		},
	}

	e = dbTest.Master().Create(&tmp1).Error
	if e != nil {
		t.Fatalf("create: %v", e)
	}

	tmp2 := TempTable{
		Name:  "tmp2",
		State: 2,
		TestJSON: &TestJSON{
			V1: 0,
			V2: "0",
			TestArray: []TestJSON{
				{
					V1:        1,
					V2:        "1",
					TestArray: nil,
				},
			},
		},
	}

	e = dbTest.Master().Create(&tmp2).Error
	if e != nil {
		t.Fatalf("create: %v", e)
	}

	var out []TempTable
	e = dbTest.Slave().Table("tmp_table").Where("name = ?", "tmp1").Find(&out).Error

	if e != nil {
		t.Errorf("find tmp1: %v", e)
	}

	if len(out) != 1 {
		t.Errorf("select len is not 1: %v", len(out))
	}
	tmp1Data, _ := json.Marshal(tmp1)
	out1Data, _ := json.Marshal(out[0])

	if string(tmp1Data) != string(out1Data) {
		t.Errorf("tmp1 save and check out not equal: \ntmp: %v\nout: %v\n", string(tmp1Data), string(out1Data))
	}

	e = dbTest.Slave().Table("tmp_table").Where("state in (?)", []uint8{1, 2}).Order("id").Find(&out).Error

	if e != nil {
		t.Errorf("find state in (1,2): %v", e)
	}

	if len(out) != 2 {
		t.Errorf("select len is not 2: %v", len(out))
	}
	tmp2Data, _ := json.Marshal(tmp2)
	out1Data, _ = json.Marshal(out[0])
	out2Data, _ := json.Marshal(out[1])
	if string(tmp1Data) != string(out1Data) {
		t.Errorf("tmp1 save and check out not equal: \ntmp: %v\nout: %v\n", string(tmp1Data), string(out1Data))
	}
	if string(tmp2Data) != string(out2Data) {
		t.Errorf("tmp2 save and check out not equal: \ntmp: %v\nout: %v\n", string(tmp2Data), string(out1Data))
	}
}

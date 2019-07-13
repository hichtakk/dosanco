package model

/*
//DataCenter
type DataCenter struct {
	gorm.Model
	Name    string  `gorm:"type:varchar(10);unique_index" json:"name"`
	Address string  `gorm:"type:varchar(255)" json:"address"`
	Floors  []Floor `json:"floors"`
}

type Floor struct {
	gorm.Model
	Name    string
	Address string
	Halls   []Hall
}

type Hall struct {
	gorm.Model
	Name    string
	Type    string // Data or Network
	RackRow []RackRow
}

type RackRow struct {
	gorm.Model
	Name  string
	Racks []Rack
}

type Rack struct {
	gorm.Model
	Name     string
	RackPDUs []RackPDU
}

type UPS struct {
	gorm.Model
	Name string
}

type PDU struct {
	gorm.Model
	Name         string
	Address      string
	PrimaryUPS   UPS
	SecondaryUPS UPS
}

type RackPDU struct {
	gorm.Model
	Name       string
	SourcePDUs []PDU
}
*/

/*
func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&DataCenter{})
	db.AutoMigrate(&Floor{})
	db.AutoMigrate(&Hall{})
	db.AutoMigrate(&RackRow{})
	db.AutoMigrate(&Rack{})
	db.AutoMigrate(&UPS{})
	db.AutoMigrate(&PDU{})
	db.AutoMigrate(&RackPDU{})

	db.Create(&DataCenter{Name: "APL", Address: "207 N United Sakura Drive, East Wenatchee, WA 98802"})

	// Read
	var dc DataCenter
	db.First(&dc, 1)                 // find product with id 1
	db.First(&dc, "name = ?", "APL") // find product with code l1212

	// Update - update product's price to 2000
	//db.Model(&product).Update("Price", 2000)

	// Delete - delete product
	//db.Delete(&product)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello")
	})
	e.Logger.Fatal(e.Start(":8080"))
}
*/

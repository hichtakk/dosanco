package model

// DataCenter represents datacenter building data.
type DataCenter struct {
	Model
	Name    string  `gorm:"type:varchar(10);unique_index" json:"name"`
	Address string  `gorm:"type:varchar(255)" json:"address"`
	Floors  []Floor `json:"floors,omitempty"`
}

// Floor represents datacenter floor or area.
type Floor struct {
	Model
	Name  string `gorm:"type:varchar(16);unique_index" json:"name"`
	Halls []Hall `json:"halls,omitempty"`
}

// Hall represents data hall in datacenter.
type Hall struct {
	Model
	Name     string    `gorm:"type:varchar(16);unique_index" json:"name"`
	Type     string    `gorm:"type:varchar(10)" json:"type"`
	RackRows []RackRow `json:"rows,omitempty"`
}

// RackRow represents row of racks in data hall.
type RackRow struct {
	Model
	Name  string `gorm:"type:varchar(16);unique_index" json:"name"`
	Racks []Rack `json:"racks,omitempty"`
}

// Rack represents each rack in row.
type Rack struct {
	Model
	Name        string    `gorm:"type:varchar(16);unique_index" json:"name"`
	RackPDUs    []RackPDU `json:"rack_pdus,omitempty"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
}

// UPS represents redundant power source
type UPS struct {
	Model
	Name        string `gorm:"type:varchar(16);unique_index" json:"name"`
	Description string `gorm:"type:varchar(255)" json:"description"`
}

// PDU represents power distribution unit on data hall
type PDU struct {
	Model
	Name         string `gorm:"type:varchar(16);unique_index" json:"name"`
	PrimaryUPS   UPS    `json:"primary_ups,omitempty"`
	SecondaryUPS UPS    `json:"secondary_ups,omitempty"`
	Description  string `gorm:"type:varchar(255)" json:"description"`
}

// RackPDU represents power distribution unit installed inside of rack
type RackPDU struct {
	Model
	Name        string `gorm:"type:varchar(16);unique_index" json:"name"`
	Address     string `gorm:"type:varchar(15)" json:"address"`
	Description string `gorm:"type:varchar(255)" json:"description"`
	SourcePDUs  []PDU  `json:"source_pdus,omitempty"`
}

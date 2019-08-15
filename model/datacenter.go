package model

//DataCenter
type DataCenter struct {
	Model
	Name    string  `gorm:"type:varchar(10);unique_index" json:"name"`
	Address string  `gorm:"type:varchar(255)" json:"address"`
	Floors  []Floor `json:"floors,omitempty"`
}

type Floor struct {
	Model
	Name  string `gorm:"type:varchar(16);unique_index" json:"name"`
	Halls []Hall `json:"halls,omitempty"`
}

type Hall struct {
	Model
	Name     string    `gorm:"type:varchar(16);unique_index" json:"name"`
	Type     string    `gorm:"type:varchar(10)" json:"type"`
	RackRows []RackRow `json:"rows,omitempty"`
}

type RackRow struct {
	Model
	Name  string `gorm:"type:varchar(16);unique_index" json:"name"`
	Racks []Rack `json:"racks,omitempty"`
}

type Rack struct {
	Model
	Name        string    `gorm:"type:varchar(16);unique_index" json:"name"`
	RackPDUs    []RackPDU `json:"rack_pdus,omitempty"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
}

type UPS struct {
	Model
	Name        string `gorm:"type:varchar(16);unique_index" json:"name"`
	Description string `gorm:"type:varchar(255)" json:"description"`
}

type PDU struct {
	Model
	Name         string `gorm:"type:varchar(16);unique_index" json:"name"`
	PrimaryUPS   UPS    `json:"primary_ups,omitempty"`
	SecondaryUPS UPS    `json:"secondary_ups,omitempty"`
	Description  string `gorm:"type:varchar(255)" json:"description"`
}

type RackPDU struct {
	Model
	Name        string `gorm:"type:varchar(16);unique_index" json:"name"`
	Address     string `gorm:"type:varchar(15)" json:"address"`
	Description string `gorm:"type:varchar(255)" json:"description"`
	SourcePDUs  []PDU  `json:"source_pdus,omitempty"`
}

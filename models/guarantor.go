package models

type Guarantor struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	CreatedAt  *LocalTime `json:"created_at"`
	UpdatedAt  *LocalTime `json:"updated_at"`
	DeletedAt  *LocalTime `gorm:"index" json:"deleted_at"`
	LoanID     int32      `json:"loan_id" gorm:"type:int"` // Foreign Key
	SyncStatus string     `json:"sync_status" gorm:"default:'WAITING'"`

	// Personal Info
	GuarantorType string `json:"guarantor_type" gorm:"default:'individual'"` // individual, juristic
	Title         string `json:"title"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	// Juristic Fields
	TradeRegistrationID string  `json:"trade_registration_id"`
	RegistrationDate    *string `json:"registration_date" gorm:"type:date;default:null"`
	TaxID               string  `json:"tax_id"`
	Gender              string  `json:"gender"`
	IdCard              string  `json:"id_card" gorm:"type:varchar(20)"`
	IdCardIssueDate     *string `json:"id_card_issue_date" gorm:"type:date;default:null"`
	IdCardExpiryDate    *string `json:"id_card_expiry_date" gorm:"type:date;default:null"`
	DateOfBirth         *string `json:"date_of_birth" gorm:"type:date;default:null"`
	Ethnicity           string  `json:"ethnicity"`
	Nationality         string  `json:"nationality"`
	Religion            string  `json:"religion"`
	MaritalStatus       string  `json:"marital_status"`
	MobilePhone         string  `json:"mobile_phone"`

	// Address - House Registration
	HouseRegNo          string `json:"house_reg_no"`
	HouseRegBuilding    string `json:"house_reg_building"`
	HouseRegFloor       string `json:"house_reg_floor"`
	HouseRegRoom        string `json:"house_reg_room"`
	HouseRegMoo         string `json:"house_reg_moo"`
	HouseRegSoi         string `json:"house_reg_soi"`
	HouseRegRoad        string `json:"house_reg_road"`
	HouseRegProvince    string `json:"house_reg_province"`
	HouseRegDistrict    string `json:"house_reg_district"`
	HouseRegSubdistrict string `json:"house_reg_subdistrict"`
	HouseRegZipcode     string `json:"house_reg_zipcode"`

	// Address - Current
	SameAsHouseReg     bool   `json:"same_as_house_reg"`
	CurrentNo          string `json:"current_no"`
	CurrentBuilding    string `json:"current_building"`
	CurrentFloor       string `json:"current_floor"`
	CurrentRoom        string `json:"current_room"`
	CurrentMoo         string `json:"current_moo"`
	CurrentSoi         string `json:"current_soi"`
	CurrentRoad        string `json:"current_road"`
	CurrentProvince    string `json:"current_province"`
	CurrentDistrict    string `json:"current_district"`
	CurrentSubdistrict string `json:"current_subdistrict"`
	CurrentZipcode     string `json:"current_zipcode"`

	// Work Info
	CompanyName     string  `json:"company_name"`
	Occupation      string  `json:"occupation"`
	Position        string  `json:"position"`
	Salary          float64 `json:"salary"`
	OtherIncome     float64 `json:"other_income"`
	IncomeSource    string  `json:"income_source"`
	WorkPhone       string  `json:"work_phone"`
	WorkNo          string  `json:"work_no"`
	WorkBuilding    string  `json:"work_building"`
	WorkFloor       string  `json:"work_floor"`
	WorkRoom        string  `json:"work_room"`
	WorkMoo         string  `json:"work_moo"`
	WorkSoi         string  `json:"work_soi"`
	WorkRoad        string  `json:"work_road"`
	WorkProvince    string  `json:"work_province"`
	WorkDistrict    string  `json:"work_district"`
	WorkSubdistrict string  `json:"work_subdistrict"`
	WorkZipcode     string  `json:"work_zipcode"`

	// Other Card Info
	OtherCardType       string  `json:"other_card_type"`
	OtherCardNumber     string  `json:"other_card_number"`
	OtherCardIssueDate  *string `json:"other_card_issue_date" gorm:"type:date;default:null"`
	OtherCardExpiryDate *string `json:"other_card_expiry_date" gorm:"type:date;default:null"`

	// Address - Document Delivery
	DocDeliveryType string `json:"doc_delivery_type"`
	DocNo           string `json:"doc_no"`
	DocBuilding     string `json:"doc_building"`
	DocFloor        string `json:"doc_floor"`
	DocRoom         string `json:"doc_room"`
	DocMoo          string `json:"doc_moo"`
	DocSoi          string `json:"doc_soi"`
	DocRoad         string `json:"doc_road"`
	DocProvince     string `json:"doc_province"`
	DocDistrict     string `json:"doc_district"`
	DocSubdistrict  string `json:"doc_subdistrict"`
	DocZipcode      string `json:"doc_zipcode"`

	// Relationship to borrower
	RelationshipHelper string `json:"relationship_helper"`
}

func (Guarantor) TableName() string {
	return "loan_applications_guarantors"
}

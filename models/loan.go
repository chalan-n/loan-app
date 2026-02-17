package models

type LoanApplication struct {
	ID int `gorm:"primaryKey"`

	// Step 1: Borrower Info
	RefCode             string  `json:"ref_code" gorm:"unique"`
	Status              string  `json:"status" gorm:"default:'D'"` // Draft, Pending, Approved, Rejected, Conditional
	SubmittedDate       string  `json:"submitted_date"`
	LastUpdateDate      string  `json:"last_update_date"`
	SyncStatus          string  `json:"sync_status" gorm:"default:'WAITING'"` // WAITING, PENDING, SYNCED
	StaffID             string  `json:"staff_id"`
	Title               string  `json:"title"`
	FirstName           string  `json:"first_name"`
	LastName            string  `json:"last_name"`
	BorrowerType        string  `json:"borrower_type" gorm:"default:'individual'"` // individual, juristic
	TradeRegistrationID string  `json:"trade_registration_id"`                     // Juristic
	RegistrationDate    *string `json:"registration_date" gorm:"type:date"`        // Juristic
	TaxID               string  `json:"tax_id"`                                    // Juristic
	Gender              string  `json:"gender"`
	IdCard              string  `json:"id_card" gorm:"type:varchar(20)"`
	IdCardIssueDate     *string `json:"id_card_issue_date" gorm:"type:date"`
	IdCardExpiryDate    *string `json:"id_card_expiry_date" gorm:"type:date"`
	DateOfBirth         *string `json:"date_of_birth" gorm:"type:date"`
	Ethnicity           string  `json:"ethnicity"`
	Nationality         string  `json:"nationality"`
	Religion            string  `json:"religion"`
	MaritalStatus       string  `json:"marital_status"`
	MobilePhone         string  `json:"mobile_phone"`
	OtherCardType       string  `json:"other_card_type"`
	OtherCardNumber     string  `json:"other_card_number"`
	OtherCardIssueDate  *string `json:"other_card_issue_date" gorm:"type:date"`
	OtherCardExpiryDate *string `json:"other_card_expiry_date" gorm:"type:date"`
	CompanyName         string  `json:"company_name"`
	Occupation          string  `json:"occupation"`
	Position            string  `json:"position"`
	Salary              float64 `json:"salary"`
	OtherIncome         float64 `json:"other_income"`
	IncomeSource        string  `json:"income_source"`
	CreditBureauStatus  string  `json:"credit_bureau_status"`

	// --- New Address Fields (Step 1) ---

	// House Registration Address
	HouseRegNo          string `json:"house_reg_no"`
	HouseRegBuilding    string `json:"house_reg_building"` // New
	HouseRegFloor       string `json:"house_reg_floor"`    // New
	HouseRegRoom        string `json:"house_reg_room"`     // New
	HouseRegMoo         string `json:"house_reg_moo"`
	HouseRegSoi         string `json:"house_reg_soi"`
	HouseRegRoad        string `json:"house_reg_road"`
	HouseRegProvince    string `json:"house_reg_province"`
	HouseRegDistrict    string `json:"house_reg_district"`
	HouseRegSubdistrict string `json:"house_reg_subdistrict"`
	HouseRegZipcode     string `json:"house_reg_zipcode"`
	SameAsHouseReg      bool   `json:"same_as_house_reg"`

	// Current Address
	CurrentCompany     string `json:"current_company"` // New
	CurrentNo          string `json:"current_no"`
	CurrentBuilding    string `json:"current_building"` // New
	CurrentFloor       string `json:"current_floor"`    // New
	CurrentRoom        string `json:"current_room"`     // New
	CurrentMoo         string `json:"current_moo"`
	CurrentSoi         string `json:"current_soi"`
	CurrentRoad        string `json:"current_road"`
	CurrentProvince    string `json:"current_province"`
	CurrentDistrict    string `json:"current_district"`
	CurrentSubdistrict string `json:"current_subdistrict"`
	CurrentZipcode     string `json:"current_zipcode"`

	// Work Address
	WorkNo          string `json:"work_no"`
	WorkBuilding    string `json:"work_building"` // New
	WorkFloor       string `json:"work_floor"`    // New
	WorkRoom        string `json:"work_room"`     // New
	WorkMoo         string `json:"work_moo"`
	WorkSoi         string `json:"work_soi"`
	WorkRoad        string `json:"work_road"`
	WorkProvince    string `json:"work_province"`
	WorkDistrict    string `json:"work_district"`
	WorkSubdistrict string `json:"work_subdistrict"`
	WorkZipcode     string `json:"work_zipcode"`
	WorkPhone       string `json:"work_phone"`

	// Document Delivery Address
	DocDeliveryType string `json:"doc_delivery_type"`
	DocNo           string `json:"doc_no"`
	DocBuilding     string `json:"doc_building"` // New
	DocFloor        string `json:"doc_floor"`    // New
	DocRoom         string `json:"doc_room"`     // New
	DocMoo          string `json:"doc_moo"`
	DocSoi          string `json:"doc_soi"`
	DocRoad         string `json:"doc_road"`
	DocProvince     string `json:"doc_province"`
	DocDistrict     string `json:"doc_district"`
	DocSubdistrict  string `json:"doc_subdistrict"`
	DocZipcode      string `json:"doc_zipcode"`

	// Step 2: Car Info
	CarType         string  `json:"car_type"`
	CarCode         string  `json:"car_code"`
	CarBrand        string  `json:"car_brand"`
	CarRegisterDate *string `json:"car_register_date"`
	CarModel        string  `json:"car_model"`
	CarYear         string  `json:"car_year"`
	CarColor        string  `json:"car_color"`
	CarWeight       float64 `json:"car_weight"`
	CarCC           int     `json:"car_cc"`
	CarMileage      float64 `json:"car_mileage"`
	CarChassisNo    string  `json:"car_chassis_no"`
	CarGear         string  `json:"car_gear"`
	CarEngineNo     string  `json:"car_engine_no"`
	CarCondition    string  `json:"car_condition"`
	LicensePlate    string  `json:"license_plate"`
	LicenseProvince string  `json:"license_province"`
	VatRate         float64 `json:"vat_rate"`
	CarPrice        float64 `json:"car_price"`
	IsRefinance     bool    `json:"is_refinance"`

	// Step 3: Contract Info
	ContractSignDate  *string `json:"contract_sign_date" gorm:"type:date"`
	LoanType          string  `json:"loan_type"`
	LoanAmount        float64 `json:"loan_amount"`
	InterestRate      float64 `json:"interest_rate"`
	Installments      int     `json:"installments"`
	InstallmentAmount float64 `json:"installment_amount"`
	DownPayment       float64 `json:"down_payment"`
	ContractStartDate *string `json:"contract_start_date" gorm:"type:date"`
	FirstPaymentDate  *string `json:"first_payment_date" gorm:"type:date"`
	TransferType      string  `json:"transfer_type"`
	TransferFee       float64 `json:"transfer_fee"`
	TaxFee            float64 `json:"tax_fee"`
	DutyFee           float64 `json:"duty_fee"`
	PaymentDay        int     `json:"payment_day"`

	// Step 4: Guarantor & Staff
	NoGuarantor       bool    `json:"no_guarantor"`
	Guarantor1Name    string  `json:"guarantor1_name"`
	Guarantor1Contact string  `json:"guarantor1_contact"`
	Guarantor2Name    string  `json:"guarantor2_name"`
	Guarantor2Contact string  `json:"guarantor2_contact"`
	Guarantor3Name    string  `json:"guarantor3_name"`
	Guarantor3Contact string  `json:"guarantor3_contact"`
	LoanOfficer       string  `json:"loan_officer"`
	CompanySellerID   string  `json:"company_seller_id"` // New Field
	CompanySeller     string  `json:"company_seller"`
	ShowroomStaff     string  `json:"showroom_staff"`
	Commission        float64 `json:"commission"`
	ScoreOfficer      float64 `json:"score_officer"`
	ScoreManager      float64 `json:"score_manager"`

	// Step 5: Life Insurance
	HasLifeInsurance     bool    `json:"has_life_insurance"`
	LifeInsuranceCompany string  `json:"life_insurance_company" gorm:"column:life_insurance_company"`
	LifeLoanPrincipal    float64 `json:"life_loan_principal"`
	LifeInterestRate     float64 `json:"life_interest_rate"`
	LifeInstallments     int     `json:"life_installments"`
	LifeGender           string  `json:"life_gender"`
	LifeDob              *string `json:"life_dob"`
	LifeSignDate         *string `json:"life_sign_date"`
	LifeInsuranceRate    float64 `json:"life_insurance_rate" gorm:"column:life_insurance_rate"`
	LifePremium          float64 `json:"life_premium" gorm:"column:life_premium"`
	Beneficiary1Name     string  `json:"beneficiary1_name"`
	Beneficiary1Relation string  `json:"beneficiary1_relation"`
	Beneficiary1Address  string  `json:"beneficiary1_address"`
	Beneficiary2Name     string  `json:"beneficiary2_name"`
	Beneficiary2Relation string  `json:"beneficiary2_relation"`
	Beneficiary2Address  string  `json:"beneficiary2_address"`
	Beneficiary3Name     string  `json:"beneficiary3_name"`
	Beneficiary3Relation string  `json:"beneficiary3_relation"`
	Beneficiary3Address  string  `json:"beneficiary3_address"`
	InsuranceSeller      string  `json:"insurance_seller"`
	InsuranceAgentEmpId  string  `json:"insurance_agent_empid"`
	InsuranceLicenseNo   string  `json:"insurance_license_no"`

	// Step 6: Car Insurance
	CarInsuranceType         string  `json:"car_insurance_type"`
	CarInsuranceCompany      string  `json:"car_insurance_company"`
	CarInsuranceClass        string  `json:"car_insurance_class"`
	CarInsuranceNotifyDate   *string `json:"car_insurance_notify_date"`
	CarInsuranceNotifyNo     string  `json:"car_insurance_notify_no"`
	CarInsuranceStartDate    *string `json:"car_insurance_start_date"`
	CarInsurancePremium      float64 `json:"car_insurance_premium"`
	CarInsuranceBeginning    float64 `json:"car_insurance_beginning"`
	CarInsuranceRefinanceFee float64 `json:"car_insurance_refinance_fee"`
	CarInsuranceAvoidanceFee float64 `json:"car_insurance_avoidance_fee"`
	CarInsuranceFile         string  `json:"car_insurance_file"`

	// Step 7: Withholding Tax
	TaxPayerType   string `json:"tax_payer_type"`
	TaxIdCard      string `json:"tax_id_card"`
	TaxPrefix      string `json:"tax_prefix"`
	TaxFirstName   string `json:"tax_first_name"`
	TaxLastName    string `json:"tax_last_name"`
	TaxHouseNo     string `json:"tax_house_no"`
	TaxBuilding    string `json:"tax_building"`
	TaxFloor       string `json:"tax_floor"`
	TaxRoom        string `json:"tax_room"`
	TaxVillage     string `json:"tax_village"`
	TaxMoo         string `json:"tax_moo"`
	TaxSoi         string `json:"tax_soi"`
	TaxRoad        string `json:"tax_road"`
	TaxProvince    string `json:"tax_province"`
	TaxDistrict    string `json:"tax_district"`
	TaxSubdistrict string `json:"tax_sub_district"`
	TaxZipcode     string `json:"tax_zipcode"`

	Guarantors []Guarantor `json:"guarantors" gorm:"foreignKey:LoanID"`
}

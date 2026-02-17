package models

type LoanProtectLicense struct {
	EmpId         string `gorm:"column:EmpId;primaryKey"`
	EmpName       string `gorm:"column:EmpName"`
	BrokerLicense string `gorm:"column:BrokerLicense"`
	LicenseStatus string `gorm:"column:LicenseStatus"`
}

func (LoanProtectLicense) TableName() string {
	return "loanprotectlicense"
}

package handlers

import (
	"fmt"
	"loan-app/config"
	"loan-app/models"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AddGuarantorGet(c *fiber.Ctx) error {
	loanID := c.Query("loan_id")
	if loanID == "" {
		return c.Redirect("/")
	}
	return c.Render("add_guarantor", fiber.Map{
		"LoanID": loanID,
	})
}

func AddGuarantorPost(c *fiber.Ctx) error {
	loanIDStr := c.FormValue("loan_id")
	loanID, _ := strconv.ParseUint(loanIDStr, 10, 64)

	// Ensure Schema is up-to-date - Removed to prevent schema conflicts
	// config.DB.AutoMigrate(&models.Guarantor{})

	// Debug log to confirm code update
	fmt.Printf("DEBUG: Processing AddGuarantorPost for loanID: %d\n", loanID)
	fmt.Printf("DEBUG: IdCardIssueDate raw: '%s'\n", c.FormValue("id_card_issue_date"))
	fmt.Printf("DEBUG: OtherCardIssueDate raw: '%s'\n", c.FormValue("other_card_issue_date"))

	guarantor := models.Guarantor{
		LoanID: int32(loanID),

		GuarantorType:       c.FormValue("guarantor_type"),
		TradeRegistrationID: c.FormValue("trade_registration_id"),
		RegistrationDate:    parseDate(c.FormValue("registration_date")),
		TaxID:               c.FormValue("tax_id"),

		Title:            c.FormValue("title"),
		FirstName:        c.FormValue("first_name"),
		LastName:         c.FormValue("last_name"),
		Gender:           c.FormValue("gender"),
		IdCard:           c.FormValue("id_card"),
		IdCardIssueDate:  parseDate(c.FormValue("id_card_issue_date")),
		IdCardExpiryDate: parseDate(c.FormValue("id_card_expiry_date")),
		DateOfBirth:      parseDate(c.FormValue("date_of_birth")),
		Ethnicity:        c.FormValue("ethnicity"),
		Nationality:      c.FormValue("nationality"),
		Religion:         c.FormValue("religion"),
		MaritalStatus:    c.FormValue("marital_status"),
		MobilePhone:      c.FormValue("mobile_phone"),

		HouseRegNo:          c.FormValue("house_reg_no"),
		HouseRegBuilding:    c.FormValue("house_reg_building"),
		HouseRegFloor:       c.FormValue("house_reg_floor"),
		HouseRegRoom:        c.FormValue("house_reg_room"),
		HouseRegMoo:         c.FormValue("house_reg_moo"),
		HouseRegSoi:         c.FormValue("house_reg_soi"),
		HouseRegRoad:        c.FormValue("house_reg_road"),
		HouseRegProvince:    c.FormValue("house_reg_province"),
		HouseRegDistrict:    c.FormValue("house_reg_district"),
		HouseRegSubdistrict: c.FormValue("house_reg_subdistrict"),
		HouseRegZipcode:     c.FormValue("house_reg_zipcode"),

		SameAsHouseReg:     c.FormValue("same_as_house_reg") == "on",
		CurrentNo:          c.FormValue("current_no"),
		CurrentBuilding:    c.FormValue("current_building"),
		CurrentFloor:       c.FormValue("current_floor"),
		CurrentRoom:        c.FormValue("current_room"),
		CurrentMoo:         c.FormValue("current_moo"),
		CurrentSoi:         c.FormValue("current_soi"),
		CurrentRoad:        c.FormValue("current_road"),
		CurrentProvince:    c.FormValue("current_province"),
		CurrentDistrict:    c.FormValue("current_district"),
		CurrentSubdistrict: c.FormValue("current_subdistrict"),
		CurrentZipcode:     c.FormValue("current_zipcode"),

		// Work Info
		CompanyName:     c.FormValue("company_name"),
		Occupation:      c.FormValue("occupation"),
		Position:        c.FormValue("position"),
		Salary:          parseMoney(c.FormValue("salary")),
		OtherIncome:     parseMoney(c.FormValue("other_income")),
		IncomeSource:    c.FormValue("income_source"),
		WorkPhone:       c.FormValue("work_phone"),
		WorkNo:          c.FormValue("work_no"),
		WorkBuilding:    c.FormValue("work_building"),
		WorkFloor:       c.FormValue("work_floor"),
		WorkRoom:        c.FormValue("work_room"),
		WorkMoo:         c.FormValue("work_moo"),
		WorkSoi:         c.FormValue("work_soi"),
		WorkRoad:        c.FormValue("work_road"),
		WorkProvince:    c.FormValue("work_province"),
		WorkDistrict:    c.FormValue("work_district"),
		WorkSubdistrict: c.FormValue("work_subdistrict"),
		WorkZipcode:     c.FormValue("work_zipcode"),

		// Other Card
		OtherCardType:       c.FormValue("other_card_type"),
		OtherCardNumber:     c.FormValue("other_card_number"),
		OtherCardIssueDate:  parseDate(c.FormValue("other_card_issue_date")),
		OtherCardExpiryDate: parseDate(c.FormValue("other_card_expiry_date")),

		// Doc Address
		DocDeliveryType: c.FormValue("doc_delivery_type"),
		DocNo:           c.FormValue("doc_no"),
		DocBuilding:     c.FormValue("doc_building"),
		DocFloor:        c.FormValue("doc_floor"),
		DocRoom:         c.FormValue("doc_room"),
		DocMoo:          c.FormValue("doc_moo"),
		DocSoi:          c.FormValue("doc_soi"),
		DocRoad:         c.FormValue("doc_road"),
		DocProvince:     c.FormValue("doc_province"),
		DocDistrict:     c.FormValue("doc_district"),
		DocSubdistrict:  c.FormValue("doc_subdistrict"),
		DocZipcode:      c.FormValue("doc_zipcode"),
	}

	// Un-comment AutoMigrate to ensure new columns are created
	config.DB.AutoMigrate(&models.Guarantor{})

	// Use Raw SQL with NULLIF to explicitly handle empty strings as NULL
	query := `INSERT INTO loan_applications_guarantors (
		created_at, updated_at, loan_id,
		guarantor_type, trade_registration_id, registration_date, tax_id,
		title, first_name, last_name, gender, id_card,
		id_card_issue_date, id_card_expiry_date, date_of_birth,
		ethnicity, nationality, religion, marital_status, mobile_phone,
		house_reg_no, house_reg_building, house_reg_floor, house_reg_room, house_reg_moo, house_reg_soi, house_reg_road, house_reg_province, house_reg_district, house_reg_subdistrict, house_reg_zipcode,
		same_as_house_reg,
		current_no, current_building, current_floor, current_room, current_moo, current_soi, current_road, current_province, current_district, current_subdistrict, current_zipcode,
		company_name, occupation, position, salary, other_income, income_source,
		work_phone, work_no, work_building, work_floor, work_room, work_moo, work_soi, work_road, work_province, work_district, work_subdistrict, work_zipcode,
		other_card_type, other_card_number, other_card_issue_date, other_card_expiry_date,
		doc_delivery_type, doc_no, doc_building, doc_floor, doc_room, doc_moo, doc_soi, doc_road, doc_province, doc_district, doc_subdistrict, doc_zipcode
	) VALUES (
		NOW(), NOW(), ?,
		?, ?, NULLIF(?, ''), ?,
		?, ?, ?, ?, ?,
		NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''),
		?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		?,
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, NULLIF(?, ''), NULLIF(?, ''),
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
	)`

	err := config.DB.Exec(query,
		guarantor.LoanID,
		guarantor.GuarantorType, guarantor.TradeRegistrationID, c.FormValue("registration_date"), guarantor.TaxID,
		guarantor.Title, guarantor.FirstName, guarantor.LastName, guarantor.Gender, guarantor.IdCard,
		c.FormValue("id_card_issue_date"), c.FormValue("id_card_expiry_date"), c.FormValue("date_of_birth"),
		guarantor.Ethnicity, guarantor.Nationality, guarantor.Religion, guarantor.MaritalStatus, guarantor.MobilePhone,
		// House
		guarantor.HouseRegNo, guarantor.HouseRegBuilding, guarantor.HouseRegFloor, guarantor.HouseRegRoom, guarantor.HouseRegMoo, guarantor.HouseRegSoi, guarantor.HouseRegRoad, guarantor.HouseRegProvince, guarantor.HouseRegDistrict, guarantor.HouseRegSubdistrict, guarantor.HouseRegZipcode,
		guarantor.SameAsHouseReg,
		// Current
		guarantor.CurrentNo, guarantor.CurrentBuilding, guarantor.CurrentFloor, guarantor.CurrentRoom, guarantor.CurrentMoo, guarantor.CurrentSoi, guarantor.CurrentRoad, guarantor.CurrentProvince, guarantor.CurrentDistrict, guarantor.CurrentSubdistrict, guarantor.CurrentZipcode,
		// Work
		guarantor.CompanyName, guarantor.Occupation, guarantor.Position, guarantor.Salary, guarantor.OtherIncome, guarantor.IncomeSource,
		guarantor.WorkPhone, guarantor.WorkNo, guarantor.WorkBuilding, guarantor.WorkFloor, guarantor.WorkRoom, guarantor.WorkMoo, guarantor.WorkSoi, guarantor.WorkRoad, guarantor.WorkProvince, guarantor.WorkDistrict, guarantor.WorkSubdistrict, guarantor.WorkZipcode,
		guarantor.OtherCardType, guarantor.OtherCardNumber, c.FormValue("other_card_issue_date"), c.FormValue("other_card_expiry_date"),
		// Doc
		guarantor.DocDeliveryType, guarantor.DocNo, guarantor.DocBuilding, guarantor.DocFloor, guarantor.DocRoom, guarantor.DocMoo, guarantor.DocSoi, guarantor.DocRoad, guarantor.DocProvince, guarantor.DocDistrict, guarantor.DocSubdistrict, guarantor.DocZipcode,
	).Error

	if err != nil {
		fmt.Printf("DEBUG: Error executing raw SQL: %v\n", err)
		return c.Status(500).SendString("Error saving guarantor (Raw SQL): " + err.Error())
	}
	fmt.Println("DEBUG: Successfully inserted via Raw SQL")

	// if err := config.DB.Create(&guarantor).Error; err != nil {
	// 	return c.Status(500).SendString("Error saving guarantor: " + err.Error())
	// }

	return c.Redirect("/step4?id=" + loanIDStr)
}

func parseMoney(amount string) float64 {
	amount = strings.ReplaceAll(amount, ",", "")
	val, _ := strconv.ParseFloat(amount, 64)
	return val
}

func parseDate(dateStr string) *string {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return nil
	}
	return &dateStr
}

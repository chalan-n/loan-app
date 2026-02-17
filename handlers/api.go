package handlers

import (
	"fmt"
	"loan-app/config"
	"loan-app/models"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
)

// Request struct for search
type SearchCarRequest struct {
	CarCode string `json:"car_code"`
}

func SearchCar(c *fiber.Ctx) error {
	var req SearchCarRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.CarCode == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Car Code is required",
		})
	}

	var car models.RedbookII
	result := config.DB.Where("carCode = ?", req.CarCode).First(&car)

	if result.Error != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Car not found",
		})
	}

	return c.JSON(car)
}

// Struct for Insurance Calculation Request
type CalculateInsuranceReq struct {
	LoanID           uint   `json:"loan_id"`
	InsuranceCompany string `json:"insurance_company"`
	Installments     uint   `json:"installments"`
	Age              int    `json:"age"`                // Optional: If provided, usage overrides calculation
	ContractSignDate string `json:"contract_sign_date"` // วันที่เซ็นสัญญาจากฟอร์ม
}

func CalculateInsuranceRate(c *fiber.Ctx) error {
	var req CalculateInsuranceReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var loan models.LoanApplication
	if err := config.DB.First(&loan, req.LoanID).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Loan not found"})
	}

	// 1. Gender: Male=1, Female=2
	genderCode := 0
	if loan.Gender == "male" {
		genderCode = 1
	} else {
		genderCode = 0
	}

	// 2. Age Calculation
	var age int
	var signDate time.Time
	var err error

	// Priority: 1. req.ContractSignDate (ฟอร์ม) → 2. loan.ContractSignDate (DB) → 3. time.Now()
	if req.ContractSignDate != "" {
		// ลองแปลงหลายรูปแบบจาก Frontend
		formats := []string{"2006-01-02", "02-01-2006", "02/01/2006"}
		parsed := false
		for _, f := range formats {
			signDate, err = time.Parse(f, req.ContractSignDate)
			if err == nil {
				parsed = true
				break
			}
		}
		if !parsed {
			signDate = time.Now()
		}
	} else if loan.ContractSignDate != nil && *loan.ContractSignDate != "" {
		signDate, err = time.Parse("2006-01-02", *loan.ContractSignDate)
		if err != nil {
			signDate = time.Now()
		}
	} else {
		signDate = time.Now()
	}

	if req.Age > 0 {
		age = req.Age
	} else {
		// Calculate Server Side if Age not provided
		// DateOfBirth is now *string
		dobStr := ""
		if loan.DateOfBirth != nil {
			dobStr = *loan.DateOfBirth
		}

		if dobStr != "" {
			dob, err := time.Parse("2006-01-02", dobStr)
			if err == nil {
				age = int(signDate.Year() - dob.Year())
				if signDate.YearDay() < dob.YearDay() {
					age--
				}
				if age < 0 {
					age = 0
				}
			} else {
				age = 0
			}
		} else {
			age = 0
		}
	}

	// 3. Call DB Function
	var rate float64

	// Format Date to YYYYMMDD
	signDateStr := signDate.Format("20060102")

	// Convert numbers to strings as per example: Fnc_LoanProtect_Rate('03','1','48','41','20250226')
	genderStr := fmt.Sprintf("%d", genderCode)
	installmentsStr := fmt.Sprintf("%d", req.Installments)
	ageStr := fmt.Sprintf("%d", age)

	// Construct Debug SQL (Strictly matching user example)
	debugSQL := fmt.Sprintf("SELECT Fnc_LoanProtect_Rate('%s', '%s', '%s', '%s', '%s') AS RATE",
		req.InsuranceCompany, genderStr, installmentsStr, ageStr, signDateStr)

	// Log to Server Console
	log.Println("-------- INSURANCE CALC DEBUG --------")
	log.Println(debugSQL)
	log.Println("--------------------------------------")

	// SQL: SELECT Fnc_LoanProtect_Rate(?, ?, ?, ?, ?) AS RATE
	query := "SELECT Fnc_LoanProtect_Rate(?, ?, ?, ?, ?) AS RATE"

	// Use Row() for scalar scan to avoid GORM slice issues
	err = config.DB.Raw(query,
		req.InsuranceCompany,
		genderStr,
		installmentsStr,
		ageStr,
		signDateStr,
	).Row().Scan(&rate)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Calculation failed: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"rate":         rate,
		"debug_age":    age,
		"debug_gender": genderCode,
		"debug_id":     loan.ID,
		"debug_sql":    debugSQL,
	})
}

// Request struct for Agent Search
type SearchAgentRequest struct {
	Query string `json:"query"`
}

func SearchAgent(c *fiber.Ctx) error {
	var req SearchAgentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Query == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Query is required",
		})
	}

	var agents []models.LoanProtectLicense
	// Search by EmpId or EmpName using LIKE
	// Note: Adjust wildcards as per DB requirement (e.g. %query%)
	searchTerm := "%" + req.Query + "%"
	result := config.DB.Where("EmpId LIKE ? OR EmpName LIKE ?", searchTerm, searchTerm).Limit(20).Find(&agents)

	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error: " + result.Error.Error(),
		})
	}

	return c.JSON(agents)
}

// Request struct for Showroom Search
type SearchShowroomRequest struct {
	Query string `json:"query"`
}

func SearchShowroom(c *fiber.Ctx) error {
	var req SearchShowroomRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Query == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Query is required",
		})
	}

	var showrooms []models.Showroom
	searchTerm := "%" + req.Query + "%"
	// Search by ID or Name
	result := config.DB.Where("ShowRoomId LIKE ? OR ShowRoomName LIKE ?", searchTerm, searchTerm).Limit(20).Find(&showrooms)

	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error: " + result.Error.Error(),
		})
	}

	return c.JSON(showrooms)
}

// GetTitles fetches all titles from the database
func GetTitles(c *fiber.Ctx) error {
	var titles []models.Title
	if err := config.DB.Find(&titles).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch titles",
		})
	}
	return c.JSON(titles)
}

// GetRaces fetches all races from the database
func GetRaces(c *fiber.Ctx) error {
	var races []models.Race
	if err := config.DB.Find(&races).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch races",
		})
	}
	return c.JSON(races)
}

// GetNations fetches all nations from the database
func GetNations(c *fiber.Ctx) error {
	var nations []models.Nation
	if err := config.DB.Find(&nations).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch nations",
		})
	}
	return c.JSON(nations)
}

// GetReligions fetches all religions from the database
func GetReligions(c *fiber.Ctx) error {
	var religions []models.Religion
	if err := config.DB.Find(&religions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch religions",
		})
	}
	return c.JSON(religions)
}

// GetSituations fetches all marital statuses from the database
func GetSituations(c *fiber.Ctx) error {
	var situations []models.Situation
	if err := config.DB.Find(&situations).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch situations",
		})
	}
	return c.JSON(situations)
}

// GetOccupations fetches all occupations from the database
func GetOccupations(c *fiber.Ctx) error {
	var occupations []models.Occupy
	if err := config.DB.Find(&occupations).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch occupations",
		})
	}
	return c.JSON(occupations)
}

// GetInsuComps fetches all insurance companies from the database
func GetInsuComps(c *fiber.Ctx) error {
	var companies []models.InsuComp
	if err := config.DB.Find(&companies).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch insurance companies",
		})
	}
	return c.JSON(companies)
}

// GetInsuClasses fetches all insurance classes from the database
func GetInsuClasses(c *fiber.Ctx) error {
	var classes []models.InsuClass
	if err := config.DB.Find(&classes).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch insurance classes",
		})
	}
	return c.JSON(classes)
}

// Cloudflare R2 Config
const (
	R2AccountId       = "1c3c174d2225cf15e7006b4a55f7a4a3"
	R2AccessKeyId     = "b7b136ea56fc886041267387cd6d211a"
	R2SecretAccessKey = "0edec23269dcdef1010706dcdeecdf44b0e174be1226d2a97e2ba40b63e22f10"
	R2BucketName      = "cmo-loan-app"
	R2Endpoint        = "https://1c3c174d2225cf15e7006b4a55f7a4a3.r2.cloudflarestorage.com"
)

// Helper to get S3 Client
func getR2Client() *s3.S3 {
	creds := credentials.NewStaticCredentials(R2AccessKeyId, R2SecretAccessKey, "")
	cfg := aws.NewConfig().
		WithRegion("auto").
		WithEndpoint(R2Endpoint).
		WithCredentials(creds)

	sess := session.Must(session.NewSession(cfg))
	return s3.New(sess)
}

// Upload Insurance File to Cloudflare R2
func UploadInsuranceFile(c *fiber.Ctx) error {
	loanID := c.Cookies("loan_id")
	if loanID == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "File upload failed"})
	}

	// Open file content
	src, err := fileHeader.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open file"})
	}
	defer src.Close()

	// Generate safe key (filename)
	filename := fmt.Sprintf("%s_%d_%s", loanID, time.Now().UnixNano(), fileHeader.Filename)

	// Upload to R2
	svc := getR2Client()
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(R2BucketName),
		Key:         aws.String(filename),
		Body:        src,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})

	if err != nil {
		log.Printf("R2 Upload Error: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload to Cloud Storage"})
	}

	// Update Database
	var loan models.LoanApplication
	if err := config.DB.First(&loan, loanID).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Loan not found"})
	}

	if loan.CarInsuranceFile == "" {
		loan.CarInsuranceFile = filename
	} else {
		loan.CarInsuranceFile = loan.CarInsuranceFile + "," + filename
	}

	config.DB.Save(&loan)

	return c.JSON(fiber.Map{
		"message":  "Upload success",
		"filename": filename,
		// Return API URL for viewing (via presigned redirect)
		"url": "/file/" + filename,
	})
}

// GetFile Redirects to Presigned URL
func GetFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" {
		return c.Status(400).SendString("Filename required")
	}

	svc := getR2Client()
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(R2BucketName),
		Key:    aws.String(filename),
	})

	// Presign for 15 minutes
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		log.Printf("Presign Error: %v", err)
		return c.Status(500).SendString("Failed to generate download link")
	}

	return c.Redirect(urlStr)
}

// Delete Insurance File
type DeleteFileRequest struct {
	Filename string `json:"filename"`
}

func DeleteInsuranceFile(c *fiber.Ctx) error {
	loanID := c.Cookies("loan_id")
	if loanID == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req DeleteFileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var loan models.LoanApplication
	if err := config.DB.First(&loan, loanID).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Loan not found"})
	}

	// Remove filename from DB field (comma separated)
	files := strings.Split(loan.CarInsuranceFile, ",")
	newFiles := []string{}
	found := false
	for _, f := range files {
		if f == req.Filename {
			found = true

			// Delete from R2
			svc := getR2Client()
			_, err := svc.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(R2BucketName),
				Key:    aws.String(f),
			})
			if err != nil {
				log.Printf("R2 Delete Error for %s: %v", f, err)
				// Continue to remove from DB anyway? Yes.
			}

		} else {
			if f != "" {
				newFiles = append(newFiles, f)
			}
		}
	}

	if found {
		loan.CarInsuranceFile = strings.Join(newFiles, ",")
		config.DB.Save(&loan)
		return c.JSON(fiber.Map{"message": "File deleted"})
	}

	return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "File not found in record"})
}

// DeleteLoan deletes a loan application and its guarantors
func DeleteLoan(c *fiber.Ctx) error {
	type Request struct {
		ID int `json:"id"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// 1. Fetch Loan Data to get File Paths
	var loan models.LoanApplication
	if err := config.DB.First(&loan, req.ID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Loan not found"})
	}

	// 2. Delete Uploaded Files
	if loan.CarInsuranceFile != "" {
		files := strings.Split(loan.CarInsuranceFile, ",")
		svc := getR2Client()
		for _, f := range files {
			if f != "" {
				// Delete from R2
				_, err := svc.DeleteObject(&s3.DeleteObjectInput{
					Bucket: aws.String(R2BucketName),
					Key:    aws.String(f),
				})
				if err != nil {
					log.Printf("Failed to delete file %s from R2: %v", f, err)
				}
			}
		}
	}

	// Transaction to ensure atomicity
	tx := config.DB.Begin()

	// 3. Delete Guarantors associated with this LoanID
	if err := tx.Where("loan_id = ?", req.ID).Unscoped().Delete(&models.Guarantor{}).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete guarantors"})
	}

	// 4. Delete Loan Application
	if err := tx.Delete(&models.LoanApplication{}, req.ID).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete loan application"})
	}

	tx.Commit()
	return c.JSON(fiber.Map{"success": true})
}

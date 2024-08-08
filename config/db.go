package config

import (
	"log"
	"os"
	"sea-study/api/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// checkEnumExists checks if an enum type exists in the database
func checkEnumExists(db *gorm.DB, enumName string) bool {
    var exists bool
    query := `SELECT EXISTS (
        SELECT 1 
        FROM pg_type 
        WHERE typname = ?
    )`
    db.Raw(query, enumName).Scan(&exists)
    return exists
}

// createEnums creates the necessary enum types in the database
func createEnums(db *gorm.DB) {
    enums := map[string]string{
        "role_enum":              "CREATE TYPE role_enum AS ENUM ('user', 'author');",
        "topup_status_enum":      "CREATE TYPE topup_status_enum AS ENUM ('pending', 'completed');",
        "course_difficulty_enum": "CREATE TYPE course_difficulty_enum AS ENUM ('beginner', 'intermediate', 'advanced');",
        "course_status_enum":     "CREATE TYPE course_status_enum AS ENUM ('active', 'inactive');",
        "progress_status_enum":   "CREATE TYPE progress_status_enum AS ENUM ('in_progress', 'completed');",
        "submission_status_enum": "CREATE TYPE submission_status_enum AS ENUM ('submitted', 'graded');",
    }

    for enumName, createQuery := range enums {
        if !checkEnumExists(db, enumName) {
            if err := db.Exec(createQuery).Error; err != nil {
                log.Printf("Error creating enum %s: %v\n", enumName, err)
            }
        } else {
            log.Printf("Enum %s already exists\n", enumName)
        }
    }
}

// InitDB connects to the database and performs auto-migration.
func InitDB() *gorm.DB {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    dsn := "host=" + os.Getenv("SUPABASE_HOST") +
        " user=" + os.Getenv("SUPABASE_USER") +
        " password=" + os.Getenv("SUPABASE_PASSWORD") +
        " dbname=" + os.Getenv("SUPABASE_DBNAME") +
        " port=" + os.Getenv("SUPABASE_PORT") +
        " sslmode=require"


    db, err := gorm.Open(postgres.New(postgres.Config{
        DSN:                  dsn,
        PreferSimpleProtocol: true, 
    }), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }


    createEnums(db)

    // AutoMigrate all models
    db.AutoMigrate(
        &models.User{},
        &models.TopupHistory{},
        &models.Course{},
        &models.Syllabus{},
        &models.ForumPost{},
        &models.SyllabusMaterial{},
        &models.UserProgress{},
        &models.Enrollment{},
        &models.Assignment{},
        &models.Submission{},
        &models.CourseReview{},
        &models.UserAssignment{},
    )

    return db
}

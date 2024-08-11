package models

// User role enums
type RoleEnum string

const (
    UserRole   RoleEnum = "user"
    AuthorRole RoleEnum = "author"
    InstructorRole RoleEnum = "instructor"
)

// Topup status enums
type TopupStatusEnum string

const (
    PendingStatus   TopupStatusEnum = "pending"
    CompletedStatus TopupStatusEnum = "completed"
)

// Course difficulty enums
type DifficultyEnum string

const (
    BeginnerLevel    DifficultyEnum = "beginner"
    IntermediateLevel DifficultyEnum = "intermediate"
    AdvancedLevel     DifficultyEnum = "advanced"
)

// Course status enums
type CourseStatusEnum string

const (
    ActiveStatus   CourseStatusEnum = "active"
    InactiveStatus CourseStatusEnum = "inactive"
)

// User progress status enums
type ProgressStatusEnum string

const (
    InProgress   ProgressStatusEnum = "in_progress"
    Completed    ProgressStatusEnum = "completed"
)

// Submission status enums
type SubmissionStatusEnum string

const (
    Submitted      SubmissionStatusEnum = "submitted"
    Graded         SubmissionStatusEnum = "graded"
)


// Course category enums
type CategoryEnum string

const (
    Android          CategoryEnum = "Android"
    Web              CategoryEnum = "Web"
    Desktop          CategoryEnum = "Desktop"
    Multiplatform    CategoryEnum = "Multiplatform"
    Game             CategoryEnum = "Game"
    IOS              CategoryEnum = "IOS"
    MachineLearning  CategoryEnum = "Machine Learning"
    DataScientist    CategoryEnum = "Data Scientist"
    React            CategoryEnum = "React"
    DevOps           CategoryEnum = "DevOps"
    Cloud            CategoryEnum = "Cloud"
)
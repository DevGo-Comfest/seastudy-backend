package constants

// User Errors
const (
	ErrUnauthorized  = "unauthorized"
	ErrInvalidUserID = "invalid user ID format"
	ErrUserNotFound  = "user not found"
	ErrAuthorNotFound = "author not found"
	ErrFailedToHash = "failed to hash password"
	ErrFailedToCreateUser = "failed to create user"
	ErrInvalidCredentials = "invalid email or password"
	ErrFailedToRetrieveUser = "failed to retrieve user"
	ErrFailedToRetrieveAuthors = "failed to retrieve authors"
	ErrFailedToGenerateToken = "failed to generate token"
	ErrUserNotAuthenticated = "user not authenticated"
	ErrFailedToRetrieveUserProfile = "failed to retrieve user profile"
)

// File Upload Errors
const (
	ErrFailedToUploadImage = "failed to upload image"
	ErrFailedToSaveImage   = "failed to save image"
	ErrFailedToSaveFile   = "failed to save file"
)

// Generic Errors
const (
	ErrInvalidInput = "invalid input data"
)

// Course Errors
const (
	ErrInvalidCourseID         = "invalid course ID"
	ErrCourseNotFound          = "course not found"
	ErrNoInstructorsFound      = "no instructors found for this course"
	ErrFailedToCreateCourse    = "failed to create course"
	ErrFailedToUpdateCourse    = "failed to update course"
	ErrFailedToDeleteCourse    = "failed to delete course"
	ErrFailedToRetrieveCourse  = "failed to retrieve course"
	ErrInvalidRating           = "invalid rating"
	ErrFailedToRetrieveCourses = "failed to retrieve courses"
	ErrFailedToRetrieveInstructors = "failed to retrieve instructors"
	ErrFailedToOpenFile = "failed to open file"
	ErrFailedToReadFile = "failed to read file"

)

// Syllabus Errors
const (
	ErrInvalidSyllabusID        = "invalid syllabus ID"
	ErrSyllabusNotFound         = "syllabus not found"
	ErrFailedToCreateSyllabus   = "failed to create syllabus"
	ErrFailedToUpdateSyllabus   = "failed to update syllabus"
	ErrFailedToDeleteSyllabus   = "failed to delete syllabus"
	ErrUnauthorizedSyllabus     = "unauthorized to modify this syllabus"
	ErrFailedToRetrieveSyllabus = "failed to retrieve syllabus"
)

// Enrollment Errors
const (
	ErrUserAlreadyEnrolled             = "user is already enrolled in the course"
	ErrInsufficientBalance             = "insufficient balance to enroll in the course"
	ErrFailedToCreateEnrollment        = "failed to create enrollment"
	ErrFailedToRetrieveEnrolledCourses = "failed to retrieve enrolled courses"
)

// Forum Post Errors
const (
	ErrFailedToCreateForumPost = "failed to create forum post"
	ErrFailedToRetrievePosts   = "failed to retrieve forum posts"
)

// User Progress Errors
const (
	ErrFailedToUpdateProgress       = "failed to update user progress"
	ErrIncompletePreviousSyllabus   = "complete all previous syllabuses to open this one"
	ErrFailedToRetrieveUserProgress = "failed to retrieve user course progress"
	ErrNoSyllabusesFound            = "no syllabuses found for this course"
)

// Review Errors
const (
	ErrFailedToCreateReview       = "failed to create review"
	ErrFailedToRetrieveReviews    = "failed to retrieve course reviews"
	ErrUserNotEnrolledInCourse    = "user is not enrolled in the course"
	ErrInvalidRate                = "rate must be between 1 and 5"
	ErrUserAlreadySubmittedReview = "user has already submitted a review for this course"
	ErrIncompleteCourseProgress   = "complete all course materials to leave a review"
	ErrNoReviewsFound             = "no reviews found"
	ErrFailedToUpdateRating       = "failed to update rating"
)

// Syllabus Material Errors
const (
	ErrInvalidSyllabusMaterialID        = "invalid syllabus material ID"
	ErrUnauthorizedSyllabusAction       = "unauthorized to perform this action on the syllabus material"
	ErrFailedToCreateSyllabusMaterial   = "failed to create syllabus material"
	ErrFailedToUpdateSyllabusMaterial   = "failed to update syllabus material"
	ErrFailedToDeleteSyllabusMaterial   = "failed to delete syllabus material"
	ErrFailedToRetrieveSyllabusMaterial = "failed to retrieve syllabus material"
)

// Topup Errors
const (
	ErrFailedToCreateTopup       = "failed to create top-up"
	ErrFailedToUpdateUserBalance = "failed to update user balance"
	ErrFailedToRetrieveTopupHistory = "fialed to retriver top-up history"
)

// Assignment Errors
const (
	ErrInvalidAssignmentID          = "invalid assignment ID"
	ErrAssignmentNotFound           = "assignment not found"
	ErrUnauthorizedAssignmentAction = "unauthorized to perform this action on the assignment"
	ErrFailedToCreateAssignment     = "failed to create assignment"
	ErrFailedToUpdateAssignment     = "failed to update assignment"
	ErrFailedToDeleteAssignment     = "failed to delete assignment"
	ErrFailedToRetrieveAssignment   = "failed to retrieve assignment"
	ErrFailedToCreateUserAssignment = "failed to create user assignment"
)

// Submission Errors
const (
	ErrInvalidSubmissionID            = "invalid submission ID"
	ErrSubmissionNotFound             = "submission not found"
	ErrUserAlreadySubmmitedAssignment = "user has already submitted for this assignment"
	ErrFailedToCreateSubmission       = "failed to submit assignment"
	ErrFailedToUpdateSubmission       = "failed to update submission"
	ErrFailedToDeleteSubmission       = "failed to delete submission"
)

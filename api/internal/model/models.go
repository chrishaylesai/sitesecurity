package model

import "time"

type Company struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Address   *string   `json:"address,omitempty" db:"address"`
	Phone     *string   `json:"phone,omitempty" db:"phone"`
	Email     *string   `json:"email,omitempty" db:"email"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type Worksite struct {
	ID        string   `json:"id" db:"id"`
	CompanyID string   `json:"companyId" db:"company_id"`
	Name      string   `json:"name" db:"name"`
	Address   *string  `json:"address,omitempty" db:"address"`
	Latitude  *float64 `json:"latitude,omitempty" db:"latitude"`
	Longitude *float64 `json:"longitude,omitempty" db:"longitude"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type Worker struct {
	ID          string    `json:"id" db:"id"`
	AuthSubject string    `json:"authSubject" db:"auth_subject"`
	FirstName   string    `json:"firstName" db:"first_name"`
	LastName    string    `json:"lastName" db:"last_name"`
	Email       string    `json:"email" db:"email"`
	Phone       *string   `json:"phone,omitempty" db:"phone"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type WorkerRole string

const (
	RoleWorker       WorkerRole = "worker"
	RoleCompanyAdmin WorkerRole = "company_admin"
	RoleSiteAdmin    WorkerRole = "site_admin"
)

type MembershipStatus string

const (
	MembershipActive   MembershipStatus = "active"
	MembershipInactive MembershipStatus = "inactive"
)

type WorkerCompany struct {
	WorkerID  string           `json:"workerId" db:"worker_id"`
	CompanyID string           `json:"companyId" db:"company_id"`
	Role      WorkerRole       `json:"role" db:"role"`
	Status    MembershipStatus `json:"status" db:"status"`
	JoinedAt  time.Time        `json:"joinedAt" db:"joined_at"`
}

type Certificate struct {
	ID                string     `json:"id" db:"id"`
	WorkerID          string     `json:"workerId" db:"worker_id"`
	Name              string     `json:"name" db:"name"`
	IssuingBody       *string    `json:"issuingBody,omitempty" db:"issuing_body"`
	CertificateNumber *string    `json:"certificateNumber,omitempty" db:"certificate_number"`
	IssuedDate        *string    `json:"issuedDate,omitempty" db:"issued_date"`
	ExpiryDate        *string    `json:"expiryDate,omitempty" db:"expiry_date"`
	CreatedAt         time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time  `json:"updatedAt" db:"updated_at"`
}

type ShiftStatus string

const (
	ShiftOpen       ShiftStatus = "open"
	ShiftAssigned   ShiftStatus = "assigned"
	ShiftInProgress ShiftStatus = "in_progress"
	ShiftCompleted  ShiftStatus = "completed"
	ShiftCancelled  ShiftStatus = "cancelled"
)

type Shift struct {
	ID          string      `json:"id" db:"id"`
	WorksiteID  string      `json:"worksiteId" db:"worksite_id"`
	CreatedBy   string      `json:"createdBy" db:"created_by"`
	Title       string      `json:"title" db:"title"`
	Description *string     `json:"description,omitempty" db:"description"`
	StartTime   time.Time   `json:"startTime" db:"start_time"`
	EndTime     time.Time   `json:"endTime" db:"end_time"`
	Status      ShiftStatus `json:"status" db:"status"`
	CreatedAt   time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time   `json:"updatedAt" db:"updated_at"`
}

type AssignmentStatus string

const (
	AssignmentOffered   AssignmentStatus = "offered"
	AssignmentAccepted  AssignmentStatus = "accepted"
	AssignmentDeclined  AssignmentStatus = "declined"
	AssignmentCompleted AssignmentStatus = "completed"
)

type ShiftAssignment struct {
	ID          string           `json:"id" db:"id"`
	ShiftID     string           `json:"shiftId" db:"shift_id"`
	WorkerID    string           `json:"workerId" db:"worker_id"`
	Status      AssignmentStatus `json:"status" db:"status"`
	AssignedAt  time.Time        `json:"assignedAt" db:"assigned_at"`
	RespondedAt *time.Time       `json:"respondedAt,omitempty" db:"responded_at"`
}

type ShiftReportTemplate struct {
	ID        string    `json:"id" db:"id"`
	CompanyID string    `json:"companyId" db:"company_id"`
	Name      string    `json:"name" db:"name"`
	Fields    string    `json:"fields" db:"fields"` // JSONB stored as string
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type ShiftReport struct {
	ID          string    `json:"id" db:"id"`
	ShiftID     string    `json:"shiftId" db:"shift_id"`
	WorkerID    string    `json:"workerId" db:"worker_id"`
	TemplateID  *string   `json:"templateId,omitempty" db:"template_id"`
	Data        string    `json:"data" db:"data"` // JSONB stored as string
	SubmittedAt time.Time `json:"submittedAt" db:"submitted_at"`
}

type LocationCheckIn struct {
	ID         string    `json:"id" db:"id"`
	WorkerID   string    `json:"workerId" db:"worker_id"`
	ShiftID    *string   `json:"shiftId,omitempty" db:"shift_id"`
	Latitude   float64   `json:"latitude" db:"latitude"`
	Longitude  float64   `json:"longitude" db:"longitude"`
	RecordedAt time.Time `json:"recordedAt" db:"recorded_at"`
}

type AlarmStatus string

const (
	AlarmRaised       AlarmStatus = "raised"
	AlarmAcknowledged AlarmStatus = "acknowledged"
	AlarmResolved     AlarmStatus = "resolved"
)

type Alarm struct {
	ID             string      `json:"id" db:"id"`
	WorkerID       string      `json:"workerId" db:"worker_id"`
	ShiftID        *string     `json:"shiftId,omitempty" db:"shift_id"`
	Latitude       *float64    `json:"latitude,omitempty" db:"latitude"`
	Longitude      *float64    `json:"longitude,omitempty" db:"longitude"`
	Message        *string     `json:"message,omitempty" db:"message"`
	Status         AlarmStatus `json:"status" db:"status"`
	RaisedAt       time.Time   `json:"raisedAt" db:"raised_at"`
	AcknowledgedAt *time.Time  `json:"acknowledgedAt,omitempty" db:"acknowledged_at"`
	ResolvedAt     *time.Time  `json:"resolvedAt,omitempty" db:"resolved_at"`
}

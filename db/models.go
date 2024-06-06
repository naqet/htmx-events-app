package db

import (
	"time"

	"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        string    `json:"id" sql:"type:uuid;primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	b.ID = uuid.NewV4().String()
	return nil
}

type User struct {
	Base
	Name              string   `json:"name"`
	Email             string   `json:"email" gorm:"unique"`
	Password          string   `json:"password"`
	PasswordConfirmed bool     `json:"passwordConfirmed" gorm:"default:false"`
	HostedEvents      []*Event `json:"hostedEvents" gorm:"many2many:hosted_events"`
	AttendedEvents    []*Event `json:"attendedEvents" gorm:"many2many:attended_events"`
}

type Event struct {
	Base
	Title       string `json:"title" gorm:"unique:not null;check:title <> ''"`
	Description string `json:"description"`
	// TODO: think about changing it to Place type with proper model and mock data
	Place     string        `json:"place"`
	StartDate time.Time     `json:"startDate"`
	EndDate   time.Time     `json:"endDate"`
	Hosts     []*User       `json:"hosts" gorm:"many2many:hosted_events;References:Email"`
	Attendees []*User       `json:"attendees" gorm:"many2many:attended_events;References:Email"`
	Agenda    []AgendaPoint `json:"agenda"`
}

type AgendaPoint struct {
	Base
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"startTime"`
	Event       Event     `json:"event"`
	EventId     string    `json:"eventId"`
}

type Invitation struct {
	Base
	From       User   `json:"from" gorm:"foreignKey:FromEmail;References:Email"`
	FromEmail  string `json:"fromEmail"`
	To         User   `json:"to" gorm:"foreignKey:ToEmail;References:Email"`
	ToEmail    string `json:"toEmail"`
	Event      Event  `json:"event" gorm:"foreignKey:EventTitle;References:Title"`
	EventTitle string `json:"eventTitle"`
	Message    string `json:"message"`
}

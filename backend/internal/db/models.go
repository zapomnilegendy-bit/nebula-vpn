package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Device struct {
	ID                 uuid.UUID  `json:"id"`
	UserID             uuid.UUID  `json:"user_id"`
	Name               string     `json:"name"`
	Platform           string     `json:"platform"`
	PublicKey          string     `json:"public_key"`
	PrivateKeyEncrypted string   `json:"-"`
	IPAddress          string     `json:"ip_address"`
	IsActive           bool       `json:"is_active"`
	LastConnectedAt    *time.Time `json:"last_connected_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
}

type VPNServer struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Host              string    `json:"host"`
	IPv4              string    `json:"ipv4,omitempty"`
	IPv6              string    `json:"ipv6,omitempty"`
	Port              int       `json:"port"`
	Protocol          string    `json:"protocol"`
	PublicKey         string    `json:"public_key"`
	Country           string    `json:"country"`
	City              string    `json:"city,omitempty"`
	Latitude          float64   `json:"latitude,omitempty"`
	Longitude         float64   `json:"longitude,omitempty"`
	LoadPercent       int       `json:"load_percent"`
	IsActive          bool      `json:"is_active"`
	IsHealthy         bool      `json:"is_healthy"`
	BandwidthMbps     int       `json:"bandwidth_mbps"`
	CurrentConnections int      `json:"current_connections"`
	MaxConnections    int       `json:"max_connections"`
	Tags              []string  `json:"tags,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type Subscription struct {
	ID                   uuid.UUID  `json:"id"`
	UserID               uuid.UUID  `json:"user_id"`
	PlanID               uuid.UUID  `json:"plan_id"`
	Status               string     `json:"status"`
	StartedAt            time.Time  `json:"started_at"`
	ExpiresAt            time.Time  `json:"expires_at"`
	AutoRenew            bool       `json:"auto_renew"`
	StripeSubscriptionID string     `json:"stripe_subscription_id,omitempty"`
	AppleReceipt         string     `json:"-"`
	GooglePurchaseToken  string     `json:"-"`
	CreatedAt            time.Time  `json:"created_at"`
}

type Plan struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	PriceCents      int       `json:"price_cents"`
	Currency        string    `json:"currency"`
	MaxDevices      int       `json:"max_devices"`
	MaxSpeedMbps    *int      `json:"max_speed_mbps,omitempty"`
	ProtocolsEnabled []string `json:"protocols_enabled"`
	LocationsEnabled []string `json:"locations_enabled"`
	TrialPeriodDays int       `json:"trial_period_days"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
}

type Connection struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	DeviceID  uuid.UUID  `json:"device_id"`
	ServerID  uuid.UUID  `json:"server_id"`
	Protocol  string     `json:"protocol"`
	BytesRX   int64      `json:"bytes_rx"`
	BytesTX   int64      `json:"bytes_tx"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
	Reason    string     `json:"reason,omitempty"`
}

type Promocode struct {
	ID              uuid.UUID  `json:"id"`
	Code            string     `json:"code"`
	DiscountPercent int        `json:"discount_percent"`
	MaxUses         *int       `json:"max_uses,omitempty"`
	CurrentUses     int        `json:"current_uses"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*User, error)
}

type DeviceRepository interface {
	Create(ctx context.Context, device *Device) error
	GetByID(ctx context.Context, id uuid.UUID) (*Device, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*Device, error)
	Update(ctx context.Context, device *Device) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByUserAndPlatform(ctx context.Context, userID uuid.UUID, platform string) (*Device, error)
}

type VPNServerRepository interface {
	Create(ctx context.Context, server *VPNServer) error
	GetByID(ctx context.Context, id uuid.UUID) (*VPNServer, error)
	List(ctx context.Context, country, protocol string, limit int) ([]*VPNServer, error)
	Update(ctx context.Context, server *VPNServer) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetOptimal(ctx context.Context, userCountry string, protocol string) (*VPNServer, error)
	UpdateLoad(ctx context.Context, id uuid.UUID, load int, connections int) error
}

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*Subscription, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Subscription, error)
	Update(ctx context.Context, sub *Subscription) error
	Cancel(ctx context.Context, id uuid.UUID) error
}

type PlanRepository interface {
	Create(ctx context.Context, plan *Plan) error
	GetByID(ctx context.Context, id uuid.UUID) (*Plan, error)
	List(ctx context.Context) ([]*Plan, error)
	Update(ctx context.Context, plan *Plan) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ConnectionRepository interface {
	Create(ctx context.Context, conn *Connection) error
	EndConnection(ctx context.Context, id uuid.UUID, reason string) error
	ListByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*Connection, error)
	GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*Connection, error)
	GetStats(ctx context.Context, userID uuid.UUID, since time.Time) (*ConnectionStats, error)
}

type ConnectionStats struct {
	TotalBytesRX int64 `json:"total_bytes_rx"`
	TotalBytesTX int64 `json:"total_bytes_tx"`
	TotalSessions int64 `json:"total_sessions"`
}

type PromocodeRepository interface {
	Create(ctx context.Context, code *Promocode) error
	GetByCode(ctx context.Context, code string) (*Promocode, error)
	IncrementUses(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*Promocode, error)
}

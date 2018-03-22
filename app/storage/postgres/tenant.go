package postgres

import (
	"strings"
	"time"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/env"
)

type dbTenant struct {
	ID             int    `db:"id"`
	Name           string `db:"name"`
	Subdomain      string `db:"subdomain"`
	CNAME          string `db:"cname"`
	Invitation     string `db:"invitation"`
	WelcomeMessage string `db:"welcome_message"`
	Status         int    `db:"status"`
}

func (t *dbTenant) toModel() *models.Tenant {
	if t == nil {
		return nil
	}

	return &models.Tenant{
		ID:             t.ID,
		Name:           t.Name,
		Subdomain:      t.Subdomain,
		CNAME:          t.CNAME,
		Invitation:     t.Invitation,
		WelcomeMessage: t.WelcomeMessage,
		Status:         t.Status,
	}
}

type dbEmailVerification struct {
	ID         int                          `db:"id"`
	Name       string                       `db:"name"`
	Email      string                       `db:"email"`
	Key        string                       `db:"key"`
	Kind       models.EmailVerificationKind `db:"kind"`
	UserID     dbx.NullInt                  `db:"user_id"`
	CreatedOn  time.Time                    `db:"created_on"`
	ExpiresOn  time.Time                    `db:"expires_on"`
	VerifiedOn dbx.NullTime                 `db:"verified_on"`
}

func (t *dbEmailVerification) toModel() *models.EmailVerification {
	model := &models.EmailVerification{
		Name:       t.Name,
		Email:      t.Email,
		Key:        t.Key,
		Kind:       t.Kind,
		CreatedOn:  t.CreatedOn,
		ExpiresOn:  t.ExpiresOn,
		VerifiedOn: nil,
	}

	if t.VerifiedOn.Valid {
		model.VerifiedOn = &t.VerifiedOn.Time
	}

	if t.UserID.Valid {
		model.UserID = int(t.UserID.Int64)
	}

	return model
}

// TenantStorage contains read and write operations for tenants
type TenantStorage struct {
	trx     *dbx.Trx
	current *models.Tenant
	user    *models.User
}

// NewTenantStorage creates a new TenantStorage
func NewTenantStorage(trx *dbx.Trx) *TenantStorage {
	return &TenantStorage{trx: trx}
}

// SetCurrentTenant to current context
func (s *TenantStorage) SetCurrentTenant(tenant *models.Tenant) {
	s.current = tenant
}

// SetCurrentUser to current context
func (s *TenantStorage) SetCurrentUser(user *models.User) {
	s.user = user
}

// Add given tenant to tenant list
func (s *TenantStorage) Add(name string, subdomain string, status int) (*models.Tenant, error) {
	var id int
	err := s.trx.Get(&id,
		`INSERT INTO tenants (name, subdomain, created_on, cname, invitation, welcome_message, status) 
		 VALUES ($1, $2, $3, '', '', '', $4) 
		 RETURNING id`, name, subdomain, time.Now(), status)
	if err != nil {
		return nil, err
	}

	return s.GetByDomain(subdomain)
}

// First returns first tenant
func (s *TenantStorage) First() (*models.Tenant, error) {
	tenant := dbTenant{}

	err := s.trx.Get(&tenant, "SELECT id, name, subdomain, cname, invitation, welcome_message, status FROM tenants ORDER BY id LIMIT 1")
	if err != nil {
		return nil, err
	}

	return tenant.toModel(), nil
}

// GetByDomain returns a tenant based on its domain
func (s *TenantStorage) GetByDomain(domain string) (*models.Tenant, error) {
	tenant := dbTenant{}

	err := s.trx.Get(&tenant, "SELECT id, name, subdomain, cname, invitation, welcome_message, status FROM tenants WHERE subdomain = $1 OR cname = $2 ORDER BY cname DESC", extractSubdomain(domain), domain)
	if err != nil {
		return nil, err
	}

	return tenant.toModel(), nil
}

// UpdateSettings of given tenant
func (s *TenantStorage) UpdateSettings(settings *models.UpdateTenantSettings) error {
	query := "UPDATE tenants SET name = $1, invitation = $2, welcome_message = $3, cname = $4 WHERE id = $5"
	_, err := s.trx.Execute(query, settings.Title, settings.Invitation, settings.WelcomeMessage, settings.CNAME, s.current.ID)
	return err
}

// IsSubdomainAvailable returns true if subdomain is available to use
func (s *TenantStorage) IsSubdomainAvailable(subdomain string) (bool, error) {
	exists, err := s.trx.Exists("SELECT id FROM tenants WHERE subdomain = $1", subdomain)
	return !exists, err
}

// IsCNAMEAvailable returns true if cname is available to use
func (s *TenantStorage) IsCNAMEAvailable(cname string) (bool, error) {
	exists, err := s.trx.Exists("SELECT id FROM tenants WHERE cname = $1 AND id <> $2", cname, s.current.ID)
	return !exists, err
}

// Activate given tenant
func (s *TenantStorage) Activate(id int) error {
	query := "UPDATE tenants SET status = $1 WHERE id = $2"
	_, err := s.trx.Execute(query, models.TenantActive, id)
	return err
}

// SaveVerificationKey used by email verification process
func (s *TenantStorage) SaveVerificationKey(key string, duration time.Duration, request models.NewEmailVerification) error {
	var userID interface{}
	if request.GetUser() != nil {
		userID = request.GetUser().ID
	}
	query := "INSERT INTO email_verifications (tenant_id, email, created_on, expires_on, key, name, kind, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := s.trx.Execute(query, s.current.ID, request.GetEmail(), time.Now(), time.Now().Add(duration), key, request.GetName(), request.GetKind(), userID)
	return err
}

// FindVerificationByKey based on current tenant
func (s *TenantStorage) FindVerificationByKey(kind models.EmailVerificationKind, key string) (*models.EmailVerification, error) {
	verification := dbEmailVerification{}

	query := "SELECT id, email, name, key, created_on, verified_on, expires_on, kind, user_id FROM email_verifications WHERE key = $1 AND kind = $2 LIMIT 1"
	err := s.trx.Get(&verification, query, key, kind)
	if err != nil {
		return nil, err
	}

	return verification.toModel(), nil
}

// SetKeyAsVerified so that it cannot be used anymore
func (s *TenantStorage) SetKeyAsVerified(key string) error {
	query := "UPDATE email_verifications SET verified_on = $1 WHERE tenant_id = $2 AND key = $3"
	_, err := s.trx.Execute(query, time.Now(), s.current.ID, key)
	return err
}

func extractSubdomain(hostname string) string {
	domain := env.MultiTenantDomain()
	if domain == "" {
		return hostname
	}

	return strings.Replace(hostname, domain, "", -1)
}

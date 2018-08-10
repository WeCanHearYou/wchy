package worker

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/uuid"
)

//Context holds references to services available for jobs
type Context struct {
	workerID string
	taskName string
	services *app.Services
	logger   log.Logger
	db       *dbx.Database
	baseURL  string
	user     *models.User
	tenant   *models.Tenant
}

//NewContext creates a new context
func NewContext(workerID, taskName string, db *dbx.Database, logger log.Logger) *Context {
	ctxLogger := logger.New()
	contextID := uuid.NewV4().String()
	ctxLogger.SetProperty(log.PropertyKeyContextID, contextID)

	return &Context{
		workerID: workerID,
		taskName: taskName,
		db:       db,
		logger:   ctxLogger,
	}
}

//SetBaseURL on context
func (c *Context) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

//SetUser on context
func (c *Context) SetUser(user *models.User) {
	c.user = user
	if user != nil {
		c.logger.SetProperty(log.PropertyKeyUserID, user.ID)
	}
	if c.services != nil {
		c.services.SetCurrentUser(user)
	}
}

//SetTenant on context
func (c *Context) SetTenant(tenant *models.Tenant) {
	c.tenant = tenant
	if tenant != nil {
		c.logger.SetProperty(log.PropertyKeyTenantID, tenant.ID)
	}
	if c.services != nil {
		c.services.SetCurrentTenant(tenant)
	}
}

//SetServices on current context
func (c *Context) SetServices(services *app.Services) {
	c.services = services
}

//WorkerID executing current context
func (c *Context) WorkerID() string {
	return c.workerID
}

//TaskName from current context
func (c *Context) TaskName() string {
	return c.taskName
}

//BaseURL from current context
func (c *Context) BaseURL() string {
	return c.baseURL
}

//User from current context
func (c *Context) User() *models.User {
	return c.user
}

//Tenant from current context
func (c *Context) Tenant() *models.Tenant {
	return c.tenant
}

//Services from current context
func (c *Context) Services() *app.Services {
	return c.services
}

//Logger from current context
func (c *Context) Logger() log.Logger {
	return c.logger
}

//Database from current context
func (c *Context) Database() *dbx.Database {
	return c.db
}

//Failure logs details of error
func (c *Context) Failure(err error) error {
	err = errors.StackN(err, 1)
	c.logger.Error(err)
	return err
}

package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"

	"github.com/getfider/fider/app/handlers"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/mock"
)

func TestIndexHandler(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.OnTenant(mock.DemoTenant).AsUser(mock.JonSnow).Execute(handlers.Index())

	Expect(code).Equals(http.StatusOK)
}

func TestDetailsHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	post, _ := services.Posts.Add("My Post", "My Post Description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		Execute(handlers.PostDetails())

	Expect(code).Equals(http.StatusOK)
}

func TestDetailsHandler_NotFound(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", "99").
		Execute(handlers.PostDetails())

	Expect(code).Equals(http.StatusNotFound)
}

func TestUpdatePostHandler_TenantStaff(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.AryaStark)
	post, _ := services.Posts.Add("My First Post", "With a description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(handlers.UpdatePost(), `{ "title": "the new title", "description": "new description" }`)

	post, _ = services.Posts.GetByNumber(post.Number)
	Expect(code).Equals(http.StatusOK)
	Expect(post.Title).Equals("the new title")
	Expect(post.Description).Equals("new description")
}

func TestUpdatePostHandler_NonAuthorized(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	post, _ := services.Posts.Add("My First Post", "With a description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", post.Number).
		ExecutePost(handlers.UpdatePost(), `{ "title": "the new title", "description": "new description" }`)

	Expect(code).Equals(http.StatusForbidden)
}

func TestUpdatePostHandler_InvalidTitle(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	post, _ := services.Posts.Add("My First Post", "With a description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(handlers.UpdatePost(), `{ "title": "", "description": "" }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestUpdatePostHandler_InvalidPost(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", 999).
		ExecutePost(handlers.UpdatePost(), `{ "title": "This is a good title!", "description": "And description too..." }`)

	Expect(code).Equals(http.StatusNotFound)
}

func TestUpdatePostHandler_DuplicateTitle(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	post1, _ := services.Posts.Add("My First Post", "With a description")
	services.Posts.Add("My Second Post", "With a description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post1.Number).
		ExecutePost(handlers.UpdatePost(), `{ "title": "My Second Post", "description": "And description too..." }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestPostCommentHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	post, _ := services.Posts.Add("My First Post", "With a description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(handlers.PostComment(), `{ "content": "This is a comment!" }`)

	Expect(code).Equals(http.StatusOK)
}

func TestPostCommentHandler_WithoutContent(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	post, _ := services.Posts.Add("My First Post", "With a description")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(handlers.PostComment(), `{ "content": "" }`)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestAddSupporterHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	first, _ := services.Posts.Add("The Post #1", "The Description #1")
	second, _ := services.Posts.Add("The Post #2", "The Description #2")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", second.Number).
		Execute(handlers.AddSupporter())

	first, _ = services.Posts.GetByNumber(1)
	second, _ = services.Posts.GetByNumber(2)

	Expect(code).Equals(http.StatusOK)
	Expect(first.TotalSupporters).Equals(0)
	Expect(second.TotalSupporters).Equals(1)
}

func TestAddSupporterHandler_InvalidPost(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", 999).
		Execute(handlers.AddSupporter())

	Expect(code).Equals(http.StatusNotFound)
}

func TestRemoveSupporterHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	post, _ := services.Posts.Add("The Post #1", "The Description #1")
	services.Posts.AddSupporter(post, mock.JonSnow)
	services.Posts.AddSupporter(post, mock.AryaStark)

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", post.ID).
		Execute(handlers.RemoveSupporter())

	post, _ = services.Posts.GetByNumber(post.Number)

	Expect(code).Equals(http.StatusOK)
	Expect(post.TotalSupporters).Equals(1)
}

func TestSetResponseHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.AryaStark)
	post, _ := services.Posts.Add("The Post #1", "The Description #1")
	services.Posts.AddSupporter(post, mock.AryaStark)

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.ID).
		ExecutePost(handlers.SetResponse(), fmt.Sprintf(`{ "status": %d, "text": "Done!" }`, models.PostCompleted))

	post, _ = services.Posts.GetByNumber(post.Number)

	Expect(code).Equals(http.StatusOK)
	Expect(post.Status).Equals(models.PostCompleted)
	Expect(post.Response.Text).Equals("Done!")
	Expect(post.Response.User.ID).Equals(mock.JonSnow.ID)
}

func TestSetResponseHandler_Unauthorized(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.AryaStark)
	post, _ := services.Posts.Add("The Post #1", "The Description #1")
	services.Posts.AddSupporter(post, mock.AryaStark)

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", post.ID).
		ExecutePost(handlers.SetResponse(), fmt.Sprintf(`{ "status": %d, "text": "Done!" }`, models.PostCompleted))

	post, _ = services.Posts.GetByNumber(post.Number)

	Expect(code).Equals(http.StatusForbidden)
	Expect(post.Status).Equals(models.PostOpen)
}

func TestSetResponseHandler_Duplicate(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.AryaStark)
	post1, _ := services.Posts.Add("The Post #1", "The Description #1")
	post2, _ := services.Posts.Add("The Post #2", "The Description #2")

	body := fmt.Sprintf(`{ "status": %d, "originalNumber": %d }`, models.PostDuplicate, post2.Number)
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post1.ID).
		ExecutePost(handlers.SetResponse(), body)
	Expect(code).Equals(http.StatusOK)

	post1, _ = services.Posts.GetByNumber(post1.Number)
	Expect(post1.Status).Equals(models.PostDuplicate)

	post2, _ = services.Posts.GetByNumber(post2.Number)
	Expect(post2.Status).Equals(models.PostOpen)
}

func TestSetResponseHandler_Duplicate_NotFound(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.AryaStark)
	post1, _ := services.Posts.Add("The Post #1", "The Description #1")

	body := fmt.Sprintf(`{ "status": %d, "originalNumber": 9999 }`, models.PostDuplicate)
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post1.ID).
		ExecutePost(handlers.SetResponse(), body)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestSetResponseHandler_Duplicate_Itself(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.AryaStark)
	post, _ := services.Posts.Add("The Post #1", "The Description #1")

	body := fmt.Sprintf(`{ "status": %d, "originalNumber": %d }`, models.PostDuplicate, post.Number)
	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.ID).
		ExecutePost(handlers.SetResponse(), body)

	Expect(code).Equals(http.StatusBadRequest)
}

func TestAddCommentHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.AryaStark)
	post, _ := services.Posts.Add("The Post #1", "The Description #1")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(handlers.PostComment(), `{ "content": "My first comment" }`)

	Expect(code).Equals(http.StatusOK)
	comments, _ := services.Posts.GetCommentsByPost(post)
	Expect(comments).HasLen(1)
}

func TestUpdateCommentHandler_Authorized(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.AryaStark)
	post, _ := services.Posts.Add("The Post #1", "The Description #1")
	commentId, _ := services.Posts.AddComment(post, "My first comment")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", post.Number).
		AddParam("id", commentId).
		ExecutePost(handlers.UpdateComment(), `{ "content": "My first comment has been edited" }`)

	Expect(code).Equals(http.StatusOK)
	comment, _ := services.Posts.GetCommentByID(commentId)
	Expect(comment.Content).Equals("My first comment has been edited")
}

func TestUpdateCommentHandler_Unauthorized(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	post, _ := services.Posts.Add("The Post #1", "The Description #1")
	commentId, _ := services.Posts.AddComment(post, "My first comment")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.AryaStark).
		AddParam("number", post.Number).
		AddParam("id", commentId).
		ExecutePost(handlers.UpdateComment(), `{ "content": "My first comment has been edited" }`)

	Expect(code).Equals(http.StatusForbidden)
	comment, _ := services.Posts.GetCommentByID(commentId)
	Expect(comment.Content).Equals("My first comment")
}

func TestDeletePostHandler_Authorized(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)
	post, _ := services.Posts.Add("The Post #1", "The Description #1")

	code, _ := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("number", post.Number).
		ExecutePost(handlers.DeletePost(), `{ }`)

	Expect(code).Equals(http.StatusOK)
	post, err := services.Posts.GetByNumber(post.Number)
	Expect(post).IsNil()
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
}

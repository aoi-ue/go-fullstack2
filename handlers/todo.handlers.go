package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/emarifer/gofiber-templ-htmx/models"
	"github.com/emarifer/gofiber-templ-htmx/views/todo_views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/sujit-baniya/flash"
)

/********** Handlers for Todo Views **********/

// Render List Page with success/error messages
func HandleViewList(c *fiber.Ctx) error {
	todo := new(models.Todo)
	todo.CreatedBy = c.Locals("userId").(uint64)

	fm := fiber.Map{
		"type": "error",
	}

	todosSlice, err := todo.GetAllTodos()
	if err != nil {
		fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		return flash.WithError(c, fm).Redirect("/todo/list")
	}

	tindex := todo_views.TodoIndex(todosSlice)
	tlist := todo_views.TodoList(
		" | Tasks List",
		fromProtected,
		flash.Get(c),
		c.Locals("username").(string),
		tindex,
	)

	handler := adaptor.HTTPHandler(templ.Handler(tlist))

	return handler(c)
}

// Render Create Todo Page with success/error messages
func HandleViewCreatePage(c *fiber.Ctx) error {

	if c.Method() == "POST" {
		fm := fiber.Map{
			"type": "error",
		}

		newTodo := new(models.Todo)
		newTodo.CreatedBy = c.Locals("userId").(uint64)
		newTodo.Title = strings.Trim(c.FormValue("title"), " ")
		newTodo.Description = strings.Trim(c.FormValue("description"), " ")

		if _, err := newTodo.CreateTodo(); err != nil {
			fm["message"] = fmt.Sprintf("something went wrong: %s", err)

			return flash.WithError(c, fm).Redirect("/todo/list")
		}

		return c.Redirect("/todo/list")
	}

	cindex := todo_views.CreateIndex()
	create := todo_views.Create(
		" | Create Todo",
		fromProtected,
		flash.Get(c),
		c.Locals("username").(string),
		cindex,
	)

	handler := adaptor.HTTPHandler(templ.Handler(create))

	return handler(c)
}

// Render Edit Todo Page with success/error messages
func HandleViewEditPage(c *fiber.Ctx) error {
	idParams, _ := strconv.Atoi(c.Params("id"))
	todoId := uint64(idParams)

	todo := new(models.Todo)
	todo.ID = todoId
	todo.CreatedBy = c.Locals("userId").(uint64)

	fm := fiber.Map{
		"type": "error",
	}

	recoveredTodo, err := todo.GetNoteById()
	if err != nil {
		fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		return flash.WithError(c, fm).Redirect("/todo/list")
	}

	if c.Method() == "POST" {
		todo.Title = strings.Trim(c.FormValue("title"), " ")
		todo.Description = strings.Trim(c.FormValue("description"), " ")
		if c.FormValue("status") == "on" {
			todo.Status = true
		} else {
			todo.Status = false
		}

		_, err := todo.UpdateTodo()
		if err != nil {
			fm["message"] = fmt.Sprintf("something went wrong: %s", err)

			return flash.WithError(c, fm).Redirect("/todo/list")
		}

		fm = fiber.Map{
			"type":    "success",
			"message": "Task successfully updated!!",
		}

		return flash.WithSuccess(c, fm).Redirect("/todo/list")
	}

	uindex := todo_views.UpdateIndex(recoveredTodo)
	update := todo_views.Update(
		fmt.Sprintf(" | Edit Todo #%d", recoveredTodo.ID),
		fromProtected,
		flash.Get(c),
		c.Locals("username").(string),
		uindex,
	)

	handler := adaptor.HTTPHandler(templ.Handler(update))

	return handler(c)
}

// Handler Remove Todo
func HandleDeleteTodo(c *fiber.Ctx) error {
	idParams, _ := strconv.Atoi(c.Params("id"))
	todoId := uint64(idParams)

	todo := new(models.Todo)
	todo.ID = todoId
	todo.CreatedBy = c.Locals("userId").(uint64)

	fm := fiber.Map{
		"type": "error",
	}

	if err := todo.DeleteTodo(); err != nil {
		fm["message"] = fmt.Sprintf("something went wrong: %s", err)

		return flash.WithError(c, fm).Redirect(
			"/todo/list",
			fiber.StatusSeeOther,
		)
	}

	fm = fiber.Map{
		"type":    "success",
		"message": "Task successfully deleted!!",
	}

	return flash.WithSuccess(c, fm).Redirect("/todo/list", fiber.StatusSeeOther)
}

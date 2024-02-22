package repository

import (
	"fmt"
	todo "github.com/h1lton/mz-course"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listId int, item todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) values ($1, $2) RETURNING id", todoItemsTable)

	var itemId int
	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	if err = row.Scan(&itemId); err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) values ($1, $2)", itemsListsTable)
	_, err = tx.Exec(createListItemsQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem
	query := fmt.Sprintf(
		`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s il on il.item_id = ti.id
		INNER JOIN %s ul on ul.list_id = il.list_id WHERE ul.list_id = $1 AND ul.user_id = $2`,
		todoItemsTable, itemsListsTable, usersListsTable,
	)

	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItemPostgres) GetById(userId, itemId int) (todo.TodoItem, error) {
	var item todo.TodoItem
	query := fmt.Sprintf(
		`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s il on il.item_id = ti.id
		INNER JOIN %s ul on ul.list_id = il.list_id WHERE ti.id = $1 AND ul.user_id = $2`,
		todoItemsTable, itemsListsTable, usersListsTable,
	)

	err := r.db.Get(&item, query, itemId, userId)

	return item, err
}

func (r *TodoItemPostgres) Delete(userId, itemId int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s ti USING %s il, %s ul 
       WHERE ti.id = il.item_id AND il.list_id = ul.list_id AND ti.id = $1 AND ul.user_id = $2`,
		todoItemsTable, itemsListsTable, usersListsTable,
	)

	_, err := r.db.Exec(query, itemId, userId)

	return err
}

func (r *TodoItemPostgres) Update(userId, itemId int, input todo.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]any, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done = $%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(
		`UPDATE %s ti SET %s FROM %s il, %s ul WHERE ti.id = il.item_id AND ul.list_id = il.list_id 
				  AND ti.id = $%d AND ul.user_id = $%d`,
		todoItemsTable, setQuery, itemsListsTable, usersListsTable, argId, argId+1)

	args = append(args, itemId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)

	return err
}

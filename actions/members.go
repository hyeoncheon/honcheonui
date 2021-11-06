package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"

	"github.com/hyeoncheon/honcheonui/models"
)

// MembersResource is the resource for the Member model
type MembersResource struct {
	buffalo.Resource
}

// List gets all Members.
//! MANAGER ONLY, currently not protected
func (v MembersResource) List(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	members := &models.Members{}
	q := tx.PaginateFromParams(c.Params())
	if err := q.All(members); err != nil {
		return errors.WithStack(err)
	}

	c.Set("pagination", q.Paginator)
	return c.Render(200, r.Auto(c, members))
}

// Show gets the data for one Member.
//! MANAGER ONLY, currently not protected
func (v MembersResource) Show(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	member := &models.Member{}
	if err := tx.Find(member, c.Param("member_id")); err != nil {
		return c.Error(404, err)
	}

	// TODO: which services
	// TODO: which resources
	// TODO: which providers

	return c.Render(200, r.Auto(c, member))
}

// Edit renders a edit form for a Member.
// for editing preferences, currently not used.
func (v MembersResource) Edit(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	member := &models.Member{}
	if err := tx.Find(member, c.Param("member_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, member))
}

// Update changes a Member in the DB.
// for updating preferences, currently not used.
func (v MembersResource) Update(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	member := &models.Member{}
	if err := tx.Find(member, c.Param("member_id")); err != nil {
		return c.Error(404, err)
	}

	if err := c.Bind(member); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(member)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(422, r.Auto(c, member))
	}

	c.Flash().Add("success", "Member was updated successfully")
	return c.Render(200, r.Auto(c, member))
}

// Destroy deletes a Member from the DB.
//! it does not revoke users grant and administrator's acceptance of access.
func (v MembersResource) Destroy(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	member := &models.Member{}
	if err := tx.Find(member, c.Param("member_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(member); err != nil {
		return errors.WithStack(err)
	}

	c.Flash().Add("success", "Member was destroyed successfully")
	return c.Render(200, r.Auto(c, member))
}

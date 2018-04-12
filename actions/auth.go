package actions

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/cloudfoundry"
	"github.com/pkg/errors"

	"github.com/hyeoncheon/honcheonui/models"
)

func init() {
	gothic.Store = App().SessionStore

	uartProvider := cloudfoundry.New(
		os.Getenv("UART_URL"),
		os.Getenv("UART_KEY"),
		os.Getenv("UART_SECRET"),
		fmt.Sprintf("%s%s", os.Getenv("HCU_URL"), "/auth/uart/callback"),
		"profile")
	uartProvider.SetName("uart")

	goth.UseProviders(
		uartProvider,
	)
}

// AuthCallback is universal callback handler for goth authorization
func AuthCallback(c buffalo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}
	c.Logger().Debugf("raw userdata: %v", r.JSON(user))

	// reach here means, user granted access and success OAuth2 sequence.
	// anyway, we need to check the person has the right for this app.
	if err := validateMembership(&user); err != nil {
		c.Logger().Warnf("user validation failed: %v", err)
		c.Flash().Add("danger", t(c, err.Error()))
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	member := &models.Member{}
	if err := tx.Find(member, user.UserID); err != nil {
		member, err = createMember(tx, &user)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	// name, avatar icon, and roles are not stored on database.
	// always refresh with information from authorization provider.
	member.Name = user.Name
	member.Avatar = user.RawData["picture"].(string)
	for _, v := range user.RawData["roles"].([]interface{}) {
		if r, ok := v.(string); ok {
			member.Roles = append(member.Roles, r)
		}
	}

	// NOTE: set initial session data for this login session
	sess := c.Session()
	sess.Set("member_id", member.ID)
	sess.Set("member_mail", member.Email)
	sess.Set("member_name", member.Name)
	sess.Set("member_icon", member.Avatar)
	sess.Set("member_roles", member.Roles)
	// Do something with the user, maybe register them/sign them in
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func validateMembership(u *goth.User) error {
	if u.Email == "" {
		return errors.New("invalid.membership..email.is.not.provided")
	}
	if u.Name == "" {
		return errors.New("invalid.membership..name.is.not.provided")
	}
	roles, ok := u.RawData["roles"].([]interface{})
	if !ok || len(roles) < 1 {
		return errors.New("invalid.membership..not.enough.roles")
	}
	return nil
}

func createMember(tx *pop.Connection, u *goth.User) (*models.Member, error) {
	member := &models.Member{Email: u.Email}
	if id, err := uuid.FromString(u.UserID); err == nil {
		member.ID = id
	} else {
		return &models.Member{}, errors.New("cannot.create.member..invalid.userid")
	}

	verrs, err := tx.ValidateAndCreate(member)
	if err != nil {
		return &models.Member{}, err
	}
	if verrs.HasAny() {
		return &models.Member{}, errors.WithStack(errors.New("cannot.create.member..validation.error"))
	}
	// TODO: add HTTP provider as default
	return member, nil
}

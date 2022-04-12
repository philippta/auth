package auth

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

// userKey is the type used to store the user in the context.
type userKey struct{}

// LookupFunc is a function that returns a user for a given id.
type LookupFunc[User any, ID any] func(id ID) (*User, bool)

// New returns a new Auth.
func New[User any, ID any](f LookupFunc[User, ID]) *Auth[User, ID] {
	return &Auth[User, ID]{
		user:           f,
		SessionManager: scs.New(),
	}
}

// Auth is an authentication middleware.
type Auth[User any, ID any] struct {
	user           LookupFunc[User, ID]
	SessionManager *scs.SessionManager
}

// Handler returns a handler that wraps h and authenticates requests.
func (a *Auth[User, ID]) Handler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		userID, ok := a.SessionManager.Get(r.Context(), "user_id").(ID)
		if !ok {
			h.ServeHTTP(w, r)
			return
		}

		user, ok := a.user(userID)
		if !ok {
			a.SessionManager.Remove(r.Context(), "user_id")
			h.ServeHTTP(w, r)
			return
		}

		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userKey{}, user)))
	}
	return a.SessionManager.LoadAndSave(http.HandlerFunc(fn))
}

// User returns the user associated with the request.
func (a *Auth[User, ID]) User(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userKey{}).(*User)
	return user, ok
}

// Ok returns true if the request is authenticated.
func (a *Auth[User, ID]) Ok(ctx context.Context) bool {
	_, ok := a.User(ctx)
	return ok
}

// Login logs in the user with the given id.
func (a *Auth[User, ID]) Login(ctx context.Context, id ID) {
	a.SessionManager.Put(ctx, "user_id", id)
}

// Logout logs out the user.
func (a *Auth[User, ID]) Logout(ctx context.Context) {
	a.SessionManager.Remove(ctx, "user_id")
}

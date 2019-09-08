package sessions

import (
	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("n1XZu33ku4V"))
package app

import (
	"net/http"
	"time"

	"github.com/acoshift/hime"
	"github.com/acoshift/httprouter"
	"github.com/acoshift/middleware"
	"github.com/acoshift/session"
	redisstore "github.com/acoshift/session/store/redis"
	"github.com/asaskevich/govalidator"
)

// New creates new app
func New(config Config) hime.HandlerFactory {
	return func(app hime.App) http.Handler {
		app.
			TemplateDir("template").
			TemplateRoot("root").
			Component("_layout.tmpl").
			Template("index", "index.tmpl").
			Template("signin", "signin.tmpl").
			BeforeRender(beforeRender).
			Routes(hime.Routes{
				"index":   "/",
				"signin":  "/signin",
				"signout": "/signout",
			})

		mux := http.NewServeMux()

		r := httprouter.New()
		r.Get(app.Route("index"), hime.H(indexHandler))
		r.Get(app.Route("signin"), hime.H(signInGetHandler))
		r.Post(app.Route("signin"), hime.H(signInPostHandler))
		r.Get(app.Route("signout"), hime.H(signOutHandler))

		mux.Handle("/", r)

		return middleware.Chain(
			session.Middleware(session.Config{
				Store: redisstore.New(redisstore.Config{
					Pool:   config.RedisPool,
					Prefix: config.RedisPrefix,
				}),
				HTTPOnly:         true,
				Secure:           session.PreferSecure,
				Path:             "/",
				Rolling:          true,
				SameSite:         session.SameSiteLax,
				Proxy:            true,
				MaxAge:           7 * 24 * time.Hour,
				DeleteOldSession: true,
			}),
		)(mux)
	}
}

func beforeRender(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getSession(r.Context()).Flash().Clear()
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		h.ServeHTTP(w, r)
	})
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func indexHandler(ctx hime.Context) hime.Result {
	return ctx.View("index", newPage(ctx))
}

func signInGetHandler(ctx hime.Context) hime.Result {
	return ctx.View("signin", newPage(ctx))
}

func signInPostHandler(ctx hime.Context) hime.Result {
	sess := getSession(ctx)
	f := sess.Flash()

	email := ctx.PostFormValueTrimSpace("email")
	password := ctx.PostFormValue("password")

	if email == "" {
		f.Add("Errors", "email required")
	} else if !govalidator.IsEmail(email) {
		f.Add("Errors", "invalid email")
	}
	if password == "" {
		f.Add("Errors", "password required")
	}

	if f.Has("Errors") {
		return ctx.RedirectToGet()
	}

	sess.Regenerate()
	sess.Set("user_id", "1")
	return ctx.RedirectTo("index")
}

func signOutHandler(ctx hime.Context) hime.Result {
	getSession(ctx).Destroy()
	return ctx.RedirectTo("index")
}

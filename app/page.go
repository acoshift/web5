package app

import "github.com/acoshift/hime"

func newPage(ctx hime.Context) map[string]interface{} {
	sess := getSession(ctx)

	return map[string]interface{}{
		"Flash":    sess.Flash().Values(),
		"IsSignIn": sess.GetString("user_id") != "",
	}
}

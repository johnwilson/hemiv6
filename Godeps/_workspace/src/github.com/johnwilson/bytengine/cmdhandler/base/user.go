package base

import (
	"github.com/johnwilson/bytengine"
)

// handler for: user.new
func UserNew(cmd bytengine.Command, user *bytengine.User, eng *bytengine.Engine) (interface{}, error) {
	usr := cmd.Args["username"].(string)
	pw := cmd.Args["password"].(string)
	err := eng.Authentication.NewUser(usr, pw, false)
	if err != nil {
		return nil, err
	}
	return true, nil
}

// handler for: user.all
func UserAll(cmd bytengine.Command, user *bytengine.User, eng *bytengine.Engine) (interface{}, error) {
	rgx := "."
	val, ok := cmd.Options["regex"]
	if ok {
		rgx = val.(string)
	}
	users, err := eng.Authentication.ListUser(rgx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// handler for: user.about
func UserAbout(cmd bytengine.Command, user *bytengine.User, eng *bytengine.Engine) (interface{}, error) {
	usr := cmd.Args["username"].(string)
	info, err := eng.Authentication.UserInfo(usr)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// handler for: user.delete
func UserDelete(cmd bytengine.Command, user *bytengine.User, eng *bytengine.Engine) (interface{}, error) {
	usr := cmd.Args["username"].(string)
	err := eng.Authentication.RemoveUser(usr)
	if err != nil {
		return nil, err
	}
	return true, nil
}

// handler for: user.passw
func UserPassw(cmd bytengine.Command, user *bytengine.User, eng *bytengine.Engine) (interface{}, error) {
	usr := cmd.Args["username"].(string)
	pw := cmd.Args["password"].(string)
	err := eng.Authentication.ChangeUserPassword(usr, pw)
	if err != nil {
		return nil, err
	}
	return true, nil
}

// handler for: user.access
func UserAccess(cmd bytengine.Command, user *bytengine.User, eng *bytengine.Engine) (interface{}, error) {
	usr := cmd.Args["username"].(string)
	grant := cmd.Args["grant"].(bool)
	err := eng.Authentication.ChangeUserStatus(usr, grant)
	if err != nil {
		return nil, err
	}
	return true, nil
}

// handler for: user.db
func UserDb(cmd bytengine.Command, user *bytengine.User, eng *bytengine.Engine) (interface{}, error) {
	usr := cmd.Args["username"].(string)
	grant := cmd.Args["grant"].(bool)
	db := cmd.Args["database"].(string)
	err := eng.Authentication.ChangeUserDbAccess(usr, db, grant)
	if err != nil {
		return nil, err
	}
	return true, nil
}

// handler for: user.whoami
func UserWhoami(cmd bytengine.Command, user *bytengine.User, eng *bytengine.Engine) (interface{}, error) {
	val := map[string]interface{}{
		"username":  user.Username,
		"databases": user.Databases,
		"root":      user.Root,
	}
	return val, nil
}

func init() {
	bytengine.RegisterCommandHandler("user.new", UserNew)
	bytengine.RegisterCommandHandler("user.all", UserAll)
	bytengine.RegisterCommandHandler("user.about", UserAbout)
	bytengine.RegisterCommandHandler("user.delete", UserDelete)
	bytengine.RegisterCommandHandler("user.passw", UserPassw)
	bytengine.RegisterCommandHandler("user.access", UserAccess)
	bytengine.RegisterCommandHandler("user.db", UserDb)
	bytengine.RegisterCommandHandler("user.whoami", UserWhoami)
}

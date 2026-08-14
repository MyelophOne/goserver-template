package app

/*
!! DO NOT DELETE OR MODIFY THIS FILE !!

!! It is used by template and will be overwritten by updates. !!

!! It shows an example of using @myelophone/goserver inside the project. !!
*/

import goserver "github.com/myelophone/goserver"

func init() {
	goserver.RegisterHook(func(s *goserver.Server) {
		if goserver.IsDev() {
			s.Logger.Print("app hooks registered")
		}
	})
}

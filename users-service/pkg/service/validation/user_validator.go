package validation

import (
	"github.com/MaxGoodfella/ivanberries-max/users-service/pkg/util"
	"strings"
)

func ValidateUser(email, password, roleID string) error {
	if strings.TrimSpace(email) == "" {
		return util.ErrUserEmailEmpty
	}
	if strings.TrimSpace(password) == "" {
		return util.ErrUserPasswordEmpty
	}
	if !util.ContainsLetter.MatchString(email) {
		return util.ErrUserEmailNoLetter
	}
	if !util.EmailRegexp.MatchString(email) {
		return util.ErrUserEmailInvalid
	}
	if !util.PasswordRegexp.MatchString(password) {
		return util.ErrUserPasswordInvalid
	}
	if strings.TrimSpace(roleID) == "" {
		return util.ErrUserRoleIDInvalid
	}
	return nil
}

package iredmail

import (
	"fmt"
)

func (s *Server) UserAliasAdd(alias, email string) error {
	_, domain := parseEmail(email)
	a := fmt.Sprintf("%v@%v", alias, domain)

	userExists, err := s.userExists(a)
	if err != nil {
		return err
	}
	if userExists {
		return fmt.Errorf("An user with %v already exists", a)
	}

	aliasExists, err := s.aliasExists(a)
	if err != nil {
		return err
	}
	if aliasExists {
		return fmt.Errorf("An alias with %v already exists", a)
	}

	_, err = s.DB.Exec(`
		INSERT INTO forwardings (address, forwarding, domain, dest_domain, is_alias, active)
		VALUES ('` + a + `', '` + email + `', '` + domain + `', '` + domain + `', 1, 1)
	`)

	return err
}

func (s *Server) UserAliasDelete(aliasEmail string) error {
	aliasExists, err := s.aliasExists(aliasEmail)
	if err != nil {
		return err
	}
	if !aliasExists {
		return fmt.Errorf("An alias with %v doesn't exists", aliasEmail)
	}

	_, err = s.DB.Exec(`
		DELETE FROM forwardings WHERE address = '` + aliasEmail + `' AND is_alias = 1
	`)

	return err
}

func (s *Server) UserAliasDeleteAll(userEmail string) error {
	_, err := s.DB.Exec(`
		DELETE FROM forwardings WHERE forwarding = '` + userEmail + `' AND is_alias = 1
	`)

	return err
}
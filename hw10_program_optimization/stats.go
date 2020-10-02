package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type User struct {
	ID       int    `json:"-"`
	Name     string `json:"-"`
	Username string `json:"-"`
	Email    string `json:"email"`
	Phone    string `json:"-"`
	Password string `json:"-"`
	Address  string `json:"-"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	reader := bufio.NewReader(r)
	domStat := make(DomainStat)
	domain = "." + domain
	u := &User{}

	for {
		line, _, err := reader.ReadLine()
		if errors.Is(err, io.EOF) {
			return domStat, nil
		} else if err != nil {
			return nil, fmt.Errorf("reading error: %w", err)
		}

		if err := json.Unmarshal(line, u); err != nil {
			return nil, fmt.Errorf(`string(%s) unmarshal error: %w`, line, err)
		}

		if !strings.HasSuffix(u.Email, domain) {
			continue
		}

		email := strings.Split(u.Email, "@")
		if len(email) != 2 {
			return nil, fmt.Errorf(`invalid email: %s`, u.Email)
		}

		emailDomain := strings.ToLower(email[1])
		domStat[emailDomain]++
	}
}

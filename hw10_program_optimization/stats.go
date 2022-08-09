package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strings"

	"github.com/mailru/easyjson"
)

var ErrReaderIsNil = errors.New("reader instance must not be nil")

const separator = "@"

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if r == nil {
		return nil, ErrReaderIsNil
	}

	result := make(DomainStat)
	reg := regexp.MustCompile("\\." + domain)

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	user := &User{}
	for scanner.Scan() {
		err := easyjson.Unmarshal(scanner.Bytes(), user)
		if err != nil {
			return nil, err
		}

		if reg.Match([]byte(user.Email)) {
			f := strings.ToLower(strings.SplitN(user.Email, separator, 2)[1])
			result[f]++
		}
	}

	return result, nil
}

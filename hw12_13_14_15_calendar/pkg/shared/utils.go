package shared

import (
	"fmt"
	"strings"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/configuration"
)

func CreatePgConnectionString(conf configuration.StorageConf) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("host=%s port=%d", conf.Host, conf.Port))

	if conf.User != "" {
		sb.WriteString(fmt.Sprintf(" user=%s", conf.User))
	}
	if conf.Password != "" {
		sb.WriteString(fmt.Sprintf(" password=%s", conf.Password))
	}
	if conf.Name != "" {
		sb.WriteString(fmt.Sprintf(" dbname=%s", conf.Name))
	}

	sb.WriteString(" sslmode=disable")
	return sb.String()
}

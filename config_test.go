package main

import "testing"

func Test_UnmarshallConfig(t *testing.T) {

	configdata := `
query: in:inbox is:unread
filters:
 - filter:
   subject: socradar
   label: socradar_incidend
 - filter:
   subject: ".+Mysql Daily Full Backup.+"
   label: mysql_backup
`

	cfg, err := UnMarshallConfig([]byte(configdata))
	if err != nil {
		t.Error(err)
	}

	t.Run("correct query value?", func(t *testing.T) {
		if cfg.Query != "in:inbox is:unread" {
			t.Errorf("want: (in:inbox is:unread), got: (%s)", cfg.Query)
		}
	})

	t.Run("filters array lenght", func(t *testing.T) {
		if len(cfg.Filters) != 2 {
			t.Errorf("want: 2, got: %d", len(cfg.Filters))
		}
	})

	t.Run("chech array item value", func(t *testing.T) {
		l := cfg.Filters[1].Label
		if l != "mysql_backup" {
			t.Errorf("want: mysql_backup, got: %s", l)
		}
	})
}

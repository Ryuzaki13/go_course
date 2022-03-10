package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Menu struct {
	Parent int32
	Name   []string
	Link   []string
}

func prepareMenu() {
	var e error
	query["SelectMenu"], e = Link.Prepare(`
SELECT "Parent",
       jsonb_agg("Name" ORDER BY "Position"),
       jsonb_agg("Link" ORDER BY "Position")
FROM "Menu"
GROUP BY "Parent"
ORDER BY "Parent"`)
	if e != nil {
		fmt.Println(e)
	}
}

func SelectMenu() (list []Menu) {
	rows, e := query["SelectMenu"].Query()
	if e != nil {
		fmt.Println(e)
		return nil
	}

	defer rows.Close()

	list = make([]Menu, 0)

	m := Menu{}

	var parent sql.NullInt32
	var name, link json.RawMessage

	for rows.Next() {
		e = rows.Scan(&parent, &name, &link)
		if e != nil {
			fmt.Println(e)
			return nil
		}

		m.Parent = parent.Int32

		e = json.Unmarshal(name, &m.Name)
		if e != nil {
			fmt.Println(e)
		}

		e = json.Unmarshal(link, &m.Link)
		if e != nil {
			fmt.Println(e)
		}

		list = append(list, m)
	}

	return list
}

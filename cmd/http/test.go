package main

import (
	"fmt"

	"angel-castaneda.com/dblog"
)

func injectTag(name, cat, desc string) {
	test := dblog.Tag{
		Name:        name,
		Category:    cat,
		Description: desc,
	}
	fmt.Printf("%s\n%s\n%s\n", test.Name, test.Category, test.Description)

	if err := dblog.AddTag(test); err != nil {
		fmt.Println(err.Error())
	}
}

func injectPost() {
	test := dblog.Post{
		Title:    "new site... again",
		FileName: "new-site-again",
		Content: `
    cool
`,
		Description: "Justifying why I spent two months remaking my site as a cover to learn go.",
		PubDate:     "2023-08-07",
		UpdateDate:  "2023-08-07",
		Thumbnail: dblog.Img{
			Src: "",
		},
		Tags: []dblog.Tag{
			{Name: "code"},
			{Name: "articles"},
			{Name: "updates"},
			{Name: "english"},
		},
	}
	fmt.Printf("%s\n%s\n%s\n%s\n%s\n", test.Title, test.FileName, test.Content, test.Description, test.PubDate)

	if err := dblog.AddPost(test); err != nil {
		fmt.Println(err.Error())
	}
}

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// HTMLInfo contains all page info.
type HTMLInfo struct {
	NicknameInfo       template.HTML
	NicknameInfoString string
	CheckBoxInfo       template.HTML
	CheckBoxInfoString string
	LinkInfo           template.HTML
	LinkInfoString     string
}

// Builds page.
func page(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		addAnswers(rw, r)
		return
	}

	var formPage *template.Template = template.Must(template.ParseFiles("templates/page.html"))
	pageInfo := addCheckBoxesAndNickname(NewRequesterContainer(""), "")
	formPage.Execute(rw, pageInfo)
}

// Adds check boxes on page.
func addCheckBoxesAndNickname(container *RequesterContainer, nickname string) *HTMLInfo {
	var pageInfo *HTMLInfo = new(HTMLInfo)
	var isChecked string
	for _, page := range Pages {
		isChecked = ""
		if container.Requesters[page.ID].Available {
			isChecked = "checked"
		}
		pageInfo.CheckBoxInfoString += fmt.Sprintf(`
            <li>
                <label for="%s">%s</label>
                <input type="checkbox" name="%s" id="%s"%s/>
            </li>
		`, page.ID, page.Name, page.ID, page.ID, isChecked)
	}
	log.Println(pageInfo)
	pageInfo.CheckBoxInfo = template.HTML(pageInfo.CheckBoxInfoString)
	pageInfo.NicknameInfo = template.HTML(nickname)
	return pageInfo
}

// Checking which textboxes are set on and creating container of user info then getting answer.
func getUsedLinks(r *http.Request, container *RequesterContainer) []*UserInfo {
	for key, _ := range r.Form {
		if _, ok := container.Requesters[key]; ok {
			container.Requesters[key] = &RequesterAvailability{
				container.Requesters[key].requester,
				true,
			}
			fmt.Println("OK")
		}
	}

	return container.GetLinks()
}

// Adds answers to page.
func addAnswers(rw http.ResponseWriter, r *http.Request) {
	var answerPage *template.Template = template.Must(template.ParseFiles("templates/page.html"))

	r.ParseForm()
	nick := r.FormValue("nickname")
	log.Println(r.Form)

	container := NewRequesterContainer(nick)
	users := getUsedLinks(r, container)
	log.Println(users)

	pageInfo := addCheckBoxesAndNickname(container, nick)
	if nick == "" {
		pageInfo.LinkInfoString = "<h3>Looks like the nickname is invalid...</h3>\n\t\t<ul>\n"
		log.Println(pageInfo)
		pageInfo.LinkInfo = template.HTML(pageInfo.LinkInfoString)
		answerPage.Execute(rw, pageInfo)
		log.Println(answerPage)
		return
	}

	if len(users) == 0 {
		pageInfo.LinkInfoString = "<h3>Looks like you didn't select any pages...</h3>\n\t\t<ul>\n"
		log.Println(pageInfo)
		pageInfo.LinkInfo = template.HTML(pageInfo.LinkInfoString)
		answerPage.Execute(rw, pageInfo)
		log.Println(answerPage)
		return
	}

	pageInfo.LinkInfoString = fmt.Sprintf("<h3>Results for nickname \"%s\":</h3>\n\t\t<ul>\n", nick)

	for _, user := range users {
		if user.IsAvailable {
			pageInfo.LinkInfoString += fmt.Sprintf("\t\t\t<li>\n\t\t\t\t<a name=\"%s\" href=\"%s\">%s: %s</a>\n\t\t\t</li>\t\n", user.SocialNetwork, user.Link, user.SocialNetwork, user.Name)
		} else {
			pageInfo.LinkInfoString += fmt.Sprintf("\t\t\t<li>\n\t\t\t\t<a name=\"%s\">%s: %s</a>\n\t\t\t</li>\t\n", user.SocialNetwork, user.SocialNetwork, user.Link)
		}
	}

	pageInfo.LinkInfoString += "\t\t</ul>"
	log.Println(pageInfo)
	pageInfo.LinkInfo = template.HTML(pageInfo.LinkInfoString)
	answerPage.Execute(rw, pageInfo)
	log.Println(answerPage)
}

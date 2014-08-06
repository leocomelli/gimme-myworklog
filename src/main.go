package main

import (
  "flag"
  "log"
  "net/http"
  "io/ioutil"
  "encoding/xml"
  "encoding/json"
  "fmt"
)

func main(){

  /* Worklog */
  type Worklog []struct {
    IssueId string `xml:"issue_id"`
    IssueKey string `xml:"issue_key"`
    Username string `xml:"username"`
    StaffId string `xml:"staff_id"`
    WorkDescription string `xml:"work_description"`
    Reporter string `xml:"reporter"`
    Hours float64 `xml:"hours"`
    WorkDate string `xml:"work_date"`
    WorkDateTime string `xml:"work_date_time"`
  }

  type Worklogs struct {
    Worklogs Worklog `xml:"worklog"`
  }

  /* Issue */
  type CustomerField struct {
    Value string
  }

  type Project struct {
    Name string
  }

  type Field struct {
    IssueProject Project `json:"project"`
    Customer CustomerField `json:"customfield_10204"`
    Budget string `json:"customfield_10019"`
  }

  type Issue struct {
    Fields Field `json:"fields"`
  }

  var url, username, password, startDate, endDate, tempoApiToken string
  flag.StringVar(&url, "url", "", "Jira url")
  flag.StringVar(&username, "username", "", "Jira username")
  flag.StringVar(&password, "password", "", "Jira password")
  flag.StringVar(&tempoApiToken, "tempoApiToken", "", "Token of Jira Templo plugin")
  flag.StringVar(&startDate, "startDate", "", "Start date")
  flag.StringVar(&endDate, "endDate", "", "End date")
  flag.Parse()

  w := Worklogs{}
  worklogsAsXml := listWorklog(url, startDate, endDate, tempoApiToken)
  xml.Unmarshal(worklogsAsXml, &w)
  
  for _, worklog := range w.Worklogs {    
    issueData := getIssue(url, username, password, worklog.IssueKey)
    fmt.Printf("\t%s\n", worklog)
    //fmt.Printf("\t%s\n", issueData)

   issue := Issue{}
   json.Unmarshal(issueData, &issue)
   fmt.Printf("Issue: %v\n", issue)
  }
}

func listWorklog(url, startDate, endDate, tempoApiToken string) []byte {
  wls, err := http.Get(url +"/plugins/servlet/tempo-getWorklog/?dateFrom=" + startDate + "&dateTo=" + endDate + "&format=xml&diffOnly=false&tempoApiToken=" + tempoApiToken)
  if err != nil {
	log.Fatal(err)
  }
  content, err := ioutil.ReadAll(wls.Body)
  wls.Body.Close()
  if err != nil {
	log.Fatal(err)
  }
  return []byte(content)
}

func getIssue(url, username, password, issueKey string) []byte {
  client := &http.Client{}

  req, _ := http.NewRequest("GET", url + "/rest/api/2/issue/" + issueKey, nil)
  req.SetBasicAuth(username, password)
  resp, _ := client.Do(req)

  issue, err := ioutil.ReadAll(resp.Body)
  resp.Body.Close()
  if err != nil {
	log.Fatal(err)
  }

  return []byte(issue)
}
package main

import (
  "os"
  "fmt"
  "log"
  "flag"  
  "net/http"
  "io/ioutil"  
  "encoding/xml"
  "encoding/json"  
  "encoding/csv"  
)

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
  Project string
  Customer string
  Budget string
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

func main(){

  var url, username, password, startDate, endDate, tempoApiToken, output string
  flag.StringVar(&url, "url", "", "Jira url")
  flag.StringVar(&username, "username", "", "Jira username")
  flag.StringVar(&password, "password", "", "Jira password")
  flag.StringVar(&tempoApiToken, "tempoApiToken", "", "Token of Jira Templo plugin")
  flag.StringVar(&startDate, "startDate", "", "Start date")
  flag.StringVar(&endDate, "endDate", "", "End date")
  flag.StringVar(&output, "output", "worklog.csv", "Output filename")
  flag.Parse()

  w := Worklogs{}
  worklogsAsXml := listWorklog(url, startDate, endDate, tempoApiToken)
  xml.Unmarshal(worklogsAsXml, &w)
  
  //for _, worklog := range w.Worklogs {    
  for i := range w.Worklogs {	
    issueData := getIssue(url, username, password, w.Worklogs[i].IssueKey)
    
    issue := Issue{}
    json.Unmarshal(issueData, &issue)
    
    w.Worklogs[i].Project = issue.Fields.IssueProject.Name
    w.Worklogs[i].Customer = issue.Fields.Customer.Value
    w.Worklogs[i].Budget = issue.Fields.Budget
  }

  w.writeCsv(output)
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

func (wls *Worklogs) writeCsv(filename string){
  file, err := os.Create(filename)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()
  
  writer := csv.NewWriter(file)
  writer.Comma = ';'

  for _, obj := range wls.Worklogs {
    var line []string
    line = append(line, obj.IssueId)
    line = append(line, obj.IssueKey)
    line = append(line, obj.Username)
    line = append(line, obj.StaffId)
    line = append(line, obj.WorkDescription)
    line = append(line, obj.Reporter)
    line = append(line, fmt.Sprintf("%g", obj.Hours))
    line = append(line, obj.WorkDate)
    line = append(line, obj.WorkDateTime)
    line = append(line, obj.Project)
    line = append(line, obj.Customer)
    line = append(line, obj.Budget)
    writer.Write(line)
  }
  writer.Flush()
}



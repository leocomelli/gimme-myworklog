package main

import (
  //"os"
  "flag"
  "log"
  "net/http"
  "io/ioutil"
  "encoding/xml"
  "fmt"
)

func main(){

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

  var url, startDate, endDate, tempoApiToken string
  flag.StringVar(&url, "url", "", "Jira url")
  flag.StringVar(&tempoApiToken, "tempoApiToken", "", "Token of Jira Templo plugin")
  flag.StringVar(&startDate, "startDate", "", "Start date")
  flag.StringVar(&endDate, "endDate", "", "End date")
  flag.Parse()

  w := Worklogs{}
  worklogsAsXml := listWorklog(url, startDate, endDate, tempoApiToken)
  xml.Unmarshal(worklogsAsXml, &w)
  
  for _, worklog := range w.Worklogs {
    fmt.Printf("\t%s\n", worklog)
  }


}

func listWorklog(url, startDate, endDate, tempoApiToken string) ([]byte) {
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
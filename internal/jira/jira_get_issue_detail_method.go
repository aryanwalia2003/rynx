package jira

type Jira_issue_result_struct struct {
	Key         string
	Summary     string
	Status      string
	Description string
	Assignee    string
	Reporter    string
}

func (j *Jira_client_struct) Jira_get_issue_detail_method(key string) (*Jira_issue_result_struct, error) {
	detail, err := j.fetch_issue_detail_method(key)
	if err != nil {
		return nil, err
	}
	return &Jira_issue_result_struct{
		Key:         detail.Key,
		Summary:     detail.Fields.Summary,
		Status:      detail.Fields.Status.Name,
		Description: detail.extract_description_util(),
		Assignee:    detail.Fields.Assignee.DisplayName,
		Reporter:    detail.Fields.Reporter.DisplayName,
	}, nil
}

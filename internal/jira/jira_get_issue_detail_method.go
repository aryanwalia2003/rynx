package jira

type Jira_issue_result_struct struct {
	Key     string
	Summary string
	Status  string
}

func (j *Jira_client_struct) Jira_get_issue_detail_method(key string) (*Jira_issue_result_struct, error) {
	detail, err := j.fetch_issue_detail_method(key)
	if err != nil {
		return nil, err
	}
	return &Jira_issue_result_struct{
		Key:     detail.Key,
		Summary: detail.Fields.Summary,
		Status:  detail.Fields.Status.Name,
	}, nil
}

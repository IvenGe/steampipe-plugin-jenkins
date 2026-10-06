package jenkins

import (
	"context"

	"github.com/IvenGe/gojenkins"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableJenkinsJob() *plugin.Table {
	return &plugin.Table{
		Name:        "jenkins_job",
		Description: "A user-configured description of work which Jenkins should perform. This table is a generic representation of a Jenkins Job.",
		Get: &plugin.GetConfig{
			Hydrate:    getJenkinsJob,
			KeyColumns: plugin.SingleColumn("full_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError(nil),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listJenkinsFolders,
			Hydrate:       listJenkinsJobs,
		},

		Columns: []*plugin.Column{
			{Name: "full_name", Type: proto.ColumnType_STRING, Description: "Unique key for the job."},
			{Name: "full_display_name", Type: proto.ColumnType_STRING, Description: "Human readable name of the job, including folder."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the job."},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "Human readable name of the job."},
			{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromField("URL"), Description: "Full URL to the job."},
			{Name: "class", Type: proto.ColumnType_STRING, Description: "Java class of the job type."},
			{Name: "properties", Type: proto.ColumnType_JSON, Description: "Properties of the job.", Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func listJenkinsJobs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	folder, ok := h.Item.(*gojenkins.Job)
	if !ok || folder == nil {
		return nil, nil
	}

	jobs, err := folder.GetInnerJobs(ctx)
	if err != nil {
		logger.Error("jenkins_job.listJenkinsJobs", "query_error", err)
		if isNotFoundErr(err) {
			return nil, nil
		}
		return nil, err
	}

	for _, job := range jobs {
		if job == nil || job.Raw == nil {
			continue
		}
		d.StreamListItem(ctx, job.Raw)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getJenkinsJob(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Debug("jenkins_job.getJenkinsJob")
	jobFullName := d.EqualsQualString("full_name")

	// Empty check for jobFullName
	if jobFullName == "" {
		return nil, nil
	}

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("jenkins_job.getJenkinsJob", "connect_error", err)
		return nil, err
	}

	jobName, jobParentNames := splitFullName(jobFullName)

	job, err := client.GetJob(ctx, jobName, jobParentNames...)
	if err != nil {
		logger.Error("jenkins_job.getJenkinsJob", "query_error", err)
		return nil, err
	}

	if job == nil || job.Raw == nil {
		return nil, nil
	}

	return job.Raw, nil
}

package jenkins

import (
	"context"

	"github.com/IvenGe/gojenkins"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableJenkinsFolder() *plugin.Table {
	return &plugin.Table{
		Name:        "jenkins_folder",
		Description: "A container that stores job projects in it. Unlike view, which is just a filter, a folder creates a separate namespace, so you can have multiple things of the same name as long as they are in different folders.",
		Get: &plugin.GetConfig{
			Hydrate:    getJenkinsFolder,
			KeyColumns: plugin.SingleColumn("full_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError(nil),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listJenkinsFolders,
		},

		Columns: []*plugin.Column{
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Raw.Description"), Description: "An optional description that can be added to the folder."},
			{Name: "display_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Raw.DisplayName"), Description: "Human readable name of the folder."},
			{Name: "full_display_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Raw.FullDisplayName"), Description: "Human readable name of the folder, including parent folder."},
			{Name: "full_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Raw.FullName"), Description: "Unique key for the folder."},
			{Name: "jobs", Type: proto.ColumnType_JSON, Transform: transform.FromField("Raw.Jobs"), Description: "Child jobs."},
			{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Raw.Name"), Description: "Name of the folder, without its parent folder."},
			{Name: "primary_view", Type: proto.ColumnType_JSON, Transform: transform.FromField("Raw.PrimaryView"), Description: "Main view of this folder."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Raw.DisplayName"), Description: titleDescription},
			{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Raw.URL"), Description: "Full URL to the folder."},
			{Name: "views", Type: proto.ColumnType_JSON, Transform: transform.FromField("Raw.Views"), Description: "Views set at this folder."},
		},
	}
}

//// Handle Jenkins root level as a folder

func handleRootFolder(client *gojenkins.Jenkins) gojenkins.Job {
	raw := &gojenkins.JobResponse{
		Description:     "Jenkins root",
		DisplayName:     "Root",
		FullDisplayName: "Root",
		FullName:        "/",
		Name:            "/",
		URL:             client.Server + "/",
	}
	if client.Raw != nil {
		raw.Jobs = client.Raw.Jobs
		raw.PrimaryView = (*gojenkins.ViewData)(&client.Raw.PrimaryView)
		raw.Views = client.Raw.Views
	}
	return gojenkins.Job{
		Jenkins: client,
		Base:    "/",
		Raw:     raw,
	}
}

//// Recursively find sub folders

func handleFolders(folders []*gojenkins.Job, ctx context.Context, d *plugin.QueryData) error {
	for _, folder := range folders {
		if folder == nil || folder.Raw == nil {
			continue
		}
		// Filter to Folder job type only
		if folder.Raw.Class != "com.cloudbees.hudson.plugins.folder.Folder" {
			continue
		}
		d.StreamListItem(ctx, folder)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil
		}

		childJobs, err := folder.GetInnerJobs(ctx)
		if err != nil {
			if isNotFoundErr(err) {
				continue
			}
			return err
		}
		if err := handleFolders(childJobs, ctx, d); err != nil {
			return err
		}
	}
	return nil
}

//// LIST FUNCTION

func listJenkinsFolders(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("jenkins_folder.listJenkinsFolders", "connect_error", err)
		return nil, err
	}

	rootFolder := handleRootFolder(client)
	d.StreamListItem(ctx, &rootFolder)

	if d.RowsRemaining(ctx) == 0 {
		return nil, nil
	}

	folders, err := client.GetAllJobs(ctx)
	if err != nil {
		logger.Error("jenkins_folder.listJenkinsFolders", "list_folders_error", err)
		if isNotFoundErr(err) {
			return nil, nil
		}
		return nil, err
	}

	if err := handleFolders(folders, ctx, d); err != nil {
		logger.Error("jenkins_folder.listJenkinsFolders", "list_child_folders_error", err)
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getJenkinsFolder(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Debug("jenkins_folder.getJenkinsFolder")
	folderFullName := d.EqualsQualString("full_name")

	// Empty check for folderFullName
	if folderFullName == "" {
		return nil, nil
	}

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("jenkins_folder.getJenkinsFolder", "connect_error", err)
		return nil, err
	}

	if folderFullName == "/" {
		rootFolder := handleRootFolder(client)
		return &rootFolder, nil
	}

	folderName, folderParentNames := splitFullName(folderFullName)

	folder, err := client.GetJob(ctx, folderName, folderParentNames...)
	if err != nil {
		logger.Error("jenkins_folder.getJenkinsFolder", "query_error", err)
		return nil, err
	}

	if folder == nil || folder.Raw == nil {
		return nil, nil
	}

	// Filter to Folder job type only
	if folder.Raw.Class != "com.cloudbees.hudson.plugins.folder.Folder" {
		return nil, nil
	}

	return folder, nil
}

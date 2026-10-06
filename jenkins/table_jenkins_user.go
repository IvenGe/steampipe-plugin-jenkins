package jenkins

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableJenkinsUser() *plugin.Table {
	return &plugin.Table{
		Name:        "jenkins_user",
		Description: "A Jenkins user account known to the controller.",
		List: &plugin.ListConfig{
			Hydrate: listJenkinsUsers,
		},

		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique user ID."},
			{Name: "full_name", Type: proto.ColumnType_STRING, Description: "User's full name."},
			{Name: "absolute_url", Type: proto.ColumnType_STRING, Description: "Absolute URL to the user profile."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("FullName"), Description: titleDescription},
		},
	}
}

//// LIST FUNCTION

func listJenkinsUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("jenkins_user.listJenkinsUsers", "connect_error", err)
		return nil, err
	}

	users, err := client.GetAllUsers(ctx)
	if err != nil {
		logger.Error("jenkins_user.listJenkinsUsers", "query_error", err)
		if isNotFoundErr(err) {
			return nil, nil
		}
		return nil, err
	}

	if users == nil || users.Raw == nil {
		return nil, nil
	}

	for _, user := range users.Raw.Users {
		d.StreamListItem(ctx, map[string]interface{}{
			"ID":          user.User.ID,
			"FullName":    user.User.FullName,
			"AbsoluteURL": user.User.AbsoluteURL,
		})

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package index

import (
	"context"

	errs "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	resourceconfig "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resource/config"
	mdutils "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/string/markdown"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IndexDataSource{}

func NewIndexDataSource() datasource.DataSource {
	return &IndexDataSource{}
}

// IndexDataSource defines the data source implementation.
type IndexDataSource struct {
	config *resourceconfig.ResourceConfig
}

// IndexDataSourceModel describes the data source data model.
type IndexDataSourceModel struct {
	Id         types.String `tfsdk:"id"`
	Database   types.String `tfsdk:"database"`
	Collection types.String `tfsdk:"collection"`
	IndexName  types.String `tfsdk:"index_name"`
	Field      types.String `tfsdk:"field"`
	Direction  types.Int64  `tfsdk:"direction"`
	Unique     types.Bool   `tfsdk:"unique"`
}

func (d *IndexDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_index"
}

func (d *IndexDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: mdutils.FormatResourceDescription(`
			This data source reads an index for single field in a collection 
			in a database on the MongoDB server.
		`),

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: mdutils.FormatSchemaDescription(
					`
						Resource identifier. 
						
						ID has a value with a format of the following: 
						
						%s
					`,
					mdutils.CodeBlock("", "databases/<database>/collections/<collection>/indexes/<index_name>"),
				),
			},
			"database": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the database to read the collection in.",
			},
			"collection": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the collection to read the index in.",
			},
			"index_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the index.",
			},
			"field": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the field to read the index on.",
			},
			"direction": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Direction of the index. 1 for ascending, -1 for descending.",
			},
			"unique": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If true, this index has a unique constraint.",
			},
		},
	}
}

func (d *IndexDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	config, diags := resourceconfig.FromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	d.config = config
}

func (d *IndexDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	client := mongoclient.New(ctx, d.config.ClientConfig).WithLogger(d.config.Logger)
	client.Run(func(client *mongoclient.MongoClient, err error) {
		if err != nil {
			resp.Diagnostics.Append(
				errs.NewMongoClientError(err).ToDiagnostic(),
			)
			return
		}

		var data IndexDataSourceModel

		// Read Terraform prior state data into the model
		resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Perform read operation
		resp.Diagnostics.Append(dataSourceRead(client, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Save updated data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	})
}
